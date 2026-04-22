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
 * Xray test ID registry.
 *
 * Use these constants in test titles for filtering with --grep:
 *   test(`should start API ${XRAY.API_LIFECYCLE.START_STOP}`, ...)
 *
 * Run a single Xray test:
 *   npm run e2e -- --grep @GKO-1464
 *
 * Run regression pack:
 *   npm run e2e:regression
 */

export const XRAY = {
  API_LIFECYCLE: {
    DEPLOY_V4_SYNC_K8S: "@GKO-69",
    DEPLOY_V4_SYNC_FROM_MGMT: "@GKO-71",
    DEPLOY_V4_DB_LESS: "@GKO-81",
    DEPLOY_V4_WITH_LABELS_CATEGORIES: "@GKO-83",
    STOP_V4_API: "@GKO-124",
    START_STOPPED_V4_API: "@GKO-126",
    DELETE_V4_API: "@GKO-140",
    UPDATE_V4_MESSAGE_API: "@GKO-141",
    UPDATE_V4_MESSAGE_MISSING_FIELDS: "@GKO-142",
    RECREATE_DELETED_V4_CLOSED_PLAN: "@GKO-159",
    MISSING_REQUIRED_FIELDS_V4_PROXY: "@GKO-165",
    NO_DEPLOY_WHEN_NO_CHANGES: "@GKO-176",
    REDEPLOY_AFTER_DELETE: "@GKO-212",
    CONTEXT_PATH_CONFLICT_V4: "@GKO-469",
    NO_PLANS_STARTED_V4: "@GKO-502",
    NO_PLANS_STOPPED_V4: "@GKO-503",
    FAILOVER_V4_PROXY: "@GKO-859",
    V4_SURVIVES_UPGRADE: "@GKO-1061",
    START_STOP_V2_V4_NATIVE: "@GKO-1464",
    POLICY_ON_API_WITHOUT_PLANS: "@GKO-1465",
    ENTRYPOINT_POLICY_MATRIX: "@GKO-1474",
    V2_V4_COMPATIBILITY: "@GKO-1448",
    // Batch 8 — visibility & lifecycleState
    V4_VISIBILITY_PRIVATE: "@GKO-172",
    V4_VISIBILITY_PUBLIC: "@GKO-173",
    V4_PUBLISHED_IN_PORTAL: "@GKO-179",
    V4_UNPUBLISHED_NOT_IN_PORTAL: "@GKO-180",
    V4_PORTAL_VISIBILITY_RULES: "@GKO-1466",
    // GKO-1220 (auto-associated groups preserved across updates) — dropped
    // from batch 8: depends on an APIM env-level "auto-associate group to
    // new APIs" setting that is not exposed via mAPI or the test harness.
    // Originally skipped in batch 7. Re-evaluate when the setting is
    // configurable through the test cluster bootstrap.
  },
  MESSAGE_APIS: {
    DEPLOY_V4_MSG_SYNC_MGMT: "@GKO-72",
    DEPLOY_V4_MSG_SYNC_K8S: "@GKO-73",
    HTTP_GET_ENTRYPOINT: "@GKO-129",
    HTTP_POST_ENTRYPOINT: "@GKO-130",
    SSE_ENTRYPOINT: "@GKO-132",
    WEBHOOK_ENTRYPOINT: "@GKO-133",
    WEBSOCKET_ENTRYPOINT: "@GKO-134",
    MOCK_ENDPOINT: "@GKO-136",
    MQTT_ENDPOINT: "@GKO-137",
    MSG_API_WITH_POLICY: "@GKO-164",
  },
  PLANS: {
    KEYLESS_PLAN_V4: "@GKO-110",
    DELETE_PLAN_SYNC_MGMT: "@GKO-113",
    DELETE_PLAN_SYNC_K8S: "@GKO-117",
    MULTIPLE_PLANS: "@GKO-160",
    APIKEY_PLAN_V4: "@GKO-161",
    OAUTH2_PLAN_V4: "@GKO-162",
    JWT_PLAN_V4: "@GKO-163",
    PUBLISHED_PLAN: "@GKO-169",
    FAIL_NON_PUBLISHED_STATUS: "@GKO-170",
    FAIL_CHANGE_PUBLISHED_STATUS: "@GKO-171",
    PLAN_PUBLISHED_WHEN_API_STOPPED: "@GKO-174",
    GENERAL_CONDITIONS: "@GKO-238",
    PLAN_LIFECYCLE_VIA_CR: "@GKO-1459",
    NATIVE_KAFKA_KEYLESS_PLAN: "@GKO-856",
    NATIVE_ADD_PLAN: "@GKO-918",
    NATIVE_REMOVE_PLAN: "@GKO-919",
    NATIVE_UPDATE_PLAN: "@GKO-920",
    NATIVE_MULTIPLE_PLANS: "@GKO-921",
    NATIVE_NO_PLAN: "@GKO-922",
  },
  VALIDATION: {
    V4_CONTEXT_PATH_CONFLICT: "@GKO-1476",
    V4_OAS_COMPLIANCE_WEBHOOK: "@GKO-1479",
    V4_DEFAULT_VALUES: "@GKO-1480",
    NON_EXISTING_GROUP_MESSAGE: "@GKO-1478",
  },
  V2_API_LIFECYCLE: {
    V2_UPDATE_API_PATH: "@GKO-1065",
    V2_MEMBER_ROLE_CHANGE_DUPLICATE_KEY: "@GKO-260",
    V2_MGMT_CTX_VALID_ON_CREATE: "@GKO-594",
    V2_MGMT_CTX_VALID_ON_UPDATE: "@GKO-597",
    V2_IMPORT_NON_EXISTING_CATEGORY_DRYRUN: "@GKO-605",
    V2_NO_PLANS_STARTED: "@GKO-606",
    V2_NO_PLANS_STOPPED: "@GKO-607",
    // GKO-653 (V2 exported read-only round-trip) — removed from batch 5:
    // APIM does not support CRD export for V2 APIs (/management/v2/.../_export/crd
    // returns 400 "definition version 2.0.0 is not supported anymore").
    // Tracked in "Batch 5 - Skipped Tests.md".
  },
  APPLICATIONS_MEMBERS: {
    APP_NON_EXISTING_MEMBER: "@GKO-533",
    APP_MEMBER_NO_ROLE: "@GKO-535",
    APP_MEMBER_NO_SOURCE: "@GKO-536",
    APP_NON_EXISTING_GROUP: "@GKO-548",
    APP_MEMBER_NON_EXISTING_ROLE: "@GKO-555",
    APP_WEB_REQUIRES_AUTH_CODE: "@GKO-581",
    APP_NON_EXISTING_ROLE: "@GKO-531",
    APP_REMOVE_MEMBER: "@GKO-534",
    APP_ADD_MEMBER_ROLE_NAME: "@GKO-538",
    APP_CHANGE_MEMBER_ROLE: "@GKO-539",
  },
  WEBHOOKS: {
    REJECT_INVALID_CRS: "@GKO-1447",
    NON_OAS_ERRORS_V4: "@GKO-77",
    NON_OAS_COMPLIANT_V4: "@GKO-76",
    INVALID_CREDENTIALS_CONTEXT: "@GKO-78",
    V2_PARENT_PATH_NOT_FOUND: "@GKO-153",
    MISSING_FIELDS_V4_MESSAGE: "@GKO-166",
    V4_PARENT_PATH_NOT_FOUND: "@GKO-281",
    NON_EXISTING_CONTEXT_V4: "@GKO-414",
    NON_EXISTING_MGMT_CONTEXT: "@GKO-465",
    MGMT_CONTEXT_INVALID_CREDS: "@GKO-474",
    RESOURCE_NO_NAME: "@GKO-515",
    RESOURCE_NO_TYPE: "@GKO-516",
    RESOURCE_NO_CONFIG: "@GKO-518",
    RESOURCE_INVALID_CONFIG: "@GKO-519",
    V4_INVALID_CRON: "@GKO-520",
    V2_CONTEXT_PATH_DUPLICATE: "@GKO-590",
    V2_CONTEXT_PATH_CONFLICT_V4: "@GKO-591",
    V2_CONTEXT_PATH_EXISTS_LOCAL_FALSE: "@GKO-609",
    V2_INVALID_CRON: "@GKO-614",
    CROSS_VERSION_SCHEDULERS_FETCHERS: "@GKO-1475",
    CROSS_VERSION_PARENT_PATH: "@GKO-1477",
  },
  APPLICATIONS: {
    APP_WITH_METADATA: "@GKO-194",
    CREATE_APP: "@GKO-335",
    UPDATE_APP: "@GKO-336",
    DELETE_APP: "@GKO-337",
    APP_NO_MGMT_CONTEXT: "@GKO-526",
    APP_NEW_CRD_ATTRIBUTES: "@GKO-527",
    APP_BOTH_SETTINGS_ERROR: "@GKO-550",
    APP_CONFIGURE_SETTINGS: "@GKO-552",
    APP_PO_IN_MEMBERS_ERROR: "@GKO-558",
    APP_SWITCH_SETTINGS_TYPE: "@GKO-561",
    APP_UPDATE_OAUTH_TYPE: "@GKO-562",
    APP_CLIENT_ID_OPTIONAL: "@GKO-563",
    APP_CLIENT_ID_UNIQUE: "@GKO-564",
    APP_PO_ROLE_OVERWRITE: "@GKO-567",
    // Batch 8 — admission edge cases & lifecycle
    APP_READ_ONLY_IN_APIM: "@GKO-505",
    APP_BROWSER_VALID_URIS: "@GKO-578",
    APP_SPA_GRANT_TYPES: "@GKO-579",
    APP_NAME_LENGTH_EDGE: "@GKO-1382",
    APP_DELETE_SUCCESS: "@GKO-1383",
  },
  SUBSCRIPTIONS: {
    ENDING_BEFORE_START: "@GKO-807",
    SYNC_FROM_K8S_ERROR_V4: "@GKO-816",
    AUTO_VALIDATE_V2: "@GKO-797",
    V2_JWT_SUBSCRIPTION: "@GKO-799",
    V4_JWT_SUBSCRIPTION: "@GKO-800",
    V2_GATEWAY_JWT_CALL: "@GKO-808",
    AUTO_VALIDATE_V4: "@GKO-815",
    V4_GATEWAY_JWT_CALL: "@GKO-817",
    V2_OAUTH2_SUBSCRIPTION: "@GKO-818",
    V4_OAUTH2_SUBSCRIPTION: "@GKO-819",
    API_MUST_BE_STARTED: "@GKO-840",
    PLAN_MUST_MATCH_V2: "@GKO-842",
    PLAN_MUST_MATCH_V4: "@GKO-843",
    SECURITY_TYPE_JWT_OAUTH2: "@GKO-844",
    API_KIND_DEFAULT: "@GKO-845",
    ERROR_UPDATE_PLAN_WITH_SUB: "@GKO-848",
    ERROR_DELETE_API_WITH_SUB: "@GKO-849",
    ERROR_DELETE_APP_WITH_SUB: "@GKO-853",
    DELETE_API_WITH_OTHER_PLAN: "@GKO-854",
    MTLS_PLAN_V4: "@GKO-869",
    V4_SUBSCRIPTION_READ_ONLY: "@GKO-795",
    SUBSCRIPTION_IMMUTABILITY: "@GKO-1460",
    V4_JWT_PLAN_DELETION_WITH_SUB: "@GKO-822",
    V4_OAUTH2_PLAN_DELETION_WITH_SUB: "@GKO-826",
    CROSS_MGMT_CONTEXT_ERROR: "@GKO-796",
    V2_LOCAL_SUBSCRIPTION_ERROR: "@GKO-798",
    V2_JWT_DELETE: "@GKO-821",
    V2_OAUTH2_DELETE: "@GKO-825",
    API_APP_SYNCED_LAST_RECONCILE: "@GKO-839",
  },
  // NATIVE_APIS group removed — 10 tests were never committed due to the
  // APIM native-plan serialization bug (see "Batch 3 - Skipped Tests.md").
  // GKO-874 and GKO-875 were listed in batch 2 but their test files were
  // also never committed; they are tracked as not-yet-covered.
  GROUPS: {
    CREATE_WITH_MEMBER: "@GKO-983",
    CREATE_NON_EXISTING_USER: "@GKO-984",
    DELETE_GROUP: "@GKO-985",
    MODIFY_GROUP: "@GKO-986",
    CREATE_WITHOUT_ROLES: "@GKO-987",
    PREVENT_PO_GROUP_AS_MEMBER: "@GKO-974",
  },
  SHARED_POLICY_GROUPS: {
    ADD_SPG_TO_API: "@GKO-976",
    REMOVE_SPG_FROM_API: "@GKO-980",
    UPDATE_SPG: "@GKO-981",
    SPG_LIFECYCLE: "@GKO-1462",
  },
  DEPLOYMENT_RECONCILIATION: {
    RECONCILE_API_CONFIG: "@GKO-1444",
    ACCEPTED_FALSE_ON_FAILURE: "@GKO-1387",
    ACCEPTED_NOT_FALSE_ON_SUCCESS: "@GKO-1388",
    ACCEPTED_UPDATES_ON_CHANGE: "@GKO-1389",
    PROCESSING_STATUS_PRESENT: "@GKO-1390",
    IDEMPOTENT_RECONCILIATION: "@GKO-1445",
    STATUS_CONDITIONS_REFLECT_STATE: "@GKO-1446",
    MGMT_CTX_CONDITION_VOCABULARY: "@GKO-1282",
    CONSISTENT_CONDITION_STRUCTURE: "@GKO-1283",
    PROCESSING_STATUS_DEPRECATED: "@GKO-1391",
    AUDITABILITY_EVENTS: "@GKO-1463",
    RECOVERY_REAPPLY: "@GKO-1808",
    CR_MANAGED_READ_ONLY: "@GKO-1456",
    OPERATOR_RESTART_RECOVERY: "@GKO-1451",
  },
  MANAGEMENT_CONTEXT: {
    NON_EXISTING_ENV: "@GKO-472",
    NON_EXISTING_ORG: "@GKO-473",
    INVALID_CREDENTIALS: "@GKO-474",
    DELETE_WITH_V2_API_REF: "@GKO-892",
    DELETE_WITH_V4_API_REF: "@GKO-893",
    DELETE_WITH_APP_REF: "@GKO-894",
    DELETE_NO_REFS: "@GKO-895",
  },
  TEMPLATING: {
    V2_MISSING_CONFIGMAP: "@GKO-676",
    V2_MISSING_KEY: "@GKO-677",
    V4_MISSING_CONFIGMAP: "@GKO-678",
    V4_MISSING_KEY: "@GKO-679",
    V2_SECRET_VALUE: "@GKO-682",
    V4_CONFIGMAP_VALUE: "@GKO-683",
    APP_CONFIGMAP_VALUE: "@GKO-684",
    MGMT_CONTEXT_BEARER_TOKEN: "@GKO-781",
  },
  IMPORT_EXPORT: {
    EXPORTED_POLICIES_V4: "@GKO-93",
    IMPORT_V4_CRD: "@GKO-218",
    EXPORT_V4_CRD: "@GKO-229",
    K8S_COMPLIANT_NAMES: "@GKO-231",
    EXPORT_V2_CRD: "@GKO-301",
    IMPORT_V2_CRD: "@GKO-303",
    V4_NO_EMAIL_ON_EXPORT: "@GKO-237",
    V4_METADATA_ON_IMPORT: "@GKO-239",
    V4_EXPORT_IMPORT_ROUND_TRIP: "@GKO-1472",
    V4_TERRAFORM_IMPORT_EXPORT_PARITY: "@GKO-1927",
    // GKO-1471 (V2 import/export round-trip) — dropped from batch 5:
    // APIM does not support CRD export for V2 APIs. Tracked in
    // "Batch 5 - Skipped Tests.md".
  },
  MEMBERS: {
    V4_NON_EXISTING_MEMBER: "@GKO-251",
    V4_NON_EXISTING_GROUP: "@GKO-252",
    V4_REMOVE_MEMBER: "@GKO-253",
    V4_MEMBER_NO_ROLE: "@GKO-254",
    V4_MEMBER_NO_SOURCE: "@GKO-255",
    V4_NON_EXISTING_MEMBERS_CRD: "@GKO-470",
    V4_PO_NOT_ALLOWED: "@GKO-569",
    V4_PO_CANT_OVERWRITE: "@GKO-571",
    V4_REMOVE_MEMBER_VARIANT: "@GKO-213",
    V4_PO_DEFINED_IN_CRD: "@GKO-244",
    V4_ADD_MEMBER_WITH_ROLE_NAME: "@GKO-247",
    V4_ADD_MEMBER_NO_ROLE: "@GKO-249",
    V4_CREATE_NON_EXISTING_GROUP: "@GKO-256",
    V4_CREATE_EXISTING_GROUP: "@GKO-257",
    V4_DUPLICATE_KEY_ON_ROLE_CHANGE: "@GKO-259",
    V4_PO_VIA_MGMT_CONTEXT: "@GKO-306",
    V4_TRANSFER_PRIMARY_OWNER: "@GKO-307",
    V4_ADD_GROUP_REFS: "@GKO-314",
    V4_NOTIFY_MEMBERS_ENABLED: "@GKO-402",
    V4_TAKE_OVER_PO_VIA_MGMT_CTX: "@GKO-658",
    V4_ADD_GROUP_REFS_VARIANT: "@GKO-1004",
    V2_ADD_MEMBER_WITH_ROLE: "@GKO-202",
    V2_NON_EXISTING_ROLE: "@GKO-204",
    V2_NON_EXISTING_MEMBER: "@GKO-205",
    V2_NON_EXISTING_GROUP: "@GKO-207",
    V2_EXISTING_GROUP: "@GKO-208",
    V2_REMOVE_MEMBER: "@GKO-216",
    V2_PO_IN_MEMBERS: "@GKO-258",
    V2_CHANGE_MEMBER_ROLE: "@GKO-308",
    V2_MEMBER_NO_ROLE: "@GKO-393",
    V2_ADD_GROUP_HRID: "@GKO-398",
    V2_MULTIPLE_GROUPS: "@GKO-399",
    V2_REMOVE_GROUP: "@GKO-400",
    V2_NOTIFY_MEMBERS: "@GKO-401",
    V2_PO_NOT_OVERWRITEABLE: "@GKO-601",
    // GKO-602 (V2 API with a different PRIMARY_OWNER is rejected) — dropped
    // from batch 5: GKO admission does not enforce this (product gap). The
    // webhook accepts the CR and the API is created in APIM. Tracked in
    // "Batch 5 - Skipped Tests.md" as a GKO product bug.
    V2_TAKE_OVER_PO: "@GKO-657",
    // GKO-659 (adding PO to members has no effect) — dropped from batch 5
    // because the companion GKO-602 scenario is also dropped, and GKO-601
    // already verifies that re-declaring the mgmt-ctx user as PO is a no-op.
    // Tracked in "Batch 5 - Skipped Tests.md".
    V2_ADD_GROUP_REFS: "@GKO-1003",
    PRIMARY_OWNER_VISIBILITY: "@GKO-1457",
  },
  DEFAULTS: {
    NAMESPACE_DEFAULT: "@GKO-463",
    VALID_NAME_NAMESPACE: "@GKO-466",
    V2_LOCAL_FALSE_DEFAULT: "@GKO-765",
    V4_SYNC_FROM_MGMT_DEFAULT: "@GKO-770",
  },
  POLICIES: {
    DEPLOY_V4_WITH_POLICY: "@GKO-94",
    REMOVE_POLICY: "@GKO-95",
    UPDATE_POLICY: "@GKO-96",
  },
  CATEGORIES: {
    VALID_CATEGORY_V4: "@GKO-267",
    NON_EXISTING_CATEGORY_V4: "@GKO-269",
    REMOVE_CATEGORY_V4: "@GKO-270",
    V4_MANY_CATEGORIES: "@GKO-268",
    V4_CATEGORY_REMOVED_FROM_APIM: "@GKO-271",
    V4_CATEGORY_RENAME_REDEPLOY: "@GKO-272",
    V4_NON_EXISTING_GROUP_REF: "@GKO-471",
    V4_DEPLOY_NON_EXISTING_CATEGORY: "@GKO-412",
    V4_IMPORT_NON_EXISTING_CATEGORY_DRYRUN: "@GKO-415",
    V4_IMPORT_NON_EXISTING_CATEGORY_APPLY: "@GKO-416",
    V2_VALID_CATEGORY: "@GKO-187",
    V2_MANY_CATEGORIES: "@GKO-189",
    V2_NON_EXISTING_CATEGORY: "@GKO-190",
    V2_REMOVE_CATEGORY: "@GKO-191",
    V2_CATEGORY_REMOVED_FROM_APIM: "@GKO-192",
    V2_CATEGORY_RENAME_REDEPLOY: "@GKO-261",
    V4_LABELS_LIFECYCLE: "@GKO-1473",
  },
  // TCP_FAILOVER group removed — tcp-failover.test.ts was never committed
  // (GKO-79 depended on TCP proxy setup that the test cluster does not
  // provide). Tracked in "Batch 3 - Skipped Tests.md".
  TERRAFORM: {
    APPLY_COMPLEX_CONFIG: "@GKO-1926",
    APIM_CONTAINS_ALL_ENTITIES: "@GKO-1929",
    PAGE_HIERARCHY_PRESERVED: "@GKO-1931",
    IDEMPOTENT_CONFIG: "@GKO-1932",
    ADD_APPLICATION: "@GKO-1373",
    ADD_SUBSCRIPTION: "@GKO-1374",
    APP_MISSING_FIELDS: "@GKO-1375",
    INVALID_SUBSCRIPTION_FORMAT: "@GKO-1376",
    REMOVE_APPLICATION: "@GKO-1378",
    CREATE_APP_AND_SUBSCRIPTION: "@GKO-1379",
    VALID_AND_MALFORMED_HCL: "@GKO-1453",
    GENERAL_CONDITIONS_PAGE: "@GKO-1930",
    // Batch 8 — error handling in TF + delete-via-TF lifecycle.
    // GKO-1381 (Role-specific access for managing Apps and subscriptions
    // via Terraform) is dropped from batch 8: the test harness has no
    // mechanism to provision a non-admin APIM user, so role scoping cannot
    // be exercised. Re-evaluate when the bootstrap supports multi-user
    // accounts (the same blocker as GKO-1541).
    INVALID_SUBSCRIPTION_CONFIG: "@GKO-1380",
    DELETE_APPLICATION_TF: "@GKO-1383",
  },
  PAGES: {
    MARKDOWN_PAGE_CRUD_V4: "@GKO-277",
    MARKDOWN_PAGE_UPDATE_V4: "@GKO-278",
    FETCHER_PAGE_V4: "@GKO-279",
    AUTOFETCH_PRESERVED: "@GKO-1933",
    DOC_CRUD_ACROSS_VERSIONS: "@GKO-1468",
    V4_DOC_OPERATIONS: "@GKO-236",
    V4_READ_ONLY_DOC: "@GKO-280",
    V4_DOC_VISIBILITY_PUBLIC: "@GKO-282",
    V4_DOC_RECONCILED: "@GKO-1470",
    V2_DOC_CRUD: "@GKO-146",
    V2_DOC_OVERSIZE: "@GKO-147",
    V2_DOC_INLINE_UPDATE: "@GKO-148",
    V2_DOC_FETCHER: "@GKO-151",
    V2_DOC_PUBLIC: "@GKO-199",
    V2_DOC_PRIVATE_NO_GROUPS: "@GKO-200",
    V2_DOC_PRIVATE_GROUPS: "@GKO-315",
    V2_DOC_PRIVATE_EXCLUDED: "@GKO-316",
    // GKO-662 (delete fetched ROOT pages) — dropped from batch 5: APIM
    // rejects V2 ROOT fetchers backed by http-fetcher with
    // "The plugin does not support to import a directory". ROOT pages need
    // a directory-listing fetcher (e.g. github-fetcher) which in turn
    // requires real GitHub credentials. Tracked in "Batch 5 - Skipped
    // Tests.md".
    V2_DOC_RENAME: "@GKO-699",
    V2_WEB_FETCHER_NO_URL: "@GKO-620",
    V2_WEB_FETCHER_WARNING: "@GKO-621",
    V2_GITHUB_FETCHER_REQUIRED_FIELDS: "@GKO-622",
    V2_GITHUB_FETCHER_WARNING: "@GKO-623",
    V4_WEB_FETCHER_NO_URL: "@GKO-629",
    V4_WEB_FETCHER_WARNING: "@GKO-628",
    V4_GITHUB_FETCHER_REQUIRED_FIELDS: "@GKO-636",
    V4_GITHUB_FETCHER_WARNING: "@GKO-637",
    V4_DOC_RENAME: "@GKO-1469",
    DOC_PUBLIC_ACCESS_CROSS_VERSION: "@GKO-1467",
    // Batch 8 — folder + page rename within a deployed API
    V4_DOC_FOLDER_RENAME: "@GKO-700",
    // GKO-626, 675, 689, 692 — dropped from batch 5 because the GKO
    // admission webhook pre-fetches github-fetcher pages at apply time
    // and the test cluster has no real GitHub credentials, so any positive
    // github-fetcher test is rejected by admission with "Page cannot be
    // fetched, this can come from either invalid / missing github
    // credentials or an invalid file path". Tracked in "Batch 5 - Skipped
    // Tests.md".
  },
  NOTIFICATIONS: {
    REMOVE_NOTIFICATION: "@GKO-1238",
    NOTIFICATION_HOOKS_GROUPS: "@GKO-1231",
    API_REF_NOTIFICATION: "@GKO-1232",
    NOTIFICATIONS_VIA_CRS: "@GKO-1461",
    // GKO-1236 (cannot delete Notification CR if referenced) — dropped from
    // batch 5: GKO does not enforce an in-use protection on Notification
    // CRs. The delete succeeds even when an API references it. Tracked in
    // "Batch 5 - Skipped Tests.md" as a product gap.
    WORKS_WITH_V2_AND_V4: "@GKO-1237",
    NOT_IN_EXPORT: "@GKO-1233",
    DUPLICATE_CONSOLE_REJECTED: "@GKO-1235",
    CR_READONLY_VIA_MAPI: "@GKO-1234",
    // Batch 8 — recipient & label assertions via mAPI
    VIEW_NOTIFICATION_SETTINGS: "@GKO-1194",
    NOTIFICATION_LABEL: "@GKO-1195",
    DEFAULT_RECIPIENT_OWNER: "@GKO-1196",
    PORTAL_NOTIFIER_TARGET_USER: "@GKO-1219",
    GROUP_MEMBERS_NOTIFIED: "@GKO-1239",
  },
  LOCAL_CONFIGMAP: {
    LOCAL_FALSE_NO_CONFIGMAP: "@GKO-765",
    DELETION_FINALIZER_CLEANUP: "@GKO-1452",
  },
  MTLS_CERTIFICATES: {
    ADD_MULTIPLE_CERTS: "@GKO-2243",
    DEPRECATED_FIELD: "@GKO-2244",
    CERT_ROTATION: "@GKO-2231",
    MTLS_SUBSCRIPTION: "@GKO-2248",
    CERT_VALID_DATES: "@GKO-2255",
    CERT_END_DATE: "@GKO-2221",
    ADD_SINGLE_CERT: "@GKO-2212",
    REMOVE_CERT: "@GKO-2247",
    REMOVE_MULTI_CERTS: "@GKO-2250",
    DEPENDENCY_RESOLUTION: "@GKO-1449",
    // ── Batch 8 — bucket I ─────────────────────────────────────
    // Application clientCertificates admission rejections.
    // GKO has no standalone MTLSCertificate CRD; these scenarios are
    // exercised through the Application CR's spec.settings.app.tls.
    CRD_BAD_PEM: "@GKO-2117",
    CRD_FORBIDDEN_FIELD_UPDATE: "@GKO-2118",
    CRD_START_EQ_END: "@GKO-2122",
    CRD_NAME_TOO_LONG: "@GKO-2124",
    CRD_MISSING_FIELDS: "@GKO-2125",
    CRD_EXPIRED_REJECTED: "@GKO-2131",
    CRD_INVALID_CHARS: "@GKO-2133",
    CRD_END_BEFORE_START: "@GKO-2135",
    CRD_MISSING_NAME: "@GKO-2143",
    CRD_MISSING_CERT_FIELD: "@GKO-2146",
    CRD_INVALID_DATA_DATES: "@GKO-2148",
  },
} as const;

export const TAGS = {
  REGRESSION: "@regression",
} as const;
