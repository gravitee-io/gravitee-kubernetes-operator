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
const apiVersion = argv.api_version;

if (!apiId) {
    console.error('Error: --api_id parameter is not provided.');
    console.error(`Usage: ${path.basename(process.argv[1])} --api_id API_ID --api_version v2|v4`);
    process.exit(1);
}

if (!apiVersion || !['v2', 'v4'].includes(apiVersion)) {
    console.error('Error: --api_version parameter is not provided or is invalid. Must be v2 or v4.');
    console.error(`Usage: ${path.basename(process.argv[1])} --api_id API_ID --api_version v2|v4`);
    process.exit(1);
}

let apiCrdExportPath;
if (apiVersion === 'v2') {
    apiCrdExportPath = `/management/organizations/DEFAULT/environments/DEFAULT/apis/${apiId}/crd`;
} else if (apiVersion === 'v4') {
    apiCrdExportPath = `/management/v2/environments/DEFAULT/apis/${apiId}/_export/crd`;
}

try {
    const crdExport = await mapiClient.get(apiCrdExportPath);
    console.log(crdExport);
} catch (error) {
    console.error(`Error: Failed to fetch CRD for API ID: ${apiId}.`);
    console.error(error.message);
    process.exit(1);
}
