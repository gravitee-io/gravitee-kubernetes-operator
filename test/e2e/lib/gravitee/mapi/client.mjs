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

import { httpClient } from "../../http/client.mjs";

async function post(path, body) {
  const url = new URL(path, httpClient.BASE_URL_API).toString();
  const res = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': httpClient.AUTH_HEADER,
    },
    body: JSON.stringify(body),
  });

  if (!res.ok) {
    const err = await res.text().catch(() => '');
    throw new Error(`Failed request to ${url}: ${res.status} ${err}`);
  }

  return res.json();
}

async function get(path) {
  const url = new URL(path, httpClient.BASE_URL_API).toString();
  const res = await fetch(url, {
    method: 'GET',
    headers: {
      'Authorization': httpClient.AUTH_HEADER,
    },
  });

  if (!res.ok) {
    const err = await res.text().catch(() => '');
    throw new Error(`Failed request to ${url}: ${res.status} ${err}`);
  }

  return res.text();
}

async function del(path) {
  const url = new URL(path, httpClient.BASE_URL_API).toString();
  const res = await fetch(url, {
    method: 'DELETE',
    headers: {
      'Authorization': httpClient.AUTH_HEADER,
    },
  });

  if (!res.ok) {
    const err = await res.text().catch(() => '');
    throw new Error(`Failed request to ${url}: ${res.status} ${err}`);
  }

  return res.text().catch(() => '');
}

export const mapiClient = {
  post,
  get,
  del,
};
