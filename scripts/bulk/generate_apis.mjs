#!/usr/bin/env zx
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

const { log } = console;
const newLogger =
    (color) =>
        (...args) =>
            log(color(...args));

const blue = newLogger(chalk.blue);
const white = newLogger(chalk.whiteBright);
const red = newLogger(chalk.red);

const timer = ms => new Promise(res => setTimeout(res, ms))

// enter the context name
const contextName = await question(blue(`Enter the context name:`));

// enter the number of APIs to generate
const nbApis = await question(white(`Number of APIs to generate:`));

// enter the time to wait between each API creation
const waitTime = await question(red(`Time to wait between each API creation (in seconds):`));


// loop to create APIs
for (let i = 0; i < nbApis; i++) {
    // create API using helm chart
    const apiName = `api-${i}`;

    await $`
cat <<EOF | kubectl apply -f -
apiVersion: gravitee.io/v1alpha1
kind: ApiDefinition
metadata:
  name: "${apiName}"
  namespace: "gko-perf-tests"
spec:
  name: "K8s Basic Example With Management Context"
  contextRef: 
    name: "${contextName}"
    namespace: "default"
  version: "1.1"
  description: "Basic api managed by Gravitee Kubernetes Operator"
  proxy:
    virtual_hosts:
      - path: "/${apiName}"
    groups:
      - endpoints:
          - name: "Default"
            target: "https://api.gravitee.io/echo"
EOF
`;
    await timer(waitTime * 1000);
}
