# API Reference

<<<<<<< HEAD
## Versions



<table>
  <thead>
        <tr>
            <th>Version</th>
            <th>Description</th>
        </tr>
  </thead>
  <tbody>
      <tr>
          <td><a href="#graviteeiov1alpha1">gravitee.io/v1alpha1</a></td>
          <td>This version is compatible with gravitee APIM version 3.x and 4.x but v4 API features are not supported.</td>
      </tr>
  </tbody>
</table>

# gravitee.io/v1alpha1

Resources

<table>
  <thead>
  </thead>
    <tbody>
        <tr>
            <td><a href="#managementcontext">ManagementContext</a></td>
            <td></td>
        </tr>
        <tr>
            <td><a href="#apidefinition">ApiDefinition</a></td>
            <td>ApiDefinition is the Schema for the apidefinitions API.</td>
        </tr>
        <tr>
            <td><a href="#apiv4definition">ApiV4Definition</a></td>
            <td>ApiV4Definition is the Schema for the v4 apidefinitions API.</td>
        </tr>
        <tr>
            <td><a href="#apiresource">ApiResource</a></td>
            <td></td>
        </tr>
        <tr>
            <td><a href="#application">Application</a></td>
            <td></td>
        </tr>
        <tr>
            <td><a href="#subscription">Subscription</a></td>
            <td></td>
        </tr>
        <tr>
            <td><a href="#sharedpolicygroup">SharedPolicyGroup</a></td>
            <td>SharedPolicyGroup</td>
        </tr>
        <tr>
            <td><a href="#notification">Notification</a></td>
            <td>Notification defines notification settings in Gravitee</td>
        </tr>
        <tr>
            <td><a href="#group">Group</a></td>
            <td></td>
        </tr></tbody>
</table>



## ManagementContext

[gravitee.io/v1alpha1](#graviteeiov1alpha1)








<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#managementcontextspec">spec</a></b></td>
        <td>object</td>
        <td>
          ManagementContext represents the configuration for a specific environment<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>object</td>
        <td>
          ManagementContextStatus defines the observed state of an API Context.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagementContext.spec
[Go to parent definition](#managementcontext)



ManagementContext represents the configuration for a specific environment

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#managementcontextspecauth">auth</a></b></td>
        <td>object</td>
        <td>
          Auth defines the authentication method used to connect to the API Management.
Can be either basic authentication credentials, a bearer token
or a reference to a kubernetes secret holding one of these two configurations.
This is optional when this context targets Gravitee Cloud.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>baseUrl</b></td>
        <td>string</td>
        <td>
          The URL of a management API instance.
This is optional when this context targets Gravitee Cloud otherwise it is required.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#managementcontextspeccloud">cloud</a></b></td>
        <td>object</td>
        <td>
          Cloud when set (token or secretRef) this context will target Gravitee Cloud.
BaseUrl will be defaulted from token data if not set,
Auth is defaulted to use the token (bearerToken),
OrgID is extracted from the token,
EnvID is defaulted when the token contains exactly one environment.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>environmentId</b></td>
        <td>string</td>
        <td>
          An existing environment id targeted by the context within the organization.
This is optional when this context targets Gravitee Cloud
and your cloud token contains only one environment ID, otherwise it is required.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>organizationId</b></td>
        <td>string</td>
        <td>
          An existing organization id targeted by the context on the management API instance.
This is optional when this context targets Gravitee Cloud otherwise it is required.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Allows to override the context path that will be appended to the baseURL.
This can be used when reverse proxying APIM with URL rewrite<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagementContext.spec.auth
[Go to parent definition](#managementcontextspec)



Auth defines the authentication method used to connect to the API Management.
Can be either basic authentication credentials, a bearer token
or a reference to a kubernetes secret holding one of these two configurations.
This is optional when this context targets Gravitee Cloud.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>bearerToken</b></td>
        <td>string</td>
        <td>
          The bearer token used to authenticate against the API Management instance
(must be generated from an admin account)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#managementcontextspecauthcredentials">credentials</a></b></td>
        <td>object</td>
        <td>
          The Basic credentials used to authenticate against the API Management instance.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#managementcontextspecauthsecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          A secret reference holding either a "bearerToken" key for bearer token authentication
or "username" and "password" keys for basic authentication<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagementContext.spec.auth.credentials
[Go to parent definition](#managementcontextspecauth)



The Basic credentials used to authenticate against the API Management instance.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>password</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>username</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagementContext.spec.auth.secretRef
[Go to parent definition](#managementcontextspecauth)



A secret reference holding either a "bearerToken" key for bearer token authentication
or "username" and "password" keys for basic authentication

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagementContext.spec.cloud
[Go to parent definition](#managementcontextspec)



Cloud when set (token or secretRef) this context will target Gravitee Cloud.
BaseUrl will be defaulted from token data if not set,
Auth is defaulted to use the token (bearerToken),
OrgID is extracted from the token,
EnvID is defaulted when the token contains exactly one environment.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#managementcontextspeccloudsecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          SecretRef secret reference holding the Gravitee cloud token in the "cloudToken" key<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>token</b></td>
        <td>string</td>
        <td>
          Token plain text Gravitee cloud token (JWT)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagementContext.spec.cloud.secretRef
[Go to parent definition](#managementcontextspeccloud)



SecretRef secret reference holding the Gravitee cloud token in the "cloudToken" key

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## ApiDefinition

[gravitee.io/v1alpha1](#graviteeiov1alpha1)
=======
## Packages
- [gravitee.io/v1alpha1](#graviteeiov1alpha1)
- [gravitee.io/v1alpha1/application](#graviteeiov1alpha1application)
- [gravitee.io/v1alpha1/base](#graviteeiov1alpha1base)
- [gravitee.io/v1alpha1/gateway](#graviteeiov1alpha1gateway)
- [gravitee.io/v1alpha1/group](#graviteeiov1alpha1group)
- [gravitee.io/v1alpha1/kafka](#graviteeiov1alpha1kafka)
- [gravitee.io/v1alpha1/management](#graviteeiov1alpha1management)
- [gravitee.io/v1alpha1/notification](#graviteeiov1alpha1notification)
- [gravitee.io/v1alpha1/policygroups](#graviteeiov1alpha1policygroups)
- [gravitee.io/v1alpha1/refs](#graviteeiov1alpha1refs)
- [gravitee.io/v1alpha1/status](#graviteeiov1alpha1status)
- [gravitee.io/v1alpha1/subscription](#graviteeiov1alpha1subscription)
- [gravitee.io/v1alpha1/utils](#graviteeiov1alpha1utils)
- [gravitee.io/v1alpha1/v2](#graviteeiov1alpha1v2)
- [gravitee.io/v1alpha1/v4](#graviteeiov1alpha1v4)


## gravitee.io/v1alpha1

Package v1alpha1 contains API Schema definitions for the  v1alpha1 API group

### Resource Types
- [ApiDefinition](#apidefinition)
- [ApiResource](#apiresource)
- [ApiV4Definition](#apiv4definition)
- [Application](#application)
- [GatewayClassParameters](#gatewayclassparameters)
- [Group](#group)
- [KafkaRoute](#kafkaroute)
- [ManagementContext](#managementcontext)
- [Notification](#notification)
- [SharedPolicyGroup](#sharedpolicygroup)
- [Subscription](#subscription)
>>>>>>> 56d297b (docs: generate docs using crd-ref-docs)



#### ApiDefinition



ApiDefinition is the Schema for the apidefinitions API.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `gravitee.io/v1alpha1` | | |
| `kind` _string_ | `ApiDefinition` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ApiDefinitionV2Spec](#apidefinitionv2spec)_ |  |  |  |
| `status` _[ApiDefinitionStatus](#apidefinitionstatus)_ |  |  |  |


#### ApiDefinitionStatus



ApiDefinitionStatus defines the observed state of API Definition.



_Appears in:_
- [ApiDefinition](#apidefinition)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `organizationId` _string_ | The organization ID, if a management context has been defined to sync with an APIM instance |  | Optional: \{\} <br /> |
| `environmentId` _string_ | The environment ID, if a management context has been defined to sync with an APIM instance |  | Optional: \{\} <br /> |
| `id` _string_ | The ID of the API definition in the Gravitee API Management instance (if an API context has been configured). |  | Optional: \{\} <br /> |
| `crossId` _string_ | The Cross ID is used to identify an API that has been promoted from one environment to another. |  |  |
| `processingStatus` _[ProcessingStatus](#processingstatus)_ | The processing status of the API definition. |  |  |
| `state` _[ApiState](#apistate)_ | The state of the API. Can be either STARTED or STOPPED. |  | Enum: [STARTED STOPPED] <br /> |
| `plans` _object (keys:string, values:string)_ | This field is used to store the list of plans that have been created<br />for the API definition if a management context has been defined<br />to sync the API with an APIM instance |  | Optional: \{\} <br /> |
| `subscriptions` _integer_ | The number of subscriptions that reference the API |  |  |
| `errors` _[Errors](#errors)_ | When API has been created regardless of errors, this field is<br />used to persist the error message encountered during admission |  |  |


#### ApiDefinitionV2Spec



The API definition is the main resource handled by the Kubernetes Operator
Most of the configuration properties defined here are already documented
in the APIM Console API Reference.
See https://docs.gravitee.io/apim/3.x/apim_installguide_rest_apis_documentation.html


<<<<<<< HEAD

If true, the Operator will create the ConfigMaps for the Gateway and pushes the API to the Management API
but without setting the update flag in the datastore.


If false, the Operator will not create the ConfigMaps for the Gateway.
Instead, it pushes the API to the Management API and forces it to update the event in the datastore.
This will cause Gateways to fetch the APIs from the datastore<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecmembersindex">members</a></b></td>
        <td>[]object</td>
        <td>
          List of members associated with the API<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecmetadataindex">metadata</a></b></td>
        <td>[]object</td>
        <td>
          List of API metadata entries<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecnotificationsrefsindex">notificationsRefs</a></b></td>
        <td>[]object</td>
        <td>
          References to Notification custom resources to setup notifications.
For an API Notification CRD `eventType` field must be set to `api`
and only events set via `apiEvents` attributes are used.
Only one notification with `target` equals to `console` is admitted.<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>notifyMembers</b></td>
        <td>boolean</td>
        <td>
          If true, new members added to the API spec will
be notified when the API is synced with APIM.<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecpageskey">pages</a></b></td>
        <td>map[string]object</td>
        <td>
          A map of pages objects.


Keys uniquely identify pages and are used to keep them in sync
with APIM when using a management context.


Renaming a key is the equivalent of deleting the page and recreating
it holding a new ID in APIM.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path_mappings</b></td>
        <td>[]string</td>
        <td>
          API Path mapping<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecplansindex">plans</a></b></td>
        <td>[]object</td>
        <td>
          API plans<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecpropertiesindex">properties</a></b></td>
        <td>[]object</td>
        <td>
          List of Properties for the API<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxy">proxy</a></b></td>
        <td>object</td>
        <td>
          The proxy of the API that specifies its VirtualHosts and Groups.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecresourcesindex">resources</a></b></td>
        <td>[]object</td>
        <td>
          Resources can be either inlined or reference the namespace and name
of an <a href="#apiresource">existing API resource definition</a>.<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecresponse_templateskeykey">response_templates</a></b></td>
        <td>map[string]map[string]object</td>
        <td>
          A list of Response Templates for the API<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecservices">services</a></b></td>
        <td>object</td>
        <td>
          Contains different services for the API (EndpointDiscovery, HealthCheck ...)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>state</b></td>
        <td>enum</td>
        <td>
          The state of API (setting the value to `STOPPED` will make the API un-reachable from the gateway)<br/>
          <br/>
            <i>Enum</i>: STARTED, STOPPED<br/>
            <i>Default</i>: STARTED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tags</b></td>
        <td>[]string</td>
        <td>
          List of Tags of the API<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>visibility</b></td>
        <td>enum</td>
        <td>
          Should the API be publicly available from the portal or not ?<br/>
          <br/>
            <i>Enum</i>: PUBLIC, PRIVATE<br/>
            <i>Default</i>: PRIVATE<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
=======

_Appears in:_
- [ApiDefinition](#apidefinition)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `description` _string_ | API description |  |  |
| `definition_context` _[DefinitionContext](#definitioncontext)_ | The definition context is used to inform a management API instance that this API definition<br />is managed using a kubernetes operator |  | Optional: \{\} <br /> |
| `lifecycle_state` _[LifecycleState](#lifecyclestate)_ | API life cycle state can be one of the values CREATED, PUBLISHED, UNPUBLISHED, DEPRECATED, ARCHIVED | CREATED | Enum: [CREATED PUBLISHED UNPUBLISHED DEPRECATED ARCHIVED] <br /> |
| `deployedAt` _integer_ | Shows the time that the API is deployed |  | Optional: \{\} <br /> |
| `gravitee` _[DefinitionVersion](#definitionversion)_ | The definition version of the API. For v1alpha1 resources, this field should always set to `2.0.0`. | 2.0.0 |  |
| `flow_mode` _[FlowMode](#flowmode)_ | The flow mode of the API. The value is either `DEFAULT` or `BEST_MATCH`. | DEFAULT | Enum: [DEFAULT BEST_MATCH] <br /> |
| `proxy` _[Proxy](#proxy)_ | The proxy of the API that specifies its VirtualHosts and Groups. |  |  |
| `services` _[Services](#services)_ | Contains different services for the API (EndpointDiscovery, HealthCheck ...) |  |  |
| `flows` _[Flow](#flow) array_ | The flow of the API | \{  \} | Optional: \{\} <br /> |
| `path_mappings` _string array_ | API Path mapping |  | Optional: \{\} <br /> |
| `plans` _[Plan](#plan) array_ | API plans | \{  \} | Optional: \{\} <br /> |
| `response_templates` _[ResponseTemplate](#responsetemplate)_ | A list of Response Templates for the API |  | Optional: \{\} <br /> |
| `members` _Member array_ | List of members associated with the API |  | Optional: \{\} <br /> |
| `pages` _[map[string]*Page](#map[string]*page)_ | A map of pages objects.<br />Keys uniquely identify pages and are used to keep them in sync<br />with APIM when using a management context.<br />Renaming a key is the equivalent of deleting the page and recreating<br />it holding a new ID in APIM. |  | Optional: \{\} <br /> |
| `execution_mode` _string_ | Execution mode that eventually runs the API in the gateway | v4-emulation-engine | Enum: [v3 v4-emulation-engine] <br /> |
| `contextRef` _[NamespacedName](#namespacedname)_ |  |  |  |
| `local` _boolean_ | local defines if the api is local or not.<br />If true, the Operator will create the ConfigMaps for the Gateway and pushes the API to the Management API<br />but without setting the update flag in the datastore.<br />If false, the Operator will not create the ConfigMaps for the Gateway.<br />Instead, it pushes the API to the Management API and forces it to update the event in the datastore.<br />This will cause Gateways to fetch the APIs from the datastore | false | Optional: \{\} <br /> |

>>>>>>> 56d297b (docs: generate docs using crd-ref-docs)

#### ApiResource









| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `gravitee.io/v1alpha1` | | |
| `kind` _string_ | `ApiResource` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ApiResourceSpec](#apiresourcespec)_ |  |  |  |
| `status` _[ApiResourceStatus](#apiresourcestatus)_ |  |  |  |


<<<<<<< HEAD

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.definition_context
[Go to parent definition](#apidefinitionspec)



The definition context is used to inform a management API instance that this API definition
is managed using a kubernetes operator

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>mode</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: fully_managed<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>origin</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: kubernetes<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>syncFrom</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: kubernetes<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.flows[index]
[Go to parent definition](#apidefinitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this flow is enabled or disabled<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          Flow condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecflowsindexconsumersindex">consumers</a></b></td>
        <td>[]object</td>
        <td>
          List of the consumers of this Flow<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Flow ID<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>methods</b></td>
        <td>[]enum</td>
        <td>
          A list of methods  for this flow (GET;POST;PUT;PATCH;DELETE;OPTIONS;HEAD;CONNECT;TRACE;OTHER)<br/>
          <br/>
            <i>Enum</i>: GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD, CONNECT, TRACE, OTHER<br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Flow name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecflowsindexpath-operator">path-operator</a></b></td>
        <td>object</td>
        <td>
          List of path operators<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecflowsindexpostindex">post</a></b></td>
        <td>[]object</td>
        <td>
          Flow post step<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecflowsindexpreindex">pre</a></b></td>
        <td>[]object</td>
        <td>
          Flow pre step<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.flows[index].consumers[index]
[Go to parent definition](#apidefinitionspecflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>consumerId</b></td>
        <td>string</td>
        <td>
          Consumer ID<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>consumerType</b></td>
        <td>integer</td>
        <td>
          Consumer type (possible values TAG)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.flows[index].path-operator
[Go to parent definition](#apidefinitionspecflowsindex)



List of path operators

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>operator</b></td>
        <td>enum</td>
        <td>
          Operator (possible values STARTS_WITH or EQUALS)<br/>
          <br/>
            <i>Enum</i>: STARTS_WITH, EQUALS<br/>
            <i>Default</i>: STARTS_WITH<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Operator path<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.flows[index].post[index]
[Go to parent definition](#apidefinitionspecflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.flows[index].pre[index]
[Go to parent definition](#apidefinitionspecflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.groupRefs[index]
[Go to parent definition](#apidefinitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.members[index]
[Go to parent definition](#apidefinitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>source</b></td>
        <td>string</td>
        <td>
          Member source<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>sourceId</b></td>
        <td>string</td>
        <td>
          Member source ID<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>role</b></td>
        <td>string</td>
        <td>
          The API role associated with this Member<br/>
          <br/>
            <i>Default</i>: USER<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.metadata[index]
[Go to parent definition](#apidefinitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>format</b></td>
        <td>enum</td>
        <td>
          Metadata Format<br/>
          <br/>
            <i>Enum</i>: STRING, NUMERIC, BOOLEAN, DATE, MAIL, URL<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Metadata Key<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Metadata Name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>defaultValue</b></td>
        <td>string</td>
        <td>
          Metadata Default value<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Metadata Value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.notificationsRefs[index]
[Go to parent definition](#apidefinitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.pages[key]
[Go to parent definition](#apidefinitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          The type of the documentation page or folder.<br/>
          <br/>
            <i>Enum</i>: MARKDOWN, SWAGGER, ASYNCAPI, ASCIIDOC, FOLDER, SYSTEM_FOLDER, ROOT<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecpageskeyaccesscontrolsindex">accessControls</a></b></td>
        <td>[]object</td>
        <td>
          If the page is private, defines a set of user groups with access<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>api</b></td>
        <td>string</td>
        <td>
          The API of the page. If empty, will be set automatically to the generated ID of the API.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>map[string]string</td>
        <td>
          Custom page configuration (e.g. page rendering can be changed to use Redoc instead of Swagger ui)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>content</b></td>
        <td>string</td>
        <td>
          The content of the page, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>crossId</b></td>
        <td>string</td>
        <td>
          CrossID is designed to identified a page across environments.
If not set, this ID will be generated in a predictable manner based on
the map key associated to this entry in the API.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>excludedAccessControls</b></td>
        <td>boolean</td>
        <td>
          if true, the references defined in the accessControls list will be
denied access instead of being granted<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>homepage</b></td>
        <td>boolean</td>
        <td>
          If true, this page will be displayed as the homepage of your API documentation.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the page. This field is mostly required when you are applying
an API exported from APIM to make the operator take control over it.
If not set, this ID will be generated in a predictable manner based on
the map key associated to this entry in the API.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          This is the display name of the page in APIM and on the portal.
This field can be edited safely if you want to rename a page.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>order</b></td>
        <td>integer</td>
        <td>
          The order used to display the page in APIM and on the portal.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>parent</b></td>
        <td>string</td>
        <td>
          If your page contains a folder, setting this field to the map key associated to the
folder entry will be reflected into APIM by making the page a child of this folder.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>parentId</b></td>
        <td>string</td>
        <td>
          The parent ID of the page. This field is mostly required when you are applying
an API exported from APIM to make the operator take control over it. Use `Parent`
in any other case.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>published</b></td>
        <td>boolean</td>
        <td>
          If true, the page will be accessible from the portal (default is false)<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecpageskeysource">source</a></b></td>
        <td>object</td>
        <td>
          Source allow you to fetch pages from various external sources, overriding page content
each time the source is fetched.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>visibility</b></td>
        <td>enum</td>
        <td>
          The visibility of the page.<br/>
          <br/>
            <i>Enum</i>: PUBLIC, PRIVATE<br/>
            <i>Default</i>: PUBLIC<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.pages[key].accessControls[index]
[Go to parent definition](#apidefinitionspecpageskey)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>referenceId</b></td>
        <td>string</td>
        <td>
          The ID denied or granted by the access control (currently only group names are supported)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>referenceType</b></td>
        <td>enum</td>
        <td>
          The type of reference denied or granted by the access control
Currently only GROUP is supported<br/>
          <br/>
            <i>Enum</i>: GROUP<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.pages[key].source
[Go to parent definition](#apidefinitionspecpageskey)



Source allow you to fetch pages from various external sources, overriding page content
each time the source is fetched.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.plans[index]
[Go to parent definition](#apidefinitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          Plan Description<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Plan name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>security</b></td>
        <td>string</td>
        <td>
          Plan Security<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>api</b></td>
        <td>string</td>
        <td>
          Specify the API associated with this plan<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>characteristics</b></td>
        <td>[]string</td>
        <td>
          List of plan characteristics<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>comment_required</b></td>
        <td>boolean</td>
        <td>
          Indicate of comment is required for this plan or not<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>crossId</b></td>
        <td>string</td>
        <td>
          The plan Cross ID.
This field is used to identify plans defined for an API
that has been promoted between different environments.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>excluded_groups</b></td>
        <td>[]string</td>
        <td>
          List of excluded groups for this plan<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecplansindexflowsindex">flows</a></b></td>
        <td>[]object</td>
        <td>
          List of different flows for this Plan<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Plan ID<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>order</b></td>
        <td>integer</td>
        <td>
          Plan order<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecplansindexpathskeyindex">paths</a></b></td>
        <td>map[string][]object</td>
        <td>
          A map of different paths (alongside their Rules) for this Plan<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>securityDefinition</b></td>
        <td>string</td>
        <td>
          Plan Security definition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>selection_rule</b></td>
        <td>string</td>
        <td>
          Plan selection rule<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>enum</td>
        <td>
          The plan status<br/>
          <br/>
            <i>Enum</i>: PUBLISHED, DEPRECATED, STAGING<br/>
            <i>Default</i>: PUBLISHED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tags</b></td>
        <td>[]string</td>
        <td>
          List of plan tags<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Plan type<br/>
          <br/>
            <i>Enum</i>: API, CATALOG<br/>
            <i>Default</i>: API<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>validation</b></td>
        <td>enum</td>
        <td>
          Plan validation strategy<br/>
          <br/>
            <i>Enum</i>: AUTO, MANUAL<br/>
            <i>Default</i>: AUTO<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.plans[index].flows[index]
[Go to parent definition](#apidefinitionspecplansindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this flow is enabled or disabled<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          Flow condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecplansindexflowsindexconsumersindex">consumers</a></b></td>
        <td>[]object</td>
        <td>
          List of the consumers of this Flow<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Flow ID<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>methods</b></td>
        <td>[]enum</td>
        <td>
          A list of methods  for this flow (GET;POST;PUT;PATCH;DELETE;OPTIONS;HEAD;CONNECT;TRACE;OTHER)<br/>
          <br/>
            <i>Enum</i>: GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD, CONNECT, TRACE, OTHER<br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Flow name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecplansindexflowsindexpath-operator">path-operator</a></b></td>
        <td>object</td>
        <td>
          List of path operators<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecplansindexflowsindexpostindex">post</a></b></td>
        <td>[]object</td>
        <td>
          Flow post step<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecplansindexflowsindexpreindex">pre</a></b></td>
        <td>[]object</td>
        <td>
          Flow pre step<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.plans[index].flows[index].consumers[index]
[Go to parent definition](#apidefinitionspecplansindexflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>consumerId</b></td>
        <td>string</td>
        <td>
          Consumer ID<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>consumerType</b></td>
        <td>integer</td>
        <td>
          Consumer type (possible values TAG)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.plans[index].flows[index].path-operator
[Go to parent definition](#apidefinitionspecplansindexflowsindex)



List of path operators

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>operator</b></td>
        <td>enum</td>
        <td>
          Operator (possible values STARTS_WITH or EQUALS)<br/>
          <br/>
            <i>Enum</i>: STARTS_WITH, EQUALS<br/>
            <i>Default</i>: STARTS_WITH<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Operator path<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.plans[index].flows[index].post[index]
[Go to parent definition](#apidefinitionspecplansindexflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.plans[index].flows[index].pre[index]
[Go to parent definition](#apidefinitionspecplansindexflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.plans[index].paths[key][index]
[Go to parent definition](#apidefinitionspecplansindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          Rule description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if the Rule is enabled or not<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>methods</b></td>
        <td>[]enum</td>
        <td>
          List of http methods for this Rule (GET;POST;PUT;PATCH;DELETE;OPTIONS;HEAD;CONNECT;TRACE;OTHER)<br/>
          <br/>
            <i>Enum</i>: GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD, CONNECT, TRACE, OTHER<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecplansindexpathskeyindexpolicy">policy</a></b></td>
        <td>object</td>
        <td>
          Rule policy<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.plans[index].paths[key][index].policy
[Go to parent definition](#apidefinitionspecplansindexpathskeyindex)



Rule policy

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Policy configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Policy name<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.properties[index]
[Go to parent definition](#apidefinitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>dynamic</b></td>
        <td>boolean</td>
        <td>
          Property is dynamic or not?<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>encryptable</b></td>
        <td>boolean</td>
        <td>
          Property is encryptable or not?<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>encrypted</b></td>
        <td>boolean</td>
        <td>
          Property Encrypted or not?<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Property Key<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Property Value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy
[Go to parent definition](#apidefinitionspec)



The proxy of the API that specifies its VirtualHosts and Groups.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#apidefinitionspecproxycors">cors</a></b></td>
        <td>object</td>
        <td>
          Proxy Cors<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxyfailover">failover</a></b></td>
        <td>object</td>
        <td>
          Proxy Failover<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindex">groups</a></b></td>
        <td>[]object</td>
        <td>
          List of endpoint groups of the proxy<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxylogging">logging</a></b></td>
        <td>object</td>
        <td>
          Logging<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>preserve_host</b></td>
        <td>boolean</td>
        <td>
          Preserve Host<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>strip_context_path</b></td>
        <td>boolean</td>
        <td>
          Strip Context Path<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxyvirtual_hostsindex">virtual_hosts</a></b></td>
        <td>[]object</td>
        <td>
          list of Virtual hosts fot the proxy<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.cors
[Go to parent definition](#apidefinitionspecproxy)



Proxy Cors

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>allowCredentials</b></td>
        <td>boolean</td>
        <td>
          Access Control - Allow credentials or not<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if the cors enabled or not<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>maxAge</b></td>
        <td>integer</td>
        <td>
          Access Control -  Max age<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>allowHeaders</b></td>
        <td>[]string</td>
        <td>
          Access Control - List of allowed headers<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>allowMethods</b></td>
        <td>[]string</td>
        <td>
          Access Control - List of allowed methods<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>allowOrigin</b></td>
        <td>[]string</td>
        <td>
          Access Control -  List of Allowed origins<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>exposeHeaders</b></td>
        <td>[]string</td>
        <td>
          Access Control - List of Exposed Headers<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runPolicies</b></td>
        <td>boolean</td>
        <td>
          Run policies or not<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.failover
[Go to parent definition](#apidefinitionspecproxy)



Proxy Failover

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>cases</b></td>
        <td>[]string</td>
        <td>
          List of Failover cases<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxAttempts</b></td>
        <td>integer</td>
        <td>
          Maximum number of attempts<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>retryTimeout</b></td>
        <td>integer</td>
        <td>
          Retry timeout<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index]
[Go to parent definition](#apidefinitionspecproxy)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindex">endpoints</a></b></td>
        <td>[]object</td>
        <td>
          List of Endpoints belonging to this group<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>headers</b></td>
        <td>map[string]string</td>
        <td>
          List of headers needed for this EndpointGroup<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexhttp">http</a></b></td>
        <td>object</td>
        <td>
          Custom HTTP SSL client options used for this EndpointGroup<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexload_balancing">load_balancing</a></b></td>
        <td>object</td>
        <td>
          The LoadBalancer Type<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          EndpointGroup name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexproxy">proxy</a></b></td>
        <td>object</td>
        <td>
          Configure the HTTP Proxy settings for this EndpointGroup if needed<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexservices">services</a></b></td>
        <td>object</td>
        <td>
          Specify different Endpoint Services<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexssl">ssl</a></b></td>
        <td>object</td>
        <td>
          Custom HTTP SSL client options used for this EndpointGroup<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index]
[Go to parent definition](#apidefinitionspecproxygroupsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>backup</b></td>
        <td>boolean</td>
        <td>
          Indicate that this ia a back-end endpoint<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexheadersindex">headers</a></b></td>
        <td>[]object</td>
        <td>
          List of headers for this endpoint<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexhealthcheck">healthcheck</a></b></td>
        <td>object</td>
        <td>
          Specify EndpointHealthCheck service settings<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexhttp">http</a></b></td>
        <td>object</td>
        <td>
          Custom HTTP client options used for this endpoint<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>inherit</b></td>
        <td>boolean</td>
        <td>
          Is endpoint inherited or not<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the endpoint<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexproxy">proxy</a></b></td>
        <td>object</td>
        <td>
          Configure the HTTP Proxy settings to reach target if needed<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexssl">ssl</a></b></td>
        <td>object</td>
        <td>
          Custom HTTP SSL client options used for this endpoint<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>target</b></td>
        <td>string</td>
        <td>
          The end target of this endpoint (backend)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tenants</b></td>
        <td>[]string</td>
        <td>
          The endpoint tenants<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          The type of endpoint (HttpEndpointType or GrpcEndpointType)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>weight</b></td>
        <td>integer</td>
        <td>
          Endpoint weight used for load-balancing<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].headers[index]
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The HTTP header name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The HTTP header value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].healthcheck
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindex)



Specify EndpointHealthCheck service settings

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is service enabled or not?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>inherit</b></td>
        <td>boolean</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Service name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>schedule</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexhealthcheckstepsindex">steps</a></b></td>
        <td>[]object</td>
        <td>
          List of health check steps<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].healthcheck.steps[index]
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindexhealthcheck)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Health Check Step Name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexhealthcheckstepsindexrequest">request</a></b></td>
        <td>object</td>
        <td>
          Health Check Step Request<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexhealthcheckstepsindexresponse">response</a></b></td>
        <td>object</td>
        <td>
          Health Check Step Response<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].healthcheck.steps[index].request
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindexhealthcheckstepsindex)



Health Check Step Request

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fromRoot</b></td>
        <td>boolean</td>
        <td>
          If true, the health check request will be issued without prepending the context path of the API.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>body</b></td>
        <td>string</td>
        <td>
          Health Check Request Body<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexhealthcheckstepsindexrequestheadersindex">headers</a></b></td>
        <td>[]object</td>
        <td>
          List of HTTP headers to include in the health check request<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>method</b></td>
        <td>enum</td>
        <td>
          The HTTP method to use when issuing the health check request<br/>
          <br/>
            <i>Enum</i>: GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD, CONNECT, TRACE, OTHER<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          The path of the endpoint handling the health check request<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].healthcheck.steps[index].request.headers[index]
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindexhealthcheckstepsindexrequest)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The HTTP header name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The HTTP header value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].healthcheck.steps[index].response
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindexhealthcheckstepsindex)



Health Check Step Response

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>assertions</b></td>
        <td>[]string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].http
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindex)



Custom HTTP client options used for this endpoint

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>followRedirects</b></td>
        <td>boolean</td>
        <td>
          Should HTTP redirects be followed or not ?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>pipelining</b></td>
        <td>boolean</td>
        <td>
          Should HTTP/1.1 pipelining be used for the connection or not ?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>useCompression</b></td>
        <td>boolean</td>
        <td>
          Should compression be used or not ?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>clearTextUpgrade</b></td>
        <td>boolean</td>
        <td>
          Should HTTP/2 clear text upgrade be used or not ?<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>connectTimeout</b></td>
        <td>integer</td>
        <td>
          Connection timeout of the http connection<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>idleTimeout</b></td>
        <td>integer</td>
        <td>
           Idle Timeout for the http connection<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keepAlive</b></td>
        <td>boolean</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keepAliveTimeout</b></td>
        <td>integer</td>
        <td>
          Should keep alive be used for the HTTP connection ?<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Default</i>: 30000<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxConcurrentConnections</b></td>
        <td>integer</td>
        <td>
          HTTP max concurrent connections<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>propagateClientAcceptEncoding</b></td>
        <td>boolean</td>
        <td>
          Propagate Client Accept-Encoding header<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readTimeout</b></td>
        <td>integer</td>
        <td>
          Read timeout<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>version</b></td>
        <td>enum</td>
        <td>
          HTTP Protocol Version (Possible values Http1 or Http2)<br/>
          <br/>
            <i>Enum</i>: HTTP_1_1, HTTP_2<br/>
            <i>Default</i>: HTTP_1_1<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].proxy
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindex)



Configure the HTTP Proxy settings to reach target if needed

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Specifies that the HTTP connection will be established through a proxy<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Proxy host name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>password</b></td>
        <td>string</td>
        <td>
          The HTTP proxy password (if the proxy requires authentication)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          The HTTP proxy port<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          The HTTP proxy type (possible values Http, Socks4, Socks5)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>useSystemProxy</b></td>
        <td>boolean</td>
        <td>
          If true, the proxy defined at the system level will be used<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>username</b></td>
        <td>string</td>
        <td>
          The HTTP proxy username (if the proxy requires authentication)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].ssl
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindex)



Custom HTTP SSL client options used for this endpoint

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>hostnameVerifier</b></td>
        <td>boolean</td>
        <td>
          Verify Hostname when establishing connection<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>trustAll</b></td>
        <td>boolean</td>
        <td>
          Whether to trust all issuers or not<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexsslheadersindex">headers</a></b></td>
        <td>[]object</td>
        <td>
          Http headers<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexsslkeystore">keyStore</a></b></td>
        <td>object</td>
        <td>
          KeyStore type (possible values PEM, PKCS12, JKS)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexendpointsindexssltruststore">trustStore</a></b></td>
        <td>object</td>
        <td>
          TrustStore type (possible values PEM, PKCS12, JKS)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].ssl.headers[index]
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindexssl)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The HTTP header name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The HTTP header value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].ssl.keyStore
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindexssl)



KeyStore type (possible values PEM, PKCS12, JKS)

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>certContent</b></td>
        <td>string</td>
        <td>
          KeyStore cert content (Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>certPath</b></td>
        <td>string</td>
        <td>
          KeyStore cert path (Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>content</b></td>
        <td>string</td>
        <td>
          The base64 encoded trustStore content, if not relying on a path to a file<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keyContent</b></td>
        <td>string</td>
        <td>
          The base64 encoded trustStore content, if not relying on a path to a file
(Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keyPath</b></td>
        <td>string</td>
        <td>
          KeyStore key path (Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>password</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          KeyStore path<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          The KeyStore type to use (possible values are PEM, PKCS12, JKS)<br/>
          <br/>
            <i>Enum</i>: PEM, PKCS12, JKS<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].endpoints[index].ssl.trustStore
[Go to parent definition](#apidefinitionspecproxygroupsindexendpointsindexssl)



TrustStore type (possible values PEM, PKCS12, JKS)

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>content</b></td>
        <td>string</td>
        <td>
          The base64 encoded trustStore content, if not relying on a path to a file<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>password</b></td>
        <td>string</td>
        <td>
          TrustStore password (Not applicable for PEM TrustStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          The TrustStore type to use (possible values are PEM, PKCS12, JKS)<br/>
          <br/>
            <i>Enum</i>: PEM, PKCS12, JKS<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].http
[Go to parent definition](#apidefinitionspecproxygroupsindex)



Custom HTTP SSL client options used for this EndpointGroup

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>followRedirects</b></td>
        <td>boolean</td>
        <td>
          Should HTTP redirects be followed or not ?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>pipelining</b></td>
        <td>boolean</td>
        <td>
          Should HTTP/1.1 pipelining be used for the connection or not ?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>useCompression</b></td>
        <td>boolean</td>
        <td>
          Should compression be used or not ?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>clearTextUpgrade</b></td>
        <td>boolean</td>
        <td>
          Should HTTP/2 clear text upgrade be used or not ?<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>connectTimeout</b></td>
        <td>integer</td>
        <td>
          Connection timeout of the http connection<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>idleTimeout</b></td>
        <td>integer</td>
        <td>
           Idle Timeout for the http connection<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keepAlive</b></td>
        <td>boolean</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keepAliveTimeout</b></td>
        <td>integer</td>
        <td>
          Should keep alive be used for the HTTP connection ?<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Default</i>: 30000<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxConcurrentConnections</b></td>
        <td>integer</td>
        <td>
          HTTP max concurrent connections<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>propagateClientAcceptEncoding</b></td>
        <td>boolean</td>
        <td>
          Propagate Client Accept-Encoding header<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readTimeout</b></td>
        <td>integer</td>
        <td>
          Read timeout<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>version</b></td>
        <td>enum</td>
        <td>
          HTTP Protocol Version (Possible values Http1 or Http2)<br/>
          <br/>
            <i>Enum</i>: HTTP_1_1, HTTP_2<br/>
            <i>Default</i>: HTTP_1_1<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].load_balancing
[Go to parent definition](#apidefinitionspecproxygroupsindex)



The LoadBalancer Type

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of the LoadBalancer (RoundRobin, Random, WeightedRoundRobin, WeightedRandom)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].proxy
[Go to parent definition](#apidefinitionspecproxygroupsindex)



Configure the HTTP Proxy settings for this EndpointGroup if needed

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Specifies that the HTTP connection will be established through a proxy<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Proxy host name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>password</b></td>
        <td>string</td>
        <td>
          The HTTP proxy password (if the proxy requires authentication)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          The HTTP proxy port<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          The HTTP proxy type (possible values Http, Socks4, Socks5)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>useSystemProxy</b></td>
        <td>boolean</td>
        <td>
          If true, the proxy defined at the system level will be used<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>username</b></td>
        <td>string</td>
        <td>
          The HTTP proxy username (if the proxy requires authentication)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].services
[Go to parent definition](#apidefinitionspecproxygroupsindex)



Specify different Endpoint Services

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexservicesdiscovery">discovery</a></b></td>
        <td>object</td>
        <td>
          Endpoint Discovery Service<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexservicesdynamic-property">dynamic-property</a></b></td>
        <td>object</td>
        <td>
          Dynamic Property Service<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexserviceshealth-check">health-check</a></b></td>
        <td>object</td>
        <td>
          Health Check Service<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].services.discovery
[Go to parent definition](#apidefinitionspecproxygroupsindexservices)



Endpoint Discovery Service

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Configuration, arbitrary map of key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is service enabled or not?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Service name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>provider</b></td>
        <td>string</td>
        <td>
          Provider name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>secondary</b></td>
        <td>boolean</td>
        <td>
          Is it secondary or not?<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tenants</b></td>
        <td>[]string</td>
        <td>
          List of tenants<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].services.dynamic-property
[Go to parent definition](#apidefinitionspecproxygroupsindexservices)



Dynamic Property Service

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Configuration, arbitrary map of key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is service enabled or not?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Service name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>provider</b></td>
        <td>enum</td>
        <td>
          <br/>
          <br/>
            <i>Enum</i>: HTTP<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>schedule</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].services.health-check
[Go to parent definition](#apidefinitionspecproxygroupsindexservices)



Health Check Service

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is service enabled or not?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Service name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>schedule</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexserviceshealth-checkstepsindex">steps</a></b></td>
        <td>[]object</td>
        <td>
          List of health check steps<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].services.health-check.steps[index]
[Go to parent definition](#apidefinitionspecproxygroupsindexserviceshealth-check)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Health Check Step Name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexserviceshealth-checkstepsindexrequest">request</a></b></td>
        <td>object</td>
        <td>
          Health Check Step Request<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexserviceshealth-checkstepsindexresponse">response</a></b></td>
        <td>object</td>
        <td>
          Health Check Step Response<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].services.health-check.steps[index].request
[Go to parent definition](#apidefinitionspecproxygroupsindexserviceshealth-checkstepsindex)



Health Check Step Request

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fromRoot</b></td>
        <td>boolean</td>
        <td>
          If true, the health check request will be issued without prepending the context path of the API.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>body</b></td>
        <td>string</td>
        <td>
          Health Check Request Body<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexserviceshealth-checkstepsindexrequestheadersindex">headers</a></b></td>
        <td>[]object</td>
        <td>
          List of HTTP headers to include in the health check request<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>method</b></td>
        <td>enum</td>
        <td>
          The HTTP method to use when issuing the health check request<br/>
          <br/>
            <i>Enum</i>: GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD, CONNECT, TRACE, OTHER<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          The path of the endpoint handling the health check request<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].services.health-check.steps[index].request.headers[index]
[Go to parent definition](#apidefinitionspecproxygroupsindexserviceshealth-checkstepsindexrequest)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The HTTP header name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The HTTP header value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].services.health-check.steps[index].response
[Go to parent definition](#apidefinitionspecproxygroupsindexserviceshealth-checkstepsindex)



Health Check Step Response

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>assertions</b></td>
        <td>[]string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].ssl
[Go to parent definition](#apidefinitionspecproxygroupsindex)



Custom HTTP SSL client options used for this EndpointGroup

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>hostnameVerifier</b></td>
        <td>boolean</td>
        <td>
          Verify Hostname when establishing connection<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>trustAll</b></td>
        <td>boolean</td>
        <td>
          Whether to trust all issuers or not<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexsslheadersindex">headers</a></b></td>
        <td>[]object</td>
        <td>
          Http headers<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexsslkeystore">keyStore</a></b></td>
        <td>object</td>
        <td>
          KeyStore type (possible values PEM, PKCS12, JKS)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecproxygroupsindexssltruststore">trustStore</a></b></td>
        <td>object</td>
        <td>
          TrustStore type (possible values PEM, PKCS12, JKS)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].ssl.headers[index]
[Go to parent definition](#apidefinitionspecproxygroupsindexssl)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The HTTP header name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The HTTP header value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].ssl.keyStore
[Go to parent definition](#apidefinitionspecproxygroupsindexssl)



KeyStore type (possible values PEM, PKCS12, JKS)

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>certContent</b></td>
        <td>string</td>
        <td>
          KeyStore cert content (Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>certPath</b></td>
        <td>string</td>
        <td>
          KeyStore cert path (Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>content</b></td>
        <td>string</td>
        <td>
          The base64 encoded trustStore content, if not relying on a path to a file<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keyContent</b></td>
        <td>string</td>
        <td>
          The base64 encoded trustStore content, if not relying on a path to a file
(Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keyPath</b></td>
        <td>string</td>
        <td>
          KeyStore key path (Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>password</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          KeyStore path<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          The KeyStore type to use (possible values are PEM, PKCS12, JKS)<br/>
          <br/>
            <i>Enum</i>: PEM, PKCS12, JKS<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.groups[index].ssl.trustStore
[Go to parent definition](#apidefinitionspecproxygroupsindexssl)



TrustStore type (possible values PEM, PKCS12, JKS)

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>content</b></td>
        <td>string</td>
        <td>
          The base64 encoded trustStore content, if not relying on a path to a file<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>password</b></td>
        <td>string</td>
        <td>
          TrustStore password (Not applicable for PEM TrustStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          The TrustStore type to use (possible values are PEM, PKCS12, JKS)<br/>
          <br/>
            <i>Enum</i>: PEM, PKCS12, JKS<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.logging
[Go to parent definition](#apidefinitionspecproxy)



Logging

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          The logging condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>content</b></td>
        <td>enum</td>
        <td>
          Which part of the request/response should be logged ?<br/>
          <br/>
            <i>Enum</i>: NONE, HEADERS, PAYLOADS, HEADERS_PAYLOADS<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>enum</td>
        <td>
          The logging mode.
CLIENT identifies the inbound request issued to the gateway,
while PROXY identifies the request issued to the upstream service.<br/>
          <br/>
            <i>Enum</i>: NONE, CLIENT, PROXY, CLIENT_PROXY<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scope</b></td>
        <td>enum</td>
        <td>
          The logging scope (which phase of the request roundtrip should be included in each log entry.<br/>
          <br/>
            <i>Enum</i>: NONE, REQUEST, RESPONSE, REQUEST_RESPONSE<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.proxy.virtual_hosts[index]
[Go to parent definition](#apidefinitionspecproxy)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>override_entrypoint</b></td>
        <td>boolean</td>
        <td>
          Indicate if Entrypoint should be overridden or not<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.resources[index]
[Go to parent definition](#apidefinitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Resource Configuration, arbitrary map of key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is resource enabled or not?<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Resource Name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecresourcesindexref">ref</a></b></td>
        <td>object</td>
        <td>
          Reference to a resource<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Resource Type<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.resources[index].ref
[Go to parent definition](#apidefinitionspecresourcesindex)



Reference to a resource

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.response_templates[key][key]
[Go to parent definition](#apidefinitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>body</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>headers</b></td>
        <td>map[string]string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>propagateErrorKeyToLogs</b></td>
        <td>boolean</td>
        <td>
          Propagate error key to logs<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>integer</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.services
[Go to parent definition](#apidefinitionspec)



Contains different services for the API (EndpointDiscovery, HealthCheck ...)

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#apidefinitionspecservicesdiscovery">discovery</a></b></td>
        <td>object</td>
        <td>
          Endpoint Discovery Service<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecservicesdynamic-property">dynamic-property</a></b></td>
        <td>object</td>
        <td>
          Dynamic Property Service<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecserviceshealth-check">health-check</a></b></td>
        <td>object</td>
        <td>
          Health Check Service<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.services.discovery
[Go to parent definition](#apidefinitionspecservices)



Endpoint Discovery Service

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Configuration, arbitrary map of key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is service enabled or not?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Service name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>provider</b></td>
        <td>string</td>
        <td>
          Provider name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>secondary</b></td>
        <td>boolean</td>
        <td>
          Is it secondary or not?<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tenants</b></td>
        <td>[]string</td>
        <td>
          List of tenants<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.services.dynamic-property
[Go to parent definition](#apidefinitionspecservices)



Dynamic Property Service

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Configuration, arbitrary map of key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is service enabled or not?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Service name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>provider</b></td>
        <td>enum</td>
        <td>
          <br/>
          <br/>
            <i>Enum</i>: HTTP<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>schedule</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.services.health-check
[Go to parent definition](#apidefinitionspecservices)



Health Check Service

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is service enabled or not?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Service name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>schedule</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecserviceshealth-checkstepsindex">steps</a></b></td>
        <td>[]object</td>
        <td>
          List of health check steps<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.services.health-check.steps[index]
[Go to parent definition](#apidefinitionspecserviceshealth-check)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Health Check Step Name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecserviceshealth-checkstepsindexrequest">request</a></b></td>
        <td>object</td>
        <td>
          Health Check Step Request<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecserviceshealth-checkstepsindexresponse">response</a></b></td>
        <td>object</td>
        <td>
          Health Check Step Response<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.services.health-check.steps[index].request
[Go to parent definition](#apidefinitionspecserviceshealth-checkstepsindex)



Health Check Step Request

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fromRoot</b></td>
        <td>boolean</td>
        <td>
          If true, the health check request will be issued without prepending the context path of the API.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>body</b></td>
        <td>string</td>
        <td>
          Health Check Request Body<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecserviceshealth-checkstepsindexrequestheadersindex">headers</a></b></td>
        <td>[]object</td>
        <td>
          List of HTTP headers to include in the health check request<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>method</b></td>
        <td>enum</td>
        <td>
          The HTTP method to use when issuing the health check request<br/>
          <br/>
            <i>Enum</i>: GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD, CONNECT, TRACE, OTHER<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          The path of the endpoint handling the health check request<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.services.health-check.steps[index].request.headers[index]
[Go to parent definition](#apidefinitionspecserviceshealth-checkstepsindexrequest)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The HTTP header name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The HTTP header value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.services.health-check.steps[index].response
[Go to parent definition](#apidefinitionspecserviceshealth-checkstepsindex)



Health Check Step Response

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>assertions</b></td>
        <td>[]string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.status
[Go to parent definition](#apidefinition)



ApiDefinitionStatus defines the observed state of API Definition.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>crossId</b></td>
        <td>string</td>
        <td>
          The Cross ID is used to identify an API that has been promoted from one environment to another.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>environmentId</b></td>
        <td>string</td>
        <td>
          The environment ID, if a management context has been defined to sync with an APIM instance<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionstatuserrors">errors</a></b></td>
        <td>object</td>
        <td>
          When API has been created regardless of errors, this field is
used to persist the error message encountered during admission<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the API definition in the Gravitee API Management instance (if an API context has been configured).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>organizationId</b></td>
        <td>string</td>
        <td>
          The organization ID, if a management context has been defined to sync with an APIM instance<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>plans</b></td>
        <td>map[string]string</td>
        <td>
          This field is used to store the list of plans that have been created
for the API definition if a management context has been defined
to sync the API with an APIM instance<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>processingStatus</b></td>
        <td>string</td>
        <td>
          The processing status of the API definition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>state</b></td>
        <td>enum</td>
        <td>
          The state of the API. Can be either STARTED or STOPPED.<br/>
          <br/>
            <i>Enum</i>: STARTED, STOPPED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subscriptions</b></td>
        <td>integer</td>
        <td>
          The number of subscriptions that reference the API<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.status.errors
[Go to parent definition](#apidefinitionstatus)



When API has been created regardless of errors, this field is
used to persist the error message encountered during admission

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>severe</b></td>
        <td>[]string</td>
        <td>
          severe errors do not pass admission and will block reconcile
hence, this field should always be during the admission phase
and is very unlikely to be persisted in the status<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>warning</b></td>
        <td>[]string</td>
        <td>
          warning errors do not block object reconciliation,
most of the time because the value is ignored or defaulted
when the API gets synced with APIM<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## ApiV4Definition

[gravitee.io/v1alpha1](#graviteeiov1alpha1)






ApiV4Definition is the Schema for the v4 apidefinitions API.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#apiv4definitionspec">spec</a></b></td>
        <td>object</td>
        <td>
          ApiV4DefinitionSpec defines the desired state of ApiDefinition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionstatus">status</a></b></td>
        <td>object</td>
        <td>
          ApiV4DefinitionStatus defines the observed state of API Definition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec
[Go to parent definition](#apiv4definition)



ApiV4DefinitionSpec defines the desired state of ApiDefinition.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindex">endpointGroups</a></b></td>
        <td>[]object</td>
        <td>
          List of Endpoint groups<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>listeners</b></td>
        <td>[]object</td>
        <td>
          List of listeners for this API<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          API name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Api Type (proxy or message)<br/>
          <br/>
            <i>Enum</i>: PROXY, MESSAGE, NATIVE<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>version</b></td>
        <td>string</td>
        <td>
          API version<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecanalytics">analytics</a></b></td>
        <td>object</td>
        <td>
          API Analytics (Not applicable for Native API)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>categories</b></td>
        <td>[]string</td>
        <td>
          The list of categories the API belongs to.
Categories are reflected in APIM portal so that consumers can easily find the APIs they need.<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecconsolenotificationconfiguration">consoleNotificationConfiguration</a></b></td>
        <td>object</td>
        <td>
          ConsoleNotification struct sent to the mAPI, not part of the CRD spec.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspeccontextref">contextRef</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>crossId</b></td>
        <td>string</td>
        <td>
          When promoting an API from one environment to the other,
this ID identifies the API across those different environments.
Setting this ID also allows to take control over an existing API on an APIM instance
(by setting the same value as defined in APIM).
If empty, a UUID will be generated based on the namespace and name of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecdefinitioncontext">definitionContext</a></b></td>
        <td>object</td>
        <td>
          The API Definition context is used to identify the Kubernetes origin of the API,
and define whether the API definition should be synchronized
from an API instance or from a config map created in the cluster (which is the default)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>definitionVersion</b></td>
        <td>enum</td>
        <td>
          The definition version of the API.<br/>
          <br/>
            <i>Enum</i>: V4<br/>
            <i>Default</i>: V4<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          API description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecfailover">failover</a></b></td>
        <td>object</td>
        <td>
          API Failover<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowexecution">flowExecution</a></b></td>
        <td>object</td>
        <td>
          API Flow Execution (Not applicable for Native API)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindex">flows</a></b></td>
        <td>[]object</td>
        <td>
          List of flows for the API<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecgrouprefsindex">groupRefs</a></b></td>
        <td>[]object</td>
        <td>
          List of group references associated with the API
These groups are references to Group custom resources created on the cluster.<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>groups</b></td>
        <td>[]string</td>
        <td>
          List of groups associated with the API.
This groups are id or name references to existing groups in APIM.<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The API ID. If empty, this field will take the value of the `metadata.uid`
field of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>labels</b></td>
        <td>[]string</td>
        <td>
          List of labels of the API<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lifecycleState</b></td>
        <td>enum</td>
        <td>
          API life cycle state can be one of the values PUBLISHED, UNPUBLISHED<br/>
          <br/>
            <i>Enum</i>: PUBLISHED, UNPUBLISHED<br/>
            <i>Default</i>: UNPUBLISHED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecmembersindex">members</a></b></td>
        <td>[]object</td>
        <td>
          List of members associated with the API<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecmetadataindex">metadata</a></b></td>
        <td>[]object</td>
        <td>
          List of API metadata entries<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecnotificationsrefsindex">notificationsRefs</a></b></td>
        <td>[]object</td>
        <td>
          References to Notification custom resources to setup notifications.
For an API Notification CRD `eventType` field must be set to `api`
and only events set via `apiEvents` attributes are used.
Only one notification with `target` equals to `console` is admitted.<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>notifyMembers</b></td>
        <td>boolean</td>
        <td>
          If true, new members added to the API spec will
be notified when the API is synced with APIM.<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecpageskey">pages</a></b></td>
        <td>map[string]object</td>
        <td>
          A map of pages objects.


Keys uniquely identify pages and are used to keep them in sync
with APIM when using a management context.


Renaming a key is the equivalent of deleting the page and recreating
it holding a new ID in APIM.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskey">plans</a></b></td>
        <td>map[string]object</td>
        <td>
          A map of plan identifiers to plan
Keys uniquely identify plans and are used to keep them in sync
when using a management context.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecpropertiesindex">properties</a></b></td>
        <td>[]object</td>
        <td>
          List of Properties for the API<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecresourcesindex">resources</a></b></td>
        <td>[]object</td>
        <td>
          Resources can be either inlined or reference the namespace and name
of an <a href="#apiresource">existing API resource definition</a>.<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecresponsetemplateskeykey">responseTemplates</a></b></td>
        <td>map[string]map[string]object</td>
        <td>
          A list of Response Templates for the API (Not applicable for Native API)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecservices">services</a></b></td>
        <td>object</td>
        <td>
          API Services (Not applicable for Native API)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>state</b></td>
        <td>enum</td>
        <td>
          The state of API (setting the value to `STOPPED` will make the API un-reachable from the gateway)<br/>
          <br/>
            <i>Enum</i>: STARTED, STOPPED<br/>
            <i>Default</i>: STARTED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tags</b></td>
        <td>[]string</td>
        <td>
          List of Tags of the API<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>visibility</b></td>
        <td>enum</td>
        <td>
          Should the API be publicly available from the portal or not ?<br/>
          <br/>
            <i>Enum</i>: PUBLIC, PRIVATE<br/>
            <i>Default</i>: PRIVATE<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index]
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Endpoint group name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexendpointsindex">endpoints</a></b></td>
        <td>[]object</td>
        <td>
          List of endpoint for the group<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>headers</b></td>
        <td>map[string]string</td>
        <td>
          Endpoint group headers, arbitrary map of key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexhttp">http</a></b></td>
        <td>object</td>
        <td>
          Endpoint group http client options<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexloadbalancer">loadBalancer</a></b></td>
        <td>object</td>
        <td>
          Endpoint group load balancer<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexservices">services</a></b></td>
        <td>object</td>
        <td>
          Endpoint group services<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sharedConfiguration</b></td>
        <td>object</td>
        <td>
          Endpoint group shared configuration, arbitrary map of key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexssl">ssl</a></b></td>
        <td>object</td>
        <td>
          Endpoint group http client SSL options<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Endpoint group type<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].endpoints[index]
[Go to parent definition](#apiv4definitionspecendpointgroupsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>inheritConfiguration</b></td>
        <td>boolean</td>
        <td>
          Should endpoint group configuration be inherited or not ?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>secondary</b></td>
        <td>boolean</td>
        <td>
          Endpoint is secondary or not?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Endpoint Configuration, arbitrary map of key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The endpoint name (this value should be unique across endpoints)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexendpointsindexservices">services</a></b></td>
        <td>object</td>
        <td>
          Endpoint Services<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sharedConfigurationOverride</b></td>
        <td>object</td>
        <td>
          Endpoint Configuration Override, arbitrary map of key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tenants</b></td>
        <td>[]string</td>
        <td>
          List of endpoint tenants<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Endpoint Type<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>weight</b></td>
        <td>integer</td>
        <td>
          Endpoint Weight<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].endpoints[index].services
[Go to parent definition](#apiv4definitionspecendpointgroupsindexendpointsindex)



Endpoint Services

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexendpointsindexserviceshealthcheck">healthCheck</a></b></td>
        <td>object</td>
        <td>
          Health check service<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].endpoints[index].services.healthCheck
[Go to parent definition](#apiv4definitionspecendpointgroupsindexendpointsindexservices)



Health check service

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is the service enabled or not ?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>overrideConfiguration</b></td>
        <td>boolean</td>
        <td>
          Service Override Configuration or not?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Service Configuration, a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Service Type<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].http
[Go to parent definition](#apiv4definitionspecendpointgroupsindex)



Endpoint group http client options

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>followRedirects</b></td>
        <td>boolean</td>
        <td>
          Should HTTP redirects be followed or not ?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>pipelining</b></td>
        <td>boolean</td>
        <td>
          Should HTTP/1.1 pipelining be used for the connection or not ?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>useCompression</b></td>
        <td>boolean</td>
        <td>
          Should compression be used or not ?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>clearTextUpgrade</b></td>
        <td>boolean</td>
        <td>
          Should HTTP/2 clear text upgrade be used or not ?<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>connectTimeout</b></td>
        <td>integer</td>
        <td>
          Connection timeout of the http connection<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>idleTimeout</b></td>
        <td>integer</td>
        <td>
           Idle Timeout for the http connection<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keepAlive</b></td>
        <td>boolean</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keepAliveTimeout</b></td>
        <td>integer</td>
        <td>
          Should keep alive be used for the HTTP connection ?<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Default</i>: 30000<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxConcurrentConnections</b></td>
        <td>integer</td>
        <td>
          HTTP max concurrent connections<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>propagateClientAcceptEncoding</b></td>
        <td>boolean</td>
        <td>
          Propagate Client Accept-Encoding header<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readTimeout</b></td>
        <td>integer</td>
        <td>
          Read timeout<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>version</b></td>
        <td>enum</td>
        <td>
          HTTP Protocol Version (Possible values Http1 or Http2)<br/>
          <br/>
            <i>Enum</i>: HTTP_1_1, HTTP_2<br/>
            <i>Default</i>: HTTP_1_1<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].loadBalancer
[Go to parent definition](#apiv4definitionspecendpointgroupsindex)



Endpoint group load balancer

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          <br/>
          <br/>
            <i>Enum</i>: ROUND_ROBIN, RANDOM, WEIGHTED_ROUND_ROBIN, WEIGHTED_RANDOM<br/>
            <i>Default</i>: ROUND_ROBIN<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].services
[Go to parent definition](#apiv4definitionspecendpointgroupsindex)



Endpoint group services

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexservicesdiscovery">discovery</a></b></td>
        <td>object</td>
        <td>
          Endpoint group discovery service<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexserviceshealthcheck">healthCheck</a></b></td>
        <td>object</td>
        <td>
          Endpoint group health check service<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].services.discovery
[Go to parent definition](#apiv4definitionspecendpointgroupsindexservices)



Endpoint group discovery service

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is the service enabled or not ?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>overrideConfiguration</b></td>
        <td>boolean</td>
        <td>
          Service Override Configuration or not?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Service Configuration, a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Service Type<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].services.healthCheck
[Go to parent definition](#apiv4definitionspecendpointgroupsindexservices)



Endpoint group health check service

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is the service enabled or not ?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>overrideConfiguration</b></td>
        <td>boolean</td>
        <td>
          Service Override Configuration or not?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Service Configuration, a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Service Type<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].ssl
[Go to parent definition](#apiv4definitionspecendpointgroupsindex)



Endpoint group http client SSL options

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>hostnameVerifier</b></td>
        <td>boolean</td>
        <td>
          Verify Hostname when establishing connection<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>trustAll</b></td>
        <td>boolean</td>
        <td>
          Whether to trust all issuers or not<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexsslheadersindex">headers</a></b></td>
        <td>[]object</td>
        <td>
          Http headers<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexsslkeystore">keyStore</a></b></td>
        <td>object</td>
        <td>
          KeyStore type (possible values PEM, PKCS12, JKS)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecendpointgroupsindexssltruststore">trustStore</a></b></td>
        <td>object</td>
        <td>
          TrustStore type (possible values PEM, PKCS12, JKS)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].ssl.headers[index]
[Go to parent definition](#apiv4definitionspecendpointgroupsindexssl)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The HTTP header name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The HTTP header value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].ssl.keyStore
[Go to parent definition](#apiv4definitionspecendpointgroupsindexssl)



KeyStore type (possible values PEM, PKCS12, JKS)

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>certContent</b></td>
        <td>string</td>
        <td>
          KeyStore cert content (Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>certPath</b></td>
        <td>string</td>
        <td>
          KeyStore cert path (Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>content</b></td>
        <td>string</td>
        <td>
          The base64 encoded trustStore content, if not relying on a path to a file<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keyContent</b></td>
        <td>string</td>
        <td>
          The base64 encoded trustStore content, if not relying on a path to a file
(Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keyPath</b></td>
        <td>string</td>
        <td>
          KeyStore key path (Only applicable for PEM KeyStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>password</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          KeyStore path<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          The KeyStore type to use (possible values are PEM, PKCS12, JKS)<br/>
          <br/>
            <i>Enum</i>: PEM, PKCS12, JKS<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.endpointGroups[index].ssl.trustStore
[Go to parent definition](#apiv4definitionspecendpointgroupsindexssl)



TrustStore type (possible values PEM, PKCS12, JKS)

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>content</b></td>
        <td>string</td>
        <td>
          The base64 encoded trustStore content, if not relying on a path to a file<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>password</b></td>
        <td>string</td>
        <td>
          TrustStore password (Not applicable for PEM TrustStore)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          The TrustStore type to use (possible values are PEM, PKCS12, JKS)<br/>
          <br/>
            <i>Enum</i>: PEM, PKCS12, JKS<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.analytics
[Go to parent definition](#apiv4definitionspec)



API Analytics (Not applicable for Native API)

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Analytics Enabled or not?<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecanalyticslogging">logging</a></b></td>
        <td>object</td>
        <td>
          Analytics Logging<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecanalyticssampling">sampling</a></b></td>
        <td>object</td>
        <td>
          Analytics Sampling<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecanalyticstracing">tracing</a></b></td>
        <td>object</td>
        <td>
          Analytics Tracing<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.analytics.logging
[Go to parent definition](#apiv4definitionspecanalytics)



Analytics Logging

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          The logging condition. This field is evaluated for HTTP requests and supports EL expressions.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecanalyticsloggingcontent">content</a></b></td>
        <td>object</td>
        <td>
          Defines which component of the request should be included in the log payload.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The logging message condition. This field is evaluated for messages and supports EL expressions.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecanalyticsloggingmode">mode</a></b></td>
        <td>object</td>
        <td>
          The logging mode defines which "hop" of the request roundtrip
should be included in the log payload.
This can be either the inbound request to the gateway,
the request issued by the gateway to the upstream service, or both.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecanalyticsloggingphase">phase</a></b></td>
        <td>object</td>
        <td>
          Defines which phase of the request roundtrip
should be included in the log payload.
This can be either the request phase, the response phase, or both.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.analytics.logging.content
[Go to parent definition](#apiv4definitionspecanalyticslogging)



Defines which component of the request should be included in the log payload.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>headers</b></td>
        <td>boolean</td>
        <td>
          Should HTTP headers be logged or not ?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>messageHeaders</b></td>
        <td>boolean</td>
        <td>
          Should message headers be logged or not ?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>messageMetadata</b></td>
        <td>boolean</td>
        <td>
          Should message metadata be logged or not ?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>messagePayload</b></td>
        <td>boolean</td>
        <td>
          Should message payloads be logged or not ?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>payload</b></td>
        <td>boolean</td>
        <td>
          Should HTTP payloads be logged or not ?<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.analytics.logging.mode
[Go to parent definition](#apiv4definitionspecanalyticslogging)



The logging mode defines which "hop" of the request roundtrip
should be included in the log payload.
This can be either the inbound request to the gateway,
the request issued by the gateway to the upstream service, or both.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>endpoint</b></td>
        <td>boolean</td>
        <td>
          If true, the request to the upstream service will be included in the log payload<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>entrypoint</b></td>
        <td>boolean</td>
        <td>
          If true, the inbound request to the gateway will be included in the log payload<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.analytics.logging.phase
[Go to parent definition](#apiv4definitionspecanalyticslogging)



Defines which phase of the request roundtrip
should be included in the log payload.
This can be either the request phase, the response phase, or both.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>request</b></td>
        <td>boolean</td>
        <td>
          Should the request phase of the request roundtrip be included in the log payload or not ?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>response</b></td>
        <td>boolean</td>
        <td>
          Should the response phase of the request roundtrip be included in the log payload or not ?<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.analytics.sampling
[Go to parent definition](#apiv4definitionspecanalytics)



Analytics Sampling

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          The sampling type to use<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Sampling Value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.analytics.tracing
[Go to parent definition](#apiv4definitionspecanalytics)



Analytics Tracing

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Specify if Tracing is Enabled or not<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>verbose</b></td>
        <td>boolean</td>
        <td>
          Specify if Tracing is Verbose or not<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.consoleNotificationConfiguration
[Go to parent definition](#apiv4definitionspec)



ConsoleNotification struct sent to the mAPI, not part of the CRD spec.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>config_type</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>groups</b></td>
        <td>[]string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>hooks</b></td>
        <td>[]string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>origin</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>referenceId</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>referenceType</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.contextRef
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.definitionContext
[Go to parent definition](#apiv4definitionspec)



The API Definition context is used to identify the Kubernetes origin of the API,
and define whether the API definition should be synchronized
from an API instance or from a config map created in the cluster (which is the default)

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>origin</b></td>
        <td>enum</td>
        <td>
          The definition context origin where the API definition is managed.
The value is always `KUBERNETES` for an API managed by the operator.<br/>
          <br/>
            <i>Enum</i>: KUBERNETES<br/>
            <i>Default</i>: KUBERNETES<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>syncFrom</b></td>
        <td>enum</td>
        <td>
          The syncFrom field defines where the gateways should source the API definition from.
If the value is `MANAGEMENT`, the API definition will be sourced from an APIM instance.
This means that the API definition *must* hold a context reference in that case.
Setting the value to `MANAGEMENT` allows to make an API definition available on
gateways deployed across multiple clusters / regions.
If the value is `KUBERNETES`, the API definition will be sourced from a config map.
This means that only gateways deployed in the same cluster will be able to sync the API definition.<br/>
          <br/>
            <i>Enum</i>: KUBERNETES, MANAGEMENT<br/>
            <i>Default</i>: MANAGEMENT<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.failover
[Go to parent definition](#apiv4definitionspec)



API Failover

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          API Failover is enabled?<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxFailures</b></td>
        <td>integer</td>
        <td>
          API Failover max failures<br/>
          <br/>
            <i>Default</i>: 5<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxRetries</b></td>
        <td>integer</td>
        <td>
          API Failover max retires<br/>
          <br/>
            <i>Default</i>: 2<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>openStateDuration</b></td>
        <td>integer</td>
        <td>
          API Failover  open state duration<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Default</i>: 10000<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>perSubscription</b></td>
        <td>boolean</td>
        <td>
          API Failover  per subscription<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>slowCallDuration</b></td>
        <td>integer</td>
        <td>
          API Failover slow call duration<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Default</i>: 2000<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flowExecution
[Go to parent definition](#apiv4definitionspec)



API Flow Execution (Not applicable for Native API)

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>matchRequired</b></td>
        <td>boolean</td>
        <td>
          Is match required or not ? If set to true, a 404 status response will be returned if no matching flow was found.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>string</td>
        <td>
          The flow mode to use<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index]
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is flow enabled or not?<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexconnectindex">connect</a></b></td>
        <td>[]object</td>
        <td>
          List of Connect flow steps (Only available for Native APIs)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the flow this field is mainly used for compatibility with
APIM exports and can be safely ignored.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexinteractindex">interact</a></b></td>
        <td>[]object</td>
        <td>
          List of Publish flow steps (Only available for Native APIs)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Flow name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexpublishindex">publish</a></b></td>
        <td>[]object</td>
        <td>
          List of Publish flow steps<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexrequestindex">request</a></b></td>
        <td>[]object</td>
        <td>
          List of Request flow steps (NOT available for Native APIs)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexresponseindex">response</a></b></td>
        <td>[]object</td>
        <td>
          List of Response flow steps (NOT available for Native APIs)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>selectors</b></td>
        <td>[]object</td>
        <td>
          List of Flow selectors<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexsubscribeindex">subscribe</a></b></td>
        <td>[]object</td>
        <td>
          List of Subscribe flow steps<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tags</b></td>
        <td>[]string</td>
        <td>
          List of tags<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].connect[index]
[Go to parent definition](#apiv4definitionspecflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexconnectindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].connect[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecflowsindexconnectindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].interact[index]
[Go to parent definition](#apiv4definitionspecflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexinteractindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].interact[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecflowsindexinteractindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].publish[index]
[Go to parent definition](#apiv4definitionspecflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexpublishindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].publish[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecflowsindexpublishindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].request[index]
[Go to parent definition](#apiv4definitionspecflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexrequestindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].request[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecflowsindexrequestindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].response[index]
[Go to parent definition](#apiv4definitionspecflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexresponseindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].response[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecflowsindexresponseindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].subscribe[index]
[Go to parent definition](#apiv4definitionspecflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexsubscribeindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flows[index].subscribe[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecflowsindexsubscribeindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.groupRefs[index]
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.members[index]
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>source</b></td>
        <td>string</td>
        <td>
          Member source<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>sourceId</b></td>
        <td>string</td>
        <td>
          Member source ID<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>role</b></td>
        <td>string</td>
        <td>
          The API role associated with this Member<br/>
          <br/>
            <i>Default</i>: USER<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.metadata[index]
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>format</b></td>
        <td>enum</td>
        <td>
          Metadata Format<br/>
          <br/>
            <i>Enum</i>: STRING, NUMERIC, BOOLEAN, DATE, MAIL, URL<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Metadata Key<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Metadata Name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>defaultValue</b></td>
        <td>string</td>
        <td>
          Metadata Default value<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Metadata Value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.notificationsRefs[index]
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.pages[key]
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          The type of the documentation page or folder.<br/>
          <br/>
            <i>Enum</i>: MARKDOWN, SWAGGER, ASYNCAPI, ASCIIDOC, FOLDER, SYSTEM_FOLDER, ROOT<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>api</b></td>
        <td>string</td>
        <td>
          The API of the page. If empty, will be set automatically to the generated ID of the API.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>map[string]string</td>
        <td>
          Custom page configuration (e.g. page rendering can be changed to use Redoc instead of Swagger ui)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>content</b></td>
        <td>string</td>
        <td>
          The content of the page, if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>crossId</b></td>
        <td>string</td>
        <td>
          CrossID is designed to identified a page across environments.
If not set, this ID will be generated in a predictable manner based on
the map key associated to this entry in the API.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>homepage</b></td>
        <td>boolean</td>
        <td>
          If true, this page will be displayed as the homepage of your API documentation.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the page. This field is mostly required when you are applying
an API exported from APIM to make the operator take control over it.
If not set, this ID will be generated in a predictable manner based on
the map key associated to this entry in the API.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          This is the display name of the page in APIM and on the portal.
This field can be edited safely if you want to rename a page.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>order</b></td>
        <td>integer</td>
        <td>
          The order used to display the page in APIM and on the portal.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>parent</b></td>
        <td>string</td>
        <td>
          If your page contains a folder, setting this field to the map key associated to the
folder entry will be reflected into APIM by making the page a child of this folder.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>parentId</b></td>
        <td>string</td>
        <td>
          The parent ID of the page. This field is mostly required when you are applying
an API exported from APIM to make the operator take control over it. Use `Parent`
in any other case.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>published</b></td>
        <td>boolean</td>
        <td>
          If true, the page will be accessible from the portal (default is false)<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecpageskeysource">source</a></b></td>
        <td>object</td>
        <td>
          Source allow you to fetch pages from various external sources, overriding page content
each time the source is fetched.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>visibility</b></td>
        <td>enum</td>
        <td>
          The visibility of the page.<br/>
          <br/>
            <i>Enum</i>: PUBLIC, PRIVATE<br/>
            <i>Default</i>: PUBLIC<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.pages[key].source
[Go to parent definition](#apiv4definitionspecpageskey)



Source allow you to fetch pages from various external sources, overriding page content
each time the source is fetched.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key]
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Plan display name, this will be the name displayed in the UI
if a management context is used to sync the API with APIM<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>characteristics</b></td>
        <td>[]string</td>
        <td>
          List of plan characteristics<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>comment_required</b></td>
        <td>boolean</td>
        <td>
          Indicate of comment is required for this plan or not<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>crossId</b></td>
        <td>string</td>
        <td>
          The plan Cross ID.
This field is used to identify plans defined for an API
that has been promoted between different environments.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>definitionVersion</b></td>
        <td>string</td>
        <td>
          Plan definition version<br/>
          <br/>
            <i>Default</i>: V4<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          Plan Description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>excludedGroups</b></td>
        <td>[]string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindex">flows</a></b></td>
        <td>[]object</td>
        <td>
          List of plan flows<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>generalConditions</b></td>
        <td>string</td>
        <td>
          The general conditions defined to use this plan<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Plan ID<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>string</td>
        <td>
          The plan mode<br/>
          <br/>
            <i>Default</i>: STANDARD<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>order</b></td>
        <td>integer</td>
        <td>
          Plan order<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeysecurity">security</a></b></td>
        <td>object</td>
        <td>
          Plan security<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>selectionRule</b></td>
        <td>string</td>
        <td>
          Plan selection rule<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>enum</td>
        <td>
          The plan status<br/>
          <br/>
            <i>Enum</i>: PUBLISHED, DEPRECATED, STAGING<br/>
            <i>Default</i>: PUBLISHED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tags</b></td>
        <td>[]string</td>
        <td>
          List of plan tags<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Plan type<br/>
          <br/>
            <i>Enum</i>: API, CATALOG<br/>
            <i>Default</i>: API<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>validation</b></td>
        <td>enum</td>
        <td>
          Plan validation strategy<br/>
          <br/>
            <i>Enum</i>: AUTO, MANUAL<br/>
            <i>Default</i>: AUTO<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index]
[Go to parent definition](#apiv4definitionspecplanskey)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is flow enabled or not?<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexconnectindex">connect</a></b></td>
        <td>[]object</td>
        <td>
          List of Connect flow steps (Only available for Native APIs)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the flow this field is mainly used for compatibility with
APIM exports and can be safely ignored.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexinteractindex">interact</a></b></td>
        <td>[]object</td>
        <td>
          List of Publish flow steps (Only available for Native APIs)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Flow name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexpublishindex">publish</a></b></td>
        <td>[]object</td>
        <td>
          List of Publish flow steps<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexrequestindex">request</a></b></td>
        <td>[]object</td>
        <td>
          List of Request flow steps (NOT available for Native APIs)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexresponseindex">response</a></b></td>
        <td>[]object</td>
        <td>
          List of Response flow steps (NOT available for Native APIs)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>selectors</b></td>
        <td>[]object</td>
        <td>
          List of Flow selectors<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexsubscribeindex">subscribe</a></b></td>
        <td>[]object</td>
        <td>
          List of Subscribe flow steps<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tags</b></td>
        <td>[]string</td>
        <td>
          List of tags<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].connect[index]
[Go to parent definition](#apiv4definitionspecplanskeyflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexconnectindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].connect[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecplanskeyflowsindexconnectindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].interact[index]
[Go to parent definition](#apiv4definitionspecplanskeyflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexinteractindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].interact[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecplanskeyflowsindexinteractindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].publish[index]
[Go to parent definition](#apiv4definitionspecplanskeyflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexpublishindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].publish[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecplanskeyflowsindexpublishindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].request[index]
[Go to parent definition](#apiv4definitionspecplanskeyflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexrequestindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].request[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecplanskeyflowsindexrequestindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].response[index]
[Go to parent definition](#apiv4definitionspecplanskeyflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexresponseindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].response[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecplanskeyflowsindexresponseindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].subscribe[index]
[Go to parent definition](#apiv4definitionspecplanskeyflowsindex)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Indicate if this FlowStep is enabled or not<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>condition</b></td>
        <td>string</td>
        <td>
          FlowStep condition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          FlowStep configuration is a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          FlowStep description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>messageCondition</b></td>
        <td>string</td>
        <td>
          The message condition (supports EL expressions)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          FlowStep name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>policy</b></td>
        <td>string</td>
        <td>
          FlowStep policy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexsubscribeindexsharedpolicygroupref">sharedPolicyGroupRef</a></b></td>
        <td>object</td>
        <td>
          Reference to an existing Shared Policy Group<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].flows[index].subscribe[index].sharedPolicyGroupRef
[Go to parent definition](#apiv4definitionspecplanskeyflowsindexsubscribeindex)



Reference to an existing Shared Policy Group

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.plans[key].security
[Go to parent definition](#apiv4definitionspecplanskey)



Plan security

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Plan Security type<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Plan security configuration, a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.properties[index]
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>dynamic</b></td>
        <td>boolean</td>
        <td>
          Property is dynamic or not?<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>encryptable</b></td>
        <td>boolean</td>
        <td>
          Property is encryptable or not?<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>encrypted</b></td>
        <td>boolean</td>
        <td>
          Property Encrypted or not?<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Property Key<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Property Value<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.resources[index]
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Resource Configuration, arbitrary map of key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is resource enabled or not?<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Resource Name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecresourcesindexref">ref</a></b></td>
        <td>object</td>
        <td>
          Reference to a resource<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Resource Type<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.resources[index].ref
[Go to parent definition](#apiv4definitionspecresourcesindex)



Reference to a resource

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.responseTemplates[key][key]
[Go to parent definition](#apiv4definitionspec)





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>body</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>headers</b></td>
        <td>map[string]string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>propagateErrorKeyToLogs</b></td>
        <td>boolean</td>
        <td>
          Propagate error key to logs<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>integer</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.services
[Go to parent definition](#apiv4definitionspec)



API Services (Not applicable for Native API)

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#apiv4definitionspecservicesdynamicproperty">dynamicProperty</a></b></td>
        <td>object</td>
        <td>
          API dynamic property service<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.services.dynamicProperty
[Go to parent definition](#apiv4definitionspecservices)



API dynamic property service

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Is the service enabled or not ?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>overrideConfiguration</b></td>
        <td>boolean</td>
        <td>
          Service Override Configuration or not?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>configuration</b></td>
        <td>object</td>
        <td>
          Service Configuration, a map of arbitrary key-values<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Service Type<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.status
[Go to parent definition](#apiv4definition)



ApiV4DefinitionStatus defines the observed state of API Definition.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>crossId</b></td>
        <td>string</td>
        <td>
          The Cross ID is used to identify an API that has been promoted from one environment to another.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>environmentId</b></td>
        <td>string</td>
        <td>
          The environment ID, if a management context has been defined to sync with an APIM instance<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionstatuserrors">errors</a></b></td>
        <td>object</td>
        <td>
          When API has been created regardless of errors, this field is
used to persist the error message encountered during admission<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the API definition in the Gravitee API Management instance (if an API context has been configured).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>organizationId</b></td>
        <td>string</td>
        <td>
          The organization ID, if a management context has been defined to sync with an APIM instance<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>plans</b></td>
        <td>map[string]string</td>
        <td>
          This field is used to store the list of plans that have been created
for the API definition if a management context has been defined
to sync the API with an APIM instance<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>processingStatus</b></td>
        <td>string</td>
        <td>
          The processing status of the API definition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>state</b></td>
        <td>enum</td>
        <td>
          The state of the API. Can be either STARTED or STOPPED.<br/>
          <br/>
            <i>Enum</i>: STARTED, STOPPED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subscriptions</b></td>
        <td>integer</td>
        <td>
          The number of subscriptions that reference the API<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.status.errors
[Go to parent definition](#apiv4definitionstatus)



When API has been created regardless of errors, this field is
used to persist the error message encountered during admission

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>severe</b></td>
        <td>[]string</td>
        <td>
          severe errors do not pass admission and will block reconcile
hence, this field should always be during the admission phase
and is very unlikely to be persisted in the status<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>warning</b></td>
        <td>[]string</td>
        <td>
          warning errors do not block object reconciliation,
most of the time because the value is ignored or defaulted
when the API gets synced with APIM<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## ApiResource

[gravitee.io/v1alpha1](#graviteeiov1alpha1)








<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#apiresourcespec">spec</a></b></td>
        <td>object</td>
        <td>
          ApiResourceSpec defines the desired state of ApiResource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiResource.spec
[Go to parent definition](#apiresource)
=======
#### ApiResourceSpec
>>>>>>> 56d297b (docs: generate docs using crd-ref-docs)



ApiResourceSpec defines the desired state of ApiResource.



_Appears in:_
- [ApiResource](#apiresource)



#### ApiResourceStatus







_Appears in:_
- [ApiResource](#apiresource)



#### ApiV4Definition



ApiV4Definition is the Schema for the v4 apidefinitions API.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `gravitee.io/v1alpha1` | | |
| `kind` _string_ | `ApiV4Definition` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ApiV4DefinitionSpec](#apiv4definitionspec)_ |  |  |  |
| `status` _[ApiV4DefinitionStatus](#apiv4definitionstatus)_ |  |  |  |


#### ApiV4DefinitionSpec



ApiV4DefinitionSpec defines the desired state of ApiDefinition.



_Appears in:_
- [ApiV4Definition](#apiv4definition)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `description` _string_ | API description |  | Optional: \{\} <br /> |
| `definitionVersion` _[DefinitionVersion](#definitionversion)_ | The definition version of the API. | V4 | Enum: [V4] <br /> |
| `definitionContext` _[DefinitionContext](#definitioncontext)_ | The API Definition context is used to identify the Kubernetes origin of the API,<br />and define whether the API definition should be synchronized<br />from an API instance or from a config map created in the cluster (which is the default) |  |  |
| `lifecycleState` _[ApiV4LifecycleState](#apiv4lifecyclestate)_ | API life cycle state can be one of the values PUBLISHED, UNPUBLISHED | UNPUBLISHED | Enum: [PUBLISHED UNPUBLISHED] <br />Optional: \{\} <br /> |
| `type` _[ApiType](#apitype)_ | Api Type (proxy or message) |  | Enum: [PROXY MESSAGE NATIVE] <br />Required: \{\} <br /> |
| `listeners` _[GenericListener](#genericlistener) array_ | List of listeners for this API |  | MinItems: 1 <br />Required: \{\} <br /> |
| `endpointGroups` _[EndpointGroup](#endpointgroup) array_ | List of Endpoint groups |  | MinItems: 1 <br />Required: \{\} <br /> |
| `plans` _[map[string]*Plan](#map[string]*plan)_ | A map of plan identifiers to plan<br />Keys uniquely identify plans and are used to keep them in sync<br />when using a management context. |  | Optional: \{\} <br /> |
| `flowExecution` _[FlowExecution](#flowexecution)_ | API Flow Execution (Not applicable for Native API) |  |  |
| `flows` _[Flow](#flow) array_ | List of flows for the API | \{  \} | Optional: \{\} <br /> |
| `analytics` _[Analytics](#analytics)_ | API Analytics (Not applicable for Native API) |  |  |
| `services` _[ApiServices](#apiservices)_ | API Services (Not applicable for Native API) |  |  |
| `responseTemplates` _[ResponseTemplate](#responsetemplate)_ | A list of Response Templates for the API (Not applicable for Native API) |  | Optional: \{\} <br /> |
| `members` _Member array_ | List of members associated with the API |  | Optional: \{\} <br /> |
| `pages` _[map[string]*Page](#map[string]*page)_ | A map of pages objects.<br />Keys uniquely identify pages and are used to keep them in sync<br />with APIM when using a management context.<br />Renaming a key is the equivalent of deleting the page and recreating<br />it holding a new ID in APIM. |  | Optional: \{\} <br /> |
| `failover` _[Failover](#failover)_ | API Failover |  |  |
| `contextRef` _[NamespacedName](#namespacedname)_ |  |  |  |


#### ApiV4DefinitionStatus



ApiV4DefinitionStatus defines the observed state of API Definition.



_Appears in:_
- [ApiV4Definition](#apiv4definition)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `organizationId` _string_ | The organization ID, if a management context has been defined to sync with an APIM instance |  | Optional: \{\} <br /> |
| `environmentId` _string_ | The environment ID, if a management context has been defined to sync with an APIM instance |  | Optional: \{\} <br /> |
| `id` _string_ | The ID of the API definition in the Gravitee API Management instance (if an API context has been configured). |  | Optional: \{\} <br /> |
| `crossId` _string_ | The Cross ID is used to identify an API that has been promoted from one environment to another. |  |  |
| `processingStatus` _[ProcessingStatus](#processingstatus)_ | The processing status of the API definition. |  |  |
| `state` _[ApiState](#apistate)_ | The state of the API. Can be either STARTED or STOPPED. |  | Enum: [STARTED STOPPED] <br /> |
| `plans` _object (keys:string, values:string)_ | This field is used to store the list of plans that have been created<br />for the API definition if a management context has been defined<br />to sync the API with an APIM instance |  | Optional: \{\} <br /> |
| `subscriptions` _integer_ | The number of subscriptions that reference the API |  |  |
| `errors` _[Errors](#errors)_ | When API has been created regardless of errors, this field is<br />used to persist the error message encountered during admission |  |  |


#### Application









| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `gravitee.io/v1alpha1` | | |
| `kind` _string_ | `Application` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ApplicationSpec](#applicationspec)_ |  |  |  |
| `status` _[ApplicationStatus](#applicationstatus)_ |  |  |  |


#### ApplicationSpec



Application is the main resource handled by the Kubernetes Operator



_Appears in:_
- [Application](#application)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Application name |  | Required: \{\} <br /> |
| `description` _string_ | Application Description |  | Required: \{\} <br /> |
| `id` _string_ | io.gravitee.definition.model.Application<br />Application ID |  |  |
| `background` _string_ | The base64 encoded background to use for this application when displaying it on the portal |  | Optional: \{\} <br /> |
| `domain` _string_ | Application domain |  | Optional: \{\} <br /> |
| `groups` _string array_ | Application groups |  | Optional: \{\} <br /> |
| `picture` _string_ | The base64 encoded picture to use for this application when displaying it on the portal (if not relying on an URL) |  | Optional: \{\} <br /> |
| `pictureUrl` _string_ | A URL pointing to the picture to use when displaying the application on the portal |  | Optional: \{\} <br /> |
| `settings` _[Setting](#setting)_ | Application settings |  | Required: \{\} <br /> |
| `notifyMembers` _boolean_ | Notify members when they are added to the application |  | Optional: \{\} <br /> |
| `metadata` _[Metadata](#metadata)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `members` _[Member](#member)_ | Application members |  | Optional: \{\} <br /> |
| `contextRef` _[NamespacedName](#namespacedname)_ |  |  | Required: \{\} <br /> |


#### ApplicationStatus



ApplicationStatus defines the observed state of Application.



_Appears in:_
- [Application](#application)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `organizationId` _string_ | The organization ID, if a management context has been defined to sync with an APIM instance |  |  |
| `environmentId` _string_ | The environment ID, if a management context has been defined to sync with an APIM instance |  |  |
| `id` _string_ | The ID of the Application, if a management context has been defined to sync with an APIM instance |  |  |
| `processingStatus` _[ProcessingStatus](#processingstatus)_ | The processing status of the Application.<br />The value is `Completed` if the sync with APIM succeeded, Failed otherwise. |  |  |
| `subscriptions` _integer_ | The number of subscriptions that reference the application |  |  |
| `errors` _[Errors](#errors)_ | When application has been created regardless of errors, this field is<br />used to persist the error message encountered during admission |  |  |


#### GatewayClassParameters









| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `gravitee.io/v1alpha1` | | |
| `kind` _string_ | `GatewayClassParameters` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[GatewayClassParametersSpec](#gatewayclassparametersspec)_ |  |  |  |
| `status` _[GatewayClassParametersStatus](#gatewayclassparametersstatus)_ |  |  |  |


#### GatewayClassParametersSpec



GatewayClassParametersSpec defines the desired state of GatewayClassParameters



_Appears in:_
- [GatewayClassParameters](#gatewayclassparameters)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `gravitee` _[GraviteeConfig](#graviteeconfig)_ | The gravitee section controls Gravitee specific features<br />and allows you to configure and customize our implementation<br />of the Kubernetes Gateway API. |  | Optional: \{\} <br /> |
| `kubernetes` _[KubernetesConfig](#kubernetesconfig)_ | The kubernetes section of the GatewayClassParameters<br />spec lets you customize core Kubernetes resources<br />that are part of your Gateway deployments. |  | Optional: \{\} <br /> |


#### GatewayClassParametersStatus







_Appears in:_
- [GatewayClassParameters](#gatewayclassparameters)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#condition-v1-meta) array_ |  |  |  |


#### Group









| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `gravitee.io/v1alpha1` | | |
| `kind` _string_ | `Group` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[GroupSpec](#groupspec)_ |  |  |  |
| `status` _[GroupStatus](#groupstatus)_ |  |  |  |


#### GroupSpec







_Appears in:_
- [Group](#group)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `contextRef` _[NamespacedName](#namespacedname)_ |  |  |  |


#### GroupStatus







_Appears in:_
- [Group](#group)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `id` _string_ | The ID of the Group in the Gravitee API Management instance |  | Optional: \{\} <br /> |
| `organizationId` _string_ | The organization ID defined in the management context |  | Optional: \{\} <br /> |
| `environmentId` _string_ | The environment ID defined in the management context |  | Optional: \{\} <br /> |
| `processingStatus` _[ProcessingStatus](#processingstatus)_ | The processing status of the Group. |  |  |
| `members` _integer_ | The number of members added to this group |  |  |
| `errors` _[Errors](#errors)_ | When group has been created regardless of errors, this field is<br />used to persist the error message encountered during admission |  |  |


#### KafkaRoute









| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `gravitee.io/v1alpha1` | | |
| `kind` _string_ | `KafkaRoute` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[KafkaRouteSpec](#kafkaroutespec)_ |  |  |  |
| `status` _[KafkaRouteStatus](#kafkaroutestatus)_ |  |  |  |


#### KafkaRouteSpec







_Appears in:_
- [KafkaRoute](#kafkaroute)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `parentRefs` _[ParentReference](#parentreference) array_ | ParentRefs references the resources (usually Gateways) that a Route wants<br />to be attached to. Note that the referenced parent resource needs to<br />allow this for the attachment to be complete. For Gateways, that means<br />the Gateway needs to allow attachment from Routes of this kind and<br />namespace. For Services, that means the Service must either be in the same<br />namespace for a "producer" route, or the mesh implementation must support<br />and allow "consumer" routes for the referenced Service. ReferenceGrant is<br />not applicable for governing ParentRefs to Services - it is not possible to<br />create a "producer" route for a Service in a different namespace from the<br />Route.<br />There are two kinds of parent resources with "Core" support:<br />* Gateway (Gateway conformance profile)<br />* Service (Mesh conformance profile, ClusterIP Services only)<br />This API may be extended in the future to support additional kinds of parent<br />resources.<br />ParentRefs must be _distinct_. This means either that:<br />* They select different objects.  If this is the case, then parentRef<br />  entries are distinct. In terms of fields, this means that the<br />  multi-part key defined by `group`, `kind`, `namespace`, and `name` must<br />  be unique across all parentRef entries in the Route.<br />* They do not select different objects, but for each optional field used,<br />  each ParentRef that selects the same object must set the same set of<br />  optional fields to different values. If one ParentRef sets a<br />  combination of optional fields, all must set the same combination.<br />Some examples:<br />* If one ParentRef sets `sectionName`, all ParentRefs referencing the<br />  same object must also set `sectionName`.<br />* If one ParentRef sets `port`, all ParentRefs referencing the same<br />  object must also set `port`.<br />* If one ParentRef sets `sectionName` and `port`, all ParentRefs<br />  referencing the same object must also set `sectionName` and `port`.<br />It is possible to separately reference multiple distinct objects that may<br />be collapsed by an implementation. For example, some implementations may<br />choose to merge compatible Gateway Listeners together. If that is the<br />case, the list of routes attached to those resources should also be<br />merged.<br />Note that for ParentRefs that cross namespace boundaries, there are specific<br />rules. Cross-namespace references are only valid if they are explicitly<br />allowed by something in the namespace they are referring to. For example,<br />Gateway has the AllowedRoutes field, and ReferenceGrant provides a<br />generic way to enable other kinds of cross-namespace reference.<br /> |
| `hostname` _[Hostname](#hostname)_ | Hostname is used to uniquely route clients to this API.<br />Your client must trust the certificate provided by the gateway,<br />and as there is a variable host in the proxy bootstrap server URL,<br />you likely need to request a wildcard SAN for the certificate presented by the gateway.<br />If empty, the hostname defined in the Kafka listener of the parent will be used. |  | MaxLength: 253 <br />MinLength: 1 <br />Pattern: `^(\*\.)?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$` <br /> |
| `backendRefs` _[KafkaBackendRef](#kafkabackendref) array_ | BackendRefs defines the backend(s) where matching requests should be sent. |  | MaxItems: 16 <br />MinItems: 1 <br /> |
| `filters` _[KafkaRouteFilter](#kafkaroutefilter) array_ | Filters define the filters that are applied to Kafka trafic matching this route. |  | MaxItems: 16 <br /> |
| `options` _object (keys:[AnnotationKey](#annotationkey), values:[AnnotationValue](#annotationvalue))_ | Options are a list of key/value pairs to enable extended configuration specific<br />to an |  | MaxProperties: 16 <br /> |


#### KafkaRouteStatus







_Appears in:_
- [KafkaRoute](#kafkaroute)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `parents` _RouteParentStatus array_ | Parents is a list of parent resources (usually Gateways) that are<br />associated with the route, and the status of the route with respect to<br />each parent. When this route attaches to a parent, the controller that<br />manages the parent must add an entry to this list when the controller<br />first sees the route and should update the entry as appropriate when the<br />route or gateway is modified.<br />Note that parent references that cannot be resolved by an implementation<br />of this API will not be added to this list. Implementations of this API<br />can only populate Route status for the Gateways/parent resources they are<br />responsible for.<br />A maximum of 32 Gateways will be represented in this list. An empty list<br />means the route has not been attached to any Gateway. |  | MaxItems: 32 <br /> |


#### ManagementContext









| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `gravitee.io/v1alpha1` | | |
| `kind` _string_ | `ManagementContext` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ManagementContextSpec](#managementcontextspec)_ |  |  |  |
| `status` _[ManagementContextStatus](#managementcontextstatus)_ |  |  |  |


#### ManagementContextSpec



ManagementContext represents the configuration for a specific environment



_Appears in:_
- [ManagementContext](#managementcontext)



#### ManagementContextStatus



ManagementContextStatus defines the observed state of an API Context.



_Appears in:_
- [ManagementContext](#managementcontext)



#### Notification



Notification defines notification settings in Gravitee





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `gravitee.io/v1alpha1` | | |
| `kind` _string_ | `Notification` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[NotificationSpec](#notificationspec)_ |  |  |  |
| `status` _[NotificationStatus](#notificationstatus)_ |  |  |  |


#### NotificationSpec



NotificationSpec defines the desired state of a Notification.
It is to be referenced in an API.



_Appears in:_
- [Notification](#notification)



#### NotificationStatus



NotificationStatus defines the observed state of the Notification.

<<<<<<< HEAD
<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#notificationstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions are the condition that must be met by the Notification
"Accepted" condition is used to indicate if the `Notification` can be used by another resource.
"ResolveRef" condition is used to indicate if an error occurred while resolving console groups.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Notification.status.conditions[index]
[Go to parent definition](#notificationstatus)



Condition contains details for one aspect of the current state of this API Resource.
---
This struct is intended for direct use as an array at the field path .status.conditions.  For example,


	type FooStatus struct{
	    // Represents the observations of a foo's current state.
	    // Known .status.conditions.type are: "Available", "Progressing", and "Degraded"
	    // +patchMergeKey=type
	    // +patchStrategy=merge
	    // +listType=map
	    // +listMapKey=type
	    Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`


	    // other fields
	}

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          lastTransitionTime is the last time the condition transitioned from one status to another.
This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          message is a human readable message indicating details about the transition.
This may be an empty string.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          reason contains a programmatic identifier indicating the reason for the condition's last transition.
Producers of specific condition types may define expected values and meanings for this field,
and whether the values are considered a guaranteed API.
The value should be a CamelCase string.
This field may not be empty.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>enum</td>
        <td>
          status of the condition, one of True, False, Unknown.<br/>
          <br/>
            <i>Enum</i>: True, False, Unknown<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type of condition in CamelCase or in foo.example.com/CamelCase.
---
Many .condition.type values are consistent across resources like Available, but because arbitrary conditions can be
useful (see .node.status.conditions), the ability to deconflict is important.
The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt)<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          observedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Minimum</i>: 0<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## Group

[gravitee.io/v1alpha1](#graviteeiov1alpha1)
=======


_Appears in:_
- [Notification](#notification)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#condition-v1-meta)_ | Conditions are the condition that must be met by the Notification<br />"Accepted" condition is used to indicate if the `Notification` can be used by another resource.<br />"ResolveRef" condition is used to indicate if an error occurred while resolving console groups. |  |  |


#### SharedPolicyGroup



SharedPolicyGroup





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `gravitee.io/v1alpha1` | | |
| `kind` _string_ | `SharedPolicyGroup` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[SharedPolicyGroupSpec](#sharedpolicygroupspec)_ |  |  |  |
| `status` _[SharedPolicyGroupSpecStatus](#sharedpolicygroupspecstatus)_ |  |  |  |


#### SharedPolicyGroupSpec



SharedPolicyGroupSpec



_Appears in:_
- [SharedPolicyGroup](#sharedpolicygroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `contextRef` _[NamespacedName](#namespacedname)_ |  |  | Required: \{\} <br /> |


#### SharedPolicyGroupSpecStatus



SharedPolicyGroupSpecStatus defines the observed state of an API Context.



_Appears in:_
- [SharedPolicyGroup](#sharedpolicygroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `organizationId` _string_ | The organization ID, if a management context has been defined to sync with an APIM instance |  |  |
| `environmentId` _string_ | The environment ID, if a management context has been defined to sync with an APIM instance |  |  |
| `crossId` _string_ | The Cross ID is used to identify an SharedPolicyGroup that has been promoted from one environment to another. |  |  |
| `id` _string_ | The ID is used to identify an SharedPolicyGroup which is unique in any environment. |  |  |
| `processingStatus` _[ProcessingStatus](#processingstatus)_ | The processing status of the SharedPolicyGroup.<br />The value is `Completed` if the sync with APIM succeeded, Failed otherwise. |  |  |
| `errors` _[Errors](#errors)_ | When SharedPolicyGroup has been created regardless of errors, this field is<br />used to persist the error message encountered during admission |  |  |


#### Subscription
>>>>>>> 56d297b (docs: generate docs using crd-ref-docs)








<<<<<<< HEAD
<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#groupspec">spec</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#groupstatus">status</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Group.spec
[Go to parent definition](#group)
=======

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `gravitee.io/v1alpha1` | | |
| `kind` _string_ | `Subscription` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[SubscriptionSpec](#subscriptionspec)_ |  |  |  |
| `status` _[SubscriptionStatus](#subscriptionstatus)_ |  |  |  |


#### SubscriptionSpec
>>>>>>> 56d297b (docs: generate docs using crd-ref-docs)





<<<<<<< HEAD
<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#groupspecmembersindex">members</a></b></td>
        <td>[]object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#groupspeccontextref">contextRef</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>notifyMembers</b></td>
        <td>boolean</td>
        <td>
          If true, new members added to the API spec will
be notified when the API is synced with APIM.<br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Group.spec.members[index]
[Go to parent definition](#groupspec)
=======


_Appears in:_
- [Subscription](#subscription)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `api` _[NamespacedName](#namespacedname)_ |  |  | Required: \{\} <br /> |
| `application` _[NamespacedName](#namespacedname)_ |  |  | Required: \{\} <br /> |
| `plan` _string_ |  |  | Required: \{\} <br /> |
| `endingAt` _string_ |  |  | Format: date-time <br />Optional: \{\} <br /> |


#### SubscriptionStatus
>>>>>>> 56d297b (docs: generate docs using crd-ref-docs)





<<<<<<< HEAD
<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>source</b></td>
        <td>string</td>
        <td>
          Member source<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>sourceId</b></td>
        <td>string</td>
        <td>
          Member source ID<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>roles</b></td>
        <td>map[string]string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: map[]<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Group.spec.contextRef
[Go to parent definition](#groupspec)
=======


_Appears in:_
- [Subscription](#subscription)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `id` _string_ | Subscription ID |  |  |
| `startedAt` _string_ | When the subscription was started and made available |  |  |
| `endingAt` _string_ | The expiry date for the subscription (no date means no expiry) |  |  |
| `processingStatus` _[ProcessingStatus](#processingstatus)_ | This value is `Completed` if the sync with APIM succeeded, Failed otherwise. |  |  |



## gravitee.io/v1alpha1/application




#### Application
>>>>>>> 56d297b (docs: generate docs using crd-ref-docs)





<<<<<<< HEAD
<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Group.status
[Go to parent definition](#group)
=======


_Appears in:_
- [ApplicationSpec](#applicationspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Application name |  | Required: \{\} <br /> |
| `description` _string_ | Application Description |  | Required: \{\} <br /> |
| `id` _string_ | io.gravitee.definition.model.Application<br />Application ID |  |  |
| `background` _string_ | The base64 encoded background to use for this application when displaying it on the portal |  | Optional: \{\} <br /> |
| `domain` _string_ | Application domain |  | Optional: \{\} <br /> |
| `groups` _string array_ | Application groups |  | Optional: \{\} <br /> |
| `picture` _string_ | The base64 encoded picture to use for this application when displaying it on the portal (if not relying on an URL) |  | Optional: \{\} <br /> |
| `pictureUrl` _string_ | A URL pointing to the picture to use when displaying the application on the portal |  | Optional: \{\} <br /> |
| `settings` _[Setting](#setting)_ | Application settings |  | Required: \{\} <br /> |
| `notifyMembers` _boolean_ | Notify members when they are added to the application |  | Optional: \{\} <br /> |
| `metadata` _[Metadata](#metadata)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  | Optional: \{\} <br /> |
| `members` _[Member](#member)_ | Application members |  | Optional: \{\} <br /> |


#### GrantType

_Underlying type:_ _string_



_Validation:_
- Enum: [authorization_code client_credentials refresh_token password implicit]

_Appears in:_
- [OAuthClientSettings](#oauthclientsettings)

| Field | Description |
| --- | --- |
| `client_credentials` |  |
| `authorization_code` |  |
| `refresh_token` |  |
| `password` |  |
| `implicit` |  |


#### Member
>>>>>>> 56d297b (docs: generate docs using crd-ref-docs)





<<<<<<< HEAD
<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>members</b></td>
        <td>integer</td>
        <td>
          The number of members added to this group<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>environmentId</b></td>
        <td>string</td>
        <td>
          The environment ID defined in the management context<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#groupstatuserrors">errors</a></b></td>
        <td>object</td>
        <td>
          When group has been created regardless of errors, this field is
used to persist the error message encountered during admission<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the Group in the Gravitee API Management instance<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>organizationId</b></td>
        <td>string</td>
        <td>
          The organization ID defined in the management context<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>processingStatus</b></td>
        <td>string</td>
        <td>
          The processing status of the Group.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Group.status.errors
[Go to parent definition](#groupstatus)



When group has been created regardless of errors, this field is
used to persist the error message encountered during admission

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>severe</b></td>
        <td>[]string</td>
        <td>
          severe errors do not pass admission and will block reconcile
hence, this field should always be during the admission phase
and is very unlikely to be persisted in the status<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>warning</b></td>
        <td>[]string</td>
        <td>
          warning errors do not block object reconciliation,
most of the time because the value is ignored or defaulted
when the API gets synced with APIM<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
=======


_Appears in:_
- [Application](#application)
- [ApplicationSpec](#applicationspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `source` _string_ | Member source |  | Required: \{\} <br /> |
| `sourceId` _string_ | Member source ID |  | Required: \{\} <br /> |
| `role` _string_ | The API role associated with this Member | USER |  |


#### MetaDataFormat

_Underlying type:_ _string_



_Validation:_
- Enum: [STRING NUMERIC BOOLEAN DATE MAIL URL]

_Appears in:_
- [Metadata](#metadata)



#### Metadata







_Appears in:_
- [Application](#application)
- [ApplicationSpec](#applicationspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Metadata Name |  | Required: \{\} <br /> |
| `value` _string_ | Metadata Value |  | Optional: \{\} <br /> |
| `defaultValue` _string_ | Metadata DefaultValue |  | Optional: \{\} <br /> |
| `format` _[MetaDataFormat](#metadataformat)_ | Metadata Format |  | Enum: [STRING NUMERIC BOOLEAN DATE MAIL URL] <br /> |
| `hidden` _boolean_ | Metadata is hidden or not? |  | Optional: \{\} <br /> |


#### OAuthClientSettings







_Appears in:_
- [Setting](#setting)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `applicationType` _[OauthType](#oauthtype)_ | Oauth client application type |  | Enum: [BACKEND_TO_BACKEND NATIVE BROWSER WEB] <br />Required: \{\} <br /> |
| `grantTypes` _[GrantType](#granttype) array_ | List of Oauth client grant types |  | Enum: [authorization_code client_credentials refresh_token password implicit] <br /> |
| `redirectUris` _string array_ | List of Oauth client redirect uris |  | Optional: \{\} <br /> |


#### OauthType

_Underlying type:_ _string_



_Validation:_
- Enum: [BACKEND_TO_BACKEND NATIVE BROWSER WEB]

_Appears in:_
- [OAuthClientSettings](#oauthclientsettings)



#### Setting







_Appears in:_
- [Application](#application)
- [ApplicationSpec](#applicationspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `app` _[SimpleSettings](#simplesettings)_ |  |  |  |
| `oauth` _[OAuthClientSettings](#oauthclientsettings)_ |  |  |  |
| `tls` _[TLSSettings](#tlssettings)_ |  |  |  |


#### SimpleSettings







_Appears in:_
- [Setting](#setting)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _string_ | Application Type |  |  |
| `clientId` _string_ | ClientID is the client id of the application |  |  |


#### Status







_Appears in:_
- [ApplicationStatus](#applicationstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `organizationId` _string_ | The organization ID, if a management context has been defined to sync with an APIM instance |  |  |
| `environmentId` _string_ | The environment ID, if a management context has been defined to sync with an APIM instance |  |  |
| `id` _string_ | The ID of the Application, if a management context has been defined to sync with an APIM instance |  |  |
| `processingStatus` _[ProcessingStatus](#processingstatus)_ | The processing status of the Application.<br />The value is `Completed` if the sync with APIM succeeded, Failed otherwise. |  |  |
| `subscriptions` _integer_ | The number of subscriptions that reference the application |  |  |
| `errors` _[Errors](#errors)_ | When application has been created regardless of errors, this field is<br />used to persist the error message encountered during admission |  |  |


#### TLSSettings



TLS settings are used to configure client side TLS in order
to be able to subscribe to a MTLS plan.



_Appears in:_
- [Setting](#setting)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `clientCertificate` _string_ | This client certificate is mandatory to subscribe to a TLS plan. |  | Required: \{\} <br /> |



## gravitee.io/v1alpha1/base




#### AccessControl







_Appears in:_
- [Page](#page)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `referenceId` _string_ | The ID denied or granted by the access control (currently only group names are supported) |  | Required: \{\} <br /> |
| `referenceType` _string_ | The type of reference denied or granted by the access control<br />Currently only GROUP is supported |  | Enum: [GROUP] <br />Required: \{\} <br /> |


#### ApiBase







_Appears in:_
- [Api](#api)
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)
- [ApiV4DefinitionSpec](#apiv4definitionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `id` _string_ | The API ID. If empty, this field will take the value of the `metadata.uid`<br />field of the resource. |  |  |
| `crossId` _string_ | When promoting an API from one environment to the other,<br />this ID identifies the API across those different environments.<br />Setting this ID also allows to take control over an existing API on an APIM instance<br />(by setting the same value as defined in APIM).<br />If empty, a UUID will be generated based on the namespace and name of the resource. |  |  |
| `name` _string_ | API name |  |  |
| `version` _string_ | API version |  | Required: \{\} <br /> |
| `state` _[ApiState](#apistate)_ | The state of API (setting the value to `STOPPED` will make the API un-reachable from the gateway) | STARTED | Enum: [STARTED STOPPED] <br /> |
| `tags` _string array_ | List of Tags of the API |  | Optional: \{\} <br /> |
| `labels` _string array_ | List of labels of the API | \{  \} | Optional: \{\} <br /> |
| `visibility` _[ApiVisibility](#apivisibility)_ | Should the API be publicly available from the portal or not ? | PRIVATE | Enum: [PUBLIC PRIVATE] <br /> |
| `properties` _[Property](#property) array_ | List of Properties for the API | \{  \} | Optional: \{\} <br /> |
| `metadata` _[MetadataEntry](#metadataentry) array_ | Refer to Kubernetes API documentation for fields of `metadata`. | \{  \} | Optional: \{\} <br /> |
| `resources` _[ResourceOrRef](#resourceorref) array_ | Resources can be either inlined or reference the namespace and name<br />of an <a href="#apiresource">existing API resource definition</a>. | \{  \} | Optional: \{\} <br /> |
| `groups` _string array_ | List of groups associated with the API.<br />This groups are id or name references to existing groups in APIM. | \{  \} | Optional: \{\} <br /> |
| `groupRefs` _[NamespacedName](#namespacedname) array_ | List of group references associated with the API<br />These groups are references to Group custom resources created on the cluster. | \{  \} | Optional: \{\} <br /> |
| `categories` _string array_ | The list of categories the API belongs to.<br />Categories are reflected in APIM portal so that consumers can easily find the APIs they need. | \{  \} | Optional: \{\} <br /> |
| `notifyMembers` _boolean_ | If true, new members added to the API spec will<br />be notified when the API is synced with APIM. | true | Optional: \{\} <br /> |
| `notificationsRefs` _[NamespacedName](#namespacedname) array_ | References to Notification custom resources to setup notifications.<br />For an API Notification CRD `eventType` field must be set to `api`<br />and only events set via `apiEvents` attributes are used.<br />Only one notification with `target` equals to `console` is admitted. | \{  \} | Optional: \{\} <br /> |
| `consoleNotificationConfiguration` _[ConsoleNotificationConfiguration](#consolenotificationconfiguration)_ | ConsoleNotification struct sent to the mAPI, not part of the CRD spec. |  |  |


#### ApiState

_Underlying type:_ _string_



_Validation:_
- Enum: [STARTED STOPPED]

_Appears in:_
- [ApiBase](#apibase)
- [ApiDefinitionStatus](#apidefinitionstatus)
- [ApiV4DefinitionStatus](#apiv4definitionstatus)
- [Status](#status)

| Field | Description |
| --- | --- |
| `STARTED` |  |
| `STOPPED` |  |


#### ApiVisibility

_Underlying type:_ _string_



_Validation:_
- Enum: [PUBLIC PRIVATE]

_Appears in:_
- [ApiBase](#apibase)



#### ConsoleNotificationConfiguration



ConsoleNotificationConfiguration mAPI object to update notification settings.



_Appears in:_
- [ApiBase](#apibase)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `referenceType` _string_ |  |  |  |
| `referenceId` _string_ |  |  |  |
| `hooks` _string array_ |  |  |  |
| `config_type` _string_ |  |  |  |
| `origin` _string_ |  |  |  |
| `user` _string_ |  |  |  |
| `groups` _string array_ |  |  |  |


#### Cors







_Appears in:_
- [Proxy](#proxy)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Indicate if the cors enabled or not |  |  |
| `allowOrigin` _string array_ | Access Control -  List of Allowed origins | \{  \} | Optional: \{\} <br /> |
| `exposeHeaders` _string array_ | Access Control - List of Exposed Headers | \{  \} | Optional: \{\} <br /> |
| `maxAge` _integer_ | Access Control -  Max age |  |  |
| `allowCredentials` _boolean_ | Access Control - Allow credentials or not |  |  |
| `allowMethods` _string array_ | Access Control - List of allowed methods | \{  \} | Optional: \{\} <br /> |
| `allowHeaders` _string array_ | Access Control - List of allowed headers | \{  \} | Optional: \{\} <br /> |
| `runPolicies` _boolean_ | Run policies or not | false |  |


#### DefinitionVersion

_Underlying type:_ _string_





_Appears in:_
- [Api](#api)
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)
- [ApiV4DefinitionSpec](#apiv4definitionspec)

| Field | Description |
| --- | --- |
| `1.0.0` |  |
| `2.0.0` |  |
| `V4` |  |
| `4.0.0` |  |


#### FlowStep







_Appears in:_
- [Flow](#flow)
- [FlowStep](#flowstep)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Indicate if this FlowStep is enabled or not | true |  |
| `policy` _string_ | FlowStep policy |  | Optional: \{\} <br /> |
| `name` _string_ | FlowStep name |  | Optional: \{\} <br /> |
| `description` _string_ | FlowStep description |  | Optional: \{\} <br /> |
| `configuration` _[GenericStringMap](#genericstringmap)_ | FlowStep configuration is a map of arbitrary key-values |  | Optional: \{\} <br /> |
| `condition` _string_ | FlowStep condition |  | Optional: \{\} <br /> |


#### HttpClientOptions







_Appears in:_
- [Endpoint](#endpoint)
- [EndpointGroup](#endpointgroup)
- [EndpointGroup](#endpointgroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `idleTimeout` _integer_ |  Idle Timeout for the http connection |  | Optional: \{\} <br /> |
| `connectTimeout` _integer_ | Connection timeout of the http connection |  | Optional: \{\} <br /> |
| `keepAlive` _boolean_ |  | true | Optional: \{\} <br /> |
| `keepAliveTimeout` _integer_ | Should keep alive be used for the HTTP connection ? | 30000 | Optional: \{\} <br /> |
| `readTimeout` _integer_ | Read timeout |  | Optional: \{\} <br /> |
| `pipelining` _boolean_ | Should HTTP/1.1 pipelining be used for the connection or not ? | false |  |
| `maxConcurrentConnections` _integer_ | HTTP max concurrent connections |  | Optional: \{\} <br /> |
| `useCompression` _boolean_ | Should compression be used or not ? | false |  |
| `propagateClientAcceptEncoding` _boolean_ | Propagate Client Accept-Encoding header | false | Optional: \{\} <br /> |
| `followRedirects` _boolean_ | Should HTTP redirects be followed or not ? | false |  |
| `clearTextUpgrade` _boolean_ | Should HTTP/2 clear text upgrade be used or not ? | true | Optional: \{\} <br /> |
| `version` _[ProtocolVersion](#protocolversion)_ | HTTP Protocol Version (Possible values Http1 or Http2) | HTTP_1_1 | Enum: [HTTP_1_1 HTTP_2] <br /> |


#### HttpClientSslOptions







_Appears in:_
- [Endpoint](#endpoint)
- [EndpointGroup](#endpointgroup)
- [EndpointGroup](#endpointgroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `trustAll` _boolean_ | Whether to trust all issuers or not | false |  |
| `hostnameVerifier` _boolean_ | Verify Hostname when establishing connection | true |  |
| `trustStore` _[TrustStore](#truststore)_ | TrustStore type (possible values PEM, PKCS12, JKS) |  |  |
| `keyStore` _[KeyStore](#keystore)_ | KeyStore type (possible values PEM, PKCS12, JKS) |  |  |
| `headers` _[HttpHeader](#httpheader) array_ | Http headers |  |  |


#### HttpHeader







_Appears in:_
- [Endpoint](#endpoint)
- [HealthCheckRequest](#healthcheckrequest)
- [HttpClientSslOptions](#httpclientssloptions)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | The HTTP header name |  | Optional: \{\} <br /> |
| `value` _string_ | The HTTP header value |  | Optional: \{\} <br /> |


#### HttpMethod

_Underlying type:_ _string_



_Validation:_
- Enum: [GET POST PUT PATCH DELETE OPTIONS HEAD CONNECT TRACE OTHER]

_Appears in:_
- [Flow](#flow)
- [HealthCheckRequest](#healthcheckrequest)
- [Rule](#rule)



#### HttpProxy







_Appears in:_
- [Endpoint](#endpoint)
- [EndpointGroup](#endpointgroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Specifies that the HTTP connection will be established through a proxy | false |  |
| `useSystemProxy` _boolean_ | If true, the proxy defined at the system level will be used | false |  |
| `host` _string_ | Proxy host name |  | Optional: \{\} <br /> |
| `port` _integer_ | The HTTP proxy port |  | Optional: \{\} <br /> |
| `username` _string_ | The HTTP proxy username (if the proxy requires authentication) |  | Optional: \{\} <br /> |
| `password` _string_ | The HTTP proxy password (if the proxy requires authentication) |  | Optional: \{\} <br /> |
| `type` _[HttpProxyType](#httpproxytype)_ | The HTTP proxy type (possible values Http, Socks4, Socks5) |  |  |


#### HttpProxyType

_Underlying type:_ _string_





_Appears in:_
- [HttpProxy](#httpproxy)

| Field | Description |
| --- | --- |
| `HTTP` |  |


#### KeyStore







_Appears in:_
- [HttpClientSslOptions](#httpclientssloptions)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[KeyStoreType](#keystoretype)_ | The KeyStore type to use (possible values are PEM, PKCS12, JKS) |  | Enum: [PEM PKCS12 JKS] <br /> |
| `path` _string_ | KeyStore path |  | Optional: \{\} <br /> |
| `content` _string_ | The base64 encoded trustStore content, if not relying on a path to a file |  | Optional: \{\} <br /> |
| `password` _string_ |  |  | Optional: \{\} <br /> |
| `keyPath` _string_ | KeyStore key path (Only applicable for PEM KeyStore) |  | Optional: \{\} <br /> |
| `keyContent` _string_ | The base64 encoded trustStore content, if not relying on a path to a file<br />(Only applicable for PEM KeyStore) |  | Optional: \{\} <br /> |
| `certPath` _string_ | KeyStore cert path (Only applicable for PEM KeyStore) |  | Optional: \{\} <br /> |
| `certContent` _string_ | KeyStore cert content (Only applicable for PEM KeyStore) |  | Optional: \{\} <br /> |


#### KeyStoreType

_Underlying type:_ _string_



_Validation:_
- Enum: [PEM PKCS12 JKS]

_Appears in:_
- [KeyStore](#keystore)
- [TrustStore](#truststore)



#### LifecycleState

_Underlying type:_ _string_



_Validation:_
- Enum: [CREATED PUBLISHED UNPUBLISHED DEPRECATED ARCHIVED]

_Appears in:_
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)



#### Member







_Appears in:_
- [Api](#api)
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)
- [ApiV4DefinitionSpec](#apiv4definitionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `source` _string_ | Member source |  | Required: \{\} <br /> |
| `sourceId` _string_ | Member source ID |  | Required: \{\} <br /> |
| `role` _string_ | The API role associated with this Member | USER |  |


#### MetadataEntry







_Appears in:_
- [ApiBase](#apibase)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `key` _string_ | Metadata Key |  |  |
| `name` _string_ | Metadata Name |  |  |
| `format` _[MetadataFormat](#metadataformat)_ | Metadata Format |  | Enum: [STRING NUMERIC BOOLEAN DATE MAIL URL] <br /> |
| `value` _string_ | Metadata Value |  |  |
| `defaultValue` _string_ | Metadata Default value |  | Optional: \{\} <br /> |


#### MetadataFormat

_Underlying type:_ _string_



_Validation:_
- Enum: [STRING NUMERIC BOOLEAN DATE MAIL URL]

_Appears in:_
- [MetadataEntry](#metadataentry)



#### NotificationConfigurationBase



NotificationConfigurationBase base object for notifications.



_Appears in:_
- [ConsoleNotificationConfiguration](#consolenotificationconfiguration)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `referenceType` _string_ |  |  |  |
| `referenceId` _string_ |  |  |  |
| `hooks` _string array_ |  |  |  |
| `config_type` _string_ |  |  |  |
| `origin` _string_ |  |  |  |


#### Operator

_Underlying type:_ _string_



_Validation:_
- Enum: [STARTS_WITH EQUALS]

_Appears in:_
- [PathOperator](#pathoperator)



#### Page







_Appears in:_
- [Page](#page)
- [Page](#page)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `id` _string_ | The ID of the page. This field is mostly required when you are applying<br />an API exported from APIM to make the operator take control over it.<br />If not set, this ID will be generated in a predictable manner based on<br />the map key associated to this entry in the API. |  | Optional: \{\} <br /> |
| `crossId` _string_ | CrossID is designed to identified a page across environments.<br />If not set, this ID will be generated in a predictable manner based on<br />the map key associated to this entry in the API. |  | Optional: \{\} <br /> |
| `name` _string_ | This is the display name of the page in APIM and on the portal.<br />This field can be edited safely if you want to rename a page. |  | Required: \{\} <br /> |
| `type` _string_ | The type of the documentation page or folder. |  | Enum: [MARKDOWN SWAGGER ASYNCAPI ASCIIDOC FOLDER SYSTEM_FOLDER ROOT] <br />Required: \{\} <br /> |
| `content` _string_ | The content of the page, if any. |  | Optional: \{\} <br /> |
| `order` _integer_ | The order used to display the page in APIM and on the portal. |  | Optional: \{\} <br /> |
| `published` _boolean_ | If true, the page will be accessible from the portal (default is false) | false | Optional: \{\} <br /> |
| `visibility` _string_ | The visibility of the page. | PUBLIC | Enum: [PUBLIC PRIVATE] <br />Optional: \{\} <br /> |
| `homepage` _boolean_ | If true, this page will be displayed as the homepage of your API documentation. | false | Optional: \{\} <br /> |
| `parent` _string_ | If your page contains a folder, setting this field to the map key associated to the<br />folder entry will be reflected into APIM by making the page a child of this folder. |  | Optional: \{\} <br /> |
| `parentId` _string_ | The parent ID of the page. This field is mostly required when you are applying<br />an API exported from APIM to make the operator take control over it. Use `Parent`<br />in any other case. |  | Optional: \{\} <br /> |
| `api` _string_ | The API of the page. If empty, will be set automatically to the generated ID of the API. |  | Optional: \{\} <br /> |
| `source` _[PageSource](#pagesource)_ | Source allow you to fetch pages from various external sources, overriding page content<br />each time the source is fetched. |  | Optional: \{\} <br /> |
| `configuration` _map[string]string_ | Custom page configuration (e.g. page rendering can be changed to use Redoc instead of Swagger ui) |  | Optional: \{\} <br /> |


#### PageSource







_Appears in:_
- [Page](#page)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _string_ |  |  | Required: \{\} <br /> |
| `configuration` _[GenericStringMap](#genericstringmap)_ |  |  | Required: \{\} <br /> |


#### Plan







_Appears in:_
- [Plan](#plan)
- [Plan](#plan)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `id` _string_ | Plan ID |  |  |
| `crossId` _string_ | The plan Cross ID.<br />This field is used to identify plans defined for an API<br />that has been promoted between different environments. |  |  |
| `tags` _string array_ | List of plan tags | \{  \} | Optional: \{\} <br /> |
| `status` _[PlanStatus](#planstatus)_ | The plan status | PUBLISHED | Enum: [PUBLISHED DEPRECATED STAGING] <br />Optional: \{\} <br /> |
| `characteristics` _string array_ | List of plan characteristics |  | Optional: \{\} <br /> |
| `validation` _[PlanValidation](#planvalidation)_ | Plan validation strategy | AUTO | Enum: [AUTO MANUAL] <br /> |
| `comment_required` _boolean_ | Indicate of comment is required for this plan or not |  | Optional: \{\} <br /> |
| `order` _integer_ | Plan order |  | Optional: \{\} <br /> |
| `type` _[PlanType](#plantype)_ | Plan type | API | Enum: [API CATALOG] <br /> |


#### PlanStatus

_Underlying type:_ _string_

The status of the plan.

_Validation:_
- Enum: [PUBLISHED DEPRECATED STAGING]

_Appears in:_
- [Plan](#plan)



#### PlanType

_Underlying type:_ _string_



_Validation:_
- Enum: [API CATALOG]

_Appears in:_
- [Plan](#plan)



#### PlanValidation

_Underlying type:_ _string_



_Validation:_
- Enum: [AUTO MANUAL]

_Appears in:_
- [Plan](#plan)



#### Plugin







_Appears in:_
- [PluginRevision](#pluginrevision)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `policy` _string_ | Plugin Policy |  | Optional: \{\} <br /> |
| `resource` _string_ | Plugin Resource |  | Optional: \{\} <br /> |
| `configuration` _[GenericStringMap](#genericstringmap)_ | Plugin Configuration, a map of arbitrary key-values |  | Optional: \{\} <br /> |


#### PluginReference







_Appears in:_
- [PluginRevision](#pluginrevision)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `namespace` _string_ | Plugin Reference Namespace |  | Optional: \{\} <br /> |
| `resource` _string_ | Plugin Reference Resource |  | Optional: \{\} <br /> |
| `name` _string_ | Plugin Reference Name |  | Optional: \{\} <br /> |




#### Property







_Appears in:_
- [ApiBase](#apibase)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `key` _string_ | Property Key |  | Optional: \{\} <br /> |
| `value` _string_ | Property Value |  | Optional: \{\} <br /> |
| `encrypted` _boolean_ | Property Encrypted or not? |  | Optional: \{\} <br /> |
| `dynamic` _boolean_ | Property is dynamic or not? |  | Optional: \{\} <br /> |
| `encryptable` _boolean_ | Property is encryptable or not? |  | Optional: \{\} <br /> |


#### ProtocolVersion

_Underlying type:_ _string_



_Validation:_
- Enum: [HTTP_1_1 HTTP_2]

_Appears in:_
- [HttpClientOptions](#httpclientoptions)



#### Resource







_Appears in:_
- [ApiResourceSpec](#apiresourcespec)
- [ResourceOrRef](#resourceorref)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Is resource enabled or not? | true | Optional: \{\} <br /> |
| `name` _string_ | Resource Name |  | Optional: \{\} <br /> |
| `type` _string_ | Resource Type |  | Optional: \{\} <br /> |
| `configuration` _[GenericStringMap](#genericstringmap)_ | Resource Configuration, arbitrary map of key-values |  | Optional: \{\} <br /> |


#### ResourceOrRef







_Appears in:_
- [ApiBase](#apibase)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `Resource` _[Resource](#resource)_ | Resource |  |  |
| `ref` _[NamespacedName](#namespacedname)_ | Reference to a resource |  |  |


#### ResponseTemplate







_Appears in:_
- [Api](#api)
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)
- [ApiV4DefinitionSpec](#apiv4definitionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `status` _integer_ |  |  | Optional: \{\} <br /> |
| `headers` _map[string]string_ |  |  | Optional: \{\} <br /> |
| `body` _string_ |  |  | Optional: \{\} <br /> |
| `propagateErrorKeyToLogs` _boolean_ | Propagate error key to logs |  |  |




#### Status







_Appears in:_
- [ApiDefinitionStatus](#apidefinitionstatus)
- [ApiV4DefinitionStatus](#apiv4definitionstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `organizationId` _string_ | The organization ID, if a management context has been defined to sync with an APIM instance |  | Optional: \{\} <br /> |
| `environmentId` _string_ | The environment ID, if a management context has been defined to sync with an APIM instance |  | Optional: \{\} <br /> |
| `id` _string_ | The ID of the API definition in the Gravitee API Management instance (if an API context has been configured). |  | Optional: \{\} <br /> |
| `crossId` _string_ | The Cross ID is used to identify an API that has been promoted from one environment to another. |  |  |
| `processingStatus` _[ProcessingStatus](#processingstatus)_ | The processing status of the API definition. |  |  |
| `state` _[ApiState](#apistate)_ | The state of the API. Can be either STARTED or STOPPED. |  | Enum: [STARTED STOPPED] <br /> |
| `plans` _object (keys:string, values:string)_ | This field is used to store the list of plans that have been created<br />for the API definition if a management context has been defined<br />to sync the API with an APIM instance |  | Optional: \{\} <br /> |
| `subscriptions` _integer_ | The number of subscriptions that reference the API |  |  |
| `errors` _[Errors](#errors)_ | When API has been created regardless of errors, this field is<br />used to persist the error message encountered during admission |  |  |


#### TrustStore







_Appears in:_
- [HttpClientSslOptions](#httpclientssloptions)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[KeyStoreType](#keystoretype)_ | The TrustStore type to use (possible values are PEM, PKCS12, JKS) |  | Enum: [PEM PKCS12 JKS] <br /> |
| `path` _string_ |  |  | Optional: \{\} <br /> |
| `content` _string_ | The base64 encoded trustStore content, if not relying on a path to a file |  | Optional: \{\} <br /> |
| `password` _string_ | TrustStore password (Not applicable for PEM TrustStore) |  | Optional: \{\} <br /> |



## gravitee.io/v1alpha1/gateway




#### Deployment







_Appears in:_
- [KubernetesConfig](#kubernetesconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `labels` _object (keys:string, values:string)_ |  |  | Optional: \{\} <br /> |
| `annotations` _object (keys:string, values:string)_ |  |  | Optional: \{\} <br /> |
| `replicas` _integer_ |  | 1 | Optional: \{\} <br /> |
| `strategy` _[DeploymentStrategy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#deploymentstrategy-v1-apps)_ |  |  | Optional: \{\} <br /> |
| `template` _[PodTemplateSpec](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podtemplatespec-v1-core)_ | The template.spec field uses the standard Kubernetes Pod template specification,<br />and its contents will be merged using a<br />[strategic merge patch](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/)<br />with Gravitee's default deployment configuration. |  | Optional: \{\} <br /> |






#### GatewayClassParameters



The GatewayClassParameters custom resource is
the Gravitee.io extension point that allows you to configure
our implementation of the [Kubernetes Gateway API](https://gateway-api.sigs.k8s.io/).
It defines a set of configuration options to control how
Gravitee Gateways are deployed and behave when managed via the Gateway API,
including licensing, Kafka support, and Kubernetes-specific deployment settings.



_Appears in:_
- [GatewayClassParametersSpec](#gatewayclassparametersspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `gravitee` _[GraviteeConfig](#graviteeconfig)_ | The gravitee section controls Gravitee specific features<br />and allows you to configure and customize our implementation<br />of the Kubernetes Gateway API. |  | Optional: \{\} <br /> |
| `kubernetes` _[KubernetesConfig](#kubernetesconfig)_ | The kubernetes section of the GatewayClassParameters<br />spec lets you customize core Kubernetes resources<br />that are part of your Gateway deployments. |  | Optional: \{\} <br /> |


#### GraviteeConfig







_Appears in:_
- [GatewayClassParameters](#gatewayclassparameters)
- [GatewayClassParametersSpec](#gatewayclassparametersspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `licenseRef` _[SecretObjectReference](https://gateway-api.sigs.k8s.io/reference/spec/#secretobjectreference)_ | A reference to a Kubernetes Secret that contains your Gravitee license key.<br />This license is required to unlock advanced capabilities like Kafka protocol support. |  | Optional: \{\} <br /> |
| `kafka` _[GraviteeKafkaConfig](#graviteekafkaconfig)_ | Use this field to enable Kafka support in the Gateway. | \{ enabled:false \} | Optional: \{\} <br /> |
| `yaml` _[GenericStringMap](#genericstringmap)_ | Use this field to provide custom gateway configuration,<br />giving you control over additional configuration blocks<br />available in the gateway<br />[settings](https://documentation.gravitee.io/apim/configure-apim/apim-components/gravitee-gateway). |  | Optional: \{\} <br /> |


#### GraviteeKafkaConfig







_Appears in:_
- [GraviteeConfig](#graviteeconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ |  | true | Optional: \{\} <br /> |
| `routingHostMode` _[GraviteeKafkaRoutingHostModeConfig](#graviteekafkaroutinghostmodeconfig)_ |  | \{  \} | Optional: \{\} <br /> |


#### GraviteeKafkaRoutingHostModeConfig







_Appears in:_
- [GraviteeKafkaConfig](#graviteekafkaconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `brokerDomainPattern` _string_ |  | broker-\{brokerId\}-\{apiHost\} | Optional: \{\} <br /> |
| `bootstrapDomainPattern` _string_ | You can find details about these configurations options in our<br />[documentation](https://documentation.gravitee.io/apim/kafka-gateway/configure-the-kafka-gateway-and-client). | \{apiHost\} | Optional: \{\} <br /> |


#### KubernetesConfig







_Appears in:_
- [GatewayClassParameters](#gatewayclassparameters)
- [GatewayClassParametersSpec](#gatewayclassparametersspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `deployment` _[Deployment](#deployment)_ | Use this field to modify pod labels and annotations,<br />adjust the number of replicas to control scaling,<br />specify update strategies for rolling updates,<br />and override the pod template to customize container specs,<br />security settings, or environment variables. |  | Optional: \{\} <br /> |
| `service` _[Service](#service)_ | Use this field to customize the Kubernetes Service that exposes the Gateway<br />by adding labels and annotations, choosing the service type,<br />configuring the external traffic policy, and specifying the load balancer class.` |  | Optional: \{\} <br /> |








#### Service







_Appears in:_
- [KubernetesConfig](#kubernetesconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `labels` _object (keys:string, values:string)_ |  |  | Optional: \{\} <br /> |
| `annotations` _object (keys:string, values:string)_ |  |  | Optional: \{\} <br /> |
| `type` _[ServiceType](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#servicetype-v1-core)_ |  | LoadBalancer | Optional: \{\} <br /> |
| `externalTrafficPolicy` _[ServiceExternalTrafficPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#serviceexternaltrafficpolicy-v1-core)_ |  | Cluster | Optional: \{\} <br /> |
| `loadBalancerClass` _string_ |  |  | Optional: \{\} <br /> |



## gravitee.io/v1alpha1/group




#### Member







_Appears in:_
- [Type](#type)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `source` _string_ | Member source |  | Required: \{\} <br /> |
| `sourceId` _string_ | Member source ID |  | Required: \{\} <br /> |
| `roles` _object (keys:[RoleScope](#rolescope), values:string)_ |  | \{  \} | Optional: \{\} <br /> |


#### RoleScope

_Underlying type:_ _string_



_Validation:_
- Enum: [API APPLICATION INTEGRATION]

_Appears in:_
- [Member](#member)



#### Status







_Appears in:_
- [GroupStatus](#groupstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `id` _string_ | The ID of the Group in the Gravitee API Management instance |  | Optional: \{\} <br /> |
| `organizationId` _string_ | The organization ID defined in the management context |  | Optional: \{\} <br /> |
| `environmentId` _string_ | The environment ID defined in the management context |  | Optional: \{\} <br /> |
| `processingStatus` _[ProcessingStatus](#processingstatus)_ | The processing status of the Group. |  |  |
| `members` _integer_ | The number of members added to this group |  |  |
| `errors` _[Errors](#errors)_ | When group has been created regardless of errors, this field is<br />used to persist the error message encountered during admission |  |  |


#### Type







_Appears in:_
- [GroupSpec](#groupspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `id` _string_ |  |  | Optional: \{\} <br /> |
| `name` _string_ |  |  | Required: \{\} <br /> |
| `notifyMembers` _boolean_ | If true, new members added to the API spec will<br />be notified when the API is synced with APIM. | true | Optional: \{\} <br /> |
| `members` _[Member](#member) array_ |  |  |  |



## gravitee.io/v1alpha1/kafka





#### KafKaRoute







_Appears in:_
- [KafkaRouteSpec](#kafkaroutespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `parentRefs` _[ParentReference](#parentreference) array_ | ParentRefs references the resources (usually Gateways) that a Route wants<br />to be attached to. Note that the referenced parent resource needs to<br />allow this for the attachment to be complete. For Gateways, that means<br />the Gateway needs to allow attachment from Routes of this kind and<br />namespace. For Services, that means the Service must either be in the same<br />namespace for a "producer" route, or the mesh implementation must support<br />and allow "consumer" routes for the referenced Service. ReferenceGrant is<br />not applicable for governing ParentRefs to Services - it is not possible to<br />create a "producer" route for a Service in a different namespace from the<br />Route.<br />There are two kinds of parent resources with "Core" support:<br />* Gateway (Gateway conformance profile)<br />* Service (Mesh conformance profile, ClusterIP Services only)<br />This API may be extended in the future to support additional kinds of parent<br />resources.<br />ParentRefs must be _distinct_. This means either that:<br />* They select different objects.  If this is the case, then parentRef<br />  entries are distinct. In terms of fields, this means that the<br />  multi-part key defined by `group`, `kind`, `namespace`, and `name` must<br />  be unique across all parentRef entries in the Route.<br />* They do not select different objects, but for each optional field used,<br />  each ParentRef that selects the same object must set the same set of<br />  optional fields to different values. If one ParentRef sets a<br />  combination of optional fields, all must set the same combination.<br />Some examples:<br />* If one ParentRef sets `sectionName`, all ParentRefs referencing the<br />  same object must also set `sectionName`.<br />* If one ParentRef sets `port`, all ParentRefs referencing the same<br />  object must also set `port`.<br />* If one ParentRef sets `sectionName` and `port`, all ParentRefs<br />  referencing the same object must also set `sectionName` and `port`.<br />It is possible to separately reference multiple distinct objects that may<br />be collapsed by an implementation. For example, some implementations may<br />choose to merge compatible Gateway Listeners together. If that is the<br />case, the list of routes attached to those resources should also be<br />merged.<br />Note that for ParentRefs that cross namespace boundaries, there are specific<br />rules. Cross-namespace references are only valid if they are explicitly<br />allowed by something in the namespace they are referring to. For example,<br />Gateway has the AllowedRoutes field, and ReferenceGrant provides a<br />generic way to enable other kinds of cross-namespace reference.<br /> |
| `hostname` _[Hostname](#hostname)_ | Hostname is used to uniquely route clients to this API.<br />Your client must trust the certificate provided by the gateway,<br />and as there is a variable host in the proxy bootstrap server URL,<br />you likely need to request a wildcard SAN for the certificate presented by the gateway.<br />If empty, the hostname defined in the Kafka listener of the parent will be used. |  | MaxLength: 253 <br />MinLength: 1 <br />Pattern: `^(\*\.)?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$` <br /> |
| `backendRefs` _[KafkaBackendRef](#kafkabackendref) array_ | BackendRefs defines the backend(s) where matching requests should be sent. |  | MaxItems: 16 <br />MinItems: 1 <br /> |
| `filters` _[KafkaRouteFilter](#kafkaroutefilter) array_ | Filters define the filters that are applied to Kafka trafic matching this route. |  | MaxItems: 16 <br /> |
| `options` _object (keys:[AnnotationKey](#annotationkey), values:[AnnotationValue](#annotationvalue))_ | Options are a list of key/value pairs to enable extended configuration specific<br />to an |  | MaxProperties: 16 <br /> |


#### KafkaACLFilter







_Appears in:_
- [KafkaRouteFilter](#kafkaroutefilter)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `rules` _[KafkaAccessControlRules](#kafkaaccesscontrolrules) array_ | Rules define a set of rules that can be use to group a set of resources together with<br />access control rules to be applied.<br />ACLs are restrictive because once they are applied,<br />proxy clients must be authorized to perform the actions they are taking.<br />If there is no ACL defined for the action taken by the user, the action is prohibited. |  | MaxItems: 16 <br />MinItems: 1 <br /> |


#### KafkaAcccessControlResourceType

_Underlying type:_ _string_



_Validation:_
- Enum: [Topic Cluster Group TransactionalIdentifier]

_Appears in:_
- [KafkaAccessControl](#kafkaaccesscontrol)



#### KafkaAccessControl







_Appears in:_
- [KafkaAccessControlRules](#kafkaaccesscontrolrules)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[KafkaAcccessControlResourceType](#kafkaacccesscontrolresourcetype)_ |  |  | Enum: [Topic Cluster Group TransactionalIdentifier] <br /> |
| `operations` _[KafkaAccessControlOperation](#kafkaaccesscontroloperation) array_ | Operations specifies the set of operations / verbs to allow for the resource<br />under access control. |  | Enum: [Create Read Write Delete Alter AlterConfigs Describe DescribeConfigs ClusterAction] <br /> |
| `match` _[KafkaAccessControlMatch](#kafkaaccesscontrolmatch)_ | Match describes how to select the resource that will be subject to the access control.<br />If not specified, any resource will be matched. |  |  |


#### KafkaAccessControlMatch







_Appears in:_
- [KafkaAccessControl](#kafkaaccesscontrol)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[KafkaAccessControlResourceMatchType](#kafkaaccesscontrolresourcematchtype)_ | Valid PathMatchType values, along with their support levels, are:<br />* "Exact" Resources whose name is an exact match to the specified string receive the ACL.<br />* "Prefix" Resources whose name starts with the specified string receive the ACL.<br />* "RegularExpression" Resources that match the specified expression receive the ACL. |  | Enum: [Exact Prefix RegularExpression] <br /> |
| `value` _string_ | Value of the resource to match against. |  |  |


#### KafkaAccessControlOperation

_Underlying type:_ _string_



_Validation:_
- Enum: [Create Read Write Delete Alter AlterConfigs Describe DescribeConfigs ClusterAction]

_Appears in:_
- [KafkaAccessControl](#kafkaaccesscontrol)



#### KafkaAccessControlResourceMatchType

_Underlying type:_ _string_



_Validation:_
- Enum: [Exact Prefix RegularExpression]

_Appears in:_
- [KafkaAccessControlMatch](#kafkaaccesscontrolmatch)



#### KafkaAccessControlRules







_Appears in:_
- [KafkaACLFilter](#kafkaaclfilter)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `resources` _[KafkaAccessControl](#kafkaaccesscontrol) array_ | A resource group together a type of matched resource and a set of operations<br />to be granted by the access control for that resource. |  | MinItems: 1 <br /> |
| `options` _object (keys:[AnnotationKey](#annotationkey), values:[AnnotationValue](#annotationvalue))_ | Options allow to specify implementation specific behaviours<br />for a set of rules. |  | MaxProperties: 16 <br /> |


#### KafkaBackendRef



This currently wraps the code gateway API BackendObjectReference type,
leaving room for e.g. backend security configuration.



_Appears in:_
- [KafKaRoute](#kafkaroute)
- [KafkaRouteSpec](#kafkaroutespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `group` _[Group](#group)_ | Group is the group of the referent. For example, "gateway.networking.k8s.io".<br />When unspecified or empty string, core API group is inferred. |  | MaxLength: 253 <br />Pattern: `^$\|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$` <br /> |
| `kind` _[Kind](#kind)_ | Kind is the Kubernetes resource kind of the referent. For example<br />"Service".<br />Defaults to "Service" when not specified.<br />ExternalName services can refer to CNAME DNS records that may live<br />outside of the cluster and as such are difficult to reason about in<br />terms of conformance. They also may not be safe to forward to (see<br />CVE-2021-25740 for more information). Implementations SHOULD NOT<br />support ExternalName Services.<br />Support: Core (Services with a type other than ExternalName)<br />Support: Implementation-specific (Services with type ExternalName) | Service | MaxLength: 63 <br />MinLength: 1 <br />Pattern: `^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$` <br /> |
| `name` _[ObjectName](#objectname)_ | Name is the name of the referent. |  | MaxLength: 253 <br />MinLength: 1 <br /> |
| `namespace` _[Namespace](#namespace)_ | Namespace is the namespace of the backend. When unspecified, the local<br />namespace is inferred.<br />Note that when a namespace different than the local namespace is specified,<br />a ReferenceGrant object is required in the referent namespace to allow that<br />namespace's owner to accept the reference. See the ReferenceGrant<br />documentation for details.<br />Support: Core |  | MaxLength: 63 <br />MinLength: 1 <br />Pattern: `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$` <br /> |
| `port` _[PortNumber](#portnumber)_ | Port specifies the destination port number to use for this resource.<br />Port is required when the referent is a Kubernetes Service. In this<br />case, the port number is the service port number, not the target port.<br />For other resources, destination port might be derived from the referent<br />resource or this field. |  | Maximum: 65535 <br />Minimum: 1 <br /> |


#### KafkaRouteFilter







_Appears in:_
- [KafKaRoute](#kafkaroute)
- [KafkaRouteSpec](#kafkaroutespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[KafkaRouteFilterType](#kafkaroutefiltertype)_ | Type identifies the type of filter to apply. |  | Enum: [ACL ExtensionRef] <br /> |
| `acl` _[KafkaACLFilter](#kafkaaclfilter)_ | ACL defines a schema for a filter that enforce access controls on Kafka trafic. |  |  |
| `extensionRef` _[LocalObjectReference](#localobjectreference)_ |  |  |  |


#### KafkaRouteFilterType

_Underlying type:_ _string_





_Appears in:_
- [KafkaRouteFilter](#kafkaroutefilter)




## gravitee.io/v1alpha1/management




#### Auth







_Appears in:_
- [Context](#context)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `bearerToken` _string_ | The bearer token used to authenticate against the API Management instance<br />(must be generated from an admin account) |  | Optional: \{\} <br /> |
| `credentials` _[BasicAuth](#basicauth)_ | The Basic credentials used to authenticate against the API Management instance. |  |  |
| `secretRef` _[NamespacedName](#namespacedname)_ | A secret reference holding either a "bearerToken" key for bearer token authentication<br />or "username" and "password" keys for basic authentication |  |  |


#### BasicAuth







_Appears in:_
- [Auth](#auth)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `username` _string_ |  |  | Required: \{\} <br /> |
| `password` _string_ |  |  | Required: \{\} <br /> |


#### Cloud







_Appears in:_
- [Context](#context)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `token` _string_ | Token plain text Gravitee cloud token (JWT) |  | Optional: \{\} <br /> |
| `secretRef` _[NamespacedName](#namespacedname)_ | SecretRef secret reference holding the Gravitee cloud token in the "cloudToken" key |  | Optional: \{\} <br /> |


#### Context







_Appears in:_
- [ManagementContextSpec](#managementcontextspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `baseUrl` _string_ | The URL of a management API instance.<br />This is optional when this context targets Gravitee Cloud otherwise it is required. |  | Optional: \{\} <br /> |
| `path` _string_ | Allows to override the context path that will be appended to the baseURL.<br />This can be used when reverse proxying APIM with URL rewrite |  | Optional: \{\} <br /> |
| `organizationId` _string_ | An existing organization id targeted by the context on the management API instance.<br />This is optional when this context targets Gravitee Cloud otherwise it is required. |  | Optional: \{\} <br /> |
| `environmentId` _string_ | An existing environment id targeted by the context within the organization.<br />This is optional when this context targets Gravitee Cloud<br />and your cloud token contains only one environment ID, otherwise it is required. |  | Optional: \{\} <br /> |
| `auth` _[Auth](#auth)_ | Auth defines the authentication method used to connect to the API Management.<br />Can be either basic authentication credentials, a bearer token<br />or a reference to a kubernetes secret holding one of these two configurations.<br />This is optional when this context targets Gravitee Cloud. |  |  |
| `cloud` _[Cloud](#cloud)_ | Cloud when set (token or secretRef) this context will target Gravitee Cloud.<br />BaseUrl will be defaulted from token data if not set,<br />Auth is defaulted to use the token (bearerToken),<br />OrgID is extracted from the token,<br />EnvID is defaulted when the token contains exactly one environment. |  |  |



## gravitee.io/v1alpha1/notification




#### ApiEvent

_Underlying type:_ _string_

ApiEvent defines the events that can be sent to the console.

_Validation:_
- Enum: [APIKEY_EXPIRED APIKEY_RENEWED APIKEY_REVOKED SUBSCRIPTION_NEW SUBSCRIPTION_ACCEPTED SUBSCRIPTION_CLOSED SUBSCRIPTION_PAUSED SUBSCRIPTION_RESUMED SUBSCRIPTION_REJECTED SUBSCRIPTION_TRANSFERRED SUBSCRIPTION_FAILED NEW_SUPPORT_TICKET API_STARTED API_STOPPED API_UPDATED API_DEPLOYED NEW_RATING NEW_RATING_ANSWER MESSAGE ASK_FOR_REVIEW REVIEW_OK REQUEST_FOR_CHANGES API_DEPRECATED NEW_SPEC_GENERATED]

_Appears in:_
- [Console](#console)



#### Console







_Appears in:_
- [Type](#type)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `groupRefs` _[NamespacedName](#namespacedname) array_ | List of group references associated with this console notification.<br />These groups are references to gravitee.io/Group custom resources created on the cluster.<br />All members of those groups will receive a notification for the defined events. |  | Optional: \{\} <br /> |
| `apiEvents` _[ApiEvent](#apievent) array_ | List events that will trigger a notification for an API. Recipients are the API primary owner<br />and all members of groups referenced in groupRefs<br />Notification spec attribute eventType must be set to "api". |  | Enum: [APIKEY_EXPIRED APIKEY_RENEWED APIKEY_REVOKED SUBSCRIPTION_NEW SUBSCRIPTION_ACCEPTED SUBSCRIPTION_CLOSED SUBSCRIPTION_PAUSED SUBSCRIPTION_RESUMED SUBSCRIPTION_REJECTED SUBSCRIPTION_TRANSFERRED SUBSCRIPTION_FAILED NEW_SUPPORT_TICKET API_STARTED API_STOPPED API_UPDATED API_DEPLOYED NEW_RATING NEW_RATING_ANSWER MESSAGE ASK_FOR_REVIEW REVIEW_OK REQUEST_FOR_CHANGES API_DEPRECATED NEW_SPEC_GENERATED] <br />Optional: \{\} <br /> |
| `groups` _string array_ | List of groups associated with the API.<br />These groups are id to existing groups in APIM. |  | Optional: \{\} <br /> |


#### EventType

_Underlying type:_ _string_

EventType defines the subject of those events.

_Validation:_
- Enum: [api]

_Appears in:_
- [Type](#type)

| Field | Description |
| --- | --- |
| `api` |  |


#### Target

_Underlying type:_ _string_

Target defines the target of the notification.

_Validation:_
- Enum: [console]

_Appears in:_
- [Type](#type)

| Field | Description |
| --- | --- |
| `console` |  |


#### Type







_Appears in:_
- [NotificationSpec](#notificationspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `target` _[Target](#target)_ | Target of the notification: "console" is for notifications in Gravitee console UI.<br />For each target there is an attribute of the same name to configure it. | console | Enum: [console] <br />Required: \{\} <br /> |
| `eventType` _[EventType](#eventtype)_ | EventType defines the subject of those events.<br />Notification can be used in API or Applications, each of those have different events.<br />An attribute starting with eventType value exists in the target configuration<br />to configure events: < eventType >Events (e.g apiEvents) | api | Enum: [api] <br />Required: \{\} <br /> |
| `console` _[Console](#console)_ | Console is used when the target value is "console" and is meant<br />to configure Gravitee console UI notifications. | \{  \} | Optional: \{\} <br /> |



## gravitee.io/v1alpha1/policygroups




#### ApiType

_Underlying type:_ _string_



_Validation:_
- Enum: [MESSAGE PROXY NATIVE]

_Appears in:_
- [SharedPolicyGroup](#sharedpolicygroup)



#### FlowPhase

_Underlying type:_ _string_



_Validation:_
- Enum: [REQUEST RESPONSE INTERACT CONNECT PUBLISH SUBSCRIBE]

_Appears in:_
- [SharedPolicyGroup](#sharedpolicygroup)



#### SharedPolicyGroup







_Appears in:_
- [SharedPolicyGroupSpec](#sharedpolicygroupspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `crossId` _string_ | CrossID to export SharedPolicyGroup into different environments |  |  |
| `name` _string_ | SharedPolicyGroup name |  | Required: \{\} <br /> |
| `description` _string_ | SharedPolicyGroup description |  |  |
| `prerequisiteMessage` _string_ | SharedPolicyGroup prerequisite Message |  |  |
| `apiType` _[ApiType](#apitype)_ | Specify the SharedPolicyGroup ApiType |  | Enum: [MESSAGE PROXY NATIVE] <br />Required: \{\} <br /> |
| `phase` _[FlowPhase](#flowphase)_ | SharedPolicyGroup phase (REQUEST;RESPONSE;INTERACT;CONNECT;PUBLISH;SUBSCRIBE) |  | Enum: [REQUEST RESPONSE INTERACT CONNECT PUBLISH SUBSCRIBE] <br />Required: \{\} <br /> |
| `steps` _[Step](#step) array_ | SharedPolicyGroup Steps |  |  |


#### Status







_Appears in:_
- [SharedPolicyGroupSpecStatus](#sharedpolicygroupspecstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `organizationId` _string_ | The organization ID, if a management context has been defined to sync with an APIM instance |  |  |
| `environmentId` _string_ | The environment ID, if a management context has been defined to sync with an APIM instance |  |  |
| `crossId` _string_ | The Cross ID is used to identify an SharedPolicyGroup that has been promoted from one environment to another. |  |  |
| `id` _string_ | The ID is used to identify an SharedPolicyGroup which is unique in any environment. |  |  |
| `processingStatus` _[ProcessingStatus](#processingstatus)_ | The processing status of the SharedPolicyGroup.<br />The value is `Completed` if the sync with APIM succeeded, Failed otherwise. |  |  |
| `errors` _[Errors](#errors)_ | When SharedPolicyGroup has been created regardless of errors, this field is<br />used to persist the error message encountered during admission |  |  |


#### Step







_Appears in:_
- [SharedPolicyGroup](#sharedpolicygroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Indicate if this FlowStep is enabled or not | true |  |
| `policy` _string_ | FlowStep policy |  | Optional: \{\} <br /> |
| `name` _string_ | FlowStep name |  | Optional: \{\} <br /> |
| `description` _string_ | FlowStep description |  | Optional: \{\} <br /> |
| `configuration` _[GenericStringMap](#genericstringmap)_ | FlowStep configuration is a map of arbitrary key-values |  | Optional: \{\} <br /> |
| `condition` _string_ | FlowStep condition |  | Optional: \{\} <br /> |



## gravitee.io/v1alpha1/refs




#### NamespacedName







_Appears in:_
- [ApiBase](#apibase)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)
- [ApiRef](#apiref)
- [ApiV4DefinitionSpec](#apiv4definitionspec)
- [ApplicationSpec](#applicationspec)
- [Auth](#auth)
- [Cloud](#cloud)
- [Console](#console)
- [FlowStep](#flowstep)
- [GroupSpec](#groupspec)
- [ResourceOrRef](#resourceorref)
- [SharedPolicyGroupSpec](#sharedpolicygroupspec)
- [SubscriptionSpec](#subscriptionspec)
- [Type](#type)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ |  |  |  |
| `namespace` _string_ |  |  | Optional: \{\} <br /> |
| `kind` _string_ |  |  | Optional: \{\} <br /> |



## gravitee.io/v1alpha1/status




#### Errors







_Appears in:_
- [ApiDefinitionStatus](#apidefinitionstatus)
- [ApiV4DefinitionStatus](#apiv4definitionstatus)
- [ApplicationStatus](#applicationstatus)
- [GroupStatus](#groupstatus)
- [SharedPolicyGroupSpecStatus](#sharedpolicygroupspecstatus)
- [Status](#status)
- [Status](#status)
- [Status](#status)
- [Status](#status)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `warning` _string array_ | warning errors do not block object reconciliation,<br />most of the time because the value is ignored or defaulted<br />when the API gets synced with APIM |  |  |
| `severe` _string array_ | severe errors do not pass admission and will block reconcile<br />hence, this field should always be during the admission phase<br />and is very unlikely to be persisted in the status |  |  |



## gravitee.io/v1alpha1/subscription








#### Type







_Appears in:_
- [SubscriptionSpec](#subscriptionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `api` _[NamespacedName](#namespacedname)_ |  |  | Required: \{\} <br /> |
| `application` _[NamespacedName](#namespacedname)_ |  |  | Required: \{\} <br /> |
| `plan` _string_ |  |  | Required: \{\} <br /> |
| `endingAt` _string_ |  |  | Format: date-time <br />Optional: \{\} <br /> |



## gravitee.io/v1alpha1/utils




#### GenericStringMap







_Appears in:_
- [DynamicPropertyService](#dynamicpropertyservice)
- [Endpoint](#endpoint)
- [EndpointDiscoveryService](#endpointdiscoveryservice)
- [EndpointGroup](#endpointgroup)
- [Entrypoint](#entrypoint)
- [FlowSelector](#flowselector)
- [FlowStep](#flowstep)
- [FlowStep](#flowstep)
- [GenericListener](#genericlistener)
- [GraviteeConfig](#graviteeconfig)
- [PageSource](#pagesource)
- [PlanSecurity](#plansecurity)
- [Plugin](#plugin)
- [Policy](#policy)
- [Resource](#resource)
- [Service](#service)
- [Step](#step)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `Object` _object (keys:string, values:[interface{}](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#interface{}-unstructured-v1))_ | Object is a JSON compatible map with string, float, int, bool, []interface\{\}, or<br />map[string]interface\{\}<br />children. |  |  |



## gravitee.io/v1alpha1/v2






#### Api







_Appears in:_
- [ApiDefinitionV2Spec](#apidefinitionv2spec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `description` _string_ | API description |  |  |
| `definition_context` _[DefinitionContext](#definitioncontext)_ | The definition context is used to inform a management API instance that this API definition<br />is managed using a kubernetes operator |  | Optional: \{\} <br /> |
| `lifecycle_state` _[LifecycleState](#lifecyclestate)_ | API life cycle state can be one of the values CREATED, PUBLISHED, UNPUBLISHED, DEPRECATED, ARCHIVED | CREATED | Enum: [CREATED PUBLISHED UNPUBLISHED DEPRECATED ARCHIVED] <br /> |
| `deployedAt` _integer_ | Shows the time that the API is deployed |  | Optional: \{\} <br /> |
| `gravitee` _[DefinitionVersion](#definitionversion)_ | The definition version of the API. For v1alpha1 resources, this field should always set to `2.0.0`. | 2.0.0 |  |
| `flow_mode` _[FlowMode](#flowmode)_ | The flow mode of the API. The value is either `DEFAULT` or `BEST_MATCH`. | DEFAULT | Enum: [DEFAULT BEST_MATCH] <br /> |
| `proxy` _[Proxy](#proxy)_ | The proxy of the API that specifies its VirtualHosts and Groups. |  |  |
| `services` _[Services](#services)_ | Contains different services for the API (EndpointDiscovery, HealthCheck ...) |  |  |
| `flows` _[Flow](#flow) array_ | The flow of the API | \{  \} | Optional: \{\} <br /> |
| `path_mappings` _string array_ | API Path mapping |  | Optional: \{\} <br /> |
| `plans` _[Plan](#plan) array_ | API plans | \{  \} | Optional: \{\} <br /> |
| `response_templates` _[ResponseTemplate](#responsetemplate)_ | A list of Response Templates for the API |  | Optional: \{\} <br /> |
| `members` _Member array_ | List of members associated with the API |  | Optional: \{\} <br /> |
| `pages` _[map[string]*Page](#map[string]*page)_ | A map of pages objects.<br />Keys uniquely identify pages and are used to keep them in sync<br />with APIM when using a management context.<br />Renaming a key is the equivalent of deleting the page and recreating<br />it holding a new ID in APIM. |  | Optional: \{\} <br /> |
| `execution_mode` _string_ | Execution mode that eventually runs the API in the gateway | v4-emulation-engine | Enum: [v3 v4-emulation-engine] <br /> |


#### Consumer







_Appears in:_
- [Flow](#flow)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `consumerType` _[ConsumerType](#consumertype)_ | Consumer type (possible values TAG) |  |  |
| `consumerId` _string_ | Consumer ID |  | Optional: \{\} <br /> |


#### ConsumerType

_Underlying type:_ _integer_





_Appears in:_
- [Consumer](#consumer)



#### DefinitionContext







_Appears in:_
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `origin` _string_ |  | kubernetes |  |
| `mode` _string_ |  | fully_managed |  |
| `syncFrom` _string_ |  | kubernetes |  |


#### DynamicPropertyProvider

_Underlying type:_ _string_



_Validation:_
- Enum: [HTTP]

_Appears in:_
- [DynamicPropertyService](#dynamicpropertyservice)

| Field | Description |
| --- | --- |
| `HTTP` |  |


#### DynamicPropertyService







_Appears in:_
- [Services](#services)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `provider` _[DynamicPropertyProvider](#dynamicpropertyprovider)_ |  |  | Enum: [HTTP] <br /> |
| `configuration` _[GenericStringMap](#genericstringmap)_ | Configuration, arbitrary map of key-values |  | Optional: \{\} <br /> |


#### Endpoint







_Appears in:_
- [EndpointGroup](#endpointgroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name of the endpoint |  | Optional: \{\} <br /> |
| `target` _string_ | The end target of this endpoint (backend) |  | Optional: \{\} <br /> |
| `weight` _integer_ | Endpoint weight used for load-balancing |  | Optional: \{\} <br /> |
| `backup` _boolean_ | Indicate that this ia a back-end endpoint |  | Optional: \{\} <br /> |
| `tenants` _string array_ | The endpoint tenants | \{  \} | Optional: \{\} <br /> |
| `type` _[EndpointType](#endpointtype)_ | The type of endpoint (HttpEndpointType or GrpcEndpointType) |  |  |
| `inherit` _boolean_ | Is endpoint inherited or not |  | Optional: \{\} <br /> |
| `proxy` _[HttpProxy](#httpproxy)_ | Configure the HTTP Proxy settings to reach target if needed |  |  |
| `http` _[HttpClientOptions](#httpclientoptions)_ | Custom HTTP client options used for this endpoint |  |  |
| `ssl` _[HttpClientSslOptions](#httpclientssloptions)_ | Custom HTTP SSL client options used for this endpoint |  |  |
| `headers` _[HttpHeader](#httpheader) array_ | List of headers for this endpoint | \{  \} | Optional: \{\} <br /> |
| `healthcheck` _[EndpointHealthCheckService](#endpointhealthcheckservice)_ | Specify EndpointHealthCheck service settings |  |  |


#### EndpointDiscoveryService







_Appears in:_
- [Services](#services)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `provider` _string_ | Provider name |  | Optional: \{\} <br /> |
| `configuration` _[GenericStringMap](#genericstringmap)_ | Configuration, arbitrary map of key-values |  | Optional: \{\} <br /> |
| `secondary` _boolean_ | Is it secondary or not? |  | Optional: \{\} <br /> |
| `tenants` _string array_ | List of tenants | \{  \} | Optional: \{\} <br /> |


#### EndpointGroup







_Appears in:_
- [Proxy](#proxy)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | EndpointGroup name |  | Optional: \{\} <br /> |
| `endpoints` _[Endpoint](#endpoint) array_ | List of Endpoints belonging to this group |  | Optional: \{\} <br /> |
| `load_balancing` _[LoadBalancer](#loadbalancer)_ | The LoadBalancer Type |  |  |
| `services` _[Services](#services)_ | Specify different Endpoint Services |  |  |
| `proxy` _[HttpProxy](#httpproxy)_ | Configure the HTTP Proxy settings for this EndpointGroup if needed |  |  |
| `http` _[HttpClientOptions](#httpclientoptions)_ | Custom HTTP SSL client options used for this EndpointGroup |  |  |
| `ssl` _[HttpClientSslOptions](#httpclientssloptions)_ | Custom HTTP SSL client options used for this EndpointGroup |  |  |
| `headers` _map[string]string_ | List of headers needed for this EndpointGroup |  | Optional: \{\} <br /> |


#### EndpointHealthCheckService







_Appears in:_
- [Endpoint](#endpoint)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `inherit` _boolean_ |  | false | Optional: \{\} <br /> |


#### EndpointType

_Underlying type:_ _string_





_Appears in:_
- [Endpoint](#endpoint)

| Field | Description |
| --- | --- |
| `http` |  |
| `grpc` |  |


#### Failover







_Appears in:_
- [Proxy](#proxy)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `maxAttempts` _integer_ | Maximum number of attempts |  | Optional: \{\} <br /> |
| `retryTimeout` _integer_ | Retry timeout |  | Optional: \{\} <br /> |
| `cases` _[FailoverCase](#failovercase) array_ | List of Failover cases |  | Optional: \{\} <br /> |


#### FailoverCase

_Underlying type:_ _string_





_Appears in:_
- [Failover](#failover)



#### Flow







_Appears in:_
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)
- [Plan](#plan)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `id` _string_ | Flow ID |  |  |
| `name` _string_ | Flow name |  | Optional: \{\} <br /> |
| `path-operator` _[PathOperator](#pathoperator)_ | List of path operators |  |  |
| `pre` _[FlowStep](#flowstep) array_ | Flow pre step | \{  \} | Optional: \{\} <br /> |
| `post` _[FlowStep](#flowstep) array_ | Flow post step | \{  \} | Optional: \{\} <br /> |
| `enabled` _boolean_ | Indicate if this flow is enabled or disabled | true |  |
| `methods` _[HttpMethod](#httpmethod) array_ | A list of methods  for this flow (GET;POST;PUT;PATCH;DELETE;OPTIONS;HEAD;CONNECT;TRACE;OTHER) | \{  \} | Optional: \{\} <br /> |
| `condition` _string_ | Flow condition |  | Optional: \{\} <br /> |
| `consumers` _[Consumer](#consumer) array_ | List of the consumers of this Flow | \{  \} | Optional: \{\} <br /> |


#### FlowMode

_Underlying type:_ _string_



_Validation:_
- Enum: [DEFAULT BEST_MATCH]

_Appears in:_
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)



#### HealthCheckRequest







_Appears in:_
- [HealthCheckStep](#healthcheckstep)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `path` _string_ | The path of the endpoint handling the health check request |  | Optional: \{\} <br /> |
| `method` _[HttpMethod](#httpmethod)_ | The HTTP method to use when issuing the health check request |  | Enum: [GET POST PUT PATCH DELETE OPTIONS HEAD CONNECT TRACE OTHER] <br /> |
| `headers` _[HttpHeader](#httpheader) array_ | List of HTTP headers to include in the health check request | \{  \} | Optional: \{\} <br /> |
| `body` _string_ | Health Check Request Body |  | Optional: \{\} <br /> |
| `fromRoot` _boolean_ | If true, the health check request will be issued without prepending the context path of the API. |  |  |


#### HealthCheckResponse







_Appears in:_
- [HealthCheckStep](#healthcheckstep)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `assertions` _string array_ |  |  | Optional: \{\} <br /> |


#### HealthCheckService







_Appears in:_
- [EndpointHealthCheckService](#endpointhealthcheckservice)
- [Services](#services)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `steps` _[HealthCheckStep](#healthcheckstep) array_ | List of health check steps | \{  \} | Optional: \{\} <br /> |


#### HealthCheckStep







_Appears in:_
- [HealthCheckService](#healthcheckservice)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Health Check Step Name |  | Optional: \{\} <br /> |
| `request` _[HealthCheckRequest](#healthcheckrequest)_ | Health Check Step Request |  |  |
| `response` _[HealthCheckResponse](#healthcheckresponse)_ | Health Check Step Response |  |  |


#### LoadBalancer







_Appears in:_
- [EndpointGroup](#endpointgroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[LoadBalancerType](#loadbalancertype)_ | Type of the LoadBalancer (RoundRobin, Random, WeightedRoundRobin, WeightedRandom) |  |  |


#### LoadBalancerType

_Underlying type:_ _string_





_Appears in:_
- [LoadBalancer](#loadbalancer)

| Field | Description |
| --- | --- |
| `ROUND_ROBIN` |  |
| `RANDOM` |  |
| `WEIGHTED_ROUND_ROBIN` |  |
| `WEIGHTED_RANDOM` |  |


#### Logging







_Appears in:_
- [Analytics](#analytics)
- [Proxy](#proxy)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `mode` _[LoggingMode](#loggingmode)_ | The logging mode.<br />CLIENT identifies the inbound request issued to the gateway,<br />while PROXY identifies the request issued to the upstream service. |  | Enum: [NONE CLIENT PROXY CLIENT_PROXY] <br /> |
| `scope` _[LoggingScope](#loggingscope)_ | The logging scope (which phase of the request roundtrip should be included in each log entry. |  | Enum: [NONE REQUEST RESPONSE REQUEST_RESPONSE] <br /> |
| `content` _[LoggingContent](#loggingcontent)_ | Which part of the request/response should be logged ? |  | Enum: [NONE HEADERS PAYLOADS HEADERS_PAYLOADS] <br /> |
| `condition` _string_ | The logging condition (supports EL expressions) |  | Optional: \{\} <br /> |


#### LoggingContent

_Underlying type:_ _string_



_Validation:_
- Enum: [NONE HEADERS PAYLOADS HEADERS_PAYLOADS]

_Appears in:_
- [Logging](#logging)



#### LoggingMode

_Underlying type:_ _string_



_Validation:_
- Enum: [NONE CLIENT PROXY CLIENT_PROXY]

_Appears in:_
- [Logging](#logging)



#### LoggingScope

_Underlying type:_ _string_



_Validation:_
- Enum: [NONE REQUEST RESPONSE REQUEST_RESPONSE]

_Appears in:_
- [Logging](#logging)



#### Page







_Appears in:_
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `accessControls` _[AccessControl](#accesscontrol) array_ | If the page is private, defines a set of user groups with access | \{  \} | Optional: \{\} <br /> |
| `excludedAccessControls` _boolean_ | if true, the references defined in the accessControls list will be<br />denied access instead of being granted |  | Optional: \{\} <br /> |




#### PathOperator







_Appears in:_
- [Flow](#flow)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `path` _string_ | Operator path |  | Optional: \{\} <br /> |
| `operator` _[Operator](#operator)_ | Operator (possible values STARTS_WITH or EQUALS) | STARTS_WITH | Enum: [STARTS_WITH EQUALS] <br /> |


#### Plan







_Appears in:_
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Plan name |  |  |
| `description` _string_ | Plan Description |  |  |
| `security` _string_ | Plan Security |  |  |
| `securityDefinition` _string_ | Plan Security definition |  | Optional: \{\} <br /> |
| `paths` _[map[string][]Rule](#map[string][]rule)_ | A map of different paths (alongside their Rules) for this Plan |  | Optional: \{\} <br /> |
| `api` _string_ | Specify the API associated with this plan |  | Optional: \{\} <br /> |
| `selection_rule` _string_ | Plan selection rule |  | Optional: \{\} <br /> |
| `flows` _[Flow](#flow) array_ | List of different flows for this Plan | \{  \} | Optional: \{\} <br /> |
| `excluded_groups` _string array_ | List of excluded groups for this plan | \{  \} | Optional: \{\} <br /> |


#### Policy







_Appears in:_
- [Rule](#rule)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Policy name |  | Optional: \{\} <br /> |
| `configuration` _[GenericStringMap](#genericstringmap)_ | Policy configuration is a map of arbitrary key-values |  | Optional: \{\} <br /> |


#### Proxy







_Appears in:_
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `virtual_hosts` _[VirtualHost](#virtualhost) array_ | list of Virtual hosts fot the proxy |  | Optional: \{\} <br /> |
| `groups` _[EndpointGroup](#endpointgroup) array_ | List of endpoint groups of the proxy | \{  \} | Optional: \{\} <br /> |
| `failover` _[Failover](#failover)_ | Proxy Failover |  |  |
| `cors` _[Cors](#cors)_ | Proxy Cors |  |  |
| `logging` _[Logging](#logging)_ | Logging |  |  |
| `strip_context_path` _boolean_ | Strip Context Path |  | Optional: \{\} <br /> |
| `preserve_host` _boolean_ | Preserve Host |  | Optional: \{\} <br /> |


#### Rule







_Appears in:_
- [Path](#path)
- [Plan](#plan)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `methods` _[HttpMethod](#httpmethod) array_ | List of http methods for this Rule (GET;POST;PUT;PATCH;DELETE;OPTIONS;HEAD;CONNECT;TRACE;OTHER) |  | Optional: \{\} <br /> |
| `policy` _[Policy](#policy)_ | Rule policy |  |  |
| `description` _string_ | Rule description |  | Optional: \{\} <br /> |
| `enabled` _boolean_ | Indicate if the Rule is enabled or not |  | Optional: \{\} <br /> |


#### Sampling







_Appears in:_
- [Analytics](#analytics)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[SamplingType](#samplingtype)_ | The sampling type to use |  |  |
| `value` _string_ | Sampling Value |  |  |


#### SamplingType

_Underlying type:_ _string_





_Appears in:_
- [Sampling](#sampling)



#### ScheduledService







_Appears in:_
- [DynamicPropertyService](#dynamicpropertyservice)
- [HealthCheckService](#healthcheckservice)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `schedule` _string_ |  |  | Optional: \{\} <br /> |


#### Service







_Appears in:_
- [EndpointDiscoveryService](#endpointdiscoveryservice)
- [ScheduledService](#scheduledservice)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Service name |  | Optional: \{\} <br /> |
| `enabled` _boolean_ | Is service enabled or not? | false | Optional: \{\} <br /> |


#### Services







_Appears in:_
- [Api](#api)
- [ApiDefinitionV2Spec](#apidefinitionv2spec)
- [EndpointGroup](#endpointgroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `discovery` _[EndpointDiscoveryService](#endpointdiscoveryservice)_ | Endpoint Discovery Service |  |  |
| `health-check` _[HealthCheckService](#healthcheckservice)_ | Health Check Service |  |  |
| `dynamic-property` _[DynamicPropertyService](#dynamicpropertyservice)_ | Dynamic Property Service |  |  |


#### VirtualHost







_Appears in:_
- [Proxy](#proxy)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `host` _string_ | Host name |  | Optional: \{\} <br /> |
| `path` _string_ | Path |  | Optional: \{\} <br /> |
| `override_entrypoint` _boolean_ | Indicate if Entrypoint should be overridden or not |  | Optional: \{\} <br /> |



## gravitee.io/v1alpha1/v4




#### AbstractListener







_Appears in:_
- [HttpListener](#httplistener)
- [KafkaListener](#kafkalistener)
- [SubscriptionListener](#subscriptionlistener)
- [TCPListener](#tcplistener)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[ListenerType](#listenertype)_ |  | HTTP | Enum: [HTTP SUBSCRIPTION TCP KAFKA] <br />Required: \{\} <br /> |
| `entrypoints` _[Entrypoint](#entrypoint) array_ |  |  | Required: \{\} <br /> |
| `servers` _string array_ |  |  | Optional: \{\} <br /> |


#### Analytics







_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Analytics Enabled or not? | true |  |
| `sampling` _[Sampling](#sampling)_ | Analytics Sampling |  |  |
| `logging` _[Logging](#logging)_ | Analytics Logging |  |  |
| `tracing` _[Tracing](#tracing)_ | Analytics Tracing |  |  |


#### Api







_Appears in:_
- [ApiV4DefinitionSpec](#apiv4definitionspec)
- [GatewayDefinitionApi](#gatewaydefinitionapi)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `description` _string_ | API description |  | Optional: \{\} <br /> |
| `definitionVersion` _[DefinitionVersion](#definitionversion)_ | The definition version of the API. | V4 | Enum: [V4] <br /> |
| `definitionContext` _[DefinitionContext](#definitioncontext)_ | The API Definition context is used to identify the Kubernetes origin of the API,<br />and define whether the API definition should be synchronized<br />from an API instance or from a config map created in the cluster (which is the default) |  |  |
| `lifecycleState` _[ApiV4LifecycleState](#apiv4lifecyclestate)_ | API life cycle state can be one of the values PUBLISHED, UNPUBLISHED | UNPUBLISHED | Enum: [PUBLISHED UNPUBLISHED] <br />Optional: \{\} <br /> |
| `type` _[ApiType](#apitype)_ | Api Type (proxy or message) |  | Enum: [PROXY MESSAGE NATIVE] <br />Required: \{\} <br /> |
| `listeners` _[GenericListener](#genericlistener) array_ | List of listeners for this API |  | MinItems: 1 <br />Required: \{\} <br /> |
| `endpointGroups` _[EndpointGroup](#endpointgroup) array_ | List of Endpoint groups |  | MinItems: 1 <br />Required: \{\} <br /> |
| `plans` _[map[string]*Plan](#map[string]*plan)_ | A map of plan identifiers to plan<br />Keys uniquely identify plans and are used to keep them in sync<br />when using a management context. |  | Optional: \{\} <br /> |
| `flowExecution` _[FlowExecution](#flowexecution)_ | API Flow Execution (Not applicable for Native API) |  |  |
| `flows` _[Flow](#flow) array_ | List of flows for the API | \{  \} | Optional: \{\} <br /> |
| `analytics` _[Analytics](#analytics)_ | API Analytics (Not applicable for Native API) |  |  |
| `services` _[ApiServices](#apiservices)_ | API Services (Not applicable for Native API) |  |  |
| `responseTemplates` _[ResponseTemplate](#responsetemplate)_ | A list of Response Templates for the API (Not applicable for Native API) |  | Optional: \{\} <br /> |
| `members` _Member array_ | List of members associated with the API |  | Optional: \{\} <br /> |
| `pages` _[map[string]*Page](#map[string]*page)_ | A map of pages objects.<br />Keys uniquely identify pages and are used to keep them in sync<br />with APIM when using a management context.<br />Renaming a key is the equivalent of deleting the page and recreating<br />it holding a new ID in APIM. |  | Optional: \{\} <br /> |
| `failover` _[Failover](#failover)_ | API Failover |  |  |


#### ApiServices







_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `dynamicProperty` _[Service](#service)_ | API dynamic property service |  |  |


#### ApiType

_Underlying type:_ _string_



_Validation:_
- Enum: [PROXY MESSAGE NATIVE]

_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)



#### ApiV4LifecycleState

_Underlying type:_ _string_



_Validation:_
- Enum: [PUBLISHED UNPUBLISHED]

_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)





#### DLQ







_Appears in:_
- [Entrypoint](#entrypoint)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `endpoint` _string_ | The endpoint to use when a message should be sent to the dead letter queue. |  | Optional: \{\} <br /> |


#### DefinitionContext







_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `origin` _[DefinitionContextOrigin](#definitioncontextorigin)_ | The definition context origin where the API definition is managed.<br />The value is always `KUBERNETES` for an API managed by the operator. | KUBERNETES | Enum: [KUBERNETES] <br /> |
| `syncFrom` _[DefinitionContextOrigin](#definitioncontextorigin)_ | The syncFrom field defines where the gateways should source the API definition from.<br />If the value is `MANAGEMENT`, the API definition will be sourced from an APIM instance.<br />This means that the API definition *must* hold a context reference in that case.<br />Setting the value to `MANAGEMENT` allows to make an API definition available on<br />gateways deployed across multiple clusters / regions.<br />If the value is `KUBERNETES`, the API definition will be sourced from a config map.<br />This means that only gateways deployed in the same cluster will be able to sync the API definition. | MANAGEMENT | Enum: [KUBERNETES MANAGEMENT] <br /> |




#### DefinitionContextOrigin

_Underlying type:_ _string_





_Appears in:_
- [DefinitionContext](#definitioncontext)

| Field | Description |
| --- | --- |
| `FULLY_MANAGED` |  |
| `KUBERNETES` |  |
| `MANAGEMENT` |  |


#### DefinitionVersion

_Underlying type:_ _string_





_Appears in:_
- [Plan](#plan)

| Field | Description |
| --- | --- |
| `V4` |  |


#### Endpoint







_Appears in:_
- [EndpointGroup](#endpointgroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | The endpoint name (this value should be unique across endpoints) |  | Optional: \{\} <br /> |
| `type` _string_ | Endpoint Type |  | Required: \{\} <br /> |
| `weight` _integer_ | Endpoint Weight |  | Optional: \{\} <br /> |
| `inheritConfiguration` _boolean_ | Should endpoint group configuration be inherited or not ? |  |  |
| `configuration` _[GenericStringMap](#genericstringmap)_ | Endpoint Configuration, arbitrary map of key-values |  | Optional: \{\} <br /> |
| `sharedConfigurationOverride` _[GenericStringMap](#genericstringmap)_ | Endpoint Configuration Override, arbitrary map of key-values |  | Optional: \{\} <br /> |
| `services` _[EndpointServices](#endpointservices)_ | Endpoint Services |  |  |
| `secondary` _boolean_ | Endpoint is secondary or not? |  |  |
| `tenants` _string array_ | List of endpoint tenants | \{  \} | Optional: \{\} <br /> |


#### EndpointGroup







_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Endpoint group name |  | Required: \{\} <br /> |
| `type` _string_ | Endpoint group type |  | Optional: \{\} <br /> |
| `loadBalancer` _[LoadBalancer](#loadbalancer)_ | Endpoint group load balancer |  |  |
| `sharedConfiguration` _[GenericStringMap](#genericstringmap)_ | Endpoint group shared configuration, arbitrary map of key-values |  | Optional: \{\} <br /> |
| `endpoints` _[Endpoint](#endpoint) array_ | List of endpoint for the group | \{  \} | Optional: \{\} <br /> |
| `services` _[EndpointGroupServices](#endpointgroupservices)_ | Endpoint group services |  |  |
| `http` _[HttpClientOptions](#httpclientoptions)_ | Endpoint group http client options |  |  |
| `ssl` _[HttpClientSslOptions](#httpclientssloptions)_ | Endpoint group http client SSL options |  |  |
| `headers` _map[string]string_ | Endpoint group headers, arbitrary map of key-values |  | Optional: \{\} <br /> |


#### EndpointGroupServices







_Appears in:_
- [EndpointGroup](#endpointgroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `discovery` _[Service](#service)_ | Endpoint group discovery service |  |  |
| `healthCheck` _[Service](#service)_ | Endpoint group health check service |  |  |


#### EndpointServices







_Appears in:_
- [Endpoint](#endpoint)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `healthCheck` _[Service](#service)_ | Health check service |  |  |






#### Entrypoint







_Appears in:_
- [AbstractListener](#abstractlistener)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _string_ |  |  | Required: \{\} <br /> |
| `qos` _[QosType](#qostype)_ |  | AUTO | Enum: [NONE AUTO AT_MOST_ONCE AT_LEAST_ONCE] <br />Optional: \{\} <br /> |
| `dlq` _[DLQ](#dlq)_ |  |  |  |
| `configuration` _[GenericStringMap](#genericstringmap)_ |  |  | Optional: \{\} <br /> |




#### Failover







_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | API Failover is enabled? | false |  |
| `maxRetries` _integer_ | API Failover max retires | 2 |  |
| `slowCallDuration` _integer_ | API Failover slow call duration | 2000 |  |
| `openStateDuration` _integer_ | API Failover  open state duration | 10000 |  |
| `maxFailures` _integer_ | API Failover max failures | 5 |  |
| `perSubscription` _boolean_ | API Failover  per subscription | true |  |


#### Flow







_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)
- [Plan](#plan)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `id` _string_ | The ID of the flow this field is mainly used for compatibility with<br />APIM exports and can be safely ignored. |  | Optional: \{\} <br /> |
| `name` _string_ | Flow name |  | Optional: \{\} <br /> |
| `enabled` _boolean_ | Is flow enabled or not? | true |  |
| `selectors` _[FlowSelector](#flowselector) array_ | List of Flow selectors |  | Optional: \{\} <br /> |
| `request` _[FlowStep](#flowstep) array_ | List of Request flow steps (NOT available for Native APIs) |  | Optional: \{\} <br /> |
| `response` _[FlowStep](#flowstep) array_ | List of Response flow steps (NOT available for Native APIs) |  | Optional: \{\} <br /> |
| `subscribe` _[FlowStep](#flowstep) array_ | List of Subscribe flow steps |  | Optional: \{\} <br /> |
| `publish` _[FlowStep](#flowstep) array_ | List of Publish flow steps |  | Optional: \{\} <br /> |
| `connect` _[FlowStep](#flowstep) array_ | List of Connect flow steps (Only available for Native APIs) |  | Optional: \{\} <br /> |
| `interact` _[FlowStep](#flowstep) array_ | List of Publish flow steps (Only available for Native APIs) |  | Optional: \{\} <br /> |
| `tags` _string array_ | List of tags |  | Optional: \{\} <br /> |


#### FlowExecution







_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `mode` _[FlowMode](#flowmode)_ | The flow mode to use |  |  |
| `matchRequired` _boolean_ | Is match required or not ? If set to true, a 404 status response will be returned if no matching flow was found. |  |  |


#### FlowMode

_Underlying type:_ _string_





_Appears in:_
- [FlowExecution](#flowexecution)



#### FlowSelector







_Appears in:_
- [Flow](#flow)



#### FlowStep







_Appears in:_
- [Flow](#flow)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Indicate if this FlowStep is enabled or not | true |  |
| `policy` _string_ | FlowStep policy |  | Optional: \{\} <br /> |
| `name` _string_ | FlowStep name |  | Optional: \{\} <br /> |
| `description` _string_ | FlowStep description |  | Optional: \{\} <br /> |
| `configuration` _[GenericStringMap](#genericstringmap)_ | FlowStep configuration is a map of arbitrary key-values |  | Optional: \{\} <br /> |
| `condition` _string_ | FlowStep condition |  | Optional: \{\} <br /> |
| `sharedPolicyGroupRef` _[NamespacedName](#namespacedname)_ | Reference to an existing Shared Policy Group |  | Optional: \{\} <br /> |
| `messageCondition` _string_ | The message condition (supports EL expressions) |  | Optional: \{\} <br /> |




#### GatewayDefinitionPlan







_Appears in:_
- [GatewayDefinitionApi](#gatewaydefinitionapi)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ |  |  |  |


#### GenericListener







_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)









#### ListenerType

_Underlying type:_ _string_



_Validation:_
- Enum: [HTTP SUBSCRIPTION TCP KAFKA]

_Appears in:_
- [AbstractListener](#abstractlistener)

| Field | Description |
| --- | --- |
| `HTTP` |  |
| `SUBSCRIPTION` |  |
| `TCP` |  |
| `KAFKA` |  |


#### LoadBalancer







_Appears in:_
- [EndpointGroup](#endpointgroup)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[LoadBalancerType](#loadbalancertype)_ |  | ROUND_ROBIN | Enum: [ROUND_ROBIN RANDOM WEIGHTED_ROUND_ROBIN WEIGHTED_RANDOM] <br /> |


#### LoadBalancerType

_Underlying type:_ _string_



_Validation:_
- Enum: [ROUND_ROBIN RANDOM WEIGHTED_ROUND_ROBIN WEIGHTED_RANDOM]

_Appears in:_
- [LoadBalancer](#loadbalancer)

| Field | Description |
| --- | --- |
| `ROUND_ROBIN` |  |
| `RANDOM` |  |
| `WEIGHTED_ROUND_ROBIN` |  |
| `WEIGHTED_RANDOM` |  |


#### Logging







_Appears in:_
- [Analytics](#analytics)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `condition` _string_ | The logging condition. This field is evaluated for HTTP requests and supports EL expressions. |  | Optional: \{\} <br /> |
| `messageCondition` _string_ | The logging message condition. This field is evaluated for messages and supports EL expressions. |  | Optional: \{\} <br /> |
| `content` _[LoggingContent](#loggingcontent)_ | Defines which component of the request should be included in the log payload. |  |  |
| `mode` _[LoggingMode](#loggingmode)_ | The logging mode defines which "hop" of the request roundtrip<br />should be included in the log payload.<br />This can be either the inbound request to the gateway,<br />the request issued by the gateway to the upstream service, or both. |  |  |
| `phase` _[LoggingPhase](#loggingphase)_ | Defines which phase of the request roundtrip<br />should be included in the log payload.<br />This can be either the request phase, the response phase, or both. |  |  |


#### LoggingContent







_Appears in:_
- [Logging](#logging)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `headers` _boolean_ | Should HTTP headers be logged or not ? |  |  |
| `messageHeaders` _boolean_ | Should message headers be logged or not ? |  |  |
| `payload` _boolean_ | Should HTTP payloads be logged or not ? |  |  |
| `messagePayload` _boolean_ | Should message payloads be logged or not ? |  |  |
| `messageMetadata` _boolean_ | Should message metadata be logged or not ? |  |  |


#### LoggingMode







_Appears in:_
- [Logging](#logging)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `entrypoint` _boolean_ | If true, the inbound request to the gateway will be included in the log payload |  |  |
| `endpoint` _boolean_ | If true, the request to the upstream service will be included in the log payload |  |  |


#### LoggingPhase







_Appears in:_
- [Logging](#logging)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `request` _boolean_ | Should the request phase of the request roundtrip be included in the log payload or not ? |  |  |
| `response` _boolean_ | Should the response phase of the request roundtrip be included in the log payload or not ? |  |  |


#### Page







_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)



#### Path







_Appears in:_
- [HttpListener](#httplistener)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `host` _string_ |  |  | Optional: \{\} <br /> |
| `path` _string_ |  |  | Required: \{\} <br /> |


#### Plan







_Appears in:_
- [Api](#api)
- [ApiV4DefinitionSpec](#apiv4definitionspec)
- [GatewayDefinitionPlan](#gatewaydefinitionplan)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Plan display name, this will be the name displayed in the UI<br />if a management context is used to sync the API with APIM |  |  |
| `description` _string_ | Plan Description |  | Optional: \{\} <br /> |
| `definitionVersion` _[DefinitionVersion](#definitionversion)_ | Plan definition version | V4 |  |
| `security` _[PlanSecurity](#plansecurity)_ | Plan security |  |  |
| `mode` _[PlanMode](#planmode)_ | The plan mode | STANDARD | Enum: [STANDARD PUSH] <br />Optional: \{\} <br /> |
| `selectionRule` _string_ | Plan selection rule |  | Optional: \{\} <br /> |
| `flows` _[Flow](#flow) array_ | List of plan flows | \{  \} | Optional: \{\} <br /> |
| `excludedGroups` _string array_ |  | \{  \} | Optional: \{\} <br /> |
| `generalConditions` _string_ | The general conditions defined to use this plan |  | Optional: \{\} <br /> |


#### PlanMode

_Underlying type:_ _string_



_Validation:_
- Enum: [STANDARD PUSH]

_Appears in:_
- [Plan](#plan)



#### PlanSecurity







_Appears in:_
- [Plan](#plan)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _string_ | Plan Security type |  | Required: \{\} <br /> |
| `configuration` _[GenericStringMap](#genericstringmap)_ | Plan security configuration, a map of arbitrary key-values |  | Optional: \{\} <br /> |


#### QosType

_Underlying type:_ _string_



_Validation:_
- Enum: [NONE AUTO AT_MOST_ONCE AT_LEAST_ONCE]

_Appears in:_
- [Entrypoint](#entrypoint)



#### Sampling







_Appears in:_
- [Analytics](#analytics)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[SamplingType](#samplingtype)_ | The sampling type to use |  |  |
| `value` _string_ | Sampling Value |  |  |


#### SamplingType

_Underlying type:_ _string_





_Appears in:_
- [Sampling](#sampling)



#### Service







_Appears in:_
- [ApiServices](#apiservices)
- [EndpointGroupServices](#endpointgroupservices)
- [EndpointServices](#endpointservices)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Is the service enabled or not ? |  |  |
| `type` _string_ | Service Type |  | Optional: \{\} <br /> |
| `overrideConfiguration` _boolean_ | Service Override Configuration or not? |  |  |
| `configuration` _[GenericStringMap](#genericstringmap)_ | Service Configuration, a map of arbitrary key-values |  | Optional: \{\} <br /> |






#### Tracing







_Appears in:_
- [Analytics](#analytics)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Specify if Tracing is Enabled or not |  |  |
| `verbose` _boolean_ | Specify if Tracing is Verbose or not |  |  |


>>>>>>> 56d297b (docs: generate docs using crd-ref-docs)
