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

import { PROJECT_DIR } from "./index.mjs";

const MANIFEST = await parseManifest();

async function getImageTag() {
  return MANIFEST.apim.image.version;
}

async function getImageRegistry() {
  return MANIFEST.apim.image.registry;
}

async function getChartVersion() {
  return MANIFEST.apim.chart.version;
}

async function getChartRegistry() {
  return MANIFEST.apim.chart.registry;
}

async function parseManifest() {
  const manifestFilePath = path.join(PROJECT_DIR, "hack", "apim.yaml");
  const manifestFile = await fs.readFile(manifestFilePath, "utf8");
  return await YAML.parse(manifestFile);
}

export const APIM = {
  getChartVersion,
  getChartRegistry,
  getImageTag,
  getImageRegistry,
};
