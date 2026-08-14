# Serving two variants of an API on the same context path

We want the same context path (`/weather-api`) served by two different backends, on two
different gateway pools, for two different consumer groups. This guide shows how to do it
with **virtual hosts**, and how to shift traffic between the two variants.

## Why sharding tags are not enough

Sharding tags control **where an API is deployed** — which gateway pool picks it up. They do
not create separate namespaces in APIM.

Context path uniqueness is validated **per environment**, before tags are ever considered.
Two APIs in the same environment declaring `/weather-api` are rejected with
`Path [/weather-api/] already exists`, even if they carry disjoint tags and never meet on the
same gateway.

The conflict rule is also *prefix containment*, not equality: `/weather-api` collides with
`/weather-api/v2` too. So carving out sub-paths does not help either.

## Why virtual hosts work

The uniqueness check runs **per host**. Paths declared with a `host` are only compared against
other paths on that same host; paths without a `host` are only compared against other hostless
paths. Give each variant its own hostname and the collision disappears:

| API              | Host                            | Path            | Result |
| ---------------- | ------------------------------- | --------------- | ------ |
| `weather-api-v1` | `consumer-group-a.example.com`  | `/weather-api`  | ✅      |
| `weather-api-v2` | `consumer-group-b.example.com`  | `/weather-api`  | ✅      |

At runtime the gateway matches on the `Host` header: a listener path with a `host` only accepts
requests carrying that exact host (case-insensitive), and a listener path without a `host`
accepts anything. Consumers reach their variant through their own hostname, on the identical
path.

Sharding tags stay in the picture — they still pin each API to its gateway pool. Virtual hosts
solve the *declaration* conflict, tags solve the *placement*.

## The manifests

See [`api.yaml`](api.yaml). The relevant parts:

```yaml
spec:
  tags: ["consumer-group-a"]        # which gateway pool deploys it
  listeners:
    - type: HTTP
      paths:
        - path: "/weather-api"
          host: "consumer-group-a.example.com"   # what makes it unique
```

```yaml
spec:
  tags: ["consumer-group-b"]
  listeners:
    - type: HTTP
      paths:
        - path: "/weather-api"
          host: "consumer-group-b.example.com"
```

Apply them:

```bash
kubectl apply -f api.yaml
```

Both are accepted, both expose `/weather-api`, each on its own gateway pool with its own
backend.

## Before you apply

- Both hostnames must resolve to the right gateway pool (or to a load balancer
  in front of it), and the certificate served for each host must cover it. A wildcard on
  `*.example.com` covers both.
- Once one API on a path uses a virtual host, give *every* API on that path a
  virtual host. A hostless `/weather-api` matches any Host header and will shadow the others on
  a shared gateway.

## Rolling traffic between the two variants

The discriminator is a header, so any L7 load balancer in front of the gateways can move
consumers between variants without touching APIM at all. Clients keep calling one stable
public hostname; the LB decides which variant answers by rewriting `Host`.

```
                              ┌── Host: consumer-group-a.example.com ──▶  Gateway pool A ──▶  weather-api-v1
Client ──▶ api.example.com ──▶ LB
                              └── Host: consumer-group-b.example.com ──▶  Gateway pool B ──▶  weather-api-v2
```

### Option 1 — no LB, cohort-by-hostname (simplest)

Each consumer group is handed its own base URL and stays there. Migrating a group is a
one-line change on their side, rollback is instant, and there is no shared component to
configure. This is the right default when the split follows organisational boundaries, which
is the case here.

### Option 2 — weighted shift behind one hostname

When you want a percentage rollout, put an LB in front and let it set the Host header. With
nginx:

```nginx
upstream gw_a { server gateway-pool-a.internal:8082; }
upstream gw_b { server gateway-pool-b.internal:8082; }

# Sticky per client: the same key always lands on the same variant.
split_clients "${http_x_consumer_id}" $variant {
    10%   b;
    *     a;
}

server {
    listen 443 ssl;
    server_name api.example.com;

    location /weather-api {
        proxy_set_header Host consumer-group-$variant.example.com;
        proxy_pass http://gw_$variant;
    }
}
```

Raise the percentage to shift traffic, set it back to `*  a` to roll back. `split_clients`
hashes the key, so a given consumer stays on the same variant across requests instead of
flapping — use a stable key (consumer id, API key, tenant header), not `$remote_addr`, if you
can.

An ingress controller with canary weights, or a service mesh with weighted routes, does the
same thing — the only requirement is that it rewrites `Host` to the target variant's virtual
host before forwarding.

### Option 3 — pinned cohorts, then flip

Match on a header or cookie to pin named accounts to the new variant first, let the rest fall
through to the default, then invert the default once you are confident. This gives you early
adopters without a percentage dial.

Whichever option you pick, the rollback path is the same: change the routing rule in front.
The APIs themselves stay deployed and untouched, so reverting takes effect immediately and
carries no risk of a failed redeploy.
