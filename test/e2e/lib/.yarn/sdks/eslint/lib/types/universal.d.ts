#!/usr/bin/env node
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


const {existsSync} = require(`fs`);
const {createRequire, register} = require(`module`);
const {resolve} = require(`path`);
const {pathToFileURL} = require(`url`);

const relPnpApiPath = "../../../../../.pnp.cjs";

const absPnpApiPath = resolve(__dirname, relPnpApiPath);
const absUserWrapperPath = resolve(__dirname, `./sdk.user.cjs`);
const absRequire = createRequire(absPnpApiPath);

const absPnpLoaderPath = resolve(absPnpApiPath, `../.pnp.loader.mjs`);
const isPnpLoaderEnabled = existsSync(absPnpLoaderPath);

if (existsSync(absPnpApiPath)) {
  if (!process.versions.pnp) {
    // Setup the environment to be able to require eslint/universal
    require(absPnpApiPath).setup();
    if (isPnpLoaderEnabled && register) {
      register(pathToFileURL(absPnpLoaderPath));
    }
  }
}

const wrapWithUserWrapper = existsSync(absUserWrapperPath)
  ? exports => absRequire(absUserWrapperPath)(exports)
  : exports => exports;

// Defer to the real eslint/universal your application uses
module.exports = wrapWithUserWrapper(absRequire(`eslint/universal`));
