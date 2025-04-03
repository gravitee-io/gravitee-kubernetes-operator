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

import { isEmptyString, LOG } from "./index.mjs";

const API_BASE = "https://circleci.com/api/v2";
const APP_BASE = "https://app.circleci.com";
const ORG = "gravitee-io";
const SCM = "github";
const PROJECT = "gravitee-kubernetes-operator";
const CIRCLECI_TOKEN = process.env.CIRCLECI_TOKEN;

if (isEmptyString(CIRCLECI_TOKEN)) {
  LOG.red("CIRCLECI_TOKEN cannot be found");
  process.exit(1);
}

export async function triggerPipeline(parameters, branch = "master") {
  const response = await fetch(
    `${API_BASE}/project/${SCM}/${ORG}/${PROJECT}/pipeline`,
    {
      method: "POST",
      headers: {
        "Circle-Token": CIRCLECI_TOKEN,
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({ parameters, branch }),
    },
  );

  if (response.status === 201) {
    const json = await response.json();
    return `${APP_BASE}/pipelines/${SCM}/${ORG}/${PROJECT}/${json.number}`;
  }

  throw new Error(`Unable to run pipeline (HTTP status ${response.status})`);
}
