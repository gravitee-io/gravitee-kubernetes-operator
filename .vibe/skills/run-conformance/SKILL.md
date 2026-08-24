---
name: run-conformance
description: Run the Gateway API conformance suite locally against the kind cluster - verify the cluster shape, clean up leftovers, load the gateway and coverage-instrumented operator images, check the cloud-provider-kind load balancer, bootstrap in the required order, and report real results. Use before reporting any change under controllers/gateway-api/, internal/k8s/gateway*.go or test/conformance/ as done, working or verified - build, lint and unit tests never execute a reconcile and are not evidence for gateway-api changes.
user-invocable: true
---

# run-conformance

The steps for this skill live in `.agent/skills/run-conformance.md`, relative to the
repository root. Read that file now and follow it; this wrapper carries no instructions
of its own. It also depends on "Verifying a fix for a flaky test" in
`.agent/rules/gateway-standards.md`.

@.agent/skills/run-conformance.md
