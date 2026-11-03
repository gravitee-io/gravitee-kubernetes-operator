# AMClientSecret

An `AMClientSecret` is a client secret of an `AMApplication`. Its key in AM is `{namespace}-{name}`; its value is
read from a `Secret` entry (`spec.secretRef`). GKO sends the value to AM and never writes it anywhere in the
cluster: not in the status, not in events. AM never returns it.

[client-secret.yml](client-secret.yml) creates a `Secret` and the `AMClientSecret` that references it. Replace the
placeholder value in the `Secret` first: it is too short on purpose, so AM rejects it.

## A client secret is immutable

AM refuses a new value for an existing secret. Changing the value in the referenced `Secret` is reported on the
`AMClientSecret` as `Accepted=False`, with AM's message and a hint to rotate the secret. `spec.applicationRef` and
`spec.secretRef` cannot be changed either.

## Rotate a client secret

1. Create a new `Secret` with a new value, and a new `AMClientSecret` that references it. The application now has
   both secrets, and clients can authenticate with either.
2. Move the clients to the new value.
3. Delete the old `AMClientSecret`, then its `Secret`.

## Deletion

Deleting an `AMClientSecret` deletes the secret in AM. GKO never deletes the referenced `Secret`: delete it
yourself once no `AMClientSecret` uses it.

Deleting the `AMApplication` deletes the application in AM with its secrets; the `AMClientSecret`s it owns are
garbage collected.
