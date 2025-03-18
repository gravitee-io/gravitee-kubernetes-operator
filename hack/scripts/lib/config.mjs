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
const repoName = "cloud-gko";
const repo = `gravitee-io/${repoName}`;
const branch = "main";
const gkoCommitHashAnnotationKey = "gko.gravitee.io/commit-hash";
const apimCommitHashAnnotationKey = "apim.gravitee.io/commit-hash";
const gkoValues = "gko/values-ci.yaml";
const apimValues = "apim/values-ci.yaml";

export const CONFIG = {
  repoName,
  repo,
  branch,
  commitHashAnnotationKey: gkoCommitHashAnnotationKey,
  apimCommitHashAnnotationKey,
  gkoValues,
  apimValues,
};
