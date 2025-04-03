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

import { fileURLToPath } from 'url';
import * as dotenv from 'dotenv';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

if (argv._.length !== 2) {
    console.error(`Usage: ${path.basename(__filename)} <ENDPOINT_PATH> <EXPECTED_STATUS_CODE>`);
    process.exit(1);
}

const endpointPath = argv._[0];
const expectedStatusCodeStr = argv._[1];
const expectedStatusCode = parseInt(expectedStatusCodeStr, 10);

if (isNaN(expectedStatusCode)) {
    console.error(`Error: Expected status code "${expectedStatusCodeStr}" is not a valid number.`);
    process.exit(1);
}

const envFilePath = path.resolve(__dirname, '..', '.env');

if (!fs.existsSync(envFilePath)) {
    console.error(`Error: .env file not found at ${envFilePath}.`);
    process.exit(1);
}

try {
    dotenv.config({ path: envFilePath });
} catch (error) {
    console.error(`Error loading .env file from ${envFilePath}:`, error);
    process.exit(1);
}

const apiGatewayBaseUrl = process.env.APIM_GATEWAY; 
if (!apiGatewayBaseUrl) {
    console.error("Error: APIM GATEWAY is not set in the .env file.");
    process.exit(1);
}

const cleanedBaseUrl = apiGatewayBaseUrl.replace(/\/+$/, ''); // Remove trailing slashes
const cleanedEndpointPath = endpointPath.replace(/^\/+/, ''); // Remove leading slashes
const url = `${cleanedBaseUrl}/${cleanedEndpointPath}`;

console.log(`Testing connection to: ${url}`);

try {
    const response = await fetch(url);
    const actualStatusCode = response.status;

    if (actualStatusCode !== expectedStatusCode) {
        console.error(`Test failed: Expected ${expectedStatusCode} but got ${actualStatusCode} when calling ${url}`);
        process.exit(1);
    } else {
        console.log(`Connection test passed: ${url} returned ${expectedStatusCode}`);
    }
} catch (error) {
    console.error(`Error during fetch request to ${url}:`);
    console.error(error); 
    process.exit(1);
}