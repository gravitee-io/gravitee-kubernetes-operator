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
 * Journey: make an API observable.
 *
 * As an API producer, I turn on analytics for my API so its traffic shows up in
 * the console, ship its logs and traces to OpenTelemetry, and later turn tracing
 * verbose to debug a problem without republishing the API.
 *
 * Two scenarios rather than one variant table, because the product splits the
 * surface in two and an API cannot cross the split with an update: `tracing` and
 * `logging` are documented "not for native APIs", while `reporterMetricsEnabled`
 * applies only to native ones. Each scenario asserts the WHOLE analytics block
 * it declares, then changes exactly one field and asserts it again — a nested
 * block can round-trip to APIM with defaults silently substituted, so a
 * single-field assertion would not notice.
 *
 * The reporter flag is asserted OFF first on purpose: the provider defaults it
 * to `true`, so `false` is the value that proves the declaration was actually
 * transmitted rather than defaulted.
 *
 * Fixtures are co-located in this folder. Nothing is asserted at the gateway:
 * analytics is control-plane configuration, and observing it on the data plane
 * needs an OTel collector and a log sink the test cluster does not run.
 */

import path from "node:path";
import { fileURLToPath } from "node:url";
import { test } from "../../../../setup.js";
import { XRAY, TAGS } from "../../../../helpers/tags.js";
import { forEachProvisioner } from "../../../../helpers/for-each-provisioner.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

const here = path.dirname(fileURLToPath(import.meta.url));

/** The one knob: whether OpenTelemetry tracing is verbose. */
interface AnalyticsParams {
  verboseTracing: boolean;
}

/** The one knob: whether the native connection-metrics reporter is running. */
interface ReporterParams {
  reporterMetricsEnabled: boolean;
}

/** Declared identically in both fixtures; only `tracing.verbose` is parameterized. */
function expectedAnalytics(params: AnalyticsParams) {
  return {
    enabled: true,
    otelLogs: { enabled: true },
    tracing: { enabled: true, verbose: params.verboseTracing },
  };
}

forEachProvisioner<AnalyticsParams>(
  {
    title: "Report an API's traffic with analytics, OTel logs and tracing",
    provisioners: {
      gko: gkoScenario<AnalyticsParams>({
        // The API carries the parameterized analytics block, so it is applied by
        // applyParams rather than as a static manifest.
        manifests: [],
        roles: { api: "observability-api" },
        dynamicRoles: ["api"],
        applyParams: async (k, params) => {
          const variant = params.verboseTracing ? "on" : "off";
          await k.apply(path.join(here, `gko/analytics-verbose-${variant}.yaml`));
        },
      }),
      terraform: tfScenario<AnalyticsParams>({
        fixture: path.join(here, "terraform/proxy"),
        toVars: (params) => ({ verbose_tracing: params.verboseTracing }),
      }),
    },
    xray: {
      gko: XRAY.OBSERVABILITY.ANALYTICS_TRACING,
      terraform: XRAY.TERRAFORM.API_OBSERVABILITY_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 90_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    await test.step("APIM records the whole analytics block", async () => {
      await mapi.waitForApiMatches(
        apiId,
        { state: "STARTED", analytics: expectedAnalytics({ verboseTracing: false }) },
        { timeoutMs: 30_000 },
      );
    });

    await test.step("Turning tracing verbose updates the block in place", async () => {
      await provisioned.update({ verboseTracing: true });
      await mapi.waitForApiMatches(
        apiId,
        { analytics: expectedAnalytics({ verboseTracing: true }) },
        { timeoutMs: 30_000 },
      );
    });
  },
  { verboseTracing: false },
);

forEachProvisioner<ReporterParams>(
  {
    title: "Disable the connection-metrics reporter on a native Kafka API",
    provisioners: {
      gko: gkoScenario<ReporterParams>({
        manifests: [],
        roles: { api: "observability-native-api" },
        dynamicRoles: ["api"],
        applyParams: async (k, params) => {
          const variant = params.reporterMetricsEnabled ? "on" : "off";
          await k.apply(path.join(here, `gko/native-reporter-${variant}.yaml`));
        },
      }),
      terraform: tfScenario<ReporterParams>({
        fixture: path.join(here, "terraform/native"),
        toVars: (params) => ({ reporter_metrics_enabled: params.reporterMetricsEnabled }),
      }),
    },
    xray: {
      gko: XRAY.OBSERVABILITY.NATIVE_REPORTER_METRICS,
      terraform: XRAY.TERRAFORM.NATIVE_REPORTER_METRICS_TF,
    },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },
    timeoutMs: { gko: 90_000 },
  },
  async ({ provisioned, mapi }) => {
    const apiId = await provisioned.apiId();

    await test.step("APIM records the reporter as disabled", async () => {
      await mapi.waitForApiMatches(
        apiId,
        { type: "NATIVE", analytics: { enabled: true, reporterMetricsEnabled: false } },
        { timeoutMs: 30_000 },
      );
    });

    await test.step("Re-enabling the reporter updates the API in place", async () => {
      await provisioned.update({ reporterMetricsEnabled: true });
      await mapi.waitForApiMatches(
        apiId,
        { analytics: { enabled: true, reporterMetricsEnabled: true } },
        { timeoutMs: 30_000 },
      );
    });
  },
  { reporterMetricsEnabled: false },
);
