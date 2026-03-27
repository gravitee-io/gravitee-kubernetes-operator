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

import {
  LOG,
  PROJECT_DIR,
  HELM,
  isEmptyString,
  toggleVerbosity,
  time,
} from "./lib/index.mjs";

const VERSION = argv.version;
const IMG = argv.img;
const CHANNELS = argv.channels || "alpha";
const DEFAULT_CHANNEL = argv["default-channel"] || "alpha";

toggleVerbosity(argv.verbose);

if (isEmptyString(VERSION)) {
  LOG.red("--version is required");
  process.exit(1);
}
if (isEmptyString(IMG)) {
  LOG.red("--img is required");
  process.exit(1);
}

const BUNDLE_DIR = path.join(PROJECT_DIR, "olm", "bundle");
const MANIFESTS_DIR = path.join(BUNDLE_DIR, "manifests");
const METADATA_DIR = path.join(BUNDLE_DIR, "metadata");
const SCORECARD_DIR = path.join(BUNDLE_DIR, "tests", "scorecard");

const SA_NAME = "gko-controller-manager";

const CLUSTER_ROLE_NAMES = new Set([
  `${SA_NAME}-cluster-role`,
  `${SA_NAME}-validation-webhook-cluster-role`,
  `${SA_NAME}-metrics-cluster-role`,
]);

LOG.magenta(`
  Generating OLM bundle for GKO v${VERSION}
    Image:           ${IMG}:${VERSION}
    Channels:        ${CHANNELS}
    Default channel: ${DEFAULT_CHANNEL}
`);

await time(async () => {
  LOG.blue("Rendering Helm chart ...");
  const rendered = await $`helm template gko ${HELM.chartDir} \
    -n gko-system \
    --set manager.image.repository=${IMG} \
    --set manager.image.tag=${VERSION} \
    --set gatewayAPI.controller.enabled=true`;

  const resources = YAML.parseAllDocuments(rendered.stdout)
    .map((d) => d.toJSON())
    .filter(Boolean);

  LOG.green(`  Parsed ${resources.length} resources`);

  const deployment = resources.find(
    (r) => r.kind === "Deployment" && r.metadata?.name === SA_NAME,
  );
  if (!deployment) throw new Error("Deployment not found in rendered chart");

  const clusterRoles = resources.filter(
    (r) =>
      r.kind === "ClusterRole" &&
      CLUSTER_ROLE_NAMES.has(r.metadata?.name?.trim()),
  );
  if (clusterRoles.length === 0) throw new Error("No ClusterRoles found");

  const leaderRole = resources.find(
    (r) =>
      r.kind === "Role" &&
      r.metadata?.name?.trim() === `${SA_NAME}-leader-election-role`,
  );
  if (!leaderRole) throw new Error("Leader election Role not found");

  const configMap = resources.find(
    (r) => r.kind === "ConfigMap" && r.metadata?.name === "gko-config",
  );
  if (!configMap) throw new Error("ConfigMap gko-config not found");

  const validatingWHC = resources.find(
    (r) => r.kind === "ValidatingWebhookConfiguration",
  );
  const mutatingWHC = resources.find(
    (r) => r.kind === "MutatingWebhookConfiguration",
  );

  LOG.green(
    `  Found Deployment, ${clusterRoles.length} ClusterRoles, leader election Role`,
  );
  LOG.green(
    `  Found ConfigMap, ${validatingWHC ? "ValidatingWHC" : "no ValidatingWHC"}, ${mutatingWHC ? "MutatingWHC" : "no MutatingWHC"}`,
  );

  LOG.blue("Assembling CSV ...");
  const csv = YAML.parse(
    await fs.readFile(
      path.join(PROJECT_DIR, "olm", "gko.clusterserviceversion.yaml"),
      "utf8",
    ),
  );

  csv.metadata.name = `gko.v${VERSION}`;
  csv.metadata.annotations.containerImage = `${IMG}:${VERSION}`;
  csv.spec.version = VERSION;

  const deploySpec = deployment.spec;
  const container = deploySpec.template.spec.containers.find(
    (c) => c.name === "manager",
  );

  const EXCLUDED_ENV_VARS = new Set([
    "NAMESPACE",
    "APPLY_CRDS",
    "ENABLE_GATEWAY_API",
    "APPLY_GATEWAY_API_CRDS",
    "WEBHOOK_CERT_SECRET_NAME",
    "WEBHOOK_NAMESPACE",
    "WEBHOOK_SERVICE_NAME",
    "WEBHOOK_VALIDATING_CONFIGURATION_NAME",
    "WEBHOOK_MUTATING_CONFIGURATION_NAME",
  ]);

  const cmData = configMap.data || {};
  container.env = Object.entries(cmData)
    .filter(([key]) => !EXCLUDED_ENV_VARS.has(key))
    .map(([name, value]) => ({ name, value: String(value) }));

  container.env.push({
    name: "NAMESPACE",
    valueFrom: {
      fieldRef: { fieldPath: "metadata.annotations['olm.targetNamespaces']" },
    },
  });

  delete container.envFrom;

  if (container.volumeMounts) {
    container.volumeMounts = container.volumeMounts.filter(
      (vm) => vm.name !== "webhook-cert",
    );
  }
  if (!container.volumeMounts) container.volumeMounts = [];
  container.volumeMounts.push({ name: "tmp", mountPath: "/tmp" });

  const podSpec = deploySpec.template.spec;
  if (podSpec.volumes) {
    podSpec.volumes = podSpec.volumes.filter((v) => v.name !== "webhook-cert");
  }
  if (!podSpec.volumes) podSpec.volumes = [];
  podSpec.volumes.push({ name: "tmp", emptyDir: {} });

  csv.spec.install.spec.deployments = [{ name: SA_NAME, spec: deploySpec }];
  csv.spec.install.spec.clusterPermissions = [
    {
      serviceAccountName: SA_NAME,
      rules: clusterRoles.flatMap((cr) => cr.rules),
    },
  ];
  csv.spec.install.spec.permissions = [
    { serviceAccountName: SA_NAME, rules: leaderRole.rules },
  ];

  const webhookDefs = [];
  const mapWebhooks = (whc, type) => {
    if (!whc?.webhooks) return;
    const suffix =
      type === "ValidatingAdmissionWebhook" ? "validate" : "mutate";
    for (const wh of whc.webhooks) {
      webhookDefs.push({
        type,
        deploymentName: SA_NAME,
        containerPort: 443,
        targetPort: 9443,
        generateName: `${wh.name}.${suffix}`,
        webhookPath: wh.clientConfig?.service?.path,
        admissionReviewVersions: wh.admissionReviewVersions,
        failurePolicy: wh.failurePolicy,
        sideEffects: wh.sideEffects,
        rules: wh.rules,
      });
    }
  };

  mapWebhooks(validatingWHC, "ValidatingAdmissionWebhook");
  mapWebhooks(mutatingWHC, "MutatingAdmissionWebhook");

  if (webhookDefs.length > 0) {
    csv.spec.webhookdefinitions = webhookDefs;
    LOG.green(`  Added ${webhookDefs.length} webhook definitions to CSV`);
  }

  LOG.blue("Writing bundle ...");
  await fs.ensureDir(MANIFESTS_DIR);
  await fs.ensureDir(METADATA_DIR);
  await fs.ensureDir(SCORECARD_DIR);

  await fs.writeFile(
    path.join(MANIFESTS_DIR, "gko.clusterserviceversion.yaml"),
    YAML.stringify(csv),
  );

  const crdFiles = await fs.readdir(HELM.crdDir);
  for (const file of crdFiles) {
    await fs.copy(path.join(HELM.crdDir, file), path.join(MANIFESTS_DIR, file));
  }
  LOG.green(`  Copied ${crdFiles.length} CRDs to bundle/manifests/`);

  await fs.writeFile(
    path.join(METADATA_DIR, "annotations.yaml"),
    YAML.stringify({
      annotations: {
        "operators.operatorframework.io.bundle.mediatype.v1": "registry+v1",
        "operators.operatorframework.io.bundle.manifests.v1": "manifests/",
        "operators.operatorframework.io.bundle.metadata.v1": "metadata/",
        "operators.operatorframework.io.bundle.package.v1":
          "gravitee-kubernetes-operator",
        "operators.operatorframework.io.bundle.channels.v1": CHANNELS,
        "operators.operatorframework.io.bundle.channel.default.v1":
          DEFAULT_CHANNEL,
      },
    }),
  );

  await fs.writeFile(
    path.join(SCORECARD_DIR, "config.yaml"),
    YAML.stringify({
      apiVersion: "scorecard.operatorframework.io/v1alpha3",
      kind: "Configuration",
      metadata: { name: "config" },
      stages: [
        {
          parallel: true,
          tests: [
            {
              image: "quay.io/operator-framework/scorecard-test:v1.26.0",
              entrypoint: ["scorecard-test", "basic-check-spec"],
              labels: { suite: "basic", test: "basic-check-spec-test" },
            },
            {
              image: "quay.io/operator-framework/scorecard-test:v1.26.0",
              entrypoint: ["scorecard-test", "olm-bundle-validation"],
              labels: { suite: "olm", test: "olm-bundle-validation-test" },
            },
          ],
        },
      ],
    }),
  );

  const totalFiles = (await fs.readdir(MANIFESTS_DIR)).length;
  LOG.green(`
  Bundle generated:
    ${MANIFESTS_DIR}/ (${totalFiles} files: 1 CSV + ${crdFiles.length} CRDs)
    ${METADATA_DIR}/annotations.yaml
    ${SCORECARD_DIR}/config.yaml
  `);
});
