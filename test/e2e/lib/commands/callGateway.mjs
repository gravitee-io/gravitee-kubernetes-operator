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

import fs from 'fs';
import https from 'https';
import fetch from 'node-fetch'; // native fetch is not used because it does not support https.Agent config (needed for mTLS/client certs)

const {
  endpoint: endpointPath,
  status: expectedStatusCode,
  notStatus: notExpectedStatusCode,
  gateway,
  cert,
  key,
  cacert,
  authorization,
  header: headerArgs,
} = argv;

if (!endpointPath || (!expectedStatusCode && !notExpectedStatusCode) || (expectedStatusCode && notExpectedStatusCode)) {
  console.error(
    'Usage: callGateway.mjs --endpoint <endpointPath> (--status <statusCode> | --notStatus <statusCode>) [--gateway <baseUrl>] [--cert <certPem>] [--key <keyPem>] [--cacert <caCertPem>] [--authorization <value>] [--header "Key: Value"]',
  );
  process.exit(1);
}

const defaultGateway = 'http://localhost:30082';
const baseUrl = (gateway || defaultGateway).replace(/\/$/, '');
const normalizedEndpoint = String(endpointPath).replace(/^\//, '');
const url = `${baseUrl}/${normalizedEndpoint}`;

const maxRetry = 60;
const retryDelay = 500;

const expected = expectedStatusCode ? parseInt(expectedStatusCode, 10) : undefined;
const notExpected = notExpectedStatusCode ? parseInt(notExpectedStatusCode, 10) : undefined;

// Optional headers (Authorization or custom headers)
const requestHeaders = {};

if (authorization) {
  requestHeaders.Authorization = String(authorization);
}

if (headerArgs) {
  const values = Array.isArray(headerArgs) ? headerArgs : [headerArgs];
  for (const value of values) {
    const headerText = String(value);
    const separatorIndex = headerText.indexOf(':');
    if (separatorIndex === -1) {
      console.warn(`Ignoring malformed header entry (expected "Key: Value"): ${headerText}`);
      continue;
    }
    const keyName = headerText.slice(0, separatorIndex).trim();
    const headerValue = headerText.slice(separatorIndex + 1).trim();
    if (!keyName || !headerValue) {
      console.warn(`Ignoring empty header key/value: ${headerText}`);
      continue;
    }
    if (keyName.toLowerCase() === 'authorization') {
      requestHeaders.Authorization = headerValue;
    } else {
      requestHeaders[keyName] = headerValue;
    }
  }
}

console.log(`Testing connection to: ${url}`);

// Configure HTTPS Agent for TLS (mTLS / CA)
const agentOptions = {
  rejectUnauthorized: false, // Default to allowing insecure, similar to curl -k
  keepAlive: false, // Disable keepAlive to avoid socket hang ups with self-signed certs
};

// If cert/key provided, load them (mTLS)
if (cert && key) {
  try {
    console.log(`Loading client cert from ${cert}`);
    console.log(`Loading client key from ${key}`);
    agentOptions.cert = fs.readFileSync(cert);
    agentOptions.key = fs.readFileSync(key);
    console.log(`Loaded cert size: ${agentOptions.cert.length}, key size: ${agentOptions.key.length}`);
  } catch (err) {
    console.error(`Error reading client cert/key: ${err.message}`);
    process.exit(1);
  }
} else if (cert || key) {
  console.error('Both --cert and --key must be provided together.');
  process.exit(1);
}

// If CA provided, load it and enable verification
if (cacert) {
  try {
    console.log(`Loading CA cert from ${cacert}`);
    agentOptions.ca = fs.readFileSync(cacert);
    // When CA is provided, enable certificate verification
    // This is necessary for proper mTLS handshake
    agentOptions.rejectUnauthorized = true;
    console.log(`Loaded CA cert size: ${agentOptions.ca.length}, verification enabled`);
  } catch (err) {
    console.error(`Error reading CA cert: ${err.message}`);
    process.exit(1);
  }
}

// Only create https.Agent for HTTPS URLs
const isHttps = url.startsWith('https://');
const agent = isHttps ? new https.Agent(agentOptions) : undefined;

let attempt = 0;
let success = false;

while (attempt < maxRetry && !success) {
  attempt++;
  console.log(`Attempt ${attempt} to connect to: ${url}`);

  try {
    const fetchOptions = isHttps ? { agent } : {};
    if (Object.keys(requestHeaders).length > 0) {
      fetchOptions.headers = requestHeaders;
    }
    const response = await fetch(url, fetchOptions);

    const actualStatusCode = response.status;

    if (expected !== undefined) {
      if (actualStatusCode !== expected) {
        console.error(`Test failed: Expected ${expected} but got ${actualStatusCode} when calling ${url}`);
      } else {
        console.log(`Connection test passed: ${url} returned ${expected}`);
        success = true;
      }
    } else if (notExpected !== undefined) {
      if (actualStatusCode === notExpected) {
        console.error(`Test failed: Expected status != ${notExpected} but got ${actualStatusCode} when calling ${url}`);
      } else {
        console.log(`Connection test passed: ${url} returned ${actualStatusCode} (!= ${notExpected})`);
        success = true;
      }
    }
  } catch (error) {
    if (notExpected !== undefined) {
      console.log(`Connection test passed: ${url} failed to connect as expected (error: ${error.message})`);
      success = true;
    } else {
      console.error(`Error during fetch request to ${url}: ${error.message}`);
      if (error.cause) {
        console.error('Cause:', error.cause);
      }
    }
  }

  if (!success && attempt < maxRetry) {
    console.log(`Retrying in ${retryDelay}ms...`);
    await new Promise((resolve) => setTimeout(resolve, retryDelay));
  }
}

if (!success) {
  console.error(`Failed to validate ${url} after ${maxRetry} attempts.`);
  process.exit(1);
}
