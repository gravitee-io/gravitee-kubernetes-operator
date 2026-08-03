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
 * Journey: tell a team about an API's events.
 *
 * As an API producer, I have the console notify a team's group when my API
 * starts or stops, then widen that to expiring api-keys.
 *
 * The two drivers MODEL this differently — GKO has a standalone `Notification`
 * CR the API references, Terraform inlines `apim_apiv4.console_notification` —
 * which is exactly why it is worth one shared journey: both must land the same
 * PORTAL notification setting (hooks + target groups) on the API in APIM.
 *
 * The notified group is created by the same provisioner and must reconcile
 * first, like associate-groups-with-an-api. APIM only lets a console
 * notification target groups the API itself belongs to, so the group
 * association is setup here rather than the thing under assertion.
 *
 * Turning the notification OFF is deliberately absent (GKO-3085): the Automation
 * API answers HTTP 500 ("organizationId must not be empty") to a console
 * notification with an empty `events` list, and setting the block to null is a
 * Terraform no-op that leaves the previous settings in place, so the Terraform
 * arm cannot express it. The GKO side removes the Notification reference and is
 * covered by GKO-1238 in tests/gko/notifications.
 *
 * Fixtures are co-located in this folder. Notification concerns that only exist
 * in the operator's model — a Notification CR shared by a V2 and a V4 API, two
 * console-target notifications on one API being rejected, notification refs
 * being excluded from the CRD export — stay in tests/gko/notifications.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test, expect } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

const BASE_EVENTS = ["API_STARTED", "API_STOPPED"];
const EXTENDED_EVENTS = [...BASE_EVENTS, "APIKEY_EXPIRED"];

interface NotificationParams {
  events: string[];
}

interface NotificationStage {
  label: string;
  params: NotificationParams;
}

const STAGES: NotificationStage[] = [
  {
    label: "notifies the group when the API starts or stops",
    params: { events: BASE_EVENTS },
  },
  {
    label: "also notifies the group about expiring api-keys",
    params: { events: EXTENDED_EVENTS },
  },
];

forEachProvisioner<NotificationParams>(
  {
    title: "Configure console notifications on an API",
    provisioners: {
      gko: gkoScenario<NotificationParams>({
        // The Group and the base Notification are static so teardown removes
        // them; the Notification variant and the API are applied by applyParams
        // (the event list is the knob, and the API must come after both). The
        // two Notification variants share one CR name, so deleting the static
        // manifest removes whichever is current.
        manifests: [path.join(here, "gko/group.yaml"), path.join(here, "gko/notification-base.yaml")],
        roles: { group: "notified-group", api: "notified-api" },
        dynamicRoles: ["api"],
        applyParams: async (k, params) => {
          const variant = params.events.length > BASE_EVENTS.length ? "extended" : "base";
          await k.apply(path.join(here, `gko/notification-${variant}.yaml`));
          await k.apply(path.join(here, "gko/api-notified.yaml"));
        },
      }),
      terraform: tfScenario<NotificationParams>({
        fixture: path.join(here, "terraform"),
        toVars: (params) => ({ notification_events: params.events }),
      }),
    },
    xray: {
      gko: [
        XRAY.NOTIFICATIONS.NOTIFICATION_HOOKS_GROUPS,
        XRAY.NOTIFICATIONS.API_REF_NOTIFICATION,
        XRAY.NOTIFICATIONS.NOTIFICATIONS_VIA_CRS,
        XRAY.NOTIFICATIONS.VIEW_NOTIFICATION_SETTINGS,
        XRAY.NOTIFICATIONS.NOTIFICATION_LABEL,
      ],
      terraform: XRAY.TERRAFORM.API_NOTIFICATIONS_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 120_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();
    const groupId = await provisioned.groupId();

    for (const [index, stage] of STAGES.entries()) {
      await test.step(`The API ${stage.label}`, async () => {
        // Stage 0 is what provision() already applied.
        if (index > 0) await provisioned.update(stage.params);

        await expect
          .poll(
            async () => {
              const settings = await mapi.fetchApiNotificationSettings(apiId);
              const portal = settings.find((s) => s.config_type === "PORTAL");
              return {
                hooks: [...(portal?.hooks ?? [])].sort(),
                groups: portal?.groups ?? [],
              };
            },
            { timeout: 30_000, message: `notification settings after: ${stage.label}` },
          )
          .toEqual({ hooks: [...stage.params.events].sort(), groups: [groupId] });
      });
    }
  },
  STAGES[0].params,
);
