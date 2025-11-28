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
 */

const resource = argv.resource;
const name = argv.name;
const namespace = argv.namespace ?? 'default';
const timeoutSeconds = Number(argv.timeout ?? 60);
const intervalSeconds = Number(argv.interval ?? 2);

if (!resource || !name) {
  console.error(
    'Usage: waitResourceDeletion.mjs --resource <RESOURCE> --name <NAME> [--namespace <NAMESPACE>] [--timeout <SECONDS>] [--interval <SECONDS>]',
  );
  process.exit(1);
}

const deadline = Date.now() + timeoutSeconds * 1000;

while (Date.now() < deadline) {
  const result = await $`kubectl get ${resource} -n ${namespace} ${name}`.nothrow();
  if (result.exitCode !== 0) {
    console.log(`Resource ${resource}/${name} does not exist in namespace ${namespace}.`);
    process.exit(0);
  }

  console.log(`Resource ${resource}/${name} still present. Retrying in ${intervalSeconds}s...`);
  await new Promise((resolve) => setTimeout(resolve, intervalSeconds * 1000));
}

console.error(`Timed out waiting for ${resource}/${name} deletion in namespace ${namespace}.`);
process.exit(1);
