# AM domain, identity provider, certificate and reporter: manual runbook

Live tests of `AMSecurityDomain`, `AMIdentityProvider`, `AMCertificate` and `AMReporter` against a real AM. Only
what the envtest suites and the AM mock server cannot show is here: AM's own validation and messages, the values
AM builds or resolves (system resources, masked secrets, loaded keystores, `expiresAt`), AM's cascades and
refusals, and drift against AM's real responses. Kubernetes-side behaviour (CEL, finalizers, owner references,
watches, apply order, Secret rotation, the fallback delete guard) is covered by the envtest suites and is not
repeated here.

Data comes from `domain-comprehensive` (Terraform). Keys AM sees: domain `am-runbook-comprehensive`, identity
providers `am-runbook-users-idp` (system) and `am-runbook-inline-users`, certificates `am-runbook-system-cert`
(system), `am-runbook-signing-p12` and `am-runbook-signing-jks`, reporters `am-runbook-system-reporter` (system),
`am-runbook-audit-file` and `am-runbook-audit-kafka`.

## Setup

Prerequisites: GKO from the branch under test, with the webhook enabled (`manager.webhook.enabled=true`), and
an AM whose Automation API is reachable from your laptop.

```bash
export AM_URL=http://localhost:30093    # AM Management API from your laptop, e.g. kubectl -n am port-forward svc/am-management-api 8093:83
export AM_TOKEN='NGMxNzk1OTQtNjI5ZS00YmZkLTk3OTUtOTQ2MjllOWJmZGUxLm1RS0RPelo2dFVCUm1EaTlxRmVObGotMHNFeHZGbUxVaVpmelQ5ZmlXZHM='
A=$AM_URL/automation/organizations/DEFAULT/environments/DEFAULT
D=am-runbook-comprehensive

am()     { curl -s -H "Authorization: Bearer $AM_TOKEN" "$A$1" "${@:2}"; }                          # GET, extra curl args allowed
am_put() { curl -s -X PUT -H "Authorization: Bearer $AM_TOKEN" -H 'Content-Type: application/json' --data @- "$A$1"; }
am_del() { curl -s -o /dev/null -w '%{http_code}\n' -X DELETE -H "Authorization: Bearer $AM_TOKEN" "$A$1"; }
st()     { kubectl -n am-runbook get "$@" -o jsonpath='{range .status.conditions[*]}{.type}={.status} {.message}{"\n"}{end}'; }
```

Edit `baseUrl` in `01-am-context.yaml` (AM as GKO reaches it from inside the cluster), then:

```bash
kubectl apply -f 00-namespace.yaml
kubectl -n am-runbook create secret generic am-context-credentials --from-literal=bearerToken="$AM_TOKEN"
kubectl apply -f 01-am-context.yaml
st amcontext am-ctx                      # Accepted=True
```

A real keystore for the certificates (the mock stores any bytes; AM loads them). `keytool` needs a JDK:

```bash
STOREPASS='Gko-Store-P@ss'
openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -days 30 -subj /CN=gko-runbook
openssl pkcs12 -export -in cert.pem -inkey key.pem -name signing -out keystore.p12 -passout pass:$STOREPASS
keytool -importkeystore -srckeystore keystore.p12 -srcstoretype PKCS12 -srcstorepass $STOREPASS \
  -destkeystore keystore.jks -deststoretype JKS -deststorepass $STOREPASS -destkeypass $STOREPASS
kubectl -n am-runbook create secret generic signing-keystore --from-literal=storepass=$STOREPASS \
  --from-literal=p12.b64="$(base64 < keystore.p12 | tr -d '\n')" --from-literal=jks.b64="$(base64 < keystore.jks | tr -d '\n')"
openssl x509 -in cert.pem -noout -enddate                # expiresAt AM must report
```

## 1. Create

### 1.1 Comprehensive domain

```bash
kubectl apply -f 02-domain.yaml
st amsecuritydomain comprehensive        # Accepted=True
am /domains/$D | jq '{key, name, path, enabled, tags, cors: .corsSettings.enabled, csp: .webProtectionSettings.csp.reportOnly, xframe: .webProtectionSettings.xframe.action, maxDelegationDepth: .tokenExchangeSettings.maxDelegationDepth}'
```

Check every section made it: `am /domains/$D | jq 'keys'`, and spot-check nested values against `02-domain.yaml`
(`accountSettings.resetPasswordCustomFormFields`, `oidc.workloadIdentitySettings`, `vhosts`).

### 1.2 System identity provider

```bash
kubectl apply -f 03-idp-system.yaml
st amidentityprovider users-idp          # Accepted=True
am /domains/$D/identities/am-runbook-users-idp | jq
```

AM builds it from its default settings: check `type` and `configuration` come from AM, not from the CR.

### 1.3 Inline identity provider, password from a Secret

```bash
kubectl apply -f 04-idp-inline-secret.yaml -f 05-idp-inline.yaml
st amidentityprovider inline-users       # Accepted=True
am /domains/$D/identities/am-runbook-inline-users | jq '{name, type, configuration: (.configuration | fromjson)}'
```

`configuration` is sent as a JSON string; the password must be the Secret's value (or AM's masked/hashed form),
never `[[ secret ... ]]`.

### 1.4 Registration IdP on the domain

Needs the IdP to exist in AM first, so it is a second step:

```bash
kubectl -n am-runbook patch amsecuritydomain comprehensive --type merge \
  -p '{"spec":{"accountSettings":{"defaultIdentityProviderForRegistration":"am-runbook-inline-users"}}}'
st amsecuritydomain comprehensive
am /domains/$D | jq .accountSettings.defaultIdentityProviderForRegistration
```

## 2. Update

```bash
# Domain: scalar, list and nested values
kubectl -n am-runbook patch amsecuritydomain comprehensive --type merge \
  -p '{"spec":{"description":"Updated by GKO","tags":["gko","drift-test","updated"],"passwordSettings":{"minLength":12}}}'
am /domains/$D | jq '{description, tags, minLength: .passwordSettings.minLength}'

# Domain: remove an optional section from the CR. Does AM keep the old value or reset it?
kubectl -n am-runbook patch amsecuritydomain comprehensive --type json -p '[{"op":"remove","path":"/spec/scim"}]'
am /domains/$D | jq .scim

# IdP: rename and add a user
kubectl -n am-runbook patch amidentityprovider inline-users --type merge -p '{"spec":{"name":"[GKO] Inline users v2"}}'
kubectl -n am-runbook patch amidentityprovider inline-users --type json -p '[{"op":"add","path":"/spec/configuration/users/-","value":{"firstname":"Carol","lastname":"QA","username":"carol","password":"Gko-Qa-C4rol-P@ss"}}]'
am /domains/$D/identities/am-runbook-inline-users | jq '{name, users: (.configuration | fromjson | .users | map(.username))}'

# IdP: rotate the Secret. No CR change: the templating watch must push the new password.
kubectl -n am-runbook create secret generic inline-users --from-literal=alice-password='Gko-Qa-R0tated!' --dry-run=client -o yaml | kubectl apply -f -
am /domains/$D/identities/am-runbook-inline-users | jq '.configuration | fromjson | .users[0]'
```

## 3. Dry-run rejections

Each `kubectl apply` must fail with AM's message, and nothing must exist afterwards, neither in Kubernetes
nor in AM. Note the message: it is what a user will read.

| File | Expected |
|---|---|
| `dryrun/domain-unknown-registration-idp.yaml` | rejected: registration IdP is not an IdP of the domain |
| `dryrun/domain-vhost-mode-without-vhosts.yaml` | rejected: vhostMode needs vhosts |
| `dryrun/domain-duplicate-path.yaml` | rejected: path already used by `comprehensive` |
| `dryrun/domain-unknown-data-plane.yaml` | rejected: unknown data plane |
| `dryrun/domain-bad-oidc-uri.yaml` | rejected: invalid post-logout redirect URI |
| `dryrun/idp-unknown-type.yaml` | rejected: unknown plugin type |
| `dryrun/idp-invalid-configuration.yaml` | rejected: configuration does not match the plugin's schema |
| `dryrun/certificate-wrong-storepass.yaml` | rejected: the keystore cannot be opened |
| `dryrun/certificate-unknown-alias.yaml` | rejected: no key under that alias |
| `dryrun/certificate-not-a-keystore.yaml` | rejected: not a PKCS#12 file |
| `dryrun/certificate-keystore-as-object.yaml` | rejected: the file field must be a JSON string (the AM OAS example shows an object) |
| `dryrun/certificate-unknown-type.yaml` | rejected: unknown plugin type |
| `dryrun/reporter-unknown-type.yaml` | rejected: unknown plugin type |
| `dryrun/reporter-invalid-configuration.yaml` | rejected: `retainDays` is not a number. If admitted, note it: AM does not check the reporter configuration |
| `dryrun/reporter-bad-expression.yaml` | note it: AM rejects the unparsable expression, or stores it and fails at audit time |

`dryrun/after-system/` needs the system certificate and reporter first: run it in 6.1 and 7.2.

```bash
for f in dryrun/*.yaml; do echo "== $f"; kubectl apply -f "$f"; done
kubectl -n am-runbook get amsecuritydomains,amidentityproviders,amcertificates,amreporters    # nothing from dryrun/
am /domains | jq                                                    # no new domain
```

An AM-side rule GKO lets through only at reconcile time shows up as `Accepted=False` instead: note which.

Updates are dry-run too:

```bash
kubectl -n am-runbook patch amsecuritydomain comprehensive --type merge -p '{"spec":{"vhostMode":true,"vhosts":null}}'   # rejected, AM unchanged
kubectl -n am-runbook patch amidentityprovider inline-users --type merge -p '{"spec":{"configuration":{"users":"alice"}}}' # rejected, AM unchanged
```

## 4. Apply order

An IdP applied before its domain is admitted with a warning and created once AM has the domain.

```bash
kubectl apply -f order/idp-first.yaml                      # Warning: domain ... not found or not yet created in AM
st amidentityprovider early-users                          # ResolvedRefs=False
kubectl apply -f order/domain-late.yaml
st amidentityprovider early-users                          # Accepted=True within seconds
am /domains/am-runbook-late/identities/am-runbook-early-users | jq .name
```

## 5. Drift

Drift is on for `comprehensive`, `users-idp` and `inline-users` (annotation `gravitee.io/drift-detection: "true"`).
"Out of band" below means a change made in AM directly: the AM console, or the Automation API with `am_put`.

### 5.1 Baseline: no drift without a remote change

```bash
kubectl -n am-runbook patch amsecuritydomain comprehensive --type merge -p '{"spec":{"description":"no drift"}}'  # admitted
kubectl -n am-runbook patch amidentityprovider inline-users --type merge -p '{"spec":{"name":"[GKO] no drift"}}'  # admitted
```

If either is rejected, the diff shows a false positive: an AM default or normalisation that drift must ignore.
A diff naming a masked secret (the inline password) is a regression of section 8.

### 5.2 Remote change on the domain

Out of band, enable CORS (a "drift probe" field):

```bash
am /domains/$D | jq '.corsSettings.enabled = true' | am_put /domains | jq .corsSettings.enabled
```

Then:

```bash
kubectl -n am-runbook patch amsecuritydomain comprehensive --type merge -p '{"spec":{"description":"after remote change"}}'
# rejected, the diff names corsSettings.enabled
kubectl -n am-runbook patch amsecuritydomain comprehensive --type merge -p '{"spec":{"corsSettings":{"enabled":true}}}'
# admitted: the CR realigns with AM
```

Repeat with other probes from `02-domain.yaml` (`accountSettings.loginAttemptsDetectionEnabled`,
`loginSettings.registerEnabled`, `oidc.redirectUriStrictMatching`, `webProtectionSettings.csp.reportOnly`), and
with a list (`tags`) and a nested list (`oidc.postLogoutRedirectUris`).

### 5.3 Remote change on an identity provider

```bash
am /domains/$D/identities/am-runbook-inline-users | jq '.name = "renamed in AM"' | am_put /domains/$D/identities | jq .name
kubectl -n am-runbook patch amidentityprovider inline-users --type json -p '[{"op":"replace","path":"/spec/configuration/users/0/lastname","value":"Changed"}]'
# rejected, the diff names name
```

Also try a change inside `configuration` out of band (add a user): the diff must point at `configuration`.

### 5.4 Remote missing

```bash
am_del /domains/$D/identities/am-runbook-inline-users             # 204
kubectl -n am-runbook patch amidentityprovider inline-users --type merge -p '{"spec":{"name":"[GKO] Inline users v3"}}'
# rejected: remote not found (DRIFT_DETECTION_ON_REMOTE_MISSING=deny)
kubectl -n am-runbook annotate amidentityprovider inline-users gravitee.io/drift-detection=false --overwrite
kubectl -n am-runbook patch amidentityprovider inline-users --type merge -p '{"spec":{"name":"[GKO] Inline users v3"}}'
am /domains/$D/identities/am-runbook-inline-users | jq .name    # recreated by the reconcile
kubectl -n am-runbook annotate amidentityprovider inline-users gravitee.io/drift-detection=true --overwrite
```

## 6. Certificates

### 6.1 Create: AM loads the keystore and reports what it builds

```bash
kubectl apply -f 06-certificate-system.yaml           # Warning: 'name' will be ignored when 'system' is 'true'
kubectl apply -f 07-certificate-pkcs12.yaml -f 08-certificate-jks.yaml
kubectl apply -f dryrun/after-system/certificate-second-system.yaml   # rejected: the domain already has one
kubectl -n am-runbook get amcertificates               # NAME and TYPE columns show AM's values; EXPIRES set
am /domains/$D/certificates | jq 'map({key, name, type, system, expiresAt})'
```

- `system-cert`: NAME and TYPE are AM's (built from `domains.certificates.default.*`), not `Ignored by AM`.
- `signing-p12` and `signing-jks`: `expiresAt` in the status matches `openssl x509 -enddate` above.
- `am /domains/$D/certificates/am-runbook-signing-p12 | jq '.configuration | fromjson'`: `storepass`,
  `keypass` and the file field come back as `********`, never the template nor the clear value.

### 6.2 Keystore rotation

The mock never computes `expiresAt`: only AM shows the new keystore was loaded.

```bash
openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -days 90 -subj /CN=gko-runbook
openssl pkcs12 -export -in cert.pem -inkey key.pem -name signing -out keystore.p12 -passout pass:$STOREPASS
kubectl -n am-runbook create secret generic signing-keystore --from-literal=storepass=$STOREPASS \
  --from-literal=p12.b64="$(base64 < keystore.p12 | tr -d '\n')" --from-literal=jks.b64="$(base64 < keystore.jks | tr -d '\n')" \
  --dry-run=client -o yaml | kubectl apply -f -
kubectl -n am-runbook get amcertificate signing-p12 -o jsonpath='{.status.expiresAt}{"\n"}'   # moved ~60 days later
```

### 6.3 Passwords resolved by AM (optional)

Needs AM's Kubernetes secret provider on the management API **and** the gateway
(`domains.secrets.providers: [{plugin: kubernetes, configuration: {enabled: true}}]`) and AM's service accounts
allowed to `get` Secrets in `am-runbook`.

```bash
kubectl apply -f 09-certificate-am-secrets.yaml        # admitted: the dry-run loads the keystore with the resolved passwords
am /domains/$D/certificates/am-runbook-signing-am | jq '.configuration | fromjson | {storepass, keypass}'
```

Without the provider, the apply must be rejected with AM's message: note it.

### 6.4 Referenced by the domain

AM does not validate these keys at domain PUT time; check it uses them.

```bash
kubectl -n am-runbook patch amsecuritydomain comprehensive --type merge \
  -p '{"spec":{"certificateSettings":{"fallbackCertificate":"am-runbook-signing-p12"},"saml":{"certificate":"am-runbook-signing-jks"}}}'
st amsecuritydomain comprehensive
am /domains/$D | jq '{fallback: .certificateSettings.fallbackCertificate, saml: .saml.certificate}'
```

### 6.5 Drift

```bash
# System certificate: name set in the CR, ignored by AM. Must be admitted.
kubectl -n am-runbook patch amcertificate system-cert --type merge -p '{"spec":{"name":"Still ignored"}}'

# Remote change of a field AM does not mask: the diff must name it.
am /domains/$D/certificates/am-runbook-signing-jks | jq '.configuration |= (fromjson | .algorithm = "RS512" | tojson)' \
  | am_put /domains/$D/certificates | jq .key
kubectl -n am-runbook annotate amcertificate signing-jks gravitee.io/drift-detection=true --overwrite
# rejected, the diff names configuration.algorithm only (no masked field), and the annotation is not applied

# Put AM back, so section 8.2 starts from a CR aligned with AM.
am /domains/$D/certificates/am-runbook-signing-jks | jq '.configuration |= (fromjson | .algorithm = "RS256" | tojson)' \
  | am_put /domains/$D/certificates | jq .key
```

The masked `storepass`, `keypass` and keystore file are covered with the other kinds in section 8.

### 6.6 Delete refused by AM

The admission guard only knows the fallback certificate. The SAML certificate (6.4) is a use only AM sees:

```bash
kubectl -n am-runbook delete amcertificate signing-jks --wait=false
st amcertificate signing-jks                           # Accepted=False with AM's refusal, CR Terminating
am /domains/$D/certificates/am-runbook-signing-jks | jq .key   # still there
kubectl -n am-runbook patch amsecuritydomain comprehensive --type json -p '[{"op":"remove","path":"/spec/saml/certificate"}]'
kubectl -n am-runbook get amcertificate signing-jks    # gone after the next retry
```

## 7. Reporters

### 7.1 No default reporter on a GKO domain

The mock has no default reporter at all; AM creates one for console domains only.

```bash
am /domains/$D/reporters | jq 'map(.key)'               # [] before 10-reporter-system.yaml
```

In the AM console (domain > Settings > Audit > Reporters), the domain must show no reporter either. Note
whether its audit log stays empty after a domain change (`kubectl patch` the description).

### 7.2 System reporter: AM builds it and ignores the CR's fields

```bash
kubectl apply -f 10-reporter-system.yaml
# Warning: 'name', 'attributeMappings', 'attributeMappingEventTypes' will be ignored when 'system' is 'true'.
kubectl -n am-runbook get amreporter system-reporter       # NAME and TYPE columns show AM's values
am /domains/$D/reporters/am-runbook-system-reporter | jq '{name, type, system, attributeMappings, attributeMappingEventTypes, configuration}'
kubectl apply -f dryrun/after-system/reporter-second-system.yaml   # rejected: the domain already has one
```

- NAME and TYPE are AM's (built from `domains.reporters.default.*`), not `Ignored by AM`.
- `attributeMappings` and `attributeMappingEventTypes` are empty in AM.
- The second system reporter is rejected at admission by the dry-run. If it is admitted and fails at reconcile
  (`Accepted=False`) instead, note it: AM only checks it on a real write.
- The console audit log of 7.1 now fills up.

### 7.3 File reporter with attribute mappings

```bash
kubectl apply -f 11-reporter-file.yaml
st amreporter audit-file                                   # Accepted=True
am /domains/$D/reporters/am-runbook-audit-file | jq '{name, type, enabled, attributeMappings, attributeMappingEventTypes, configuration: (.configuration | fromjson)}'
```

`configuration` comes back as an object of the CR's values (`retainDays` a number, not a string), the mappings
and the event types as sent.

### 7.4 Kafka reporter, password from a Secret

Needs the `reporter-am-kafka` plugin in AM. No broker is needed to store the reporter; note whether the dry-run
tries to reach `bootstrapServers` (a rejection naming the broker means it does).

```bash
kubectl -n am-runbook create secret generic audit-kafka --from-literal=password='Gko-Kafka-P@ss'
kubectl apply -f 12-reporter-kafka-secret.yaml
st amreporter audit-kafka                                  # Accepted=True
am /domains/$D/reporters/am-runbook-audit-kafka | jq '.configuration | fromjson | .password'   # "********"
```

### 7.5 Password resolved by AM (optional)

Same AM secret provider set-up as 6.3, plus `get` on Secrets in `am-runbook` for AM's service accounts.

```bash
kubectl apply -f 13-reporter-kafka-am-secrets.yaml         # admitted: the dry-run loads the reporter
am /domains/$D/reporters/am-runbook-audit-kafka-am | jq '.configuration | fromjson | .password'
```

Note what AM returns: the `{#secrets.get(...)}` expression as sent, or `********`. Without the provider, the
apply must be rejected with AM's message: note it.

### 7.6 Drift

```bash
# System reporter: name and attribute mappings set in the CR, ignored by AM. Must be admitted.
kubectl -n am-runbook patch amreporter system-reporter --type merge -p '{"spec":{"name":"Still ignored"}}'

# File reporter: remote change of a configuration field.
am /domains/$D/reporters/am-runbook-audit-file | jq '.configuration |= (fromjson | .retainDays = 30 | tojson)' \
  | am_put /domains/$D/reporters | jq .key
kubectl -n am-runbook patch amreporter audit-file --type merge -p '{"spec":{"name":"[GKO] Audit to file v2"}}'
# rejected, the diff names configuration.retainDays
kubectl -n am-runbook patch amreporter audit-file --type merge -p '{"spec":{"configuration":{"retainDays":30}}}'
# admitted: the CR realigns with AM

# File reporter: disabled, then attribute mappings changed, in the AM console. Each time the diff must name
# the field (enabled, attributeMappings), and realigning the CR must be admitted.
```

The Kafka password is covered with the other kinds in section 8.

## 8. Masked values and drift (all resources)

AM returns every configuration field its plugin flags as sensitive as `********`. The certificate, identity
provider and reporter `configuration` are tagged `drift:"unstructured:masked"` (am-sdk v2.3.1): a value AM
returns as exactly `********` is equivalent to a string (or nothing) in the CR, at any depth; a CR object or list
against `********` drifts, and every other configuration value is still compared. The mock does not mask, so only a real AM shows which fields are masked and in what
shape. A value masked in any other form (a hash, `****`, a partly masked JSON string) is still a drift: that
is what this section looks for.

The domain has no plugin configuration and is not concerned.

### 8.1 Which fields AM masks

```bash
am /domains/$D/identities/am-runbook-inline-users | jq '.configuration | fromjson | .users[0]'
am /domains/$D/certificates/am-runbook-signing-p12 | jq '.configuration | fromjson'
am /domains/$D/certificates/am-runbook-signing-jks | jq '.configuration | fromjson'
am /domains/$D/reporters/am-runbook-audit-kafka | jq '.configuration | fromjson'
am /domains/$D/certificates/am-runbook-signing-am | jq '.configuration | fromjson'    # optional, 6.3
am /domains/$D/reporters/am-runbook-audit-kafka-am | jq '.configuration | fromjson'   # optional, 7.5
```

Fill in the table. Every masked field must come back as exactly `********`, as a whole value.

| Resource | Field | AM returns |
|---|---|---|
| `inline-users` | `users[].password` | `********`, a hash, or the clear value? |
| `signing-p12` | `storepass`, `keypass`, keystore file (`content`) | |
| `signing-jks` | `storepass`, `keypass`, keystore file | |
| `audit-kafka` | `password` | |
| `signing-am`, `audit-kafka-am` | the `{#secrets.get(...)}` fields | the expression, or `********`? |
| system resources | their configuration (built by AM) | |

### 8.2 The false positive, one kind at a time

Drift is on for `inline-users`, `signing-p12` and `audit-kafka`; turn it on for `signing-jks`, whose keystore
file field is `jks`, not `content`. Change a field that is neither masked nor changed in AM:

```bash
kubectl -n am-runbook annotate amcertificate signing-jks gravitee.io/drift-detection=true --overwrite
kubectl -n am-runbook patch amidentityprovider inline-users --type merge -p '{"spec":{"name":"[GKO] masked v2"}}'
kubectl -n am-runbook patch amcertificate signing-p12 --type merge -p '{"spec":{"name":"[GKO] masked v2"}}'
kubectl -n am-runbook patch amcertificate signing-jks --type merge -p '{"spec":{"name":"[GKO] masked v2"}}'
kubectl -n am-runbook patch amreporter audit-kafka --type merge -p '{"spec":{"name":"[GKO] masked v2"}}'
```

All four must be admitted. A rejection whose diff names a field of 8.1 means AM masks it in a shape
other than `********`: note the value AM returns. Any other field in the diff is a new false positive: note it.

Then with the AM secret provider variants (optional):

```bash
kubectl -n am-runbook annotate amcertificate signing-am gravitee.io/drift-detection=true --overwrite
kubectl -n am-runbook annotate amreporter audit-kafka-am gravitee.io/drift-detection=true --overwrite
kubectl -n am-runbook patch amcertificate signing-am --type merge -p '{"spec":{"name":"[GKO] masked v2"}}'
kubectl -n am-runbook patch amreporter audit-kafka-am --type merge -p '{"spec":{"name":"[GKO] masked v2"}}'
```

Both must be admitted, whether AM returns the expression as sent or `********` (8.1).

### 8.3 What AM does with `********` sent back

The OAS says `********` sent on an update keeps the stored value, and is refused on a create. GKO always sends
the resolved value, so drift does not depend on it; check it with a dry-run of AM's own response to record AM's
contract:

```bash
# Update: the GET response, masked values included, sent back as is. Expect no dryRunErrors.
am /domains/$D/reporters/am-runbook-audit-kafka | am_put "/domains/$D/reporters?dryRun=true" | jq .dryRunErrors
am /domains/$D/certificates/am-runbook-signing-p12 | am_put "/domains/$D/certificates?dryRun=true" | jq .dryRunErrors
# Create: the same body under a new key. Expect a refusal.
am /domains/$D/reporters/am-runbook-audit-kafka | jq '.key = "am-runbook-masked-create"' \
  | am_put "/domains/$D/reporters?dryRun=true" | jq '.dryRunErrors // .'
```

For the certificate, also send it back for real (no `dryRun`) and check `expiresAt` is unchanged: AM still
loads the keystore with the password it kept.

### 8.4 A masked field changed out of band

In the AM console, change the Kafka reporter's password, then:

```bash
am /domains/$D/reporters/am-runbook-audit-kafka | jq '.configuration | fromjson | .password'   # still "********"
```

Drift cannot see this change: AM returns the same `********`. Note it as expected
behaviour. The next reconcile (a CR or Secret change) puts the Secret's password back.

## 9. Delete and cascade

### 9.1 One identity provider

```bash
kubectl -n am-runbook delete amidentityprovider inline-users
am /domains/$D/identities/am-runbook-inline-users -o /dev/null -w '%{http_code}\n'   # 404
```

The domain still references it as registration IdP (step 1.4): note what AM does (refuse the delete, so the
CR stays `Terminating` with AM's message, or delete and clear the reference).

### 9.2 Reporters

```bash
kubectl -n am-runbook delete amreporter audit-file
am /domains/$D/reporters/am-runbook-audit-file -o /dev/null -w '%{http_code}\n'      # 404
kubectl -n am-runbook delete amreporter system-reporter --wait=false
st amreporter system-reporter
```

Note what AM does with the system reporter: delete it (404 afterwards, and the domain is back to no reporter, as
in 7.1), or refuse it (CR stays `Terminating` with AM's message).

### 9.3 AMContext still in use

```bash
kubectl -n am-runbook delete amcontext am-ctx --wait=false   # rejected: referenced by AM security domains
```

### 9.4 Domain deleted outside Kubernetes

```bash
am_del /domains/am-runbook-late                                # AM deletes its identity providers with it
kubectl -n am-runbook patch amidentityprovider early-users --type merge -p '{"spec":{"name":"[GKO] Early v2"}}'
st amidentityprovider early-users                              # Accepted=False with AM's 404
kubectl -n am-runbook patch amsecuritydomain late --type merge -p '{"spec":{"description":"recreate"}}'
am /domains/am-runbook-late | jq .key                          # recreated by the reconcile
```

Then check whether `early-users` comes back too (it only reconciles on its own change or the next resync).

### 9.5 Domain with identity providers, certificates and reporters

```bash
kubectl apply -f 04-idp-inline-secret.yaml -f 05-idp-inline.yaml        # bring inline-users back
kubectl -n am-runbook delete amsecuritydomain comprehensive
am /domains/$D -o /dev/null -w '%{http_code}\n'                          # 404
am /domains/$D/identities -o /dev/null -w '%{http_code}\n'               # 404: gone with the domain
am /domains/$D/certificates -o /dev/null -w '%{http_code}\n'             # 404: AM deletes them with force, fallback included
am /domains/$D/reporters -o /dev/null -w '%{http_code}\n'                # 404: gone with the domain
kubectl -n am-runbook get amidentityproviders,amcertificates,amreporters # all garbage collected, none stuck Terminating
```

### 9.6 Cleanup

```bash
kubectl -n am-runbook delete amsecuritydomain late
kubectl -n am-runbook delete amcontext am-ctx                  # admitted now
kubectl delete namespace am-runbook
am /domains | jq                                                # no am-runbook-* domain left
```
