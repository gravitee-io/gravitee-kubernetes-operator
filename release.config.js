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

const config = {
  branches: ["master"],
  tagFormat: "${version}",
};
const changelogFile = "CHANGELOG.md";
const chartDirectory = "helm/gko"
const crdDirectory = "helm/gko/crds"

const plugins = [
  "@semantic-release/commit-analyzer",
  "@semantic-release/release-notes-generator",
  [
    "@semantic-release/changelog",
    {
      changelogFile,
    },
  ],
  [
    "@semantic-release/exec",
    {
      prepareCmd:
        "npx zx scripts/release-helm-chart.mjs --version ${nextRelease.version} --img graviteeio/kubernetes-operator",
    },
  ],
  [
    "@semantic-release/github",
    {
      assets: [
        { 
          path: crdDirectory, label: "Operator Custom Resource Definitions"
        },
      ],
    },
  ],
  [
    "@semantic-release/git",
    {
      assets: [changelogFile, chartDirectory],
      message: "chore(release): ${nextRelease.version} [skip ci]",
    },
  ],
];

module.exports = { ...config, plugins };
