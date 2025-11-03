#!/usr/bin/env zx
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
 *
 * Usage (npx --yes --quiet zx exportAndAssertApi.mjs [options]):
 *   --api_version     v2 or v4 (required)
 *   --api_id          API id (optional if --api_resource + --api_name provided)
 *   --api_resource    Kubernetes resource name (e.g. apidefinitions.gravitee.io)
 *   --api_name        Kubernetes resource instance to look up the API id
 *   --namespace       Namespace for the lookup (default: default)
 *   --assert          Assertion, repeatable. Format: path[:expectedValue]
 *   --assert_not      Negated assertion. Format: path[:unexpectedValue]
 *   --max_attempts    Retry attempts before failing (default: 30)
 *   --delay_ms        Delay between retries in milliseconds (default: 1000)
 *
 * Paths support dot/bracket notation (e.g. spec.pages.hello.accessControls[0].referenceId),
 * and values are optional. When omitted, the script only checks that the path exists/does not exist.
 * Values can be primitives (true, 42, text) or JSON objects/arrays.
 * Values containing spaces are supported without quoting; everything after the first ':'
 * until the next option flag is treated as the value. You may also wrap the value in
 * single or double quotes (they will be stripped automatically).
 *
 * Example:
 *   npx --yes --quiet zx exportAndAssertApi.mjs \
 *     --api_version v2 \
 *     --api_resource apidefinitions.gravitee.io \
 *     --api_name my-api \
 *     --namespace default \
 *     --assert spec.pages.overview.accessControls:'{"referenceId":"team-a","referenceType":"GROUP"}' \
 *     --assert_not spec.categories:deprecated \
 *     --assert_contains spec.pages.overview.content:This is an update \
 *     --assert_contains spec.pages.overview.content:"Another sentence with spaces"
 */

import { isDeepStrictEqual } from 'node:util';
import { parse as parseYamlDocument } from 'yaml';
import { mapiClient } from '../gravitee/mapi/client.mjs';

const apiVersion = argv['api_version'];
if (!apiVersion || !['v2', 'v4'].includes(apiVersion)) {
  console.error('Please set --api_version (v2 or v4).');
  process.exit(1);
}

const namespace = argv.namespace ?? 'default';
const apiResource = argv['api_resource'];
const apiName = argv['api_name'];
const providedApiId = argv['api_id'];

if (!providedApiId && (!apiResource || !apiName)) {
  console.error('Provide either --api_id or the trio --api_resource, --api_name, --namespace.');
  process.exit(1);
}

const maxAttempts = Number(argv['max_attempts']) || 10;
const delayMs = Number(argv['delay_ms']) || 1000;

// Parse assertion flags from raw argv tokens to preserve values containing spaces
const rawAssertionArgs = parseAssertionArgs(process.argv.slice(2));

const mustHave = collectAssertionsFromList(rawAssertionArgs.assert, undefined);
const mustNotHave = collectAssertionsFromList(rawAssertionArgs.assert_not, undefined);
const mustContain = collectAssertionsFromList(rawAssertionArgs.assert_contains, 'contains');
const mustNotContain = collectAssertionsFromList(rawAssertionArgs.assert_not_contains, 'not_contains');

if (!mustHave.length && !mustNotHave.length && !mustContain.length && !mustNotContain.length) {
  console.error('Nothing to check. Add at least one --assert or --assert_not rule.');
  process.exit(1);
}

let apiId = providedApiId;

for (let attempt = 1; attempt <= maxAttempts; attempt += 1) {
  try {
    if (!apiId) {
      apiId = await resolveApiId(apiResource, namespace, apiName);
    }

    const exportedYaml = await exportApiAsYaml(apiId, apiVersion);
    const exportedObject = await parseYaml(exportedYaml);

    const errors = [
      ...checkRules(exportedObject, mustHave, false),
      ...checkRules(exportedObject, mustNotHave, true),
      ...checkContains(exportedObject, mustContain, false),
      ...checkContains(exportedObject, mustNotContain, true),
    ];

    if (!errors.length) {
      console.log(`Assertions satisfied (attempt ${attempt}).`);
      process.exit(0);
    }

    logAttempt(attempt, maxAttempts, errors);
  } catch (error) {
    apiId = providedApiId;
    logAttempt(attempt, maxAttempts, [error.message ?? String(error)]);
  }

  if (attempt < maxAttempts) {
    await sleep(delayMs);
  }
}

console.error('Assertions were not met before the timeout.');
process.exit(1);

function collectAssertionsFromList(list, mode) {
  const entries = Array.isArray(list) ? list : list ? [list] : [];
  return entries
    .map((entry) => entry && entry.toString().trim())
    .filter(Boolean)
    .map((text) => ({ ...parseAssertion(text), mode }));
}

function parseAssertionArgs(tokens) {
  const result = {
    assert: [],
    assert_not: [],
    assert_contains: [],
    assert_not_contains: [],
  };

  const keys = new Set(['--assert', '--assert_not', '--assert_contains', '--assert_not_contains']);

  for (let i = 0; i < tokens.length; i += 1) {
    const t = tokens[i];
    if (!keys.has(t)) {
      continue;
    }

    const key = t.replace(/^--/, '');
    // Consume subsequent tokens that belong to this assertion until the next option (starts with --) or end
    let value = '';
    if (i + 1 < tokens.length) {
      i += 1;
      value = tokens[i] ?? '';
      while (i + 1 < tokens.length && !String(tokens[i + 1]).startsWith('--')) {
        i += 1;
        value += ` ${tokens[i]}`;
      }
    }

    if (value) {
      result[key].push(value);
    }
  }

  return result;
}

function parseAssertion(text) {
  const separatorIndex = text.indexOf(':');
  if (separatorIndex === -1) {
    return { path: text.trim(), value: undefined };
  }

  const path = text.slice(0, separatorIndex).trim();
  const rawValue = text.slice(separatorIndex + 1).trim();
  return {
    path,
    value: parseValue(rawValue),
  };
}

function parseValue(raw) {
  // Strip matching surrounding quotes to allow values containing spaces
  let processed = raw;
  if (
    typeof processed === 'string' &&
    processed.length >= 2 &&
    ((processed.startsWith('"') && processed.endsWith('"')) || (processed.startsWith("'") && processed.endsWith("'")))
  ) {
    processed = processed.slice(1, -1);
  }
  if (processed === '') {
    return '';
  }
  if (processed === 'true') {
    return true;
  }
  if (processed === 'false') {
    return false;
  }
  if (processed === 'null') {
    return null;
  }

  const numericValue = Number(processed);
  if (!Number.isNaN(numericValue)) {
    return numericValue;
  }

  if ((processed.startsWith('{') && processed.endsWith('}')) || (processed.startsWith('[') && processed.endsWith(']'))) {
    try {
      return JSON.parse(processed);
    } catch {
      return processed;
    }
  }

  return processed;
}

function checkRules(exportedObject, assertions, negate) {
  const issues = [];
  for (const assertion of assertions) {
    const actual = getByPath(exportedObject, assertion.path);
    const matches = valueMatches(actual, assertion.value);
    if (!negate && !matches) {
      issues.push(
        assertion.value === undefined
          ? `Expected ${assertion.path} to exist.`
          : `Expected ${assertion.path} to contain ${display(assertion.value)} but saw ${display(actual)}.`,
      );
    }
    if (negate && matches) {
      issues.push(
        assertion.value === undefined
          ? `Expected ${assertion.path} to be absent.`
          : `Expected ${assertion.path} to not contain ${display(assertion.value)}.`,
      );
    }
  }
  return issues;
}

function checkContains(exportedObject, assertions, negate) {
  const issues = [];
  for (const assertion of assertions) {
    const actual = getByPath(exportedObject, assertion.path);
    const expected = assertion.value;
    const ok = containsMatch(actual, expected);
    if (!negate && !ok) {
      issues.push(`Expected ${assertion.path} to contain substring ${display(expected)} but saw ${display(actual)}.`);
    }
    if (negate && ok) {
      issues.push(`Expected ${assertion.path} to not contain substring ${display(expected)}.`);
    }
  }
  return issues;
}

function containsMatch(actual, expected) {
  if (expected === undefined) {
    return false;
  }
  if (actual === undefined || actual === null) {
    return false;
  }
  if (typeof actual !== 'string') {
    return false;
  }
  return actual.includes(String(expected));
}

function valueMatches(actual, expected) {
  if (expected === undefined) {
    return actual !== undefined;
  }
  if (actual === undefined) {
    return expected === 0; // treat missing collections as empty when expecting zero
  }
  if (Array.isArray(actual)) {
    return arrayContains(actual, expected);
  }
  if (isPlainObject(actual) && isPlainObject(expected)) {
    return isSubset(expected, actual);
  }
  return isDeepStrictEqual(actual, expected);
}

function arrayContains(array, expected) {
  if (Array.isArray(expected)) {
    return expected.every((item) => arrayContains(array, item));
  }

  return array.some((candidate) => {
    if (Array.isArray(candidate)) {
      return Array.isArray(expected) && arrayContains(candidate, expected);
    }
    if (isPlainObject(expected) && isPlainObject(candidate)) {
      return isSubset(expected, candidate);
    }
    return isDeepStrictEqual(candidate, expected);
  });
}

function isSubset(expected, actual) {
  return Object.entries(expected).every(([key, value]) => {
    const actualValue = actual[key];
    if (isPlainObject(value) && isPlainObject(actualValue)) {
      return isSubset(value, actualValue);
    }
    if (Array.isArray(value) && Array.isArray(actualValue)) {
      return value.every((entry) => arrayContains(actualValue, entry));
    }
    if (Array.isArray(value)) {
      return false;
    }
    if (Array.isArray(actualValue)) {
      return arrayContains(actualValue, value);
    }
    return isDeepStrictEqual(actualValue, value);
  });
}

function isPlainObject(value) {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function getByPath(root, path) {
  const segments = path.replace(/\[(\d+)\]/g, '.$1').split('.');
  let current = root;

  for (const segment of segments) {
    if (!segment) {
      continue;
    }
    if (current === null || current === undefined) {
      return undefined;
    }

    if (segment === 'length') {
      if (!Array.isArray(current)) {
        return undefined;
      }
      current = current.length;
      continue;
    }

    const numeric = Number(segment);
    if (Array.isArray(current) && !Number.isNaN(numeric)) {
      current = current[numeric];
      continue;
    }

    current = current[segment];
  }

  return current;
}

async function resolveApiId(resource, resourceNamespace, name) {
  const result = await $`kubectl get ${resource} -n ${resourceNamespace} ${name} -o jsonpath='{.status.id}'`.nothrow();
  if (result.exitCode !== 0) {
    throw new Error(result.stderr.trim() || `Failed to resolve API id for ${name}.`);
  }

  const apiId = result.stdout.trim();
  if (!apiId) {
    throw new Error(`API id not yet available for ${name}.`);
  }

  return apiId;
}

async function exportApiAsYaml(apiId, version) {
  const path =
    version === 'v2'
      ? `/management/organizations/DEFAULT/environments/DEFAULT/apis/${apiId}/crd`
      : `/management/v2/environments/DEFAULT/apis/${apiId}/_export/crd`;
  const { body } = await mapiClient.get(path);
  return body;
}

async function parseYaml(yamlString) {
  if (!yamlString || !yamlString.trim()) {
    throw new Error('Exported payload was empty.');
  }

  try {
    return parseYamlDocument(yamlString);
  } catch (error) {
    throw new Error(`Could not parse exported YAML: ${error.message}`);
  }
}

function display(value) {
  if (typeof value === 'string') {
    return `"${value}"`;
  }
  if (value === undefined) {
    return 'undefined';
  }
  try {
    return JSON.stringify(value);
  } catch {
    return String(value);
  }
}

function logAttempt(attempt, maxAttempts, errors) {
  console.error(`Attempt ${attempt}/${maxAttempts} failed:`);
  errors.forEach((err) => console.error(` - ${err}`));
}
