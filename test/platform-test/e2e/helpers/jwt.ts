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
 * Minimal RS256 JWT signer for gateway data-plane assertions against a JWT plan.
 *
 * The suite has no JWT signer (existing tests only assert "401 without token"),
 * so a survival/reachability test cannot prove a legitimate caller gets through.
 * This mints a token matching exactly what APIM's GIVEN_KEY resolver expects -
 * the same shape as examples/usecase/subscribe-to-jwt-plan/pki/get_token.sh:
 * header {alg:RS256, typ:JWT}; claims iat, exp, client_id, sub. It is signed with
 * the test private key whose public half is in the jwt-secret fixture, so the
 * gateway verifies it. Uses only node:crypto - no npm dependency.
 */

import { createSign } from "node:crypto";
import { readFileSync } from "node:fs";
import { fixture } from "../setup.js";

function base64url(input: string | Buffer): string {
  return Buffer.from(input)
    .toString("base64")
    .replace(/=/g, "")
    .replace(/\+/g, "-")
    .replace(/\//g, "_");
}

/**
 * Sign an RS256 JWT for `subject` (used as both `sub` and `client_id`, so it must
 * match the subscribed application's clientId). Valid for `ttlSeconds` (default 1h).
 */
export function signJwt(subject: string, ttlSeconds = 3600): string {
  const pem = readFileSync(fixture("upgrade/jwt-private.key"), "utf8");
  const now = Math.floor(Date.now() / 1000);
  const header = base64url(JSON.stringify({ alg: "RS256", typ: "JWT" }));
  const payload = base64url(
    JSON.stringify({ iat: now, exp: now + ttlSeconds, client_id: subject, sub: subject }),
  );
  const signingInput = `${header}.${payload}`;
  const signature = base64url(createSign("RSA-SHA256").update(signingInput).sign(pem));
  return `${signingInput}.${signature}`;
}
