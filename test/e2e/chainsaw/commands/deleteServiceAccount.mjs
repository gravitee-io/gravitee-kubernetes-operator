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
  console.error('Usage: deleteServiceAccount.mjs --name <SERVICE_ACCOUNT_NAME>');
  process.exit(1);
}

try {
  const searchPath = `/management/organizations/DEFAULT/search/users?q=${encodeURIComponent(name)}`;
  const resText = await mapiClient.get(searchPath);

  let users;
  try {
    users = JSON.parse(resText);
  } catch (parseErr) {
    throw new Error(`Unexpected response while searching users (not JSON).`);
  }

  if (!Array.isArray(users)) {
    throw new Error('Unexpected response format: expected an array of users.');
  }

const foundUser = users.find((u) => u && u.lastname === name);
  if (!foundUser) {
    console.error(`No user found with lastname exactly '${name}'.`);
    process.exit(0);
  }

  await mapiClient.del(`/management/organizations/DEFAULT/users/${foundUser.id}`);
  console.log(`Deleted service account user '${name}' (id: ${foundUser.id})`);
} catch (e) {
  console.error(`Failed to delete service account user: ${e.message}`);
  process.exit(1);
}
