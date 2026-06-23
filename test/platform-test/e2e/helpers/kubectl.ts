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

/**
 * kubectl CLI wrapper. The implementation now lives in the runner-agnostic
 * `@gravitee/platform-test` library (`src/provisioners/engines/kubectl.ts`) so
 * the GKO provisioner can reuse it; this module re-exports it unchanged so the
 * `e2e/helpers/kubectl.js` import path keeps working for existing tests and the
 * Playwright `kubectl` fixture.
 */
export * from "../../src/provisioners/engines/kubectl.js";
