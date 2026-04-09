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

// ── Enums ─────────────────────────────────────────────────────

export type DefinitionVersion = "V1" | "V2" | "V4" | "FEDERATED" | "FEDERATED_AGENT";
export type ApiState = "CLOSED" | "INITIALIZED" | "STARTED" | "STOPPED" | "STOPPING";
export type ApiVisibility = "PUBLIC" | "PRIVATE";
export type ApiLifecycleState = "CREATED" | "PUBLISHED" | "UNPUBLISHED" | "DEPRECATED" | "ARCHIVED";
export type ApiDeploymentState = "NEED_REDEPLOY" | "DEPLOYED";
export type ApiType = "PROXY" | "MESSAGE" | "NATIVE" | "A2A_PROXY" | "LLM_PROXY" | "MCP_PROXY";
export type ApiWorkflowState = "DRAFT" | "IN_REVIEW" | "REQUEST_FOR_CHANGES" | "REVIEW_OK";

// ── API Types (v2 Management API) ─────────────────────────────

/**
 * Discriminated union of all supported API definition versions.
 * Use DeepPartial<Api> in assertions — partial matching handles all variants.
 */
export type Api = ApiV1 | ApiV2 | ApiV4 | ApiFederated | ApiFederatedAgent;

export interface ApiBase {
  id: string;
  name: string;
  description?: string;
  crossId?: string;
  apiVersion: string;
  definitionVersion: DefinitionVersion;
  state: ApiState;
  visibility: ApiVisibility;
  lifecycleState: ApiLifecycleState;
  deploymentState?: ApiDeploymentState;
  workflowState?: ApiWorkflowState;
  tags?: string[];
  labels?: string[];
  groups?: string[];
  categories?: string[];
  deployedAt?: string;
  createdAt: string;
  updatedAt: string;
  primaryOwner: PrimaryOwner;
  disableMembershipNotifications?: boolean;
  originContext?: OriginContext;
  responseTemplates?: Record<string, Record<string, ResponseTemplate>>;
  resources?: Resource[];
  properties?: Property[];
  _links?: ApiLinks;
}

export interface ApiV1 extends ApiBase {
  definitionVersion: "V1";
}

export interface ApiV4 extends ApiBase {
  definitionVersion: "V4";
  type: ApiType;
  listeners: Listener[];
  endpointGroups: EndpointGroupV4[];
  analytics?: Analytics;
  failover?: FailoverV4;
  flowExecution?: FlowExecution;
  flows?: FlowV4[];
  services?: ApiServices;
  allowedInApiProducts?: boolean;
}

export interface ApiV2 extends ApiBase {
  definitionVersion: "V2";
  environmentId?: string;
  executionMode?: ExecutionMode;
  contextPath?: string;
  proxy: Proxy;
  flowMode?: FlowMode;
  flows?: FlowV2[];
  services?: ApiServicesV2;
  pathMappings?: string[];
  entrypoints?: ApiEntrypoint[];
}

export interface ApiFederated extends ApiBase {
  definitionVersion: "FEDERATED";
}

/**
 * Agent-sourced federated API (A2A protocol).
 */
export interface ApiFederatedAgent extends ApiBase {
  definitionVersion: "FEDERATED_AGENT";
  url?: string;
  documentationUrl?: string;
  provider?: A2AProvider;
  defaultInputModes?: string[];
  defaultOutputModes?: string[];
  capabilities?: string[];
  securitySchemes?: Record<string, unknown>;
  security?: Record<string, unknown>;
  skills?: A2ASkill[];
}

export interface A2AProvider {
  organization?: string;
  url?: string;
}

export interface A2ASkill {
  id?: string;
  name?: string;
  description?: string;
  tags?: string[];
  examples?: string[];
  inputModes?: string[];
  outputModes?: string[];
}

// ── Origin Context ────────────────────────────────────────────

export type OriginContext = ManagementOriginContext | KubernetesOriginContext | IntegrationOriginContext;

export interface ManagementOriginContext {
  origin: "MANAGEMENT";
}

export interface KubernetesOriginContext {
  origin: "KUBERNETES";
  mode?: "FULLY_MANAGED";
  syncFrom?: "MANAGEMENT" | "KUBERNETES";
}

export interface IntegrationOriginContext {
  origin: "INTEGRATION";
  integrationId?: string;
}

// ── API Supporting Types ──────────────────────────────────────

export interface ResponseTemplate {
  statusCode?: number;
  headers?: Record<string, string>;
  body?: string;
}

export interface Resource {
  name: string;
  type: string;
  configuration?: Record<string, unknown>;
  enabled?: boolean;
}

export interface Property {
  key: string;
  value?: string;
  encrypted?: boolean;
  dynamic?: boolean;
  encryptable?: boolean;
}

export interface ApiLinks {
  pictureUrl?: string;
  backgroundUrl?: string;
}

export interface ApiEntrypoint {
  target?: string;
  host?: string;
  tags?: string[];
}

export type ExecutionMode = "V3" | "V4_EMULATION_ENGINE";

// ── Listener Types (discriminated union) ──────────────────────

export type ListenerType = "HTTP" | "SUBSCRIPTION" | "TCP" | "KAFKA";

export type Listener = HttpListener | SubscriptionListener | TcpListener | KafkaListener;

export interface BaseListener {
  type: ListenerType;
  entrypoints?: Entrypoint[];
  servers?: string[];
}

export interface HttpListener extends BaseListener {
  type: "HTTP";
  paths?: PathV4[];
  pathMappings?: string[];
  cors?: Cors;
}

export interface SubscriptionListener extends BaseListener {
  type: "SUBSCRIPTION";
}

export interface TcpListener extends BaseListener {
  type: "TCP";
  hosts?: string[];
}

export interface KafkaListener extends BaseListener {
  type: "KAFKA";
  host?: string;
  port?: number;
}

// ── Entrypoint ────────────────────────────────────────────────

export type Qos = "NONE" | "AUTO" | "AT_MOST_ONCE" | "AT_LEAST_ONCE";

export interface Dlq {
  endpoint?: string;
}

export interface Entrypoint {
  type: string;
  qos?: Qos;
  dlq?: Dlq;
  configuration?: Record<string, unknown>;
}

// ── PathV4 ────────────────────────────────────────────────────

export interface PathV4 {
  host?: string;
  path?: string;
  overrideAccess?: boolean;
}

// ── EndpointGroup / Endpoint (V4) ─────────────────────────────

export interface EndpointGroupV4 {
  name: string;
  type: string;
  loadBalancer?: LoadBalancer;
  sharedConfiguration?: Record<string, unknown>;
  endpoints?: EndpointV4[];
  services?: EndpointGroupServices;
}

export interface EndpointV4 {
  name: string;
  type: string;
  weight?: number;
  inheritConfiguration?: boolean;
  configuration?: Record<string, unknown>;
  sharedConfigurationOverride?: Record<string, unknown>;
  services?: EndpointServices;
  secondary?: boolean;
  tenants?: string[];
}

export interface EndpointGroupServices {
  healthCheck?: ServiceV4;
  discovery?: ServiceV4;
}

export interface EndpointServices {
  healthCheck?: ServiceV4;
}

export interface ServiceV4 {
  enabled?: boolean;
  type?: string;
  configuration?: Record<string, unknown>;
  overrideConfiguration?: boolean;
}

export interface LoadBalancer {
  type: LoadBalancerType;
}

export type LoadBalancerType = "RANDOM" | "ROUND_ROBIN" | "WEIGHTED_RANDOM" | "WEIGHTED_ROUND_ROBIN";

// ── Failover / FlowExecution ──────────────────────────────────

export interface FailoverV4 {
  enabled?: boolean;
  maxRetries?: number;
  slowCallDuration?: number;
  openStateDuration?: number;
  maxFailures?: number;
  perSubscription?: boolean;
}

export interface FlowExecution {
  mode?: FlowMode;
  matchRequired?: boolean;
}

export type FlowMode = "BEST_MATCH" | "DEFAULT";

// ── API Services ──────────────────────────────────────────────

export interface ApiServices {
  dynamicProperty?: ServiceV4;
}

export interface ApiServicesV2 {
  dynamicProperty?: ServiceV2;
  healthCheck?: ServiceV2;
}

export interface ServiceV2 {
  enabled?: boolean;
  type?: string;
  configuration?: Record<string, unknown>;
}

// ── Cors ──────────────────────────────────────────────────────

export interface Cors {
  enabled?: boolean;
  allowOrigin?: string[];
  allowMethods?: string[];
  allowHeaders?: string[];
  allowCredentials?: boolean;
  maxAge?: number;
  exposeHeaders?: string[];
  runPolicies?: boolean;
}

// ── Plan Types (discriminated union) ──────────────────────────

export type Plan = PlanV2 | PlanV4 | PlanFederated;

export interface GenericPlan {
  id: string;
  name: string;
  description?: string;
  apiId: string;
  apiProductId?: string;
  definitionVersion: DefinitionVersion;
  status: PlanStatus;
  security: PlanSecurity;
  mode?: PlanMode;
  type?: PlanType;
  validation: PlanValidation;
  order?: number;
  characteristics?: string[];
  commentRequired?: boolean;
  commentMessage?: string;
  generalConditions?: string;
  crossId?: string;
  excludedGroups?: string[];
  selectionRule?: string;
  tags?: string[];
  createdAt: string;
  updatedAt: string;
  publishedAt?: string;
  closedAt?: string;
}

export interface PlanV4 extends GenericPlan {
  definitionVersion: "V4";
  flows?: FlowV4[];
}

export interface PlanV2 extends GenericPlan {
  definitionVersion: "V2";
  flows?: FlowV2[];
  paths?: Record<string, Rule[]>;
}

export interface PlanFederated extends GenericPlan {
  definitionVersion: "FEDERATED";
}

export type PlanStatus = "STAGING" | "PUBLISHED" | "DEPRECATED" | "CLOSED";
export type PlanMode = "STANDARD" | "PUSH";
export type PlanType = "API" | "CATALOG";

export interface PlanSecurity {
  type: PlanSecurityType;
  configuration?: Record<string, unknown>;
}

export type PlanSecurityType = "KEY_LESS" | "API_KEY" | "JWT" | "OAUTH2" | "MTLS";
export type PlanValidation = "AUTO" | "MANUAL";

export interface Rule {
  methods?: HttpMethod[];
  description?: string;
  enabled?: boolean;
  policy?: string;
  configuration?: Record<string, unknown>;
}

export type HttpMethod = "CONNECT" | "DELETE" | "GET" | "HEAD" | "OPTIONS" | "PATCH" | "POST" | "PUT" | "TRACE" | "OTHER";

// ── Subscription Types ────────────────────────────────────────

export interface Subscription {
  id: string;
  api: BaseApi;
  apiProduct?: BaseApiProduct;
  plan: BasePlan;
  application: BaseApplication;
  status: SubscriptionStatus;
  consumerStatus: SubscriptionConsumerStatus;
  consumerMessage?: string;
  publisherMessage?: string;
  processedBy?: BaseUser;
  subscribedBy?: BaseUser;
  processedAt?: string;
  startingAt?: string;
  endingAt?: string;
  closedAt?: string;
  pausedAt?: string;
  consumerPausedAt?: string;
  metadata?: Record<string, string>;
  daysToExpirationOnLastNotification?: number;
  consumerConfiguration?: SubscriptionConsumerConfiguration;
  failureCause?: string;
  origin?: "KUBERNETES" | "MANAGEMENT";
  createdAt: string;
  updatedAt: string;
}

export type SubscriptionStatus = "PENDING" | "ACCEPTED" | "CLOSED" | "REJECTED" | "PAUSED" | "RESUMED";
export type SubscriptionConsumerStatus = "STARTED" | "STOPPED" | "FAILURE";

/** @deprecated Use SubscriptionConsumerStatus instead */
export type ConsumerStatus = SubscriptionConsumerStatus;

export interface BaseApi {
  id: string;
  name?: string;
  description?: string;
}

export interface BasePlan {
  id: string;
  name?: string;
  description?: string;
  apiId?: string;
  apiProductId?: string;
  security?: PlanSecurity;
  mode?: PlanMode;
}

export interface BaseApplication {
  id: string;
  name?: string;
  description?: string;
  domain?: string;
  type?: string;
  primaryOwner?: PrimaryOwner;
  apiKeyMode?: "SHARED" | "UNSPECIFIED" | "EXCLUSIVE";
}

/** Full Application resource from APIM management API. */
export interface Application extends BaseApplication {
  applicationType?: string;
  status?: string;
  groups?: string[];
  settings?: ApplicationSettings;
  metadata?: Record<string, ApplicationMetadataValue>;
  disableMembershipNotifications?: boolean;
  originContext?: OriginContext;
  createdAt?: string;
  updatedAt?: string;
}

export interface ApplicationSettings {
  app?: ApplicationSimpleSettings;
  oauth?: ApplicationOAuthSettings;
}

export interface ApplicationSimpleSettings {
  clientId?: string;
  type?: string;
}

export interface ApplicationOAuthSettings {
  clientId?: string;
  clientSecret?: string;
  applicationType?: string;
  grantTypes?: string[];
  redirectUris?: string[];
}

export interface ApplicationMetadataValue {
  name?: string;
  value?: string;
  format?: string;
  defaultValue?: string;
}

/** Paginated list response from APIM. */
export interface PaginatedResult<T> {
  data: T[];
  pagination?: {
    page: number;
    perPage: number;
    pageCount: number;
    pageItemsCount: number;
    totalCount: number;
  };
}

export interface BaseApiProduct {
  id: string;
  name?: string;
  description?: string;
}

export interface BaseUser {
  id: string;
  displayName?: string;
}

export interface SubscriptionConsumerConfiguration {
  entrypointId?: string;
  channel?: string;
  entrypointConfiguration?: Record<string, unknown>;
}

// ── Supporting Types ──────────────────────────────────────────

export interface PrimaryOwner {
  id: string;
  email?: string;
  displayName: string;
  type: "USER" | "GROUP";
}

// ── Flow V4 ───────────────────────────────────────────────────

export interface FlowV4 {
  id?: string;
  name?: string;
  enabled?: boolean;
  selectors?: Selector[];
  request?: StepV4[];
  response?: StepV4[];
  subscribe?: StepV4[];
  publish?: StepV4[];
  entrypointConnect?: StepV4[];
  interact?: StepV4[];
  tags?: string[];
}

export interface StepV4 {
  name?: string;
  description?: string;
  enabled?: boolean;
  policy: string;
  configuration?: Record<string, unknown>;
  condition?: string;
  messageCondition?: string;
}

// ── Selectors (discriminated union) ───────────────────────────

export type Selector = HttpSelector | ChannelSelector | ConditionSelector | McpSelector;

export interface HttpSelector {
  type: "HTTP";
  path?: string;
  pathOperator?: "STARTS_WITH" | "EQUALS";
  methods?: HttpMethod[];
}

export interface ChannelSelector {
  type: "CHANNEL";
  channel?: string;
  channelOperator?: "STARTS_WITH" | "EQUALS";
  operations?: ("SUBSCRIBE" | "PUBLISH")[];
  entrypoints?: string[];
}

export interface ConditionSelector {
  type: "CONDITION";
  condition?: string;
}

export interface McpSelector {
  type: "MCP";
  methods?: string[];
}

// ── Flow V2 ───────────────────────────────────────────────────

export interface FlowV2 {
  id?: string;
  name?: string;
  enabled?: boolean;
  pathOperator?: PathOperator;
  pre?: StepV2[];
  post?: StepV2[];
  methods?: HttpMethod[];
  condition?: string;
  consumers?: Consumer[];
  stage?: FlowStage;
}

export interface PathOperator {
  path?: string;
  operator?: "STARTS_WITH" | "EQUALS";
}

export interface StepV2 {
  name?: string;
  description?: string;
  enabled?: boolean;
  policy?: string;
  configuration?: Record<string, unknown>;
  condition?: string;
}

export interface Consumer {
  consumerId?: string;
  consumerType?: "TAG";
}

export type FlowStage = "PLATFORM" | "API" | "PLAN";

// ── Proxy (V2) ────────────────────────────────────────────────

export interface Proxy {
  virtualHosts?: VirtualHost[];
  groups?: EndpointGroupV2[];
  failover?: Failover;
  cors?: Cors;
  logging?: LoggingV2;
  stripContextPath?: boolean;
  preserveHost?: boolean;
  servers?: string[];
}

export interface VirtualHost {
  host?: string;
  path?: string;
  overrideEntrypoint?: boolean;
}

export interface Failover {
  enabled?: boolean;
  maxAttempts?: number;
  retryTimeout?: number;
}

// ── EndpointGroup / Endpoint (V2) ─────────────────────────────

export interface EndpointGroupV2 {
  name: string;
  endpoints?: EndpointV2[];
  loadBalancer?: LoadBalancer;
  services?: EndpointGroupServicesV2;
  httpProxy?: HttpProxy;
  httpClientOptions?: HttpClientOptions;
  httpClientSslOptions?: HttpClientSslOptions;
  headers?: HttpHeader[];
}

export interface EndpointV2 {
  name: string;
  target: string;
  weight?: number;
  backup?: boolean;
  status?: EndpointStatus;
  tenants?: string[];
  type?: string;
  inherit?: boolean;
  healthCheck?: EndpointHealthCheckService;
  httpProxy?: HttpProxy;
  httpClientOptions?: HttpClientOptions;
  httpClientSslOptions?: HttpClientSslOptions;
  headers?: HttpHeader[];
}

export type EndpointStatus = "UP" | "DOWN" | "TRANSITIONALLY_UP" | "TRANSITIONALLY_DOWN";

export interface EndpointGroupServicesV2 {
  healthCheck?: ServiceV2;
  discovery?: ServiceV2;
}

export interface EndpointHealthCheckService {
  enabled?: boolean;
  schedule?: string;
  configuration?: Record<string, unknown>;
}

export interface HttpProxy {
  enabled?: boolean;
  useSystemProxy?: boolean;
  host?: string;
  port?: number;
  username?: string;
  password?: string;
  type?: "HTTP" | "SOCKS4" | "SOCKS5";
}

export interface HttpClientOptions {
  connectTimeout?: number;
  idleTimeout?: number;
  keepAliveTimeout?: number;
  keepAlive?: boolean;
  readTimeout?: number;
  pipelining?: boolean;
  maxConcurrentConnections?: number;
  useCompression?: boolean;
  followRedirects?: boolean;
  propagateClientAcceptEncoding?: boolean;
  clearTextUpgrade?: boolean;
  version?: "HTTP_1_1" | "HTTP_2";
}

export interface HttpClientSslOptions {
  trustAll?: boolean;
  hostnameVerifier?: boolean;
  trustStore?: TrustStore;
  keyStore?: KeyStore;
}

export interface TrustStore {
  type?: "JKS" | "PKCS12" | "PEM" | "NONE";
  path?: string;
  content?: string;
  password?: string;
}

export interface KeyStore {
  type?: "JKS" | "PKCS12" | "PEM" | "NONE";
  path?: string;
  content?: string;
  password?: string;
  certPath?: string;
  certContent?: string;
  keyPath?: string;
  keyContent?: string;
}

export interface HttpHeader {
  name: string;
  value: string;
}

// ── Logging ───────────────────────────────────────────────────

export interface LoggingV2 {
  mode?: "NONE" | "CLIENT" | "PROXY" | "CLIENT_PROXY";
  condition?: string;
  scope?: "NONE" | "REQUEST" | "RESPONSE" | "REQUEST_RESPONSE";
  content?: "NONE" | "HEADERS" | "PAYLOADS" | "HEADERS_PAYLOADS";
}

export interface LoggingV4 {
  mode?: LoggingMode;
  phase?: LoggingPhase;
  content?: LoggingContent;
  condition?: string;
}

export interface LoggingMode {
  entrypoint?: boolean;
  endpoint?: boolean;
}

export interface LoggingPhase {
  request?: boolean;
  response?: boolean;
}

export interface LoggingContent {
  messagePayload?: boolean;
  messageHeaders?: boolean;
  messageMetadata?: boolean;
  headers?: boolean;
  payload?: boolean;
}

// ── Tracing ───────────────────────────────────────────────────

export interface TracingV4 {
  enabled?: boolean;
  verbose?: boolean;
}

// ── Analytics ─────────────────────────────────────────────────

export interface Analytics {
  enabled?: boolean;
  sampling?: Sampling;
  logging?: LoggingV4;
  tracing?: TracingV4;
}

export interface Sampling {
  type?: SamplingType;
  value?: string;
}

export type SamplingType = "PROBABILITY" | "TEMPORAL" | "COUNT" | "WINDOWED_COUNT";

// ── Notification Types ───────────────────────────────────────

export interface NotificationSetting {
  config_type: "PORTAL" | "GENERIC";
  referenceType: string;
  referenceId: string;
  hooks: string[];
  groups: string[];
  name?: string;
}
