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
 * Message APIs — Lifecycle tests.
 *
 * Xray tests:
 *   GKO-72:  Deploy V4 Message API with syncFrom Management
 *   GKO-73:  Deploy V4 Message API with syncFrom Kubernetes
 *   GKO-129: Deploy V4 message API with HTTP GET entrypoint
 *   GKO-130: Deploy V4 message API with HTTP POST entrypoint
 *   GKO-132: Deploy V4 message API with SSE entrypoint
 *   GKO-133: Deploy V4 message API with Webhooks entrypoint
 *   GKO-134: Deploy V4 message API with Websockets entrypoint
 *   GKO-135: Deploy V4 message API with Kafka endpoint
 *   GKO-136: Deploy V4 message API with Mock endpoint
 *   GKO-164: Deploy V4 message API with policy
 *
 * Preconditions:
 *   - APIM, Gateway, and GKO operator are running
 *   - A ManagementContext "dev-ctx" exists in the default namespace
 */

import { test, fixture, expect } from "../../../setup.js";
import { XRAY, TAGS } from "../../../helpers/tags.js";

test.describe("Message APIs — Lifecycle", () => {
  // ── GKO-72: Deploy V4 Message API with syncFrom Management ──

  test(`Deploy V4 Message API with syncFrom Management ${XRAY.MESSAGE_APIS.DEPLOY_V4_MSG_SYNC_MGMT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-sync-mgmt";
    const fixturePath = fixture("crds/message-apis/v4-message-api-sync-mgmt.yaml");

    await test.step("Apply CRD with syncFrom Management", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API is STARTED in APIM with type MESSAGE", async () => {
      await mapi.assertApiMatches(apiId, { name: API_NAME, state: "STARTED" });
      const api = await mapi.fetchApi(apiId);
      expect(api.type).toBe("MESSAGE");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-73: Deploy V4 Message API with syncFrom Kubernetes ──

  test(`Deploy V4 Message API with syncFrom Kubernetes ${XRAY.MESSAGE_APIS.DEPLOY_V4_MSG_SYNC_K8S} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-sync-k8s";
    const fixturePath = fixture("crds/message-apis/v4-message-api-sync-k8s.yaml");

    await test.step("Apply CRD with syncFrom Kubernetes", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API is STARTED in APIM with type MESSAGE", async () => {
      await mapi.assertApiMatches(apiId, { name: API_NAME, state: "STARTED" });
      const api = await mapi.fetchApi(apiId);
      expect(api.type).toBe("MESSAGE");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-129: Deploy V4 message API with HTTP GET entrypoint ─

  test(`Deploy V4 message API with HTTP GET entrypoint ${XRAY.MESSAGE_APIS.HTTP_GET_ENTRYPOINT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-http-get";
    const fixturePath = fixture("crds/message-apis/v4-message-api-http-get.yaml");

    await test.step("Apply CRD with HTTP GET entrypoint", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API exists in APIM with type MESSAGE", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api).toBeTruthy();
      expect(api.type).toBe("MESSAGE");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-130: Deploy V4 message API with HTTP POST entrypoint ─

  test(`Deploy V4 message API with HTTP POST entrypoint ${XRAY.MESSAGE_APIS.HTTP_POST_ENTRYPOINT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-http-post";
    const fixturePath = fixture("crds/message-apis/v4-message-api-http-post.yaml");

    await test.step("Apply CRD with HTTP POST entrypoint", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API exists in APIM with type MESSAGE", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api).toBeTruthy();
      expect(api.type).toBe("MESSAGE");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-132: Deploy V4 message API with SSE entrypoint ──────

  test(`Deploy V4 message API with SSE entrypoint ${XRAY.MESSAGE_APIS.SSE_ENTRYPOINT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-sse";
    const fixturePath = fixture("crds/message-apis/v4-message-api-sse.yaml");

    await test.step("Apply CRD with SSE entrypoint", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API exists in APIM with type MESSAGE", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api).toBeTruthy();
      expect(api.type).toBe("MESSAGE");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-133: Deploy V4 message API with Webhooks entrypoint ─

  test(`Deploy V4 message API with Webhooks entrypoint ${XRAY.MESSAGE_APIS.WEBHOOK_ENTRYPOINT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-webhook";
    const fixturePath = fixture("crds/message-apis/v4-message-api-webhook.yaml");

    await test.step("Apply CRD with Webhook entrypoint", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API exists in APIM with type MESSAGE", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api).toBeTruthy();
      expect(api.type).toBe("MESSAGE");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-134: Deploy V4 message API with Websockets entrypoint

  test(`Deploy V4 message API with Websockets entrypoint ${XRAY.MESSAGE_APIS.WEBSOCKET_ENTRYPOINT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-websocket";
    const fixturePath = fixture("crds/message-apis/v4-message-api-websocket.yaml");

    await test.step("Apply CRD with Websocket entrypoint", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API exists in APIM with type MESSAGE", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api).toBeTruthy();
      expect(api.type).toBe("MESSAGE");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-135: Deploy V4 message API with Kafka endpoint ──────

  // SKIP: APIM dev build returns 500 "required key [consumer/producer] not found"
  // for Kafka endpoint configuration. GKO/APIM serialization bug.
  test.skip(`Deploy V4 message API with Kafka endpoint ${XRAY.MESSAGE_APIS.KAFKA_ENDPOINT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-kafka";
    const fixturePath = fixture("crds/message-apis/v4-message-api-kafka.yaml");

    await test.step("Apply CRD with Kafka endpoint", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API exists in APIM with type MESSAGE", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api).toBeTruthy();
      expect(api.type).toBe("MESSAGE");
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-136: Deploy V4 message API with Mock endpoint ───────

  test(`Deploy V4 message API with Mock endpoint ${XRAY.MESSAGE_APIS.MOCK_ENDPOINT} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-http-get";
    const fixturePath = fixture("crds/message-apis/v4-message-api-http-get.yaml");

    await test.step("Apply CRD with Mock endpoint", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API is STARTED in APIM", async () => {
      await mapi.assertApiStarted(apiId);
    });

    await kubectl.del(fixturePath);
  });

  // ── GKO-164: Deploy V4 message API with policy ──────────────

  test(`Deploy V4 message API with policy ${XRAY.MESSAGE_APIS.MSG_API_WITH_POLICY} ${TAGS.REGRESSION}`, async ({
    kubectl,
    mapi,
  }) => {
    const API_NAME = "e2e-v4-msg-policy";
    const fixturePath = fixture("crds/message-apis/v4-message-api-with-policy.yaml");

    await test.step("Apply CRD with transform-headers policy", async () => {
      await kubectl.apply(fixturePath);
      await kubectl.waitForCondition("apiv4definition", API_NAME, "Accepted");
    });

    const status = await kubectl.getStatus<{ id: string }>("apiv4definition", API_NAME);
    const apiId = status.id;

    await test.step("API has flows configured in APIM", async () => {
      const api = await mapi.fetchApi(apiId);
      expect(api).toBeTruthy();
      if ("flows" in api && api.flows) {
        expect(api.flows.length).toBeGreaterThanOrEqual(1);
      }
    });

    await kubectl.del(fixturePath);
  });
});
