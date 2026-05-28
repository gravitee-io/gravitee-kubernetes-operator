#!/usr/bin/env bash
# Copyright (C) 2015 The Gravitee team (http://gravitee.io)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Push the Playwright e2e results to Xray Cloud as a Test Execution so the
# matching Jira Test issues' "Test Coverage" panel reflects each run.
#
# The JUnit report is transformed (junit-to-xray.mjs) into Xray's JSON results
# format, keyed by the Jira issue id embedded in each test title — the platform
# test suite already includes those via the XRAY constants in
# e2e/helpers/tags.ts (e.g. "... @GKO-2865 @regression"). Importing as Xray
# JSON (rather than raw JUnit) is what makes results attach to the existing
# Tests by key instead of creating duplicate Test issues.
#
# Best-effort: if Xray is unreachable, credentials are missing, or the upload
# fails for any reason, the script exits 0 with a warning. The pipeline's
# pass/fail gate is unaffected — that's CircleCI's own test-result parsing.
#
# Usage:
#   XRAY_CLIENT_ID=...  XRAY_CLIENT_SECRET=... \
#       ./test/platform-test/scripts/xray/push-results.sh [path/to/junit-results.xml]
#
# Env vars (consumed by the script):
#   XRAY_CLIENT_ID         (required) Xray Cloud API client id
#   XRAY_CLIENT_SECRET     (required) Xray Cloud API client secret
#   XRAY_TEST_PLAN_KEY     (optional) existing Test Plan issue key; if set, the
#                          new Test Execution is linked to it so the Plan view
#                          shows trend over time
#   CIRCLE_BUILD_URL,
#   CIRCLE_BUILD_NUM,
#   CIRCLE_BRANCH          (optional) populated automatically in CircleCI; used
#                          in the resulting Test Execution summary/description

set -uo pipefail

# Resolve paths relative to this script so it works both from the repo root
# and once test/platform-test/ is extracted into its own repository.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLATFORM_TEST_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

RESULTS_PATH="${1:-${PLATFORM_TEST_DIR}/playwright-results/results.xml}"
TRANSFORM="${SCRIPT_DIR}/junit-to-xray.mjs"
PAYLOAD_FILE="${TMPDIR:-/tmp}/xray-execution.json"
KEY_FILE="${TMPDIR:-/tmp}/xray-test-execution.txt"

log() { printf '[xray-push] %s\n' "$*"; }
warn_and_exit() { log "WARN: $*"; log "skipping Xray upload — pipeline not affected"; exit 0; }

# ── Preflight ────────────────────────────────────────────────────────

[[ -f "${RESULTS_PATH}" ]] \
  || warn_and_exit "no JUnit results found at ${RESULTS_PATH}"

[[ -n "${XRAY_CLIENT_ID:-}" && -n "${XRAY_CLIENT_SECRET:-}" ]] \
  || warn_and_exit "XRAY_CLIENT_ID / XRAY_CLIENT_SECRET not set"

command -v jq >/dev/null \
  || warn_and_exit "jq not on PATH"

command -v node >/dev/null \
  || warn_and_exit "node not on PATH"

[[ -f "${TRANSFORM}" ]] \
  || warn_and_exit "transform script missing at ${TRANSFORM}"

# ── 1. Authenticate → JWT ────────────────────────────────────────────
#
# The /authenticate endpoint returns the token as a JSON string (i.e. a
# bare quoted string). `jq -r .` strips the quotes; -f makes curl error on
# non-2xx so we can fall through to warn_and_exit.

TOKEN=$(curl -fsS -X POST \
          -H 'Content-Type: application/json' \
          -d "{\"client_id\":\"${XRAY_CLIENT_ID}\",\"client_secret\":\"${XRAY_CLIENT_SECRET}\"}" \
          https://xray.cloud.getxray.app/api/v2/authenticate \
        | jq -r '.') \
  || warn_and_exit "auth to Xray Cloud failed"

[[ -n "${TOKEN}" && "${TOKEN}" != "null" ]] \
  || warn_and_exit "auth returned an empty token"

# ── 2. Build the Xray JSON payload ───────────────────────────────────
#
# Convert the JUnit report into Xray's JSON results format, keyed by the
# @GKO-NNNN Test issue id embedded in each title, so each result attaches to
# the *existing* Test (the JUnit importer instead matches Tests by
# classname+name and would create duplicate Test issues). Skipped / fixme'd
# cases are dropped so they don't overwrite a Test's prior result with TODO.

BRANCH="${CIRCLE_BRANCH:-$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo local)}"
BUILD_NUM="${CIRCLE_BUILD_NUM:-local}"
BUILD_URL="${CIRCLE_BUILD_URL:-(local run)}"

XRAY_SUMMARY="GKO e2e - ${BUILD_NUM} (${BRANCH})" \
XRAY_DESCRIPTION="Automated import from CircleCI: ${BUILD_URL}" \
XRAY_TEST_PLAN_KEY="${XRAY_TEST_PLAN_KEY:-}" \
  node "${TRANSFORM}" "${RESULTS_PATH}" > "${PAYLOAD_FILE}" \
  || warn_and_exit "failed to build Xray payload from ${RESULTS_PATH}"

TEST_COUNT=$(jq '.tests | length' "${PAYLOAD_FILE}" 2>/dev/null || echo 0)
if [[ "${TEST_COUNT}" -eq 0 ]]; then
  warn_and_exit "no @GKO-tagged results to report (all skipped or untagged)"
fi
log "reporting ${TEST_COUNT} test result(s) to Xray Cloud"

# ── 3. POST the Xray JSON results ────────────────────────────────────
#
# Capture the body and HTTP status separately (no -f) so a non-2xx
# response surfaces Xray's actual error message instead of an opaque
# curl exit code — the import endpoint returns the reason in the body.

BODY_FILE="${TMPDIR:-/tmp}/xray-upload-response.json"
HTTP_CODE=$(curl -sS -o "${BODY_FILE}" -w '%{http_code}' -X POST \
              -H "Authorization: Bearer ${TOKEN}" \
              -H "Content-Type: application/json" \
              --data @"${PAYLOAD_FILE}" \
              "https://xray.cloud.getxray.app/api/v2/import/execution") \
  || warn_and_exit "could not reach Xray Cloud import endpoint"

RESP=$(cat "${BODY_FILE}" 2>/dev/null)

if [[ "${HTTP_CODE}" != 2* ]]; then
  log "WARN: import endpoint returned HTTP ${HTTP_CODE}"
  log "raw response: ${RESP}"
  warn_and_exit "upload to Xray Cloud failed"
fi

# ── 4. Capture the Test Execution key as a build artifact ────────────

TEST_EXEC_KEY=$(echo "${RESP}" | jq -r '.testExecIssue.key // .key // empty')
if [[ -z "${TEST_EXEC_KEY}" ]]; then
  log "WARN: Xray accepted the upload but no Test Execution key was returned"
  log "raw response: ${RESP}"
  exit 0
fi

echo "${TEST_EXEC_KEY}" > "${KEY_FILE}"
log "created Test Execution ${TEST_EXEC_KEY}"
log "  Jira: https://gravitee.atlassian.net/browse/${TEST_EXEC_KEY}"
log "  key written to ${KEY_FILE}"
