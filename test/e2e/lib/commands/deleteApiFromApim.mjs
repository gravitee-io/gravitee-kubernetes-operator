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

import { mapiClient } from '../gravitee/mapi/client.mjs';

const { api_id: apiId, api_version: apiVersion, org_id: orgId = 'DEFAULT', env_id: envId = 'DEFAULT' } = argv;

if (!apiId || !apiVersion) {
  console.error('Usage: deleteApiFromApim.mjs --api_id <API_ID> --api_version <v2|v4> [--org_id <ORG>] [--env_id <ENV>]');
  process.exit(1);
}

let endpoint;
switch (apiVersion) {
  case 'v2':
    endpoint = `/management/organizations/${orgId}/environments/${envId}/apis/${apiId}?closePlans=true`;
    break;
  case 'v4':
    endpoint = `/management/v2/environments/${envId}/apis/${apiId}?closePlans=true`;
    break;
  default:
    console.error(`Unsupported api_version '${apiVersion}'. Expected 'v2' or 'v4'.`);
    process.exit(1);
}

try {
  await mapiClient.del(endpoint);
  console.log(`API '${apiId}' deleted from APIM (${apiVersion}).`);
} catch (error) {
  console.error(`Failed to delete API '${apiId}': ${error.message}`);
  process.exit(1);
}
