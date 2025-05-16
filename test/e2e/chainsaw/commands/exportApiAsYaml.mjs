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


await dotenv.config(`${__dirname}/../.env`);

const { APIM_API, APIM_AUTH } = process.env;
if (!APIM_API || !APIM_AUTH) {
    console.error('Error: APIM_API or APIM_AUTH not set in .env file.');
    process.exit(1);
}

const apiId = argv.api_id;
const apiVersion = argv.api_version;

if (!apiId) {
    console.error('Error: --api_id parameter is not provided.');
    console.error(`Usage: ${path.basename(process.argv[1])} --api_id API_ID --api_version v2|v4`);
    process.exit(2);
}

if (!apiVersion || !['v2', 'v4'].includes(apiVersion)) {
    console.error('Error: --api_version parameter is not provided or is invalid. Must be v2 or v4.');
    console.error(`Usage: ${path.basename(process.argv[1])} --api_id API_ID --api_version v2|v4`);
    process.exit(3);
}

let url;
if (apiVersion === 'v2') {
    url = `${APIM_API}/management/organizations/DEFAULT/environments/DEFAULT/apis/${apiId}/crd`;
} else if (apiVersion === 'v4') {
    url = `${APIM_API}/management/v2/environments/DEFAULT/apis/${apiId}/_export/crd`;
}

try {
    const headers = {
        'Authorization': `Bearer ${APIM_AUTH}`,
        'Content-Type': 'application/yaml'
    };

    const response = await fetch(url, {
        method: 'GET',
        headers: headers,
    });

    if (!response.ok) {
        let errorBody = '';
        try {
            errorBody = await response.text();
        } catch (bodyError) {
            errorBody = '(Could not read error response body)';
        }
        console.error(`Error: API request failed with status ${response.status} ${response.statusText}`);
        if (errorBody) {
            console.error(`Response body:\n${errorBody}`);
        }
        process.exit(4);
    }

    const responseBody = await response.text();
    console.log(responseBody);

} catch (error) {
    console.error(`Error: Failed to fetch API CRD for ID: ${apiId}.`);
    console.error(error);
    process.exit(5);
}