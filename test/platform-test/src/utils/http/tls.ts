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

import * as https from "node:https";
import * as http from "node:http";
import type { FetchFn, TlsOptions } from "../../types/http.js";

/**
 * Create a fetch-compatible function backed by node:https.request, enabling
 * mTLS client certificates and custom CA verification.
 *
 * Native fetch (Node 18+) does not support HTTPS agent configuration needed
 * for client certificates. This helper fills that gap using only built-in
 * node:https / node:http APIs (zero runtime dependencies).
 *
 * HTTP URLs fall back to node:http.request automatically.
 *
 * @example
 * import { readFile } from "node:fs/promises";
 * import { createTlsFetch } from "@gravitee/platform-test/utils/http";
 *
 * const cert = await readFile("client.crt");
 * const key  = await readFile("client.key");
 * const ca   = await readFile("ca.crt");
 *
 * const gateway = apim.gateway(
 *   { baseUrl: "https://localhost:8443" },
 *   createTlsFetch({ cert, key, ca }),
 * );
 * await gateway.assertResponds("/mtls-demo", { status: 200 });
 */
export function createTlsFetch(tls: TlsOptions): FetchFn {
  return (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    const url =
      typeof input === "string"
        ? input
        : input instanceof URL
          ? input.href
          : (input as Request).url;

    const method = (init?.method ?? "GET").toUpperCase();

    // Normalise headers to a plain object
    const requestHeaders: Record<string, string> = {};
    if (init?.headers) {
      if (init.headers instanceof Headers) {
        init.headers.forEach((v, k) => {
          requestHeaders[k] = v;
        });
      } else if (Array.isArray(init.headers)) {
        for (const [k, v] of init.headers as [string, string][]) {
          requestHeaders[k] = v;
        }
      } else {
        Object.assign(requestHeaders, init.headers);
      }
    }

    const parsed = new URL(url);
    const isHttps = parsed.protocol === "https:";

    return new Promise<Response>((resolve, reject) => {
      const options: https.RequestOptions = {
        hostname: parsed.hostname,
        port: parsed.port ? parseInt(parsed.port, 10) : isHttps ? 443 : 80,
        path: `${parsed.pathname}${parsed.search}`,
        method,
        headers: requestHeaders,
        ...(isHttps
          ? {
              cert: tls.cert,
              key: tls.key,
              ca: tls.ca,
              rejectUnauthorized: tls.rejectUnauthorized ?? false,
            }
          : {}),
      };

      const lib: typeof https = isHttps ? https : (http as unknown as typeof https);
      const req = lib.request(options, (res) => {
        const chunks: Buffer[] = [];
        res.on("data", (chunk: Buffer) => chunks.push(chunk));
        res.on("end", () => {
          const body = Buffer.concat(chunks);
          const status = res.statusCode ?? 0;
          const headers = new Headers();
          for (const [k, v] of Object.entries(res.headers)) {
            if (v !== undefined) {
              const values = Array.isArray(v) ? v : [String(v)];
              for (const val of values) {
                headers.append(k, val);
              }
            }
          }
          resolve(new Response(body, { status, headers }));
        });
        res.on("error", reject);
      });

      req.on("error", reject);
      req.end();
    });
  };
}
