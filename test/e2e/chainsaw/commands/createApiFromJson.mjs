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

await dotenv.config(`${__dirname}/../.env`);

const { jsonfile } = argv;
if (!jsonfile) {
    console.error('Usage: createApiFromJson.mjs --jsonfile <JSON_FILE>');
    process.exit(1);
}

try {
    await $`test -f ${jsonfile}`;
} catch {
    console.error(`Error: File ${jsonfile} not found.`);
    process.exit(1);
}

const { APIM_API, APIM_AUTH } = process.env;
if (!APIM_API || !APIM_AUTH) {
    console.error('Error: APIM_API or APIM_AUTH not set in .env file.');
    process.exit(1);
}

const fileContent = fs.readFileSync(jsonfile, 'utf8');
const apiDefinition = JSON.parse(fileContent);

let url;
if (apiDefinition.gravitee === '2.0.0') {
    url = `${APIM_API}/management/organizations/DEFAULT/environments/DEFAULT/apis/import`;
} else if (apiDefinition.api && apiDefinition.api.definitionVersion === 'V4') {
    url = `${APIM_API}/management/v2/environments/DEFAULT/apis/_import/definition`;
} else {
    console.error('Unknown API definition version');
    process.exit(1);
}


  try {
        const response = await fetch(url, {
            method: 'POST',
            body: fileContent,
            headers: {
                'Authorization': `Bearer ${APIM_AUTH}`,
                'Content-Type': 'application/json'
            }
    });

    if (!response.ok) {
        console.error(`Error during API import: HTTP ${response.status} ${response.statusText}`);
        process.exit(1);
    }

    const responseData = await response.text();
    console.log(responseData);
} catch (error) {
    console.error('Error during API import:', error);
    process.exit(1);
}
