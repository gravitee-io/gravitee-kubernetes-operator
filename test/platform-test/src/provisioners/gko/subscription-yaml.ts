/**
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * Builders for GKO Subscription / Secret manifests, used by the GKO provisioner
 * to apply parameterized subscriptions (api-key sets, rotations) that cannot be
 * expressed as a static fixture. Generalized from the inline builders in the
 * original v4-subscriptions-apikey pilot.
 *
 * Single-quoted YAML scalars are used for key values so values containing
 * backticks (e.g. the GKO templating syntax `[[ secret `name/key` ]]`)
 * round-trip cleanly.
 */

export interface ApiKeyEntry {
  key: string;
  expireAt?: string;
}

export interface SubscriptionYamlOptions {
  /** metadata.name of the Subscription CR. */
  name: string;
  /** spec.api.name (the target API CR name). */
  apiName: string;
  /** spec.api.kind. Defaults to "ApiV4Definition". */
  apiKind?: string;
  /** spec.plan (the plan key, e.g. "ApiKey"). */
  plan: string;
  /** spec.application.name. */
  applicationName: string;
  /**
   * spec.apiKeys entries. When empty/omitted the apiKeys block is left out
   * entirely so APIM auto-generates a single key.
   */
  keys?: ApiKeyEntry[];
}

/** Build a Subscription manifest targeting a plan on an API + application. */
export function subscriptionYaml(opts: SubscriptionYamlOptions): string {
  const { name, apiName, apiKind = "ApiV4Definition", plan, applicationName, keys = [] } = opts;

  const lines = [
    "apiVersion: gravitee.io/v1alpha1",
    "kind: Subscription",
    "metadata:",
    `  name: ${name}`,
    "spec:",
    "  api:",
    `    name: "${apiName}"`,
    `    kind: "${apiKind}"`,
    `  plan: "${plan}"`,
    "  application:",
    `    name: "${applicationName}"`,
  ];

  if (keys.length > 0) {
    lines.push("  apiKeys:");
    for (const k of keys) {
      const escaped = k.key.replaceAll("'", "''");
      lines.push(`    - key: '${escaped}'`);
      if (k.expireAt) lines.push(`      expireAt: "${k.expireAt}"`);
    }
  }

  lines.push("");
  return lines.join("\n");
}

/** Build a Secret manifest carrying a single `apiKey` field. */
export function apiKeySecretYaml(name: string, value: string): string {
  return [
    "apiVersion: v1",
    "kind: Secret",
    "metadata:",
    `  name: ${name}`,
    "type: Opaque",
    "stringData:",
    `  apiKey: '${value.replaceAll("'", "''")}'`,
    "",
  ].join("\n");
}
