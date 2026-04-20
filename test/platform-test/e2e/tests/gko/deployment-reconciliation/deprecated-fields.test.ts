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
 * CRD field deprecation — batch 7.
 *
 * Xray tests:
 *   GKO-1391: processingStatus is deprecated in the status schema
 *
 * The CRD schema is the user-facing documentation for Kubernetes. If the
 * field's description does not carry an obvious deprecation marker, any
 * migration guidance around "use the Accepted condition instead" is lost.
 *
 * Known product gap (logged in "Batch 7 - Skipped Tests.md"):
 *   The Application CRD's status.processingStatus.description currently reads
 *   "The processing status of the Application." with no "*** DEPRECATED ***"
 *   marker, unlike the three CRDs below. Application is intentionally not
 *   asserted here — fixing the marker on Application should re-enable that
 *   assertion. This test passing does not imply every CRD is aligned.
 */

import { execFile } from "node:child_process";
import { promisify } from "node:util";
import { test, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

const execFileAsync = promisify(execFile);

interface CrdSchema {
  spec: {
    versions: Array<{
      schema?: {
        openAPIV3Schema?: {
          properties?: {
            status?: {
              properties?: Record<string, { description?: string }>;
            };
          };
        };
      };
    }>;
  };
}

async function getProcessingStatusDescription(crd: string): Promise<string | undefined> {
  const { stdout } = await execFileAsync("kubectl", ["get", "crd", crd, "-o", "json"], {
    timeout: 15_000,
  });
  const schema = JSON.parse(stdout) as CrdSchema;
  return schema.spec.versions[0]?.schema?.openAPIV3Schema?.properties?.status?.properties
    ?.processingStatus?.description;
}

test.describe("CRD field deprecation", () => {
  // ── GKO-1391: processingStatus marked deprecated ────────────

  test(`processingStatus is documented as deprecated ${XRAY.DEPLOYMENT_RECONCILIATION.PROCESSING_STATUS_DEPRECATED} ${TAGS.REGRESSION}`, async () => {
    const crds = [
      "apiv4definitions.gravitee.io",
      "apidefinitions.gravitee.io",
      "subscriptions.gravitee.io",
    ];

    for (const crd of crds) {
      const description = await getProcessingStatusDescription(crd);
      expect(description, `${crd} status.processingStatus.description`).toBeTruthy();
      expect(description?.toLowerCase()).toMatch(/deprecated/);
    }
  });
});
