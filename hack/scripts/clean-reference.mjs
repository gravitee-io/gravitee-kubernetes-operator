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

import { PROJECT_DIR } from "./lib/index.mjs";

const REFERENCE_FILE = path.join(PROJECT_DIR, "docs", "api", "reference.md");

async function cleanReference(reference) {
  return reference
    .replaceAll(/<\/gateway:experimental:description>/g, "")
    .replaceAll(/<gateway.*>/g, "");
}

async function readReference() {
  return await fs.readFile(REFERENCE_FILE, "utf8");
}

const dirty = await readReference();
const clean = await cleanReference(dirty);
await fs.writeFile(REFERENCE_FILE, clean);
