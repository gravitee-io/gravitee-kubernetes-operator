#!/bin/bash
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


if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <ENDPOINT_PATH> <EXPECTED_STATUS_CODE>" >&2
    exit 1
fi

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
ENV_FILE="$SCRIPT_DIR/.."/.env

if [ -f "$ENV_FILE" ]; then
  source "$ENV_FILE"
else
  echo "Error: .env file not found at $ENV_FILE. Make sure it exists in the expected location." >&2
  exit 1
fi

ENDPOINT=$1
EXPECTED_STATUS_CODE=$2
URL="${APIM_GATEWAY%/}/${ENDPOINT#/}"

ACTUAL_STATUS_CODE=$(curl $URL -sS -o /dev/null -w "%{http_code}")

if [ "$ACTUAL_STATUS_CODE" != "$EXPECTED_STATUS_CODE" ]; then
    echo "Test failed: Expected $EXPECTED_STATUS_CODE but got $ACTUAL_STATUS_CODE when calling $URL" >&2
    exit 1
else
    echo "Connection test passed: ${URL} returned ${EXPECTED_STATUS_CODE}"
fi