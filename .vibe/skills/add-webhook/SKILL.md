---
name: add-webhook
description: Add validating or defaulting admission webhooks for a CRD - internal/admission/<resource>/ctrl.go and validate.go, the generic admission.Validator/Defaulter wiring, the required order of checks (templates, context ref, spec, cluster, APIM dry-run, drift last), the Helm webhook manifests, registration in main.go under ENABLE_WEBHOOK, and where the tests go. Use when asked to validate a CRD on create or update, reject invalid specs, default missing fields, or add an admission webhook.
user-invocable: true
---

# add-webhook

The steps for this skill live in `.agent/skills/add-webhook.md`, relative to the
repository root. Read that file now and follow it; this wrapper carries no instructions
of its own.

@.agent/skills/add-webhook.md
