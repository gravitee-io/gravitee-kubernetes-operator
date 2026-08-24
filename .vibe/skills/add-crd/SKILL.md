---
name: add-crd
description: Add a new Custom Resource Definition to the operator - define the type in api/v1alpha1, register it in the SchemeBuilder, implement the internal/core interfaces, run make generate manifests reference, then wire up the APIM client layer, the controller, the webhook and search indexers. Use when asked to add a new CRD, a new gravitee.io resource kind, or a new API type to GKO.
user-invocable: true
---

# add-crd

The steps for this skill live in `.agent/skills/add-crd.md`, relative to the repository
root. Read that file now and follow it; this wrapper carries no instructions of its own.

@.agent/skills/add-crd.md
