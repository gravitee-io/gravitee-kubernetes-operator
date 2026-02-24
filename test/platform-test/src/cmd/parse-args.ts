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
 * Parse an argv array (starting after `node script.js`) into
 * a subcommand name and a flags map.
 *
 * Flag format: --key value  or  --key=value
 */
export function parseArgs(argv: string[]): { subcommand: string; flags: Record<string, string> } {
  const [subcommand = "", ...rest] = argv;
  const flags: Record<string, string> = {};

  for (let i = 0; i < rest.length; i++) {
    const arg = rest[i];

    if (arg.startsWith("--")) {
      const eqIdx = arg.indexOf("=");
      if (eqIdx !== -1) {
        // --key=value form
        const key = arg.slice(2, eqIdx);
        const value = arg.slice(eqIdx + 1);
        flags[key] = value;
      } else {
        // --key value form
        const key = arg.slice(2);
        const next = rest[i + 1];
        if (next !== undefined && !next.startsWith("--")) {
          flags[key] = next;
          i++;
        } else {
          // boolean flag (no value)
          flags[key] = "true";
        }
      }
    }
  }

  return { subcommand, flags };
}
