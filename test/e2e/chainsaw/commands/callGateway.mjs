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

await dotenv.config('.env');
const { endpoint: endpointPath, status: expectedStatusCode } = argv;

if (!endpointPath || !expectedStatusCode) {
    console.error("Usage: callGateway.mjs --endpoint <endpointPath> --status <statusCode>");
    process.exit(1);
}

const apiGatewayBaseUrl = process.env.APIM_GATEWAY;
if (!apiGatewayBaseUrl) {
    console.error("Error: APIM_GATEWAY is not set.");
    process.exit(1);
}

const base = new URL(apiGatewayBaseUrl);
base.pathname = path.posix.join(base.pathname, endpointPath);
const url = base.toString();

console.log(`Testing connection to: ${url}`);

try {
    const response = await fetch(url);
    const actualStatusCode = response.status;

    if (actualStatusCode !== parseInt(expectedStatusCode, 10)) {
        console.error(`Test failed: Expected ${expectedStatusCode} but got ${actualStatusCode} when calling ${url}`);
        process.exit(1);
    } else {
        console.log(`Connection test passed: ${url} returned ${expectedStatusCode}`);
    }
} catch (error) {
    console.error(`Error during fetch request to ${url}: ${error.message}`);
    process.exit(1);
}
