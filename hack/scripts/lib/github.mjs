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

import { LOG, PROJECT_DIR } from "./index.mjs";
import { Version } from "./version.mjs";

const githubWorflows = path.join(PROJECT_DIR, ".github", "workflows");

async function rolloutMatrix(fileName, jobName, newVersion) {
  const branch = new Version(newVersion).branch();
  const workflowFile = path.join(githubWorflows, fileName);
  LOG.blue(`Reading file ${workflowFile}`);
  const workflow = await fs.readFile(workflowFile, "utf8");
  const workflowYaml = await YAML.parse(workflow);
  LOG.blue(`Rolling out github matrix for job ${jobName}`);
  const job = workflowYaml.jobs[jobName];
  const branches = job.strategy.matrix.branch;
  const removedBranch = branches.shift();
  LOG.blue(`Removed branch ${removedBranch} from matrix`);
  branches.splice(branches.length - 1, 0, branch);
  LOG.blue(`Inserted branch ${branch} in matrix`);
  await fs.writeFile(workflowFile, YAML.stringify(workflowYaml));
}

export const GH = {
  rolloutMatrix,
};
