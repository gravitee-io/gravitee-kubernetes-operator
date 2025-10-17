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

const apiId = argv.api_id;
const orgId = argv.org_id || 'DEFAULT';
const envId = argv.env_id || 'DEFAULT';
const expectedStatus = argv.status ? Number(argv.status) : undefined;

if (!apiId) {
  console.error('Error: --api_id parameter is not provided.');
  console.error(`Usage: ${path.basename(process.argv[1])} --api_id API_ID [--org_id ORG_ID] [--env_id ENV_ID] [--status HTTP_STATUS]`);
  process.exit(1);
}

const apiPath = `/management/organizations/${orgId}/environments/${envId}/apis/${apiId}`;

try {
  // mapiClient.get returns response text and throws on non-2xx with message containing status.
  const resText = await mapiClient.get(apiPath);
  const actualStatus = 200; // success path
  if (expectedStatus !== undefined) {
    if (actualStatus !== expectedStatus) {
      console.error(`Expected HTTP status ${expectedStatus} but got ${actualStatus}.`);
      process.exit(1);
    }
  }
  console.log(resText);
} catch (error) {
  // If status checking is requested, verify it against error message.
  const match = /\b(\d{3})\b/.exec(error.message);
  const actual = match ? Number(match[1]) : undefined;
  if (expectedStatus !== undefined) {
    if (actual === expectedStatus) {
      // Print the error message to stdout for visibility when expected error occurs
      console.log(error.message);
      process.exit(0);
    }
    console.error(`Expected HTTP status ${expectedStatus}${actual ? ` but got ${actual}` : ''}.`);
  } else {
    console.error(`Error: Failed to fetch API by ID: ${apiId}.`);
    console.error(error.message);
  }
  process.exit(1);
}
