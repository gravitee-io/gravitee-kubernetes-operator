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

const { name } = argv;
if (!name) {
    console.error('Usage: createServiceAccount.mjs --name <SERVICE_ACCOUNT_NAME>');
    process.exit(1);
}

const newUser = {
  lastname: name,
  email: `${name}@graviteesource.com`,
  service: true,
};

try {
    const createdUser = await mapiClient.post('/management/organizations/DEFAULT/users', newUser);
    console.log(JSON.stringify(createdUser));
} catch (e) {
    console.error(`Failed to create service account user: ${e.message}`);
    process.exit(1);
}
