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

const { name, description, orgId = 'DEFAULT', envId = 'DEFAULT' } = argv;
if (!name || typeof name !== 'string' || name.length < 1) {
    console.error('Usage: createCategory.mjs --name <CATEGORY_NAME> [--description <DESCRIPTION>] [--orgId <ORG_ID>] [--envId <ENV_ID>]');
    process.exit(1);
}

const body = { name };

const endpoint = `/management/organizations/${orgId}/environments/${envId}/configuration/categories`;

try {
    const createdCategory = await mapiClient.post(endpoint, body);
    console.log(JSON.stringify(createdCategory));
} catch (e) {
    console.error(`Failed to create category: ${e.message}`);
    process.exit(1);
}
