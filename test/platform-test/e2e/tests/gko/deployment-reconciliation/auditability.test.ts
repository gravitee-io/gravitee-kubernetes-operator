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
 * Auditability.
 *
 * Xray tests:
 *   GKO-1463: Updates on a CR are traceable — k8s Events are emitted,
 *             status.conditions.lastTransitionTime reflects the change,
 *             and APIM-side updatedAt advances after the reconcile.
 *
 * The GKO-1463 Xray scope also mentions APIM audit logs. That API is not
 * currently exposed via the test Mapi client, so we verify traceability
 * through the k8s events (operator-side) and updatedAt (APIM-side) which
 * together give a reliable audit trail.
 */

import { execFile as execFileCb } from "node:child_process";
import { promisify } from "node:util";
import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";
import * as kubectlSafe from "../../../helpers/kubectl.js";

const execFile = promisify(execFileCb);

const V4_ORIGINAL = "crds/api-v4-definitions/v4-proxy-api-reconcile.yaml";
const V4_UPDATED = "crds/api-v4-definitions/v4-proxy-api-reconcile-updated.yaml";
const API_NAME = "e2e-v4-reconcile";

interface K8sEvent {
  metadata: { uid: string; resourceVersion?: string };
  message: string;
  reason?: string;
  count?: number;
  lastTimestamp?: string;
  eventTime?: string;
}

interface Condition {
  type: string;
  status: string;
  lastTransitionTime?: string;
}

interface StatusWithConditions {
  conditions?: Condition[];
}

async function fetchEvents(name: string): Promise<K8sEvent[]> {
  const { stdout } = await execFile("kubectl", [
    "get",
    "events",
    `--field-selector=involvedObject.name=${name}`,
    "-n",
    "default",
    "-o",
    "json",
  ]);
  return (JSON.parse(stdout) as { items: K8sEvent[] }).items;
}

/**
 * Sum of `count` across every "succeeded" event for this object. K8s
 * deduplicates identical (reason, message) events by incrementing
 * `count` on the existing event rather than creating a new UID, so the
 * UID set alone won't reveal a fresh reconcile. The total strictly
 * advances on each new "Update succeeded" emission.
 */
async function totalSucceededEventCount(name: string): Promise<number> {
  const events = await fetchEvents(name);
  return events
    .filter((e) => e.message.includes("succeeded"))
    .reduce((sum, e) => sum + (e.count ?? 1), 0);
}

test.describe("Auditability — changes traced via events and updatedAt", () => {
  test.afterEach(async () => {
    await kubectlSafe.del(fixture(V4_UPDATED)).catch(() => {});
    await kubectlSafe.del(fixture(V4_ORIGINAL)).catch(() => {});
  });

  test(`Update to CR is auditable via events and updatedAt ${XRAY.DEPLOYMENT_RECONCILIATION.AUDITABILITY_EVENTS} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    await test.step("Deploy initial CR", async () => {
      await kubectl.apply(fixture(V4_ORIGINAL));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const apiId = (
      await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME)
    ).id;
    const before = await mapi.fetchApi(apiId);
    const beforeUpdatedAt = before.updatedAt;

    // Capture baselines BEFORE applying the update so we can prove the
    // update — not the initial create — is what produced the new event
    // and the new lastTransitionTime.
    const preUpdateSucceededCount = await totalSucceededEventCount(API_NAME);
    const preUpdateStatus = await kubectl.getStatus<StatusWithConditions>(
      "apiv4definition",
      API_NAME,
    );
    const preUpdateTransition = preUpdateStatus.conditions?.find(
      (c) => c.type === "Accepted",
    )?.lastTransitionTime;
    expect(preUpdateTransition, "Accepted condition must have a lastTransitionTime").toBeTruthy();

    // Separate reconciles by a second so updatedAt and lastTransitionTime
    // can advance.
    await new Promise((r) => setTimeout(r, 1_100));

    await test.step("Update CR", async () => {
      await kubectl.apply(fixture(V4_UPDATED));
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    await test.step("APIM updatedAt has advanced", async () => {
      await expect
        .poll(
          async () => {
            const api = await mapi.fetchApi(apiId);
            return new Date(api.updatedAt).getTime();
          },
          { timeout: 30_000, intervals: [1_000] },
        )
        .toBeGreaterThan(new Date(beforeUpdatedAt).getTime());
    });

    await test.step("Reconcile fired a NEW k8s Event for the update", async () => {
      // K8s aggregates identical (reason, message) events on the same
      // object by incrementing `count` on the existing event rather
      // than minting a new UID. So the proof of a fresh reconcile is
      // the total `count` advancing — not a new UID appearing.
      await expect
        .poll(
          async () => totalSucceededEventCount(API_NAME),
          { timeout: 30_000, intervals: [1_000] },
        )
        .toBeGreaterThan(preUpdateSucceededCount);
    });

    await test.step("Accepted condition lastTransitionTime advances on update", async () => {
      await expect
        .poll(
          async () => {
            const status = await kubectl.getStatus<StatusWithConditions>(
              "apiv4definition",
              API_NAME,
            );
            const t = status.conditions?.find((c) => c.type === "Accepted")?.lastTransitionTime;
            return t ? new Date(t).getTime() : 0;
          },
          { timeout: 30_000, intervals: [1_000] },
        )
        .toBeGreaterThan(new Date(preUpdateTransition!).getTime());
    });

    await kubectl.del(fixture(V4_UPDATED));
  });
});
