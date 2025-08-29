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

import { mapiClient } from '../lib/mapiClient.mjs';

const apiId = argv.api_id;
const orgId = argv.org_id || 'DEFAULT';
const envId = argv.env_id || 'DEFAULT';

if (!apiId) {
	console.error('Error: --api_id parameter is not provided.');
	console.error(`Usage: ${path.basename(process.argv[1])} --api_id API_ID [--org_id ORG_ID] [--env_id ENV_ID]`);
	process.exit(1);
}

const notificationsPath = `/management/organizations/${orgId}/environments/${envId}/apis/${apiId}/notificationsettings`;

try {
	const notifications = await mapiClient.get(notificationsPath);
	console.log(notifications);
} catch (error) {
	console.error(`Error: Failed to fetch notification settings for API ID: ${apiId}.`);
	console.error(error.message);
	process.exit(1);
}

