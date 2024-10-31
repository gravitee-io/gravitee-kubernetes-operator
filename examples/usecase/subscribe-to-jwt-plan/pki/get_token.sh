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

dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

pem=$( cat "${dir}/private.key" )
client_id="echo-client"
sub="echo-client"

now=$( date +%s )
iat="${now}"
# shellcheck disable=SC2004
exp=$((${now} + 3600))
header_raw='{"alg":"RS256", "typ": "JWT"}'
header=$( echo -n "${header_raw}" | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n' )
payload_raw='{"iat":'"$iat"',"exp":'"$exp"',"client_id":"'"$client_id"'", "sub":"'"$sub"'"}'
payload=$( echo -n "${payload_raw}" | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n' )
header_payload="${header}"."${payload}"
signature=$( openssl dgst -sha256 -sign <(echo -n "${pem}") <(echo -n "${header_payload}") | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n' )
echo "${header}.${payload}.${signature}"
