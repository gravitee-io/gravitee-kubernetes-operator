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
  branches: [{ name: "alpha", prerelease: true }, "master"],
  tagFormat: "${version}",
};

const branch = process.env.CIRCLE_BRANCH;
const changelogFileName =
  branch === "master" ? "CHANGELOG.md" : `CHANGELOG-${branch.toUpperCase()}.md`;

const plugins = [
  "@semantic-release/commit-analyzer",
  "@semantic-release/release-notes-generator",
  [
    "@semantic-release/changelog",
    {
      changelogFile: changelogFileName,
    },
  ],
  [
    "@semantic-release/exec",
    {
      prepareCmd:
        "IMG=graviteeio/kubernetes-operator:${nextRelease.version} BUNDLE_IMG=graviteeio/kubernetes-operator-bundle:${nextRelease.version} make docker-build docker-push bundle-standalone bundle-build bundle-push lint-fix",
    },
  ],
  [
    "@semantic-release/github",
    {
      assets: [{ path: "bundle.yml", label: "Operator resources bundle" }],
    },
  ],
  [
    "@semantic-release/git",
    {
      assets: [changelogFileName],
      message: "chore(release): ${nextRelease.version} [skip ci]",
    },
  ],
];

module.exports = { ...config, plugins };
