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

import { APIM } from "./lib/apim.mjs";
import { LOG, time } from "./lib/index.mjs";

const APIM_IMAGE_REGISTRY = await APIM.getImageRegistry();
const APIM_IMAGE_TAG = await APIM.getImageTag();

async function pullAndLoadGatewayImage() {
  LOG.blue(
    `Pulling gateway image from ${APIM_IMAGE_REGISTRY}/apim-gateway:${APIM_IMAGE_TAG}...`,
  );
  await $`docker pull ${APIM_IMAGE_REGISTRY}/apim-gateway:${APIM_IMAGE_TAG}`;
  LOG.blue(`Tagging gateway image as gateway:latest`);
  await $`docker tag ${APIM_IMAGE_REGISTRY}/apim-gateway:${APIM_IMAGE_TAG} gateway:latest`;
  LOG.blue(`Loading gateway image into kind cluster`);
  await $`kind load docker-image gateway:latest --name gravitee`;
}

await time(pullAndLoadGatewayImage);
