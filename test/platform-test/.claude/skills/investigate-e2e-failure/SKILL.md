---
name: investigate-e2e-failure
description: Work a failing Gravitee platform e2e test (Playwright, test/platform-test) from symptom to root cause - leaked-resource cascade, eventual consistency, version gating, cluster/infra trouble, or a genuine bug in APIM, GKO or the Terraform provider. Use when an e2e test fails locally or in CI, when a suite times out, when several unrelated tests fail at once, or when asked why a test is flaky.
allowed-tools: Bash, Read, Edit, Glob, Grep
---

# investigate-e2e-failure

A failure in this suite is one of five things. Work them **in order**: the
earlier causes produce symptoms that look exactly like the later ones, so
skipping ahead leads to fixing a test that was never wrong.

The rules this rests on are in
[`.agent/rules/e2e-test-authoring.md`](../../../.agent/rules/e2e-test-authoring.md).
Two of them govern this whole procedure:

> **Never fix a failure by raising a timeout or skipping the test.**
> **Always consider that the root cause is a bug in the component under test.**

Work from `test/platform-test/`.

---

## Step 0: gather

```bash
# The run's own output
open playwright-report/index.html          # HTML report, per-step traces
grep -c "<failure" playwright-results/results.xml

# Cluster state at the time of the failure
kubectl get pods -A
kubectl get apiv4definitions,applications,subscriptions -A
kubectl logs deployment/gko-controller-manager -n default --tail=200

# Reproduce the single test
npm run e2e -- --grep @GKO-NNNN
npm run e2e -- --grep @GKO-NNNN --provision-with terraform
```

Note **whether one test failed or many**, and **which failed first**. That single
fact routes Steps 1 and 4.

## Step 1: leaked resource cascade

**Symptom:** a wave of generic 30s timeouts across unrelated suites; tests that
pass in isolation fail in a full run.

The suite runs serially, single-worker, against **one shared APIM**, and
API/App names are a global namespace. One test that leaks a name makes every
later apply of that name collide with stuck state.

1. Find the **first** failure in the run, not the loudest one.
2. Check that test for missing safety-net cleanup: inline `finally` does **not**
   run on a Playwright timeout. `forEachProvisioner` provides the net; a
   hand-rolled test needs its own `afterEach`/`afterAll`.
3. Check cleanup order: subscriptions, then applications, then APIs. The
   admission webhook blocks deleting an application that still has subscriptions,
   so a wrong order leaves both behind.
4. Sweep the leftovers and re-run:

```bash
kubectl get apiv4definitions,apidefinitions,applications,subscriptions -A
kubectl delete apiv4definition,application,subscription -l gravitee.io/e2e=true -A
```

**Fix = add the missing safety net to the first failing test**, not a longer
timeout on the ones that cascaded.

## Step 2: eventual consistency

**Symptom:** one test fails intermittently, on an assertion that immediately
follows an apply; the value is correct when you check by hand afterwards.

`kubectl apply` and `terraform apply` both return before APIM and the gateway
have converged.

- Convert the single-shot assertion to `mapi.waitForX(...)`, `expect.poll(...)`,
  or the `poll()` util.
- **Combine polled checks atomically:** one `expect.poll(...).toMatchObject({...})`
  over every field. Polling one field and then re-fetching for the rest can
  observe two different states.
- A reconcile is triggered by re-`kubectl apply -f` on a modified CR file, not by
  `kubectl patch` or `kubectl annotate`. A test using the latter may simply never
  have triggered the change it is waiting for.
- `provisioned.view` is **not** a convergence wait after `update()`: GKO's
  condition `observedGeneration` can lag after a re-apply (GKO-2940), so a
  post-update read can return the pre-update `Accepted=True`. Use `mapi`.

## Step 3: version gate

**Symptom:** the test fails only on an older APIM line, or only in the
`--run-up-to-version` / upgrade jobs, with a "field unknown" or 400-shaped error.

The feature is not in the APIM version under test. The APIM image comes from the
`gravitee-io/gravitee` CircleCI orb, not from this repo.

- Tag the test `since("<version the feature shipped in>")` so
  `--run-up-to-version` skips it on older clusters. A `since` tag is per
  provisioner: a feature can land in each driver at a different version.
- **`since` is not a substitute for a skip.** Only use it when the version
  boundary is real and known.
- Check the local cluster's version before concluding:

```bash
kubectl get deploy -n gravitee -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.template.spec.containers[0].image}{"\n"}{end}'
```

## Step 4: infrastructure, not the test

**Symptom:** `ECONNRESET`, `SocketError`, a mid-run collapse, or several
consecutive failures that all hit the gateway or mAPI.

Check the platform's own health **before** blaming the test:

```bash
kubectl get pods -n gravitee
kubectl get pods -n gravitee -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.containerStatuses[*].restartCount}{"\t"}{.status.containerStatuses[*].lastState.terminated.reason}{"\n"}{end}'
kubectl describe pod -n gravitee -l app.kubernetes.io/name=gateway
```

- `OOMKilled` or a climbing `restartCount` on the gateway or management API
  means the suite starved the platform, not that the test is wrong. Say so, and
  raise the pod's memory rather than the test's timeout.
- `globalSetup` runs five pre-flight checks (management API, gateway, cluster,
  CRDs, operator deployment). A failure there is an environment problem;
  `e2e/README.md` maps each symptom to its fix.
- Terraform arms need `terraform` on `PATH` at the CI-pinned version, and pick up
  a locally built provider from the mirror if `scripts/build-tf-provider.sh` was
  run.

## Step 5: a bug in the component under test

If Steps 1 to 4 come back clean, the failure is evidence, not noise. Establish:

- **Which component:** GKO (operator), APIM backend / Automation API, gateway, or
  the Terraform provider. The `origin` label does not tell you: `KUBERNETES`
  means "written via the Automation API", which is the write path for **both**
  GKO and Terraform.
- **Does it reproduce through the other provisioner?** Running the same journey
  with `--provision-with terraform` separates an operator bug from a platform bug
  in one command. This is the single most useful diagnostic the suite has.
- **Minimal repro:** the smallest manifest or HCL that shows it, saved to a file
  so it can be attached to the ticket.
- **Expected vs actual**, with the operator log line or HTTP response body.

Watch for the recurring shapes: missing validation, silent failure, an optional
field that is functionally required, schema mismatch between the CRD and the
Automation API, and a successful apply/reconcile with broken runtime behaviour.

### Filing it

- Issue type **Private Bug** in project GKO, `Team` = **DevEx**
  (`customfield_10001`).
- **No** `Quality Management` / `GKO-Quality-Management` labels: those are for
  Stories about the e2e framework, not for bugs.
- Link the failing test to the bug as **"is tested by"**.
- Show the full ticket draft and wait for approval before creating it.

### While it is open

- If the behaviour is genuinely broken, `test.fixme` the affected test **with the
  ticket in the reason**, so the gap is visible and re-enabling it is findable.
- If an arm cannot be expressed at all, use `pending: { <provisioner>: "<reason + ticket>" }`
  on the scenario rather than deleting it.
- Never green-wash by weakening the assertion.

## Reporting back

State plainly which of the five it was, the evidence, and what changed. If the
answer is "a bug in the component under test", say that first: a working test
that exposes a real defect is a success, not a failure to be worked around.
