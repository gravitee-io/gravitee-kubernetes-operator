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

import { LOG, PROJECT_DIR } from "./lib/index.mjs";
import { Version } from "./lib/version.mjs";

const mergifyConfig = path.join(PROJECT_DIR, ".mergify.yml");

export async function rolloutMergify(newVersion) {
  LOG.blue(`
    Rolling out mergify config`);
  const branch = new Version(newVersion).branch();
  const config = await fs.readFile(mergifyConfig, "utf8");
  const configYaml = await YAML.parse(config);
  const rules = configYaml.pull_request_rules;
  const rule = rules.pop();
  const removedBranch = rule.actions.backport.branches[0];
  LOG.blue(`
    Removed branch ${removedBranch} from config`);
  rule.name = rule.name.replace(removedBranch, branch);
  rule.conditions = [`label=apply-on-${branch.replaceAll(".", "-")}`];
  rule.actions.backport.branches = [branch];
  rule.actions.backport.title = rule.actions.backport.title.replace(
    removedBranch,
    branch,
  );
  configYaml.pull_request_rules.splice(1, 0, rule);
  LOG.blue(`
    Inserted branch ${branch} in config`);
  LOG.log();
  await fs.writeFile(mergifyConfig, YAML.stringify(configYaml));
}
