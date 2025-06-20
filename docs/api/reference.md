# Gravitee Kubernetes Operator API Reference

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
        </tr>
        <tr>
            <td><a href="#gatewayclassparameters">GatewayClassParameters</a></td>
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
          Auth defines the authentication method used to connect to the API Management.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>baseUrl</b></td>
        <td>string</td>
        <td>
          The URL of a management API instance.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#managementcontextspeccloud">cloud</a></b></td>
        <td>object</td>
        <td>
          Cloud when set (token or secretRef) this context will target Gravitee Cloud.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>environmentId</b></td>
        <td>string</td>
        <td>
          An existing environment id targeted by the context within the organization.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>organizationId</b></td>
        <td>string</td>
        <td>
          An existing organization id targeted by the context on the management API instance.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Allows to override the context path that will be appended to the baseURL.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ManagementContext.spec.auth
[Go to parent definition](#managementcontextspec)



Auth defines the authentication method used to connect to the API Management.

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
(must be generated from...<br/>
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
or "username"...<br/>
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
        <td>true</td>
      </tr><tr>
        <td><b>username</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ManagementContext.spec.auth.secretRef
[Go to parent definition](#managementcontextspecauth)



A secret reference holding either a "bearerToken" key for bearer token authentication
or "username"...

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






ApiDefinition is the Schema for the apidefinitions API.

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
        <td><b><a href="#apidefinitionspec">spec</a></b></td>
        <td>object</td>
        <td>
          The API definition is the main resource handled by the Kubernetes Operator
Most of the...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionstatus">status</a></b></td>
        <td>object</td>
        <td>
          ApiDefinitionStatus defines the observed state of API Definition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiDefinition.spec
[Go to parent definition](#apidefinition)



The API definition is the main resource handled by the Kubernetes Operator
Most of the...

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
          API description<br/>
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
        <td><b>version</b></td>
        <td>string</td>
        <td>
          API version<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>categories</b></td>
        <td>[]string</td>
        <td>
          The list of categories the API belongs to.<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecconsolenotificationconfiguration">consoleNotificationConfiguration</a></b></td>
        <td>object</td>
        <td>
          ConsoleNotification struct sent to the mAPI, not part of the CRD spec.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspeccontextref">contextRef</a></b></td>
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
this ID identifies the API across those...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecdefinition_context">definition_context</a></b></td>
        <td>object</td>
        <td>
          The definition context is used to inform a management API instance that this API definition
is...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deployedAt</b></td>
        <td>integer</td>
        <td>
          Shows the time that the API is deployed<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>execution_mode</b></td>
        <td>enum</td>
        <td>
          Execution mode that eventually runs the API in the gateway<br/>
          <br/>
            <i>Enum</i>: v3, v4-emulation-engine<br/>
            <i>Default</i>: v4-emulation-engine<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>flow_mode</b></td>
        <td>enum</td>
        <td>
          The flow mode of the API. The value is either `DEFAULT` or `BEST_MATCH`.<br/>
          <br/>
            <i>Enum</i>: DEFAULT, BEST_MATCH<br/>
            <i>Default</i>: DEFAULT<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecflowsindex">flows</a></b></td>
        <td>[]object</td>
        <td>
          The flow of the API<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gravitee</b></td>
        <td>string</td>
        <td>
          The definition version of the API. For v1alpha1 resources, this field should always set to `2.0.0`.<br/>
          <br/>
            <i>Default</i>: 2.0.0<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecgrouprefsindex">groupRefs</a></b></td>
        <td>[]object</td>
        <td>
          List of group references associated with the API
These groups are references to Group custom...<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>groups</b></td>
        <td>[]string</td>
        <td>
          List of groups associated with the API.<br/>
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
        <td><b>lifecycle_state</b></td>
        <td>enum</td>
        <td>
          API life cycle state can be one of the values CREATED, PUBLISHED, UNPUBLISHED, DEPRECATED, ARCHIVED<br/>
          <br/>
            <i>Enum</i>: CREATED, PUBLISHED, UNPUBLISHED, DEPRECATED, ARCHIVED<br/>
            <i>Default</i>: CREATED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>local</b></td>
        <td>boolean</td>
        <td>
          local defines if the api is local or not.<br/>
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
          References to Notification custom resources to setup notifications.<br/>
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
          A map of pages objects.<br/>
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
of an <a...<br/>
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


### ApiDefinition.spec.consoleNotificationConfiguration
[Go to parent definition](#apidefinitionspec)



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


### ApiDefinition.spec.contextRef
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


### ApiDefinition.spec.definition_context
[Go to parent definition](#apidefinitionspec)



The definition context is used to inform a management API instance that this API definition
is...

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
        <td><b>name</b></td>
        <td>string</td>
        <td>
          This is the display name of the page in APIM and on the portal.<br/>
        </td>
        <td>true</td>
      </tr><tr>
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
          CrossID is designed to identified a page across environments.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>excludedAccessControls</b></td>
        <td>boolean</td>
        <td>
          if true, the references defined in the accessControls list will be
denied access instead of being...<br/>
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
          The ID of the page.<br/>
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
folder entry...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>parentId</b></td>
        <td>string</td>
        <td>
          The parent ID of the page.<br/>
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
each time...<br/>
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
        <td>true</td>
      </tr><tr>
        <td><b>referenceType</b></td>
        <td>enum</td>
        <td>
          The type of reference denied or granted by the access control
Currently only GROUP is supported<br/>
          <br/>
            <i>Enum</i>: GROUP<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### ApiDefinition.spec.pages[key].source
[Go to parent definition](#apidefinitionspecpageskey)



Source allow you to fetch pages from various external sources, overriding page content
each time...

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
          The plan Cross ID.<br/>
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
(Only applicable for PEM...<br/>
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
(Only applicable for PEM...<br/>
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
          The logging mode.<br/>
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
used to persist the error message...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the API definition in the Gravitee API Management instance (if an API context has been...<br/>
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
for the API definition if a...<br/>
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
used to persist the error message...

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
hence, this field should always be...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>warning</b></td>
        <td>[]string</td>
        <td>
          warning errors do not block object reconciliation,
most of the time because the value is ignored or...<br/>
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
          The list of categories the API belongs to.<br/>
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
this ID identifies the API across those...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecdefinitioncontext">definitionContext</a></b></td>
        <td>object</td>
        <td>
          The API Definition context is used to identify the Kubernetes origin of the API,
and define whether...<br/>
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
These groups are references to Group custom...<br/>
          <br/>
            <i>Default</i>: []<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>groups</b></td>
        <td>[]string</td>
        <td>
          List of groups associated with the API.<br/>
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
          References to Notification custom resources to setup notifications.<br/>
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
          A map of pages objects.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskey">plans</a></b></td>
        <td>map[string]object</td>
        <td>
          A map of plan identifiers to plan
Keys uniquely identify plans and are used to keep them in sync...<br/>
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
of an <a...<br/>
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
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Endpoint Type<br/>
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
        <td><b>weight</b></td>
        <td>integer</td>
        <td>
          Endpoint Weight<br/>
          <br/>
            <i>Format</i>: int32<br/>
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
(Only applicable for PEM...<br/>
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
should be included in the log payload.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecanalyticsloggingphase">phase</a></b></td>
        <td>object</td>
        <td>
          Defines which phase of the request roundtrip
should be included in the log payload.<br/>
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
and define whether...

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
          The definition context origin where the API definition is managed.<br/>
          <br/>
            <i>Enum</i>: KUBERNETES<br/>
            <i>Default</i>: KUBERNETES<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>syncFrom</b></td>
        <td>enum</td>
        <td>
          The syncFrom field defines where the gateways should source the API definition from.<br/>
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
          Is match required or not ?<br/>
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
APIM exports and can be safely...<br/>
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
        <td><b>name</b></td>
        <td>string</td>
        <td>
          This is the display name of the page in APIM and on the portal.<br/>
        </td>
        <td>true</td>
      </tr><tr>
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
          CrossID is designed to identified a page across environments.<br/>
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
          The ID of the page.<br/>
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
folder entry...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>parentId</b></td>
        <td>string</td>
        <td>
          The parent ID of the page.<br/>
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
each time...<br/>
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
each time...

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
if a management context is used to...<br/>
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
          The plan Cross ID.<br/>
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
APIM exports and can be safely...<br/>
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
used to persist the error message...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the API definition in the Gravitee API Management instance (if an API context has been...<br/>
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
for the API definition if a...<br/>
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
used to persist the error message...

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
hence, this field should always be...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>warning</b></td>
        <td>[]string</td>
        <td>
          warning errors do not block object reconciliation,
most of the time because the value is ignored or...<br/>
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



ApiResourceSpec defines the desired state of ApiResource.

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
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Resource Type<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## Application

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
        <td><b><a href="#applicationspec">spec</a></b></td>
        <td>object</td>
        <td>
          Application is the main resource handled by the Kubernetes Operator<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#applicationstatus">status</a></b></td>
        <td>object</td>
        <td>
          ApplicationStatus defines the observed state of Application.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Application.spec
[Go to parent definition](#application)



Application is the main resource handled by the Kubernetes Operator

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
        <td><b><a href="#applicationspeccontextref">contextRef</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          Application Description<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Application name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#applicationspecsettings">settings</a></b></td>
        <td>object</td>
        <td>
          Application settings<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>background</b></td>
        <td>string</td>
        <td>
          The base64 encoded background to use for this application when displaying it on the portal<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>domain</b></td>
        <td>string</td>
        <td>
          Application domain<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>groups</b></td>
        <td>[]string</td>
        <td>
          Application groups<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          io.gravitee.definition.model.Application
Application ID<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#applicationspecmembersindex">members</a></b></td>
        <td>[]object</td>
        <td>
          Application members<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#applicationspecmetadataindex">metadata</a></b></td>
        <td>[]object</td>
        <td>
          Application metadata<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>notifyMembers</b></td>
        <td>boolean</td>
        <td>
          Notify members when they are added to the application<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>picture</b></td>
        <td>string</td>
        <td>
          The base64 encoded picture to use for this application when displaying it on the portal (if not...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>pictureUrl</b></td>
        <td>string</td>
        <td>
          A URL pointing to the picture to use when displaying the application on the portal<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Application.spec.contextRef
[Go to parent definition](#applicationspec)





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


### Application.spec.settings
[Go to parent definition](#applicationspec)



Application settings

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
        <td><b><a href="#applicationspecsettingsapp">app</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#applicationspecsettingsoauth">oauth</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#applicationspecsettingstls">tls</a></b></td>
        <td>object</td>
        <td>
          TLS settings are used to configure client side TLS in order
to be able to subscribe to a MTLS plan.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Application.spec.settings.app
[Go to parent definition](#applicationspecsettings)





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
          Application Type<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>clientId</b></td>
        <td>string</td>
        <td>
          ClientID is the client id of the application<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Application.spec.settings.oauth
[Go to parent definition](#applicationspecsettings)





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
        <td><b>applicationType</b></td>
        <td>enum</td>
        <td>
          Oauth client application type<br/>
          <br/>
            <i>Enum</i>: BACKEND_TO_BACKEND, NATIVE, BROWSER, WEB<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>grantTypes</b></td>
        <td>[]enum</td>
        <td>
          List of Oauth client grant types<br/>
          <br/>
            <i>Enum</i>: authorization_code, client_credentials, refresh_token, password, implicit<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>redirectUris</b></td>
        <td>[]string</td>
        <td>
          List of Oauth client redirect uris<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Application.spec.settings.tls
[Go to parent definition](#applicationspecsettings)



TLS settings are used to configure client side TLS in order
to be able to subscribe to a MTLS plan.

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
        <td><b>clientCertificate</b></td>
        <td>string</td>
        <td>
          This client certificate is mandatory to subscribe to a TLS plan.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Application.spec.members[index]
[Go to parent definition](#applicationspec)





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


### Application.spec.metadata[index]
[Go to parent definition](#applicationspec)





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
          Metadata Name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>defaultValue</b></td>
        <td>string</td>
        <td>
          Metadata DefaultValue<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>format</b></td>
        <td>enum</td>
        <td>
          Metadata Format<br/>
          <br/>
            <i>Enum</i>: STRING, NUMERIC, BOOLEAN, DATE, MAIL, URL<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hidden</b></td>
        <td>boolean</td>
        <td>
          Metadata is hidden or not?<br/>
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


### Application.status
[Go to parent definition](#application)



ApplicationStatus defines the observed state of Application.

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
        <td><b>environmentId</b></td>
        <td>string</td>
        <td>
          The environment ID, if a management context has been defined to sync with an APIM instance<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#applicationstatuserrors">errors</a></b></td>
        <td>object</td>
        <td>
          When application has been created regardless of errors, this field is
used to persist the error...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the Application, if a management context has been defined to sync with an APIM instance<br/>
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
        <td><b>processingStatus</b></td>
        <td>string</td>
        <td>
          The processing status of the Application.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subscriptions</b></td>
        <td>integer</td>
        <td>
          The number of subscriptions that reference the application<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Application.status.errors
[Go to parent definition](#applicationstatus)



When application has been created regardless of errors, this field is
used to persist the error...

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
hence, this field should always be...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>warning</b></td>
        <td>[]string</td>
        <td>
          warning errors do not block object reconciliation,
most of the time because the value is ignored or...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## Subscription

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
        <td><b><a href="#subscriptionspec">spec</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#subscriptionstatus">status</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Subscription.spec
[Go to parent definition](#subscription)





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
        <td><b><a href="#subscriptionspecapi">api</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#subscriptionspecapplication">application</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>plan</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>endingAt</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Subscription.spec.api
[Go to parent definition](#subscriptionspec)





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


### Subscription.spec.application
[Go to parent definition](#subscriptionspec)





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


### Subscription.status
[Go to parent definition](#subscription)





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
        <td><b>endingAt</b></td>
        <td>string</td>
        <td>
          The expiry date for the subscription (no date means no expiry)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          Subscription ID<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>processingStatus</b></td>
        <td>string</td>
        <td>
          This value is `Completed` if the sync with APIM succeeded, Failed otherwise.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>startedAt</b></td>
        <td>string</td>
        <td>
          When the subscription was started and made available<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## SharedPolicyGroup

[gravitee.io/v1alpha1](#graviteeiov1alpha1)






SharedPolicyGroup

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
        <td><b><a href="#sharedpolicygroupspec">spec</a></b></td>
        <td>object</td>
        <td>
          SharedPolicyGroupSpec<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#sharedpolicygroupstatus">status</a></b></td>
        <td>object</td>
        <td>
          SharedPolicyGroupSpecStatus defines the observed state of an API Context.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SharedPolicyGroup.spec
[Go to parent definition](#sharedpolicygroup)



SharedPolicyGroupSpec

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
        <td><b>apiType</b></td>
        <td>enum</td>
        <td>
          Specify the SharedPolicyGroup ApiType<br/>
          <br/>
            <i>Enum</i>: MESSAGE, PROXY, NATIVE<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#sharedpolicygroupspeccontextref">contextRef</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          SharedPolicyGroup name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>phase</b></td>
        <td>enum</td>
        <td>
          SharedPolicyGroup phase (REQUEST;RESPONSE;INTERACT;CONNECT;PUBLISH;SUBSCRIBE)<br/>
          <br/>
            <i>Enum</i>: REQUEST, RESPONSE, INTERACT, CONNECT, PUBLISH, SUBSCRIBE<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>crossId</b></td>
        <td>string</td>
        <td>
          CrossID to export SharedPolicyGroup into different environments<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          SharedPolicyGroup description<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>prerequisiteMessage</b></td>
        <td>string</td>
        <td>
          SharedPolicyGroup prerequisite Message<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#sharedpolicygroupspecstepsindex">steps</a></b></td>
        <td>[]object</td>
        <td>
          SharedPolicyGroup Steps<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SharedPolicyGroup.spec.contextRef
[Go to parent definition](#sharedpolicygroupspec)





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


### SharedPolicyGroup.spec.steps[index]
[Go to parent definition](#sharedpolicygroupspec)





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


### SharedPolicyGroup.status
[Go to parent definition](#sharedpolicygroup)



SharedPolicyGroupSpecStatus defines the observed state of an API Context.

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
          The Cross ID is used to identify an SharedPolicyGroup that has been promoted from one environment...<br/>
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
        <td><b><a href="#sharedpolicygroupstatuserrors">errors</a></b></td>
        <td>object</td>
        <td>
          When SharedPolicyGroup has been created regardless of errors, this field is
used to persist the...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID is used to identify an SharedPolicyGroup which is unique in any environment.<br/>
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
        <td><b>processingStatus</b></td>
        <td>string</td>
        <td>
          The processing status of the SharedPolicyGroup.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### SharedPolicyGroup.status.errors
[Go to parent definition](#sharedpolicygroupstatus)



When SharedPolicyGroup has been created regardless of errors, this field is
used to persist the...

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
hence, this field should always be...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>warning</b></td>
        <td>[]string</td>
        <td>
          warning errors do not block object reconciliation,
most of the time because the value is ignored or...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## Notification

[gravitee.io/v1alpha1](#graviteeiov1alpha1)






Notification defines notification settings in Gravitee

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
        <td><b><a href="#notificationspec">spec</a></b></td>
        <td>object</td>
        <td>
          NotificationSpec defines the desired state of a Notification.
It is to be referenced in an API.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#notificationstatus">status</a></b></td>
        <td>object</td>
        <td>
          NotificationStatus defines the observed state of the Notification.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Notification.spec
[Go to parent definition](#notification)



NotificationSpec defines the desired state of a Notification.
It is to be referenced in an API.

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
        <td><b>eventType</b></td>
        <td>enum</td>
        <td>
          EventType defines the subject of those events.<br/>
          <br/>
            <i>Enum</i>: api<br/>
            <i>Default</i>: api<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>target</b></td>
        <td>enum</td>
        <td>
          Target of the notification: "console" is for notifications in Gravitee console UI.<br/>
          <br/>
            <i>Enum</i>: console<br/>
            <i>Default</i>: console<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#notificationspecconsole">console</a></b></td>
        <td>object</td>
        <td>
          Console is used when the target value is "console" and is meant
to configure Gravitee console UI...<br/>
          <br/>
            <i>Default</i>: map[]<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Notification.spec.console
[Go to parent definition](#notificationspec)



Console is used when the target value is "console" and is meant
to configure Gravitee console UI...

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
        <td><b>apiEvents</b></td>
        <td>[]enum</td>
        <td>
          List events that will trigger a notification for an API.<br/>
          <br/>
            <i>Enum</i>: APIKEY_EXPIRED, APIKEY_RENEWED, APIKEY_REVOKED, SUBSCRIPTION_NEW, SUBSCRIPTION_ACCEPTED, SUBSCRIPTION_CLOSED, SUBSCRIPTION_PAUSED, SUBSCRIPTION_RESUMED, SUBSCRIPTION_REJECTED, SUBSCRIPTION_TRANSFERRED, SUBSCRIPTION_FAILED, NEW_SUPPORT_TICKET, API_STARTED, API_STOPPED, API_UPDATED, API_DEPLOYED, NEW_RATING, NEW_RATING_ANSWER, MESSAGE, ASK_FOR_REVIEW, REVIEW_OK, REQUEST_FOR_CHANGES, API_DEPRECATED, NEW_SPEC_GENERATED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#notificationspecconsolegrouprefsindex">groupRefs</a></b></td>
        <td>[]object</td>
        <td>
          List of group references associated with this console notification.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>groups</b></td>
        <td>[]string</td>
        <td>
          List of groups associated with the API.
These groups are id to existing groups in APIM.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Notification.spec.console.groupRefs[index]
[Go to parent definition](#notificationspecconsole)





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


### Notification.status
[Go to parent definition](#notification)



NotificationStatus defines the observed state of the Notification.

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
"Accepted" condition is used to...<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Notification.status.conditions[index]
[Go to parent definition](#notificationstatus)



Condition contains details for one aspect of the current state of this API Resource.

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
          lastTransitionTime is the last time the condition transitioned from one status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          message is a human readable message indicating details about the transition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          reason contains a programmatic identifier indicating the reason for the condition's last transition.<br/>
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
          type of condition in CamelCase or in foo.example.com/CamelCase.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          observedGeneration represents the .metadata.generation that the condition was set based upon.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Minimum</i>: 0<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## Group

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
used to persist the error message...<br/>
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
used to persist the error message...

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
hence, this field should always be...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>warning</b></td>
        <td>[]string</td>
        <td>
          warning errors do not block object reconciliation,
most of the time because the value is ignored or...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

## GatewayClassParameters

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
        <td><b><a href="#gatewayclassparametersspec">spec</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersstatus">status</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec
[Go to parent definition](#gatewayclassparameters)





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
        <td><b><a href="#gatewayclassparametersspecgravitee">gravitee</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetes">kubernetes</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.gravitee
[Go to parent definition](#gatewayclassparametersspec)





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
        <td><b><a href="#gatewayclassparametersspecgraviteekafka">kafka</a></b></td>
        <td>object</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: map[enabled:false]<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspecgraviteelicenseref">licenseRef</a></b></td>
        <td>object</td>
        <td>
          SecretObjectReference identifies an API object including its namespace,
defaulting to Secret.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>yaml</b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.gravitee.kafka
[Go to parent definition](#gatewayclassparametersspecgravitee)





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
          <br/>
          <br/>
            <i>Default</i>: true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspecgraviteekafkaroutinghostmode">routingHostMode</a></b></td>
        <td>object</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: map[]<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.gravitee.kafka.routingHostMode
[Go to parent definition](#gatewayclassparametersspecgraviteekafka)





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
        <td><b>bootstrapDomainPattern</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: {apiHost}<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>brokerDomainPattern</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: broker-{brokerId}-{apiHost}<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.gravitee.licenseRef
[Go to parent definition](#gatewayclassparametersspecgravitee)



SecretObjectReference identifies an API object including its namespace,
defaulting to Secret.

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
          Name is the name of the referent.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>group</b></td>
        <td>string</td>
        <td>
          Group is the group of the referent. For example, "gateway.networking.k8s.io".<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind is kind of the referent. For example "Secret".<br/>
          <br/>
            <i>Default</i>: Secret<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace is the namespace of the referenced object.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes
[Go to parent definition](#gatewayclassparametersspec)





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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeployment">deployment</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesservice">service</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment
[Go to parent definition](#gatewayclassparametersspeckubernetes)





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
        <td><b>annotations</b></td>
        <td>map[string]string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>labels</b></td>
        <td>map[string]string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>replicas</b></td>
        <td>integer</td>
        <td>
          <br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymentstrategy">strategy</a></b></td>
        <td>object</td>
        <td>
          DeploymentStrategy describes how to replace existing pods with new ones.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplate">template</a></b></td>
        <td>object</td>
        <td>
          PodTemplateSpec describes the data a pod should have when created from a template<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.strategy
[Go to parent definition](#gatewayclassparametersspeckubernetesdeployment)



DeploymentStrategy describes how to replace existing pods with new ones.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymentstrategyrollingupdate">rollingUpdate</a></b></td>
        <td>object</td>
        <td>
          Rolling update config params. Present only if DeploymentStrategyType =
RollingUpdate.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of deployment. Can be "Recreate" or "RollingUpdate". Default is RollingUpdate.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.strategy.rollingUpdate
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymentstrategy)



Rolling update config params. Present only if DeploymentStrategyType =
RollingUpdate.

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
        <td><b>maxSurge</b></td>
        <td>int or string</td>
        <td>
          The maximum number of pods that can be scheduled above the desired number of
pods.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxUnavailable</b></td>
        <td>int or string</td>
        <td>
          The maximum number of pods that can be unavailable during the update.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template
[Go to parent definition](#gatewayclassparametersspeckubernetesdeployment)



PodTemplateSpec describes the data a pod should have when created from a template

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
        <td><b>metadata</b></td>
        <td>object</td>
        <td>
          Standard object's metadata.
More info: https://git.k8s.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespec">spec</a></b></td>
        <td>object</td>
        <td>
          Specification of the desired behavior of the pod.
More info: https://git.k8s.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplate)



Specification of the desired behavior of the pod.
More info: https://git.k8s.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex">containers</a></b></td>
        <td>[]object</td>
        <td>
          List of containers belonging to the pod.
Containers cannot currently be added or removed.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>activeDeadlineSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod may be active on the node relative to
StartTime before the...<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinity">affinity</a></b></td>
        <td>object</td>
        <td>
          If specified, the pod's scheduling constraints<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>automountServiceAccountToken</b></td>
        <td>boolean</td>
        <td>
          AutomountServiceAccountToken indicates whether a service account token should be automatically...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecdnsconfig">dnsConfig</a></b></td>
        <td>object</td>
        <td>
          Specifies the DNS parameters of a pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>dnsPolicy</b></td>
        <td>string</td>
        <td>
          Set DNS policy for the pod.
Defaults to "ClusterFirst".<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enableServiceLinks</b></td>
        <td>boolean</td>
        <td>
          EnableServiceLinks indicates whether information about services should be injected into pod's...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex">ephemeralContainers</a></b></td>
        <td>[]object</td>
        <td>
          List of ephemeral containers run in this pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespechostaliasesindex">hostAliases</a></b></td>
        <td>[]object</td>
        <td>
          HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts
file if...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostIPC</b></td>
        <td>boolean</td>
        <td>
          Use the host's ipc namespace.
Optional: Default to false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostNetwork</b></td>
        <td>boolean</td>
        <td>
          Host networking requested for this pod. Use the host's network namespace.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostPID</b></td>
        <td>boolean</td>
        <td>
          Use the host's pid namespace.
Optional: Default to false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostUsers</b></td>
        <td>boolean</td>
        <td>
          Use the host's user namespace.
Optional: Default to true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostname</b></td>
        <td>string</td>
        <td>
          Specifies the hostname of the Pod
If not specified, the pod's hostname will be set to a...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecimagepullsecretsindex">imagePullSecrets</a></b></td>
        <td>[]object</td>
        <td>
          ImagePullSecrets is an optional list of references to secrets in the same namespace to use for...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex">initContainers</a></b></td>
        <td>[]object</td>
        <td>
          List of initialization containers belonging to the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>nodeName</b></td>
        <td>string</td>
        <td>
          NodeName indicates in which node this pod is scheduled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>nodeSelector</b></td>
        <td>map[string]string</td>
        <td>
          NodeSelector is a selector which must be true for the pod to fit on a node.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecos">os</a></b></td>
        <td>object</td>
        <td>
          Specifies the OS of the containers in the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>overhead</b></td>
        <td>map[string]int or string</td>
        <td>
          Overhead represents the resource overhead associated with running a pod for a given RuntimeClass.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>preemptionPolicy</b></td>
        <td>string</td>
        <td>
          PreemptionPolicy is the Policy for preempting pods with lower priority.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>priority</b></td>
        <td>integer</td>
        <td>
          The priority value. Various system components use this field to find the
priority of the pod.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>priorityClassName</b></td>
        <td>string</td>
        <td>
          If specified, indicates the pod's priority.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecreadinessgatesindex">readinessGates</a></b></td>
        <td>[]object</td>
        <td>
          If specified, all readiness gates will be evaluated for pod readiness.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecresourceclaimsindex">resourceClaims</a></b></td>
        <td>[]object</td>
        <td>
          ResourceClaims defines which ResourceClaims must be allocated
and reserved before the Pod is...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecresources">resources</a></b></td>
        <td>object</td>
        <td>
          Resources is the total amount of CPU and Memory resources required by all
containers in the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>restartPolicy</b></td>
        <td>string</td>
        <td>
          Restart policy for all containers within the pod.
One of Always, OnFailure, Never.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runtimeClassName</b></td>
        <td>string</td>
        <td>
          RuntimeClassName refers to a RuntimeClass object in the node.k8s.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>schedulerName</b></td>
        <td>string</td>
        <td>
          If specified, the pod will be dispatched by specified scheduler.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecschedulinggatesindex">schedulingGates</a></b></td>
        <td>[]object</td>
        <td>
          SchedulingGates is an opaque list of values that if specified will block scheduling the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecsecuritycontext">securityContext</a></b></td>
        <td>object</td>
        <td>
          SecurityContext holds pod-level security attributes and common container settings.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceAccount</b></td>
        <td>string</td>
        <td>
          DeprecatedServiceAccount is a deprecated alias for ServiceAccountName.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceAccountName</b></td>
        <td>string</td>
        <td>
          ServiceAccountName is the name of the ServiceAccount to use to run this pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>setHostnameAsFQDN</b></td>
        <td>boolean</td>
        <td>
          If true the pod's hostname will be configured as the pod's FQDN, rather than the leaf name (the...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>shareProcessNamespace</b></td>
        <td>boolean</td>
        <td>
          Share a single process namespace between all of the containers in a pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subdomain</b></td>
        <td>string</td>
        <td>
          If specified, the fully qualified Pod hostname will be "<hostname>.<subdomain>.<pod namespace>.svc.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespectolerationsindex">tolerations</a></b></td>
        <td>[]object</td>
        <td>
          If specified, the pod's tolerations.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespectopologyspreadconstraintsindex">topologySpreadConstraints</a></b></td>
        <td>[]object</td>
        <td>
          TopologySpreadConstraints describes how a group of pods ought to spread across topology
domains.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex">volumes</a></b></td>
        <td>[]object</td>
        <td>
          List of volumes that can be mounted by containers belonging to the pod.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



A single application container that you want to run within a pod.

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
          Name of the container specified as a DNS_LABEL.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>args</b></td>
        <td>[]string</td>
        <td>
          Arguments to the entrypoint.
The container image's CMD is used if this is not provided.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Entrypoint array. Not executed within a shell.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvindex">env</a></b></td>
        <td>[]object</td>
        <td>
          List of environment variables to set in the container.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvfromindex">envFrom</a></b></td>
        <td>[]object</td>
        <td>
          List of sources to populate environment variables in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>image</b></td>
        <td>string</td>
        <td>
          Container image name.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imagePullPolicy</b></td>
        <td>string</td>
        <td>
          Image pull policy.
One of Always, Never, IfNotPresent.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycle">lifecycle</a></b></td>
        <td>object</td>
        <td>
          Actions that the management system should take in response to container lifecycle events.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlivenessprobe">livenessProbe</a></b></td>
        <td>object</td>
        <td>
          Periodic probe of container liveness.
Container will be restarted if the probe fails.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexportsindex">ports</a></b></td>
        <td>[]object</td>
        <td>
          List of ports to expose from the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexreadinessprobe">readinessProbe</a></b></td>
        <td>object</td>
        <td>
          Periodic probe of container service readiness.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexresizepolicyindex">resizePolicy</a></b></td>
        <td>[]object</td>
        <td>
          Resources resize policy for the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexresources">resources</a></b></td>
        <td>object</td>
        <td>
          Compute Resources required by this container.
Cannot be updated.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>restartPolicy</b></td>
        <td>string</td>
        <td>
          RestartPolicy defines the restart behavior of individual containers in a pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexsecuritycontext">securityContext</a></b></td>
        <td>object</td>
        <td>
          SecurityContext defines the security options the container should be run with.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexstartupprobe">startupProbe</a></b></td>
        <td>object</td>
        <td>
          StartupProbe indicates that the Pod has successfully initialized.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stdin</b></td>
        <td>boolean</td>
        <td>
          Whether this container should allocate a buffer for stdin in the container runtime.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stdinOnce</b></td>
        <td>boolean</td>
        <td>
          Whether the container runtime should close the stdin channel after it has been opened by
a single...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationMessagePath</b></td>
        <td>string</td>
        <td>
          Optional: Path at which the file to which the container's termination message
will be written is...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationMessagePolicy</b></td>
        <td>string</td>
        <td>
          Indicate how the termination message should be populated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tty</b></td>
        <td>boolean</td>
        <td>
          Whether this container should allocate a TTY for itself, also requires 'stdin' to be true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexvolumedevicesindex">volumeDevices</a></b></td>
        <td>[]object</td>
        <td>
          volumeDevices is the list of block devices to be used by the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexvolumemountsindex">volumeMounts</a></b></td>
        <td>[]object</td>
        <td>
          Pod volumes to mount into the container's filesystem.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>workingDir</b></td>
        <td>string</td>
        <td>
          Container's working directory.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].env[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



EnvVar represents an environment variable present in a Container.

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
          Name of the environment variable. Must be a C_IDENTIFIER.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Variable references $(VAR_NAME) are expanded
using the previously defined environment variables in...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvindexvaluefrom">valueFrom</a></b></td>
        <td>object</td>
        <td>
          Source for the environment variable's value. Cannot be used if value is not empty.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].env[index].valueFrom
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvindex)



Source for the environment variable's value. Cannot be used if value is not empty.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvindexvaluefromconfigmapkeyref">configMapKeyRef</a></b></td>
        <td>object</td>
        <td>
          Selects a key of a ConfigMap.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvindexvaluefromfieldref">fieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvindexvaluefromresourcefieldref">resourceFieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvindexvaluefromsecretkeyref">secretKeyRef</a></b></td>
        <td>object</td>
        <td>
          Selects a key of a secret in the pod's namespace<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].env[index].valueFrom.configMapKeyRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvindexvaluefrom)



Selects a key of a ConfigMap.

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key to select.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the ConfigMap or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].env[index].valueFrom.fieldRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvindexvaluefrom)



Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.

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
        <td><b>fieldPath</b></td>
        <td>string</td>
        <td>
          Path of the field to select in the specified API version.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          Version of the schema the FieldPath is written in terms of, defaults to "v1".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].env[index].valueFrom.resourceFieldRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvindexvaluefrom)



Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.

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
        <td><b>resource</b></td>
        <td>string</td>
        <td>
          Required: resource to select<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>containerName</b></td>
        <td>string</td>
        <td>
          Container name: required for volumes, optional for env vars<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>divisor</b></td>
        <td>int or string</td>
        <td>
          Specifies the output format of the exposed resources, defaults to "1"<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].env[index].valueFrom.secretKeyRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvindexvaluefrom)



Selects a key of a secret in the pod's namespace

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key of the secret to select from.  Must be a valid secret key.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the Secret or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].envFrom[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



EnvFromSource represents the source of a set of ConfigMaps or Secrets

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvfromindexconfigmapref">configMapRef</a></b></td>
        <td>object</td>
        <td>
          The ConfigMap to select from<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>prefix</b></td>
        <td>string</td>
        <td>
          Optional text to prepend to the name of each environment variable. Must be a C_IDENTIFIER.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvfromindexsecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          The Secret to select from<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].envFrom[index].configMapRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvfromindex)



The ConfigMap to select from

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the ConfigMap must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].envFrom[index].secretRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexenvfromindex)



The Secret to select from

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the Secret must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



Actions that the management system should take in response to container lifecycle events.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecyclepoststart">postStart</a></b></td>
        <td>object</td>
        <td>
          PostStart is called immediately after a container is created.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycleprestop">preStop</a></b></td>
        <td>object</td>
        <td>
          PreStop is called immediately before a container is terminated due to an
API request or management...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stopSignal</b></td>
        <td>string</td>
        <td>
          StopSignal defines which signal will be sent to a container when it is being stopped.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.postStart
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycle)



PostStart is called immediately after a container is created.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecyclepoststartexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecyclepoststarthttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecyclepoststartsleep">sleep</a></b></td>
        <td>object</td>
        <td>
          Sleep represents a duration that the container should sleep.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecyclepoststarttcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.postStart.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecyclepoststart)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.postStart.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecyclepoststart)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecyclepoststarthttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.postStart.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecyclepoststarthttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.postStart.sleep
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecyclepoststart)



Sleep represents a duration that the container should sleep.

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
        <td><b>seconds</b></td>
        <td>integer</td>
        <td>
          Seconds is the number of seconds to sleep.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.postStart.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecyclepoststart)



Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.preStop
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycle)



PreStop is called immediately before a container is terminated due to an
API request or management...

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycleprestopexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycleprestophttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycleprestopsleep">sleep</a></b></td>
        <td>object</td>
        <td>
          Sleep represents a duration that the container should sleep.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycleprestoptcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.preStop.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycleprestop)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.preStop.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycleprestop)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycleprestophttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.preStop.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycleprestophttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.preStop.sleep
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycleprestop)



Sleep represents a duration that the container should sleep.

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
        <td><b>seconds</b></td>
        <td>integer</td>
        <td>
          Seconds is the number of seconds to sleep.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].lifecycle.preStop.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlifecycleprestop)



Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].livenessProbe
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



Periodic probe of container liveness.
Container will be restarted if the probe fails.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlivenessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlivenessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlivenessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlivenessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].livenessProbe.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlivenessprobe)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].livenessProbe.grpc
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlivenessprobe)



GRPC specifies a GRPC HealthCheckRequest.

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
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].livenessProbe.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlivenessprobe)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlivenessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].livenessProbe.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlivenessprobehttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].livenessProbe.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexlivenessprobe)



TCPSocket specifies a connection to a TCP port.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].ports[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



ContainerPort represents a network port in a single container.

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
        <td><b>containerPort</b></td>
        <td>integer</td>
        <td>
          Number of port to expose on the pod's IP address.
This must be a valid port number, 0 < x < 65536.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>hostIP</b></td>
        <td>string</td>
        <td>
          What host IP to bind the external port to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostPort</b></td>
        <td>integer</td>
        <td>
          Number of port to expose on the host.
If specified, this must be a valid port number, 0 < x < 65536.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          If specified, this must be an IANA_SVC_NAME and unique within the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>protocol</b></td>
        <td>string</td>
        <td>
          Protocol for port. Must be UDP, TCP, or SCTP.
Defaults to "TCP".<br/>
          <br/>
            <i>Default</i>: TCP<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].readinessProbe
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



Periodic probe of container service readiness.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexreadinessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexreadinessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexreadinessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexreadinessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].readinessProbe.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexreadinessprobe)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].readinessProbe.grpc
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexreadinessprobe)



GRPC specifies a GRPC HealthCheckRequest.

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
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].readinessProbe.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexreadinessprobe)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexreadinessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].readinessProbe.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexreadinessprobehttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].readinessProbe.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexreadinessprobe)



TCPSocket specifies a connection to a TCP port.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].resizePolicy[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



ContainerResizePolicy represents resource resize policy for the container.

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
        <td><b>resourceName</b></td>
        <td>string</td>
        <td>
          Name of the resource to which this resource resize policy applies.
Supported values: cpu, memory.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>restartPolicy</b></td>
        <td>string</td>
        <td>
          Restart policy to apply when specified resource is resized.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].resources
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



Compute Resources required by this container.
Cannot be updated.
More info: https://kubernetes.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexresourcesclaimsindex">claims</a></b></td>
        <td>[]object</td>
        <td>
          Claims lists the names of resources, defined in spec.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>limits</b></td>
        <td>map[string]int or string</td>
        <td>
          Limits describes the maximum amount of compute resources allowed.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>requests</b></td>
        <td>map[string]int or string</td>
        <td>
          Requests describes the minimum amount of compute resources required.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].resources.claims[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexresources)



ResourceClaim references one entry in PodSpec.ResourceClaims.

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
          Name must match the name of one entry in pod.spec.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>request</b></td>
        <td>string</td>
        <td>
          Request is the name chosen for a request in the referenced claim.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].securityContext
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



SecurityContext defines the security options the container should be run with.

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
        <td><b>allowPrivilegeEscalation</b></td>
        <td>boolean</td>
        <td>
          AllowPrivilegeEscalation controls whether a process can gain more
privileges than its parent...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexsecuritycontextapparmorprofile">appArmorProfile</a></b></td>
        <td>object</td>
        <td>
          appArmorProfile is the AppArmor options to use by this container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexsecuritycontextcapabilities">capabilities</a></b></td>
        <td>object</td>
        <td>
          The capabilities to add/drop when running containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>privileged</b></td>
        <td>boolean</td>
        <td>
          Run container in privileged mode.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>procMount</b></td>
        <td>string</td>
        <td>
          procMount denotes the type of proc mount to use for the containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnlyRootFilesystem</b></td>
        <td>boolean</td>
        <td>
          Whether this container has a read-only root filesystem.
Default is false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsGroup</b></td>
        <td>integer</td>
        <td>
          The GID to run the entrypoint of the container process.
Uses runtime default if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsNonRoot</b></td>
        <td>boolean</td>
        <td>
          Indicates that the container must run as a non-root user.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUser</b></td>
        <td>integer</td>
        <td>
          The UID to run the entrypoint of the container process.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexsecuritycontextselinuxoptions">seLinuxOptions</a></b></td>
        <td>object</td>
        <td>
          The SELinux context to be applied to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexsecuritycontextseccompprofile">seccompProfile</a></b></td>
        <td>object</td>
        <td>
          The seccomp options to use by this container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexsecuritycontextwindowsoptions">windowsOptions</a></b></td>
        <td>object</td>
        <td>
          The Windows specific settings applied to all containers.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].securityContext.appArmorProfile
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexsecuritycontext)



appArmorProfile is the AppArmor options to use by this container.

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
          type indicates which kind of AppArmor profile will be applied.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile loaded on the node that should be used.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].securityContext.capabilities
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexsecuritycontext)



The capabilities to add/drop when running containers.

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
        <td><b>add</b></td>
        <td>[]string</td>
        <td>
          Added capabilities<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>drop</b></td>
        <td>[]string</td>
        <td>
          Removed capabilities<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].securityContext.seLinuxOptions
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexsecuritycontext)



The SELinux context to be applied to the container.

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
        <td><b>level</b></td>
        <td>string</td>
        <td>
          Level is SELinux level label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>role</b></td>
        <td>string</td>
        <td>
          Role is a SELinux role label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type is a SELinux type label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          User is a SELinux user label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].securityContext.seccompProfile
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexsecuritycontext)



The seccomp options to use by this container.

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
          type indicates which kind of seccomp profile will be applied.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile defined in a file on the node should be used.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].securityContext.windowsOptions
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexsecuritycontext)



The Windows specific settings applied to all containers.

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
        <td><b>gmsaCredentialSpec</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpec is where the GMSA admission webhook
(https://github.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gmsaCredentialSpecName</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpecName is the name of the GMSA credential spec to use.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostProcess</b></td>
        <td>boolean</td>
        <td>
          HostProcess determines if a container should be run as a 'Host Process' container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUserName</b></td>
        <td>string</td>
        <td>
          The UserName in Windows to run the entrypoint of the container process.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].startupProbe
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



StartupProbe indicates that the Pod has successfully initialized.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexstartupprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexstartupprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexstartupprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexstartupprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].startupProbe.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexstartupprobe)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].startupProbe.grpc
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexstartupprobe)



GRPC specifies a GRPC HealthCheckRequest.

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
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].startupProbe.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexstartupprobe)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexstartupprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].startupProbe.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexstartupprobehttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].startupProbe.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindexstartupprobe)



TCPSocket specifies a connection to a TCP port.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].volumeDevices[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



volumeDevice describes a mapping of a raw block device within a container.

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
        <td><b>devicePath</b></td>
        <td>string</td>
        <td>
          devicePath is the path inside of the container that the device will be mapped to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          name must match the name of a persistentVolumeClaim in the pod<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.containers[index].volumeMounts[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespeccontainersindex)



VolumeMount describes a mounting of a Volume within a container.

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
        <td><b>mountPath</b></td>
        <td>string</td>
        <td>
          Path within the container at which the volume should be mounted.  Must
not contain ':'.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          This must match the Name of a Volume.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mountPropagation</b></td>
        <td>string</td>
        <td>
          mountPropagation determines how mounts are propagated from the host
to container and the other way...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          Mounted read-only if true, read-write otherwise (false or unspecified).
Defaults to false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>recursiveReadOnly</b></td>
        <td>string</td>
        <td>
          RecursiveReadOnly specifies whether read-only mounts should be handled
recursively.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subPath</b></td>
        <td>string</td>
        <td>
          Path within the volume from which the container's volume should be mounted.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subPathExpr</b></td>
        <td>string</td>
        <td>
          Expanded path within the volume from which the container's volume should be mounted.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



If specified, the pod's scheduling constraints

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinity">nodeAffinity</a></b></td>
        <td>object</td>
        <td>
          Describes node affinity scheduling rules for the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinity">podAffinity</a></b></td>
        <td>object</td>
        <td>
          Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone, etc.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinity">podAntiAffinity</a></b></td>
        <td>object</td>
        <td>
          Describes pod anti-affinity scheduling rules (e.g.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.nodeAffinity
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinity)



Describes node affinity scheduling rules for the pod.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinitypreferredduringschedulingignoredduringexecutionindex">preferredDuringSchedulingIgnoredDuringExecution</a></b></td>
        <td>[]object</td>
        <td>
          The scheduler will prefer to schedule pods to nodes that satisfy
the affinity expressions specified...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinityrequiredduringschedulingignoredduringexecution">requiredDuringSchedulingIgnoredDuringExecution</a></b></td>
        <td>object</td>
        <td>
          If the affinity requirements specified by this field are not met at
scheduling time, the pod will...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinity)



An empty preferred scheduling term matches all objects with implicit weight 0
(i.e. it's a no-op).

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinitypreferredduringschedulingignoredduringexecutionindexpreference">preference</a></b></td>
        <td>object</td>
        <td>
          A node selector term, associated with the corresponding weight.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>weight</b></td>
        <td>integer</td>
        <td>
          Weight associated with matching the corresponding nodeSelectorTerm, in the range 1-100.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].preference
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinitypreferredduringschedulingignoredduringexecutionindex)



A node selector term, associated with the corresponding weight.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinitypreferredduringschedulingignoredduringexecutionindexpreferencematchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          A list of node selector requirements by node's labels.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinitypreferredduringschedulingignoredduringexecutionindexpreferencematchfieldsindex">matchFields</a></b></td>
        <td>[]object</td>
        <td>
          A list of node selector requirements by node's fields.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].preference.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinitypreferredduringschedulingignoredduringexecutionindexpreference)



A node selector requirement is a selector that contains values, a key, and an operator
that relates...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          An array of string values. If the operator is In or NotIn,
the values array must be non-empty.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].preference.matchFields[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinitypreferredduringschedulingignoredduringexecutionindexpreference)



A node selector requirement is a selector that contains values, a key, and an operator
that relates...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          An array of string values. If the operator is In or NotIn,
the values array must be non-empty.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinity)



If the affinity requirements specified by this field are not met at
scheduling time, the pod will...

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinityrequiredduringschedulingignoredduringexecutionnodeselectortermsindex">nodeSelectorTerms</a></b></td>
        <td>[]object</td>
        <td>
          Required. A list of node selector terms. The terms are ORed.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinityrequiredduringschedulingignoredduringexecution)



A null or empty node selector term matches no objects. The requirements of
them are ANDed.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinityrequiredduringschedulingignoredduringexecutionnodeselectortermsindexmatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          A list of node selector requirements by node's labels.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinityrequiredduringschedulingignoredduringexecutionnodeselectortermsindexmatchfieldsindex">matchFields</a></b></td>
        <td>[]object</td>
        <td>
          A list of node selector requirements by node's fields.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[index].matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinityrequiredduringschedulingignoredduringexecutionnodeselectortermsindex)



A node selector requirement is a selector that contains values, a key, and an operator
that relates...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          An array of string values. If the operator is In or NotIn,
the values array must be non-empty.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[index].matchFields[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitynodeaffinityrequiredduringschedulingignoredduringexecutionnodeselectortermsindex)



A node selector requirement is a selector that contains values, a key, and an operator
that relates...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          An array of string values. If the operator is In or NotIn,
the values array must be non-empty.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinity)



Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone, etc.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinitypreferredduringschedulingignoredduringexecutionindex">preferredDuringSchedulingIgnoredDuringExecution</a></b></td>
        <td>[]object</td>
        <td>
          The scheduler will prefer to schedule pods to nodes that satisfy
the affinity expressions specified...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinityrequiredduringschedulingignoredduringexecutionindex">requiredDuringSchedulingIgnoredDuringExecution</a></b></td>
        <td>[]object</td>
        <td>
          If the affinity requirements specified by this field are not met at
scheduling time, the pod will...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinity)



The weights of all of the matched WeightedPodAffinityTerm fields are added per-node to find the...

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinityterm">podAffinityTerm</a></b></td>
        <td>object</td>
        <td>
          Required. A pod affinity term, associated with the corresponding weight.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>weight</b></td>
        <td>integer</td>
        <td>
          weight associated with matching the corresponding podAffinityTerm,
in the range 1-100.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].podAffinityTerm
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinitypreferredduringschedulingignoredduringexecutionindex)



Required. A pod affinity term, associated with the corresponding weight.

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
        <td><b>topologyKey</b></td>
        <td>string</td>
        <td>
          This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching...<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermlabelselector">labelSelector</a></b></td>
        <td>object</td>
        <td>
          A label query over a set of resources, in this case pods.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabelKeys</b></td>
        <td>[]string</td>
        <td>
          MatchLabelKeys is a set of pod label keys to select which pods will
be taken into consideration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mismatchLabelKeys</b></td>
        <td>[]string</td>
        <td>
          MismatchLabelKeys is a set of pod label keys to select which pods will
be taken into consideration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermnamespaceselector">namespaceSelector</a></b></td>
        <td>object</td>
        <td>
          A label query over the set of namespaces that the term applies to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespaces</b></td>
        <td>[]string</td>
        <td>
          namespaces specifies a static list of namespace names that the term applies to.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].podAffinityTerm.labelSelector
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinityterm)



A label query over a set of resources, in this case pods.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermlabelselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].podAffinityTerm.labelSelector.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermlabelselector)



A label selector requirement is a selector that contains values, a key, and an operator that...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].podAffinityTerm.namespaceSelector
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinityterm)



A label query over the set of namespaces that the term applies to.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermnamespaceselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].podAffinityTerm.namespaceSelector.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermnamespaceselector)



A label selector requirement is a selector that contains values, a key, and an operator that...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinity)



Defines a set of pods (namely those matching the labelSelector
relative to the given namespace(s))...

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
        <td><b>topologyKey</b></td>
        <td>string</td>
        <td>
          This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching...<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinityrequiredduringschedulingignoredduringexecutionindexlabelselector">labelSelector</a></b></td>
        <td>object</td>
        <td>
          A label query over a set of resources, in this case pods.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabelKeys</b></td>
        <td>[]string</td>
        <td>
          MatchLabelKeys is a set of pod label keys to select which pods will
be taken into consideration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mismatchLabelKeys</b></td>
        <td>[]string</td>
        <td>
          MismatchLabelKeys is a set of pod label keys to select which pods will
be taken into consideration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinityrequiredduringschedulingignoredduringexecutionindexnamespaceselector">namespaceSelector</a></b></td>
        <td>object</td>
        <td>
          A label query over the set of namespaces that the term applies to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespaces</b></td>
        <td>[]string</td>
        <td>
          namespaces specifies a static list of namespace names that the term applies to.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution[index].labelSelector
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinityrequiredduringschedulingignoredduringexecutionindex)



A label query over a set of resources, in this case pods.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinityrequiredduringschedulingignoredduringexecutionindexlabelselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution[index].labelSelector.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinityrequiredduringschedulingignoredduringexecutionindexlabelselector)



A label selector requirement is a selector that contains values, a key, and an operator that...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution[index].namespaceSelector
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinityrequiredduringschedulingignoredduringexecutionindex)



A label query over the set of namespaces that the term applies to.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinityrequiredduringschedulingignoredduringexecutionindexnamespaceselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution[index].namespaceSelector.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodaffinityrequiredduringschedulingignoredduringexecutionindexnamespaceselector)



A label selector requirement is a selector that contains values, a key, and an operator that...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinity)



Describes pod anti-affinity scheduling rules (e.g.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinitypreferredduringschedulingignoredduringexecutionindex">preferredDuringSchedulingIgnoredDuringExecution</a></b></td>
        <td>[]object</td>
        <td>
          The scheduler will prefer to schedule pods to nodes that satisfy
the anti-affinity expressions...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinityrequiredduringschedulingignoredduringexecutionindex">requiredDuringSchedulingIgnoredDuringExecution</a></b></td>
        <td>[]object</td>
        <td>
          If the anti-affinity requirements specified by this field are not met at
scheduling time, the pod...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinity)



The weights of all of the matched WeightedPodAffinityTerm fields are added per-node to find the...

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinityterm">podAffinityTerm</a></b></td>
        <td>object</td>
        <td>
          Required. A pod affinity term, associated with the corresponding weight.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>weight</b></td>
        <td>integer</td>
        <td>
          weight associated with matching the corresponding podAffinityTerm,
in the range 1-100.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].podAffinityTerm
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinitypreferredduringschedulingignoredduringexecutionindex)



Required. A pod affinity term, associated with the corresponding weight.

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
        <td><b>topologyKey</b></td>
        <td>string</td>
        <td>
          This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching...<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermlabelselector">labelSelector</a></b></td>
        <td>object</td>
        <td>
          A label query over a set of resources, in this case pods.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabelKeys</b></td>
        <td>[]string</td>
        <td>
          MatchLabelKeys is a set of pod label keys to select which pods will
be taken into consideration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mismatchLabelKeys</b></td>
        <td>[]string</td>
        <td>
          MismatchLabelKeys is a set of pod label keys to select which pods will
be taken into consideration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermnamespaceselector">namespaceSelector</a></b></td>
        <td>object</td>
        <td>
          A label query over the set of namespaces that the term applies to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespaces</b></td>
        <td>[]string</td>
        <td>
          namespaces specifies a static list of namespace names that the term applies to.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].podAffinityTerm.labelSelector
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinityterm)



A label query over a set of resources, in this case pods.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermlabelselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].podAffinityTerm.labelSelector.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermlabelselector)



A label selector requirement is a selector that contains values, a key, and an operator that...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].podAffinityTerm.namespaceSelector
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinityterm)



A label query over the set of namespaces that the term applies to.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermnamespaceselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution[index].podAffinityTerm.namespaceSelector.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinitypreferredduringschedulingignoredduringexecutionindexpodaffinitytermnamespaceselector)



A label selector requirement is a selector that contains values, a key, and an operator that...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinity)



Defines a set of pods (namely those matching the labelSelector
relative to the given namespace(s))...

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
        <td><b>topologyKey</b></td>
        <td>string</td>
        <td>
          This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching...<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinityrequiredduringschedulingignoredduringexecutionindexlabelselector">labelSelector</a></b></td>
        <td>object</td>
        <td>
          A label query over a set of resources, in this case pods.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabelKeys</b></td>
        <td>[]string</td>
        <td>
          MatchLabelKeys is a set of pod label keys to select which pods will
be taken into consideration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mismatchLabelKeys</b></td>
        <td>[]string</td>
        <td>
          MismatchLabelKeys is a set of pod label keys to select which pods will
be taken into consideration.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinityrequiredduringschedulingignoredduringexecutionindexnamespaceselector">namespaceSelector</a></b></td>
        <td>object</td>
        <td>
          A label query over the set of namespaces that the term applies to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespaces</b></td>
        <td>[]string</td>
        <td>
          namespaces specifies a static list of namespace names that the term applies to.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution[index].labelSelector
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinityrequiredduringschedulingignoredduringexecutionindex)



A label query over a set of resources, in this case pods.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinityrequiredduringschedulingignoredduringexecutionindexlabelselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution[index].labelSelector.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinityrequiredduringschedulingignoredduringexecutionindexlabelselector)



A label selector requirement is a selector that contains values, a key, and an operator that...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution[index].namespaceSelector
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinityrequiredduringschedulingignoredduringexecutionindex)



A label query over the set of namespaces that the term applies to.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinityrequiredduringschedulingignoredduringexecutionindexnamespaceselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution[index].namespaceSelector.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecaffinitypodantiaffinityrequiredduringschedulingignoredduringexecutionindexnamespaceselector)



A label selector requirement is a selector that contains values, a key, and an operator that...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.dnsConfig
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



Specifies the DNS parameters of a pod.

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
        <td><b>nameservers</b></td>
        <td>[]string</td>
        <td>
          A list of DNS name server IP addresses.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecdnsconfigoptionsindex">options</a></b></td>
        <td>[]object</td>
        <td>
          A list of DNS resolver options.
This will be merged with the base options generated from DNSPolicy.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>searches</b></td>
        <td>[]string</td>
        <td>
          A list of DNS search domains for host-name lookup.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.dnsConfig.options[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecdnsconfig)



PodDNSConfigOption defines DNS resolver options of a pod.

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
          Name is this DNS resolver option's name.
Required.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Value is this DNS resolver option's value.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



An EphemeralContainer is a temporary container that you may add to an existing Pod for...

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
          Name of the ephemeral container specified as a DNS_LABEL.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>args</b></td>
        <td>[]string</td>
        <td>
          Arguments to the entrypoint.
The image's CMD is used if this is not provided.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Entrypoint array. Not executed within a shell.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvindex">env</a></b></td>
        <td>[]object</td>
        <td>
          List of environment variables to set in the container.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvfromindex">envFrom</a></b></td>
        <td>[]object</td>
        <td>
          List of sources to populate environment variables in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>image</b></td>
        <td>string</td>
        <td>
          Container image name.
More info: https://kubernetes.io/docs/concepts/containers/images<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imagePullPolicy</b></td>
        <td>string</td>
        <td>
          Image pull policy.
One of Always, Never, IfNotPresent.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycle">lifecycle</a></b></td>
        <td>object</td>
        <td>
          Lifecycle is not allowed for ephemeral containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlivenessprobe">livenessProbe</a></b></td>
        <td>object</td>
        <td>
          Probes are not allowed for ephemeral containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexportsindex">ports</a></b></td>
        <td>[]object</td>
        <td>
          Ports are not allowed for ephemeral containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexreadinessprobe">readinessProbe</a></b></td>
        <td>object</td>
        <td>
          Probes are not allowed for ephemeral containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexresizepolicyindex">resizePolicy</a></b></td>
        <td>[]object</td>
        <td>
          Resources resize policy for the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexresources">resources</a></b></td>
        <td>object</td>
        <td>
          Resources are not allowed for ephemeral containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>restartPolicy</b></td>
        <td>string</td>
        <td>
          Restart policy for the container to manage the restart behavior of each
container within a pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexsecuritycontext">securityContext</a></b></td>
        <td>object</td>
        <td>
          Optional: SecurityContext defines the security options the ephemeral container should be run with.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexstartupprobe">startupProbe</a></b></td>
        <td>object</td>
        <td>
          Probes are not allowed for ephemeral containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stdin</b></td>
        <td>boolean</td>
        <td>
          Whether this container should allocate a buffer for stdin in the container runtime.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stdinOnce</b></td>
        <td>boolean</td>
        <td>
          Whether the container runtime should close the stdin channel after it has been opened by
a single...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>targetContainerName</b></td>
        <td>string</td>
        <td>
          If set, the name of the container from PodSpec that this ephemeral container targets.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationMessagePath</b></td>
        <td>string</td>
        <td>
          Optional: Path at which the file to which the container's termination message
will be written is...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationMessagePolicy</b></td>
        <td>string</td>
        <td>
          Indicate how the termination message should be populated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tty</b></td>
        <td>boolean</td>
        <td>
          Whether this container should allocate a TTY for itself, also requires 'stdin' to be true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexvolumedevicesindex">volumeDevices</a></b></td>
        <td>[]object</td>
        <td>
          volumeDevices is the list of block devices to be used by the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexvolumemountsindex">volumeMounts</a></b></td>
        <td>[]object</td>
        <td>
          Pod volumes to mount into the container's filesystem.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>workingDir</b></td>
        <td>string</td>
        <td>
          Container's working directory.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].env[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



EnvVar represents an environment variable present in a Container.

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
          Name of the environment variable. Must be a C_IDENTIFIER.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Variable references $(VAR_NAME) are expanded
using the previously defined environment variables in...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvindexvaluefrom">valueFrom</a></b></td>
        <td>object</td>
        <td>
          Source for the environment variable's value. Cannot be used if value is not empty.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].env[index].valueFrom
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvindex)



Source for the environment variable's value. Cannot be used if value is not empty.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvindexvaluefromconfigmapkeyref">configMapKeyRef</a></b></td>
        <td>object</td>
        <td>
          Selects a key of a ConfigMap.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvindexvaluefromfieldref">fieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvindexvaluefromresourcefieldref">resourceFieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvindexvaluefromsecretkeyref">secretKeyRef</a></b></td>
        <td>object</td>
        <td>
          Selects a key of a secret in the pod's namespace<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].env[index].valueFrom.configMapKeyRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvindexvaluefrom)



Selects a key of a ConfigMap.

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key to select.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the ConfigMap or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].env[index].valueFrom.fieldRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvindexvaluefrom)



Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.

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
        <td><b>fieldPath</b></td>
        <td>string</td>
        <td>
          Path of the field to select in the specified API version.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          Version of the schema the FieldPath is written in terms of, defaults to "v1".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].env[index].valueFrom.resourceFieldRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvindexvaluefrom)



Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.

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
        <td><b>resource</b></td>
        <td>string</td>
        <td>
          Required: resource to select<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>containerName</b></td>
        <td>string</td>
        <td>
          Container name: required for volumes, optional for env vars<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>divisor</b></td>
        <td>int or string</td>
        <td>
          Specifies the output format of the exposed resources, defaults to "1"<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].env[index].valueFrom.secretKeyRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvindexvaluefrom)



Selects a key of a secret in the pod's namespace

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key of the secret to select from.  Must be a valid secret key.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the Secret or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].envFrom[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



EnvFromSource represents the source of a set of ConfigMaps or Secrets

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvfromindexconfigmapref">configMapRef</a></b></td>
        <td>object</td>
        <td>
          The ConfigMap to select from<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>prefix</b></td>
        <td>string</td>
        <td>
          Optional text to prepend to the name of each environment variable. Must be a C_IDENTIFIER.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvfromindexsecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          The Secret to select from<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].envFrom[index].configMapRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvfromindex)



The ConfigMap to select from

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the ConfigMap must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].envFrom[index].secretRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexenvfromindex)



The Secret to select from

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the Secret must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



Lifecycle is not allowed for ephemeral containers.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecyclepoststart">postStart</a></b></td>
        <td>object</td>
        <td>
          PostStart is called immediately after a container is created.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycleprestop">preStop</a></b></td>
        <td>object</td>
        <td>
          PreStop is called immediately before a container is terminated due to an
API request or management...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stopSignal</b></td>
        <td>string</td>
        <td>
          StopSignal defines which signal will be sent to a container when it is being stopped.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.postStart
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycle)



PostStart is called immediately after a container is created.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecyclepoststartexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecyclepoststarthttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecyclepoststartsleep">sleep</a></b></td>
        <td>object</td>
        <td>
          Sleep represents a duration that the container should sleep.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecyclepoststarttcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.postStart.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecyclepoststart)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.postStart.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecyclepoststart)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecyclepoststarthttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.postStart.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecyclepoststarthttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.postStart.sleep
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecyclepoststart)



Sleep represents a duration that the container should sleep.

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
        <td><b>seconds</b></td>
        <td>integer</td>
        <td>
          Seconds is the number of seconds to sleep.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.postStart.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecyclepoststart)



Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.preStop
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycle)



PreStop is called immediately before a container is terminated due to an
API request or management...

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycleprestopexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycleprestophttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycleprestopsleep">sleep</a></b></td>
        <td>object</td>
        <td>
          Sleep represents a duration that the container should sleep.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycleprestoptcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.preStop.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycleprestop)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.preStop.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycleprestop)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycleprestophttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.preStop.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycleprestophttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.preStop.sleep
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycleprestop)



Sleep represents a duration that the container should sleep.

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
        <td><b>seconds</b></td>
        <td>integer</td>
        <td>
          Seconds is the number of seconds to sleep.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].lifecycle.preStop.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlifecycleprestop)



Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].livenessProbe
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



Probes are not allowed for ephemeral containers.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlivenessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlivenessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlivenessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlivenessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].livenessProbe.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlivenessprobe)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].livenessProbe.grpc
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlivenessprobe)



GRPC specifies a GRPC HealthCheckRequest.

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
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].livenessProbe.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlivenessprobe)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlivenessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].livenessProbe.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlivenessprobehttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].livenessProbe.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexlivenessprobe)



TCPSocket specifies a connection to a TCP port.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].ports[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



ContainerPort represents a network port in a single container.

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
        <td><b>containerPort</b></td>
        <td>integer</td>
        <td>
          Number of port to expose on the pod's IP address.
This must be a valid port number, 0 < x < 65536.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>hostIP</b></td>
        <td>string</td>
        <td>
          What host IP to bind the external port to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostPort</b></td>
        <td>integer</td>
        <td>
          Number of port to expose on the host.
If specified, this must be a valid port number, 0 < x < 65536.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          If specified, this must be an IANA_SVC_NAME and unique within the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>protocol</b></td>
        <td>string</td>
        <td>
          Protocol for port. Must be UDP, TCP, or SCTP.
Defaults to "TCP".<br/>
          <br/>
            <i>Default</i>: TCP<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].readinessProbe
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



Probes are not allowed for ephemeral containers.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexreadinessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexreadinessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexreadinessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexreadinessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].readinessProbe.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexreadinessprobe)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].readinessProbe.grpc
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexreadinessprobe)



GRPC specifies a GRPC HealthCheckRequest.

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
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].readinessProbe.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexreadinessprobe)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexreadinessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].readinessProbe.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexreadinessprobehttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].readinessProbe.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexreadinessprobe)



TCPSocket specifies a connection to a TCP port.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].resizePolicy[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



ContainerResizePolicy represents resource resize policy for the container.

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
        <td><b>resourceName</b></td>
        <td>string</td>
        <td>
          Name of the resource to which this resource resize policy applies.
Supported values: cpu, memory.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>restartPolicy</b></td>
        <td>string</td>
        <td>
          Restart policy to apply when specified resource is resized.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].resources
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



Resources are not allowed for ephemeral containers.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexresourcesclaimsindex">claims</a></b></td>
        <td>[]object</td>
        <td>
          Claims lists the names of resources, defined in spec.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>limits</b></td>
        <td>map[string]int or string</td>
        <td>
          Limits describes the maximum amount of compute resources allowed.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>requests</b></td>
        <td>map[string]int or string</td>
        <td>
          Requests describes the minimum amount of compute resources required.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].resources.claims[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexresources)



ResourceClaim references one entry in PodSpec.ResourceClaims.

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
          Name must match the name of one entry in pod.spec.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>request</b></td>
        <td>string</td>
        <td>
          Request is the name chosen for a request in the referenced claim.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].securityContext
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



Optional: SecurityContext defines the security options the ephemeral container should be run with.

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
        <td><b>allowPrivilegeEscalation</b></td>
        <td>boolean</td>
        <td>
          AllowPrivilegeEscalation controls whether a process can gain more
privileges than its parent...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexsecuritycontextapparmorprofile">appArmorProfile</a></b></td>
        <td>object</td>
        <td>
          appArmorProfile is the AppArmor options to use by this container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexsecuritycontextcapabilities">capabilities</a></b></td>
        <td>object</td>
        <td>
          The capabilities to add/drop when running containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>privileged</b></td>
        <td>boolean</td>
        <td>
          Run container in privileged mode.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>procMount</b></td>
        <td>string</td>
        <td>
          procMount denotes the type of proc mount to use for the containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnlyRootFilesystem</b></td>
        <td>boolean</td>
        <td>
          Whether this container has a read-only root filesystem.
Default is false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsGroup</b></td>
        <td>integer</td>
        <td>
          The GID to run the entrypoint of the container process.
Uses runtime default if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsNonRoot</b></td>
        <td>boolean</td>
        <td>
          Indicates that the container must run as a non-root user.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUser</b></td>
        <td>integer</td>
        <td>
          The UID to run the entrypoint of the container process.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexsecuritycontextselinuxoptions">seLinuxOptions</a></b></td>
        <td>object</td>
        <td>
          The SELinux context to be applied to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexsecuritycontextseccompprofile">seccompProfile</a></b></td>
        <td>object</td>
        <td>
          The seccomp options to use by this container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexsecuritycontextwindowsoptions">windowsOptions</a></b></td>
        <td>object</td>
        <td>
          The Windows specific settings applied to all containers.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].securityContext.appArmorProfile
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexsecuritycontext)



appArmorProfile is the AppArmor options to use by this container.

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
          type indicates which kind of AppArmor profile will be applied.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile loaded on the node that should be used.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].securityContext.capabilities
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexsecuritycontext)



The capabilities to add/drop when running containers.

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
        <td><b>add</b></td>
        <td>[]string</td>
        <td>
          Added capabilities<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>drop</b></td>
        <td>[]string</td>
        <td>
          Removed capabilities<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].securityContext.seLinuxOptions
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexsecuritycontext)



The SELinux context to be applied to the container.

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
        <td><b>level</b></td>
        <td>string</td>
        <td>
          Level is SELinux level label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>role</b></td>
        <td>string</td>
        <td>
          Role is a SELinux role label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type is a SELinux type label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          User is a SELinux user label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].securityContext.seccompProfile
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexsecuritycontext)



The seccomp options to use by this container.

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
          type indicates which kind of seccomp profile will be applied.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile defined in a file on the node should be used.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].securityContext.windowsOptions
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexsecuritycontext)



The Windows specific settings applied to all containers.

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
        <td><b>gmsaCredentialSpec</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpec is where the GMSA admission webhook
(https://github.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gmsaCredentialSpecName</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpecName is the name of the GMSA credential spec to use.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostProcess</b></td>
        <td>boolean</td>
        <td>
          HostProcess determines if a container should be run as a 'Host Process' container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUserName</b></td>
        <td>string</td>
        <td>
          The UserName in Windows to run the entrypoint of the container process.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].startupProbe
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



Probes are not allowed for ephemeral containers.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexstartupprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexstartupprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexstartupprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexstartupprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].startupProbe.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexstartupprobe)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].startupProbe.grpc
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexstartupprobe)



GRPC specifies a GRPC HealthCheckRequest.

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
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].startupProbe.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexstartupprobe)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexstartupprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].startupProbe.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexstartupprobehttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].startupProbe.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindexstartupprobe)



TCPSocket specifies a connection to a TCP port.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].volumeDevices[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



volumeDevice describes a mapping of a raw block device within a container.

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
        <td><b>devicePath</b></td>
        <td>string</td>
        <td>
          devicePath is the path inside of the container that the device will be mapped to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          name must match the name of a persistentVolumeClaim in the pod<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.ephemeralContainers[index].volumeMounts[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecephemeralcontainersindex)



VolumeMount describes a mounting of a Volume within a container.

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
        <td><b>mountPath</b></td>
        <td>string</td>
        <td>
          Path within the container at which the volume should be mounted.  Must
not contain ':'.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          This must match the Name of a Volume.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mountPropagation</b></td>
        <td>string</td>
        <td>
          mountPropagation determines how mounts are propagated from the host
to container and the other way...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          Mounted read-only if true, read-write otherwise (false or unspecified).
Defaults to false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>recursiveReadOnly</b></td>
        <td>string</td>
        <td>
          RecursiveReadOnly specifies whether read-only mounts should be handled
recursively.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subPath</b></td>
        <td>string</td>
        <td>
          Path within the volume from which the container's volume should be mounted.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subPathExpr</b></td>
        <td>string</td>
        <td>
          Expanded path within the volume from which the container's volume should be mounted.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.hostAliases[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



HostAlias holds the mapping between IP and hostnames that will be injected as an entry in the
pod's...

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
        <td><b>ip</b></td>
        <td>string</td>
        <td>
          IP address of the host file entry.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>hostnames</b></td>
        <td>[]string</td>
        <td>
          Hostnames for the above IP address.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.imagePullSecrets[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



LocalObjectReference contains enough information to let you locate the
referenced object inside the...

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



A single application container that you want to run within a pod.

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
          Name of the container specified as a DNS_LABEL.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>args</b></td>
        <td>[]string</td>
        <td>
          Arguments to the entrypoint.
The container image's CMD is used if this is not provided.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Entrypoint array. Not executed within a shell.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvindex">env</a></b></td>
        <td>[]object</td>
        <td>
          List of environment variables to set in the container.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvfromindex">envFrom</a></b></td>
        <td>[]object</td>
        <td>
          List of sources to populate environment variables in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>image</b></td>
        <td>string</td>
        <td>
          Container image name.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imagePullPolicy</b></td>
        <td>string</td>
        <td>
          Image pull policy.
One of Always, Never, IfNotPresent.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycle">lifecycle</a></b></td>
        <td>object</td>
        <td>
          Actions that the management system should take in response to container lifecycle events.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlivenessprobe">livenessProbe</a></b></td>
        <td>object</td>
        <td>
          Periodic probe of container liveness.
Container will be restarted if the probe fails.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexportsindex">ports</a></b></td>
        <td>[]object</td>
        <td>
          List of ports to expose from the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexreadinessprobe">readinessProbe</a></b></td>
        <td>object</td>
        <td>
          Periodic probe of container service readiness.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexresizepolicyindex">resizePolicy</a></b></td>
        <td>[]object</td>
        <td>
          Resources resize policy for the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexresources">resources</a></b></td>
        <td>object</td>
        <td>
          Compute Resources required by this container.
Cannot be updated.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>restartPolicy</b></td>
        <td>string</td>
        <td>
          RestartPolicy defines the restart behavior of individual containers in a pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexsecuritycontext">securityContext</a></b></td>
        <td>object</td>
        <td>
          SecurityContext defines the security options the container should be run with.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexstartupprobe">startupProbe</a></b></td>
        <td>object</td>
        <td>
          StartupProbe indicates that the Pod has successfully initialized.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stdin</b></td>
        <td>boolean</td>
        <td>
          Whether this container should allocate a buffer for stdin in the container runtime.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stdinOnce</b></td>
        <td>boolean</td>
        <td>
          Whether the container runtime should close the stdin channel after it has been opened by
a single...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationMessagePath</b></td>
        <td>string</td>
        <td>
          Optional: Path at which the file to which the container's termination message
will be written is...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationMessagePolicy</b></td>
        <td>string</td>
        <td>
          Indicate how the termination message should be populated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tty</b></td>
        <td>boolean</td>
        <td>
          Whether this container should allocate a TTY for itself, also requires 'stdin' to be true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexvolumedevicesindex">volumeDevices</a></b></td>
        <td>[]object</td>
        <td>
          volumeDevices is the list of block devices to be used by the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexvolumemountsindex">volumeMounts</a></b></td>
        <td>[]object</td>
        <td>
          Pod volumes to mount into the container's filesystem.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>workingDir</b></td>
        <td>string</td>
        <td>
          Container's working directory.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].env[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



EnvVar represents an environment variable present in a Container.

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
          Name of the environment variable. Must be a C_IDENTIFIER.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Variable references $(VAR_NAME) are expanded
using the previously defined environment variables in...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvindexvaluefrom">valueFrom</a></b></td>
        <td>object</td>
        <td>
          Source for the environment variable's value. Cannot be used if value is not empty.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].env[index].valueFrom
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvindex)



Source for the environment variable's value. Cannot be used if value is not empty.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvindexvaluefromconfigmapkeyref">configMapKeyRef</a></b></td>
        <td>object</td>
        <td>
          Selects a key of a ConfigMap.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvindexvaluefromfieldref">fieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvindexvaluefromresourcefieldref">resourceFieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvindexvaluefromsecretkeyref">secretKeyRef</a></b></td>
        <td>object</td>
        <td>
          Selects a key of a secret in the pod's namespace<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].env[index].valueFrom.configMapKeyRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvindexvaluefrom)



Selects a key of a ConfigMap.

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key to select.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the ConfigMap or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].env[index].valueFrom.fieldRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvindexvaluefrom)



Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.

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
        <td><b>fieldPath</b></td>
        <td>string</td>
        <td>
          Path of the field to select in the specified API version.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          Version of the schema the FieldPath is written in terms of, defaults to "v1".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].env[index].valueFrom.resourceFieldRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvindexvaluefrom)



Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.

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
        <td><b>resource</b></td>
        <td>string</td>
        <td>
          Required: resource to select<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>containerName</b></td>
        <td>string</td>
        <td>
          Container name: required for volumes, optional for env vars<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>divisor</b></td>
        <td>int or string</td>
        <td>
          Specifies the output format of the exposed resources, defaults to "1"<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].env[index].valueFrom.secretKeyRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvindexvaluefrom)



Selects a key of a secret in the pod's namespace

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key of the secret to select from.  Must be a valid secret key.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the Secret or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].envFrom[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



EnvFromSource represents the source of a set of ConfigMaps or Secrets

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvfromindexconfigmapref">configMapRef</a></b></td>
        <td>object</td>
        <td>
          The ConfigMap to select from<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>prefix</b></td>
        <td>string</td>
        <td>
          Optional text to prepend to the name of each environment variable. Must be a C_IDENTIFIER.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvfromindexsecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          The Secret to select from<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].envFrom[index].configMapRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvfromindex)



The ConfigMap to select from

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the ConfigMap must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].envFrom[index].secretRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexenvfromindex)



The Secret to select from

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the Secret must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



Actions that the management system should take in response to container lifecycle events.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecyclepoststart">postStart</a></b></td>
        <td>object</td>
        <td>
          PostStart is called immediately after a container is created.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycleprestop">preStop</a></b></td>
        <td>object</td>
        <td>
          PreStop is called immediately before a container is terminated due to an
API request or management...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stopSignal</b></td>
        <td>string</td>
        <td>
          StopSignal defines which signal will be sent to a container when it is being stopped.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.postStart
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycle)



PostStart is called immediately after a container is created.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecyclepoststartexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecyclepoststarthttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecyclepoststartsleep">sleep</a></b></td>
        <td>object</td>
        <td>
          Sleep represents a duration that the container should sleep.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecyclepoststarttcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.postStart.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecyclepoststart)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.postStart.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecyclepoststart)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecyclepoststarthttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.postStart.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecyclepoststarthttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.postStart.sleep
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecyclepoststart)



Sleep represents a duration that the container should sleep.

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
        <td><b>seconds</b></td>
        <td>integer</td>
        <td>
          Seconds is the number of seconds to sleep.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.postStart.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecyclepoststart)



Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.preStop
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycle)



PreStop is called immediately before a container is terminated due to an
API request or management...

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycleprestopexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycleprestophttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycleprestopsleep">sleep</a></b></td>
        <td>object</td>
        <td>
          Sleep represents a duration that the container should sleep.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycleprestoptcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.preStop.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycleprestop)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.preStop.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycleprestop)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycleprestophttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.preStop.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycleprestophttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.preStop.sleep
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycleprestop)



Sleep represents a duration that the container should sleep.

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
        <td><b>seconds</b></td>
        <td>integer</td>
        <td>
          Seconds is the number of seconds to sleep.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].lifecycle.preStop.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlifecycleprestop)



Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].livenessProbe
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



Periodic probe of container liveness.
Container will be restarted if the probe fails.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlivenessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlivenessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlivenessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlivenessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].livenessProbe.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlivenessprobe)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].livenessProbe.grpc
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlivenessprobe)



GRPC specifies a GRPC HealthCheckRequest.

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
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].livenessProbe.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlivenessprobe)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlivenessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].livenessProbe.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlivenessprobehttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].livenessProbe.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexlivenessprobe)



TCPSocket specifies a connection to a TCP port.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].ports[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



ContainerPort represents a network port in a single container.

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
        <td><b>containerPort</b></td>
        <td>integer</td>
        <td>
          Number of port to expose on the pod's IP address.
This must be a valid port number, 0 < x < 65536.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>hostIP</b></td>
        <td>string</td>
        <td>
          What host IP to bind the external port to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostPort</b></td>
        <td>integer</td>
        <td>
          Number of port to expose on the host.
If specified, this must be a valid port number, 0 < x < 65536.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          If specified, this must be an IANA_SVC_NAME and unique within the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>protocol</b></td>
        <td>string</td>
        <td>
          Protocol for port. Must be UDP, TCP, or SCTP.
Defaults to "TCP".<br/>
          <br/>
            <i>Default</i>: TCP<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].readinessProbe
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



Periodic probe of container service readiness.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexreadinessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexreadinessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexreadinessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexreadinessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].readinessProbe.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexreadinessprobe)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].readinessProbe.grpc
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexreadinessprobe)



GRPC specifies a GRPC HealthCheckRequest.

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
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].readinessProbe.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexreadinessprobe)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexreadinessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].readinessProbe.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexreadinessprobehttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].readinessProbe.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexreadinessprobe)



TCPSocket specifies a connection to a TCP port.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].resizePolicy[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



ContainerResizePolicy represents resource resize policy for the container.

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
        <td><b>resourceName</b></td>
        <td>string</td>
        <td>
          Name of the resource to which this resource resize policy applies.
Supported values: cpu, memory.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>restartPolicy</b></td>
        <td>string</td>
        <td>
          Restart policy to apply when specified resource is resized.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].resources
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



Compute Resources required by this container.
Cannot be updated.
More info: https://kubernetes.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexresourcesclaimsindex">claims</a></b></td>
        <td>[]object</td>
        <td>
          Claims lists the names of resources, defined in spec.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>limits</b></td>
        <td>map[string]int or string</td>
        <td>
          Limits describes the maximum amount of compute resources allowed.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>requests</b></td>
        <td>map[string]int or string</td>
        <td>
          Requests describes the minimum amount of compute resources required.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].resources.claims[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexresources)



ResourceClaim references one entry in PodSpec.ResourceClaims.

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
          Name must match the name of one entry in pod.spec.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>request</b></td>
        <td>string</td>
        <td>
          Request is the name chosen for a request in the referenced claim.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].securityContext
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



SecurityContext defines the security options the container should be run with.

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
        <td><b>allowPrivilegeEscalation</b></td>
        <td>boolean</td>
        <td>
          AllowPrivilegeEscalation controls whether a process can gain more
privileges than its parent...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexsecuritycontextapparmorprofile">appArmorProfile</a></b></td>
        <td>object</td>
        <td>
          appArmorProfile is the AppArmor options to use by this container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexsecuritycontextcapabilities">capabilities</a></b></td>
        <td>object</td>
        <td>
          The capabilities to add/drop when running containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>privileged</b></td>
        <td>boolean</td>
        <td>
          Run container in privileged mode.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>procMount</b></td>
        <td>string</td>
        <td>
          procMount denotes the type of proc mount to use for the containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnlyRootFilesystem</b></td>
        <td>boolean</td>
        <td>
          Whether this container has a read-only root filesystem.
Default is false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsGroup</b></td>
        <td>integer</td>
        <td>
          The GID to run the entrypoint of the container process.
Uses runtime default if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsNonRoot</b></td>
        <td>boolean</td>
        <td>
          Indicates that the container must run as a non-root user.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUser</b></td>
        <td>integer</td>
        <td>
          The UID to run the entrypoint of the container process.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexsecuritycontextselinuxoptions">seLinuxOptions</a></b></td>
        <td>object</td>
        <td>
          The SELinux context to be applied to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexsecuritycontextseccompprofile">seccompProfile</a></b></td>
        <td>object</td>
        <td>
          The seccomp options to use by this container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexsecuritycontextwindowsoptions">windowsOptions</a></b></td>
        <td>object</td>
        <td>
          The Windows specific settings applied to all containers.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].securityContext.appArmorProfile
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexsecuritycontext)



appArmorProfile is the AppArmor options to use by this container.

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
          type indicates which kind of AppArmor profile will be applied.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile loaded on the node that should be used.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].securityContext.capabilities
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexsecuritycontext)



The capabilities to add/drop when running containers.

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
        <td><b>add</b></td>
        <td>[]string</td>
        <td>
          Added capabilities<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>drop</b></td>
        <td>[]string</td>
        <td>
          Removed capabilities<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].securityContext.seLinuxOptions
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexsecuritycontext)



The SELinux context to be applied to the container.

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
        <td><b>level</b></td>
        <td>string</td>
        <td>
          Level is SELinux level label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>role</b></td>
        <td>string</td>
        <td>
          Role is a SELinux role label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type is a SELinux type label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          User is a SELinux user label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].securityContext.seccompProfile
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexsecuritycontext)



The seccomp options to use by this container.

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
          type indicates which kind of seccomp profile will be applied.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile defined in a file on the node should be used.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].securityContext.windowsOptions
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexsecuritycontext)



The Windows specific settings applied to all containers.

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
        <td><b>gmsaCredentialSpec</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpec is where the GMSA admission webhook
(https://github.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gmsaCredentialSpecName</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpecName is the name of the GMSA credential spec to use.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostProcess</b></td>
        <td>boolean</td>
        <td>
          HostProcess determines if a container should be run as a 'Host Process' container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUserName</b></td>
        <td>string</td>
        <td>
          The UserName in Windows to run the entrypoint of the container process.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].startupProbe
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



StartupProbe indicates that the Pod has successfully initialized.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexstartupprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexstartupprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexstartupprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexstartupprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].startupProbe.exec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexstartupprobe)



Exec specifies a command to execute in the container.

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
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command ...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].startupProbe.grpc
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexstartupprobe)



GRPC specifies a GRPC HealthCheckRequest.

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
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].startupProbe.httpGet
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexstartupprobe)



HTTPGet specifies an HTTP GET request to perform.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexstartupprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].startupProbe.httpGet.httpHeaders[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexstartupprobehttpget)



HTTPHeader describes a custom header to be used in HTTP probes

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
          The header field name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].startupProbe.tcpSocket
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindexstartupprobe)



TCPSocket specifies a connection to a TCP port.

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
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].volumeDevices[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



volumeDevice describes a mapping of a raw block device within a container.

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
        <td><b>devicePath</b></td>
        <td>string</td>
        <td>
          devicePath is the path inside of the container that the device will be mapped to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          name must match the name of a persistentVolumeClaim in the pod<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.initContainers[index].volumeMounts[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecinitcontainersindex)



VolumeMount describes a mounting of a Volume within a container.

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
        <td><b>mountPath</b></td>
        <td>string</td>
        <td>
          Path within the container at which the volume should be mounted.  Must
not contain ':'.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          This must match the Name of a Volume.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mountPropagation</b></td>
        <td>string</td>
        <td>
          mountPropagation determines how mounts are propagated from the host
to container and the other way...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          Mounted read-only if true, read-write otherwise (false or unspecified).
Defaults to false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>recursiveReadOnly</b></td>
        <td>string</td>
        <td>
          RecursiveReadOnly specifies whether read-only mounts should be handled
recursively.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subPath</b></td>
        <td>string</td>
        <td>
          Path within the volume from which the container's volume should be mounted.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subPathExpr</b></td>
        <td>string</td>
        <td>
          Expanded path within the volume from which the container's volume should be mounted.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.os
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



Specifies the OS of the containers in the pod.

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
          Name is the name of the operating system. The currently supported values are linux and windows.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.readinessGates[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



PodReadinessGate contains the reference to a pod condition

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
        <td><b>conditionType</b></td>
        <td>string</td>
        <td>
          ConditionType refers to a condition in the pod's condition list with matching type.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.resourceClaims[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



PodResourceClaim references exactly one ResourceClaim, either directly
or by naming a...

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
          Name uniquely identifies this resource claim inside the pod.
This must be a DNS_LABEL.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>resourceClaimName</b></td>
        <td>string</td>
        <td>
          ResourceClaimName is the name of a ResourceClaim object in the same
namespace as this pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>resourceClaimTemplateName</b></td>
        <td>string</td>
        <td>
          ResourceClaimTemplateName is the name of a ResourceClaimTemplate
object in the same namespace as...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.resources
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



Resources is the total amount of CPU and Memory resources required by all
containers in the pod.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecresourcesclaimsindex">claims</a></b></td>
        <td>[]object</td>
        <td>
          Claims lists the names of resources, defined in spec.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>limits</b></td>
        <td>map[string]int or string</td>
        <td>
          Limits describes the maximum amount of compute resources allowed.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>requests</b></td>
        <td>map[string]int or string</td>
        <td>
          Requests describes the minimum amount of compute resources required.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.resources.claims[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecresources)



ResourceClaim references one entry in PodSpec.ResourceClaims.

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
          Name must match the name of one entry in pod.spec.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>request</b></td>
        <td>string</td>
        <td>
          Request is the name chosen for a request in the referenced claim.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.schedulingGates[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



PodSchedulingGate is associated to a Pod to guard its scheduling.

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
          Name of the scheduling gate.
Each scheduling gate must have a unique name field.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.securityContext
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



SecurityContext holds pod-level security attributes and common container settings.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecsecuritycontextapparmorprofile">appArmorProfile</a></b></td>
        <td>object</td>
        <td>
          appArmorProfile is the AppArmor options to use by the containers in this pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>fsGroup</b></td>
        <td>integer</td>
        <td>
          A special supplemental group that applies to all containers in a pod.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>fsGroupChangePolicy</b></td>
        <td>string</td>
        <td>
          fsGroupChangePolicy defines behavior of changing ownership and permission of the volume
before...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsGroup</b></td>
        <td>integer</td>
        <td>
          The GID to run the entrypoint of the container process.
Uses runtime default if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsNonRoot</b></td>
        <td>boolean</td>
        <td>
          Indicates that the container must run as a non-root user.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUser</b></td>
        <td>integer</td>
        <td>
          The UID to run the entrypoint of the container process.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>seLinuxChangePolicy</b></td>
        <td>string</td>
        <td>
          seLinuxChangePolicy defines how the container's SELinux label is applied to all volumes used by the...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecsecuritycontextselinuxoptions">seLinuxOptions</a></b></td>
        <td>object</td>
        <td>
          The SELinux context to be applied to all containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecsecuritycontextseccompprofile">seccompProfile</a></b></td>
        <td>object</td>
        <td>
          The seccomp options to use by the containers in this pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>supplementalGroups</b></td>
        <td>[]integer</td>
        <td>
          A list of groups applied to the first process run in each container, in
addition to the container's...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>supplementalGroupsPolicy</b></td>
        <td>string</td>
        <td>
          Defines how supplemental groups of the first container processes are calculated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecsecuritycontextsysctlsindex">sysctls</a></b></td>
        <td>[]object</td>
        <td>
          Sysctls hold a list of namespaced sysctls used for the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecsecuritycontextwindowsoptions">windowsOptions</a></b></td>
        <td>object</td>
        <td>
          The Windows specific settings applied to all containers.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.securityContext.appArmorProfile
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecsecuritycontext)



appArmorProfile is the AppArmor options to use by the containers in this pod.

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
          type indicates which kind of AppArmor profile will be applied.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile loaded on the node that should be used.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.securityContext.seLinuxOptions
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecsecuritycontext)



The SELinux context to be applied to all containers.

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
        <td><b>level</b></td>
        <td>string</td>
        <td>
          Level is SELinux level label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>role</b></td>
        <td>string</td>
        <td>
          Role is a SELinux role label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type is a SELinux type label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          User is a SELinux user label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.securityContext.seccompProfile
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecsecuritycontext)



The seccomp options to use by the containers in this pod.

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
          type indicates which kind of seccomp profile will be applied.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile defined in a file on the node should be used.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.securityContext.sysctls[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecsecuritycontext)



Sysctl defines a kernel parameter to be set

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
          Name of a property to set<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Value of a property to set<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.securityContext.windowsOptions
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecsecuritycontext)



The Windows specific settings applied to all containers.

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
        <td><b>gmsaCredentialSpec</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpec is where the GMSA admission webhook
(https://github.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gmsaCredentialSpecName</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpecName is the name of the GMSA credential spec to use.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostProcess</b></td>
        <td>boolean</td>
        <td>
          HostProcess determines if a container should be run as a 'Host Process' container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUserName</b></td>
        <td>string</td>
        <td>
          The UserName in Windows to run the entrypoint of the container process.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.tolerations[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



The pod this Toleration is attached to tolerates any taint that matches
the triple...

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
        <td><b>effect</b></td>
        <td>string</td>
        <td>
          Effect indicates the taint effect to match. Empty means match all taint effects.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Key is the taint key that the toleration applies to. Empty means match all taint keys.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Operator represents a key's relationship to the value.
Valid operators are Exists and Equal.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tolerationSeconds</b></td>
        <td>integer</td>
        <td>
          TolerationSeconds represents the period of time the toleration (which must be
of effect NoExecute,...<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Value is the taint value the toleration matches to.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.topologySpreadConstraints[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



TopologySpreadConstraint specifies how to spread matching pods among the given topology.

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
        <td><b>maxSkew</b></td>
        <td>integer</td>
        <td>
          MaxSkew describes the degree to which pods may be unevenly distributed.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>topologyKey</b></td>
        <td>string</td>
        <td>
          TopologyKey is the key of node labels.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>whenUnsatisfiable</b></td>
        <td>string</td>
        <td>
          WhenUnsatisfiable indicates how to deal with a pod if it doesn't satisfy
the spread constraint.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespectopologyspreadconstraintsindexlabelselector">labelSelector</a></b></td>
        <td>object</td>
        <td>
          LabelSelector is used to find matching pods.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabelKeys</b></td>
        <td>[]string</td>
        <td>
          MatchLabelKeys is a set of pod label keys to select the pods over which
spreading will be...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>minDomains</b></td>
        <td>integer</td>
        <td>
          MinDomains indicates a minimum number of eligible domains.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>nodeAffinityPolicy</b></td>
        <td>string</td>
        <td>
          NodeAffinityPolicy indicates how we will treat Pod's nodeAffinity/nodeSelector
when calculating pod...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>nodeTaintsPolicy</b></td>
        <td>string</td>
        <td>
          NodeTaintsPolicy indicates how we will treat node taints when calculating
pod topology spread skew.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.topologySpreadConstraints[index].labelSelector
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespectopologyspreadconstraintsindex)



LabelSelector is used to find matching pods.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespectopologyspreadconstraintsindexlabelselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.topologySpreadConstraints[index].labelSelector.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespectopologyspreadconstraintsindexlabelselector)



A label selector requirement is a selector that contains values, a key, and an operator that...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespec)



Volume represents a named volume in a pod that may be accessed by any container in the pod.

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
          name of the volume.
Must be a DNS_LABEL and unique within the pod.
More info: https://kubernetes.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexawselasticblockstore">awsElasticBlockStore</a></b></td>
        <td>object</td>
        <td>
          awsElasticBlockStore represents an AWS Disk resource that is attached to a
kubelet's host machine...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexazuredisk">azureDisk</a></b></td>
        <td>object</td>
        <td>
          azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexazurefile">azureFile</a></b></td>
        <td>object</td>
        <td>
          azureFile represents an Azure File Service mount on the host and bind mount to the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexcephfs">cephfs</a></b></td>
        <td>object</td>
        <td>
          cephFS represents a Ceph FS mount on the host that shares a pod's lifetime.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexcinder">cinder</a></b></td>
        <td>object</td>
        <td>
          cinder represents a cinder volume attached and mounted on kubelets host machine.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexconfigmap">configMap</a></b></td>
        <td>object</td>
        <td>
          configMap represents a configMap that should populate this volume<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexcsi">csi</a></b></td>
        <td>object</td>
        <td>
          csi (Container Storage Interface) represents ephemeral storage that is handled by certain external...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexdownwardapi">downwardAPI</a></b></td>
        <td>object</td>
        <td>
          downwardAPI represents downward API about the pod that should populate this volume<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexemptydir">emptyDir</a></b></td>
        <td>object</td>
        <td>
          emptyDir represents a temporary directory that shares a pod's lifetime.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeral">ephemeral</a></b></td>
        <td>object</td>
        <td>
          ephemeral represents a volume that is handled by a cluster storage driver.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexfc">fc</a></b></td>
        <td>object</td>
        <td>
          fc represents a Fibre Channel resource that is attached to a kubelet's host machine and then...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexflexvolume">flexVolume</a></b></td>
        <td>object</td>
        <td>
          flexVolume represents a generic volume resource that is
provisioned/attached using an exec based...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexflocker">flocker</a></b></td>
        <td>object</td>
        <td>
          flocker represents a Flocker volume attached to a kubelet's host machine.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexgcepersistentdisk">gcePersistentDisk</a></b></td>
        <td>object</td>
        <td>
          gcePersistentDisk represents a GCE Disk resource that is attached to a
kubelet's host machine and...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexgitrepo">gitRepo</a></b></td>
        <td>object</td>
        <td>
          gitRepo represents a git repository at a particular revision.
Deprecated: GitRepo is deprecated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexglusterfs">glusterfs</a></b></td>
        <td>object</td>
        <td>
          glusterfs represents a Glusterfs mount on the host that shares a pod's lifetime.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexhostpath">hostPath</a></b></td>
        <td>object</td>
        <td>
          hostPath represents a pre-existing file or directory on the host
machine that is directly exposed...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindeximage">image</a></b></td>
        <td>object</td>
        <td>
          image represents an OCI object (a container image or artifact) pulled and mounted on the kubelet's...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexiscsi">iscsi</a></b></td>
        <td>object</td>
        <td>
          iscsi represents an ISCSI Disk resource that is attached to a
kubelet's host machine and then...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexnfs">nfs</a></b></td>
        <td>object</td>
        <td>
          nfs represents an NFS mount on the host that shares a pod's lifetime
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexpersistentvolumeclaim">persistentVolumeClaim</a></b></td>
        <td>object</td>
        <td>
          persistentVolumeClaimVolumeSource represents a reference to a
PersistentVolumeClaim in the same...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexphotonpersistentdisk">photonPersistentDisk</a></b></td>
        <td>object</td>
        <td>
          photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexportworxvolume">portworxVolume</a></b></td>
        <td>object</td>
        <td>
          portworxVolume represents a portworx volume attached and mounted on kubelets host machine.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojected">projected</a></b></td>
        <td>object</td>
        <td>
          projected items for all in one resources secrets, configmaps, and downward API<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexquobyte">quobyte</a></b></td>
        <td>object</td>
        <td>
          quobyte represents a Quobyte mount on the host that shares a pod's lifetime.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexrbd">rbd</a></b></td>
        <td>object</td>
        <td>
          rbd represents a Rados Block Device mount on the host that shares a pod's lifetime.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexscaleio">scaleIO</a></b></td>
        <td>object</td>
        <td>
          scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexsecret">secret</a></b></td>
        <td>object</td>
        <td>
          secret represents a secret that should populate this volume.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexstorageos">storageos</a></b></td>
        <td>object</td>
        <td>
          storageOS represents a StorageOS volume attached and mounted on Kubernetes nodes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexvspherevolume">vsphereVolume</a></b></td>
        <td>object</td>
        <td>
          vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].awsElasticBlockStore
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



awsElasticBlockStore represents an AWS Disk resource that is attached to a
kubelet's host machine...

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
        <td><b>volumeID</b></td>
        <td>string</td>
        <td>
          volumeID is unique ID of the persistent disk resource in AWS (Amazon EBS volume).<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type of the volume that you want to mount.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>partition</b></td>
        <td>integer</td>
        <td>
          partition is the partition in the volume that you want to mount.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly value true will force the readOnly setting in VolumeMounts.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].azureDisk
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.

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
        <td><b>diskName</b></td>
        <td>string</td>
        <td>
          diskName is the Name of the data disk in the blob storage<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>diskURI</b></td>
        <td>string</td>
        <td>
          diskURI is the URI of data disk in the blob storage<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>cachingMode</b></td>
        <td>string</td>
        <td>
          cachingMode is the Host Caching mode: None, Read Only, Read Write.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is Filesystem type to mount.<br/>
          <br/>
            <i>Default</i>: ext4<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          kind expected values are Shared: multiple blob disks per storage account  Dedicated: single blob...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly Defaults to false (read/write).<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].azureFile
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



azureFile represents an Azure File Service mount on the host and bind mount to the pod.

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
        <td><b>secretName</b></td>
        <td>string</td>
        <td>
          secretName is the  name of secret that contains Azure Storage Account Name and Key<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>shareName</b></td>
        <td>string</td>
        <td>
          shareName is the azure share Name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly defaults to false (read/write).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].cephfs
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



cephFS represents a Ceph FS mount on the host that shares a pod's lifetime.

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
        <td><b>monitors</b></td>
        <td>[]string</td>
        <td>
          monitors is Required: Monitors is a collection of Ceph monitors
More info: https://examples.k8s.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is Optional: Used as the mounted root, rather than the full Ceph tree, default is /<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly is Optional: Defaults to false (read/write).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>secretFile</b></td>
        <td>string</td>
        <td>
          secretFile is Optional: SecretFile is the path to key ring for User, default is /etc/ceph/user.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexcephfssecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef is Optional: SecretRef is reference to the authentication secret for User, default is...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          user is optional: User is the rados user name, default is admin
More info: https://examples.k8s.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].cephfs.secretRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexcephfs)



secretRef is Optional: SecretRef is reference to the authentication secret for User, default is...

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].cinder
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



cinder represents a cinder volume attached and mounted on kubelets host machine.

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
        <td><b>volumeID</b></td>
        <td>string</td>
        <td>
          volumeID used to identify the volume in cinder.
More info: https://examples.k8s.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly defaults to false (read/write).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexcindersecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef is optional: points to a secret object containing parameters used to connect
to OpenStack.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].cinder.secretRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexcinder)



secretRef is optional: points to a secret object containing parameters used to connect
to OpenStack.

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].configMap
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



configMap represents a configMap that should populate this volume

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
        <td><b>defaultMode</b></td>
        <td>integer</td>
        <td>
          defaultMode is optional: mode bits used to set permissions on created files by default.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexconfigmapitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          items if unspecified, each key-value pair in the Data field of the referenced
ConfigMap will be...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          optional specify whether the ConfigMap or its keys must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].configMap.items[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexconfigmap)



Maps a string key to a path within a volume.

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the key to project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is the relative path of the file to map the key to.
May not be an absolute path.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          mode is Optional: mode bits used to set permissions on this file.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].csi
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



csi (Container Storage Interface) represents ephemeral storage that is handled by certain external...

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
        <td><b>driver</b></td>
        <td>string</td>
        <td>
          driver is the name of the CSI driver that handles this volume.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType to mount. Ex. "ext4", "xfs", "ntfs".<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexcsinodepublishsecretref">nodePublishSecretRef</a></b></td>
        <td>object</td>
        <td>
          nodePublishSecretRef is a reference to the secret object containing
sensitive information to pass...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly specifies a read-only configuration for the volume.
Defaults to false (read/write).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeAttributes</b></td>
        <td>map[string]string</td>
        <td>
          volumeAttributes stores driver-specific properties that are passed to the CSI
driver.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].csi.nodePublishSecretRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexcsi)



nodePublishSecretRef is a reference to the secret object containing
sensitive information to pass...

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].downwardAPI
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



downwardAPI represents downward API about the pod that should populate this volume

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
        <td><b>defaultMode</b></td>
        <td>integer</td>
        <td>
          Optional: mode bits to use on created files by default.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexdownwardapiitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          Items is a list of downward API volume file<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].downwardAPI.items[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexdownwardapi)



DownwardAPIVolumeFile represents information to create the file containing the pod field

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
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Required: Path is  the relative path name of the file to be created.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexdownwardapiitemsindexfieldref">fieldRef</a></b></td>
        <td>object</td>
        <td>
          Required: Selects a field of the pod: only annotations, labels, name, namespace and uid are...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          Optional: mode bits used to set permissions on this file, must be an octal value
between 0000 and...<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexdownwardapiitemsindexresourcefieldref">resourceFieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].downwardAPI.items[index].fieldRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexdownwardapiitemsindex)



Required: Selects a field of the pod: only annotations, labels, name, namespace and uid are...

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
        <td><b>fieldPath</b></td>
        <td>string</td>
        <td>
          Path of the field to select in the specified API version.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          Version of the schema the FieldPath is written in terms of, defaults to "v1".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].downwardAPI.items[index].resourceFieldRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexdownwardapiitemsindex)



Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.

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
        <td><b>resource</b></td>
        <td>string</td>
        <td>
          Required: resource to select<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>containerName</b></td>
        <td>string</td>
        <td>
          Container name: required for volumes, optional for env vars<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>divisor</b></td>
        <td>int or string</td>
        <td>
          Specifies the output format of the exposed resources, defaults to "1"<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].emptyDir
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



emptyDir represents a temporary directory that shares a pod's lifetime.

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
        <td><b>medium</b></td>
        <td>string</td>
        <td>
          medium represents what type of storage medium should back this directory.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sizeLimit</b></td>
        <td>int or string</td>
        <td>
          sizeLimit is the total amount of local storage required for this EmptyDir volume.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].ephemeral
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



ephemeral represents a volume that is handled by a cluster storage driver.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplate">volumeClaimTemplate</a></b></td>
        <td>object</td>
        <td>
          Will be used to create a stand-alone PVC to provision the volume.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].ephemeral.volumeClaimTemplate
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeral)



Will be used to create a stand-alone PVC to provision the volume.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplatespec">spec</a></b></td>
        <td>object</td>
        <td>
          The specification for the PersistentVolumeClaim.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>metadata</b></td>
        <td>object</td>
        <td>
          May contain labels and annotations that will be copied into the PVC
when creating it.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].ephemeral.volumeClaimTemplate.spec
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplate)



The specification for the PersistentVolumeClaim.

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
        <td><b>accessModes</b></td>
        <td>[]string</td>
        <td>
          accessModes contains the desired access modes the volume should have.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplatespecdatasource">dataSource</a></b></td>
        <td>object</td>
        <td>
          dataSource field can be used to specify either:
* An existing VolumeSnapshot object (snapshot.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplatespecdatasourceref">dataSourceRef</a></b></td>
        <td>object</td>
        <td>
          dataSourceRef specifies the object from which to populate the volume with data, if a non-empty...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplatespecresources">resources</a></b></td>
        <td>object</td>
        <td>
          resources represents the minimum resources the volume should have.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplatespecselector">selector</a></b></td>
        <td>object</td>
        <td>
          selector is a label query over volumes to consider for binding.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storageClassName</b></td>
        <td>string</td>
        <td>
          storageClassName is the name of the StorageClass required by the claim.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeAttributesClassName</b></td>
        <td>string</td>
        <td>
          volumeAttributesClassName may be used to set the VolumeAttributesClass used by this claim.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeMode</b></td>
        <td>string</td>
        <td>
          volumeMode defines what type of volume is required by the claim.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeName</b></td>
        <td>string</td>
        <td>
          volumeName is the binding reference to the PersistentVolume backing this claim.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].ephemeral.volumeClaimTemplate.spec.dataSource
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplatespec)



dataSource field can be used to specify either:
* An existing VolumeSnapshot object (snapshot.

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
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind is the type of resource being referenced<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of resource being referenced<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiGroup</b></td>
        <td>string</td>
        <td>
          APIGroup is the group for the resource being referenced.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].ephemeral.volumeClaimTemplate.spec.dataSourceRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplatespec)



dataSourceRef specifies the object from which to populate the volume with data, if a non-empty...

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
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind is the type of resource being referenced<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of resource being referenced<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiGroup</b></td>
        <td>string</td>
        <td>
          APIGroup is the group for the resource being referenced.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace is the namespace of resource being referenced
Note that when a namespace is specified, a...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].ephemeral.volumeClaimTemplate.spec.resources
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplatespec)



resources represents the minimum resources the volume should have.

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
        <td><b>limits</b></td>
        <td>map[string]int or string</td>
        <td>
          Limits describes the maximum amount of compute resources allowed.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>requests</b></td>
        <td>map[string]int or string</td>
        <td>
          Requests describes the minimum amount of compute resources required.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].ephemeral.volumeClaimTemplate.spec.selector
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplatespec)



selector is a label query over volumes to consider for binding.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplatespecselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].ephemeral.volumeClaimTemplate.spec.selector.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexephemeralvolumeclaimtemplatespecselector)



A label selector requirement is a selector that contains values, a key, and an operator that...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].fc
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



fc represents a Fibre Channel resource that is attached to a kubelet's host machine and then...

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
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lun</b></td>
        <td>integer</td>
        <td>
          lun is Optional: FC target lun number<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly is Optional: Defaults to false (read/write).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>targetWWNs</b></td>
        <td>[]string</td>
        <td>
          targetWWNs is Optional: FC target worldwide names (WWNs)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>wwids</b></td>
        <td>[]string</td>
        <td>
          wwids Optional: FC volume world wide identifiers (wwids)
Either wwids or combination of targetWWNs...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].flexVolume
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



flexVolume represents a generic volume resource that is
provisioned/attached using an exec based...

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
        <td><b>driver</b></td>
        <td>string</td>
        <td>
          driver is the name of the driver to use for this volume.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>options</b></td>
        <td>map[string]string</td>
        <td>
          options is Optional: this field holds extra command options if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly is Optional: defaults to false (read/write).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexflexvolumesecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef is Optional: secretRef is reference to the secret object containing
sensitive information...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].flexVolume.secretRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexflexvolume)



secretRef is Optional: secretRef is reference to the secret object containing
sensitive information...

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].flocker
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



flocker represents a Flocker volume attached to a kubelet's host machine.

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
        <td><b>datasetName</b></td>
        <td>string</td>
        <td>
          datasetName is Name of the dataset stored as metadata -> name on the dataset for Flocker
should be...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>datasetUUID</b></td>
        <td>string</td>
        <td>
          datasetUUID is the UUID of the dataset. This is unique identifier of a Flocker dataset<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].gcePersistentDisk
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



gcePersistentDisk represents a GCE Disk resource that is attached to a
kubelet's host machine and...

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
        <td><b>pdName</b></td>
        <td>string</td>
        <td>
          pdName is unique name of the PD resource in GCE. Used to identify the disk in GCE.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is filesystem type of the volume that you want to mount.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>partition</b></td>
        <td>integer</td>
        <td>
          partition is the partition in the volume that you want to mount.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the ReadOnly setting in VolumeMounts.
Defaults to false.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].gitRepo
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



gitRepo represents a git repository at a particular revision.
Deprecated: GitRepo is deprecated.

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
        <td><b>repository</b></td>
        <td>string</td>
        <td>
          repository is the URL<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>directory</b></td>
        <td>string</td>
        <td>
          directory is the target directory name.
Must not contain or start with '..'.  If '.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>revision</b></td>
        <td>string</td>
        <td>
          revision is the commit hash for the specified revision.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].glusterfs
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



glusterfs represents a Glusterfs mount on the host that shares a pod's lifetime.

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
        <td><b>endpoints</b></td>
        <td>string</td>
        <td>
          endpoints is the endpoint name that details Glusterfs topology.
More info: https://examples.k8s.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is the Glusterfs volume path.
More info: https://examples.k8s.io/volumes/glusterfs/README.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the Glusterfs volume to be mounted with read-only permissions.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].hostPath
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



hostPath represents a pre-existing file or directory on the host
machine that is directly exposed...

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
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path of the directory on the host.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type for HostPath Volume
Defaults to ""
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].image
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



image represents an OCI object (a container image or artifact) pulled and mounted on the kubelet's...

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
        <td><b>pullPolicy</b></td>
        <td>string</td>
        <td>
          Policy for pulling OCI objects.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reference</b></td>
        <td>string</td>
        <td>
          Required: Image or artifact reference to be used.
Behaves in the same way as pod.spec.containers[*].<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].iscsi
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



iscsi represents an ISCSI Disk resource that is attached to a
kubelet's host machine and then...

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
        <td><b>iqn</b></td>
        <td>string</td>
        <td>
          iqn is the target iSCSI Qualified Name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lun</b></td>
        <td>integer</td>
        <td>
          lun represents iSCSI Target Lun number.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>targetPortal</b></td>
        <td>string</td>
        <td>
          targetPortal is iSCSI Target Portal.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>chapAuthDiscovery</b></td>
        <td>boolean</td>
        <td>
          chapAuthDiscovery defines whether support iSCSI Discovery CHAP authentication<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>chapAuthSession</b></td>
        <td>boolean</td>
        <td>
          chapAuthSession defines whether support iSCSI Session CHAP authentication<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type of the volume that you want to mount.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initiatorName</b></td>
        <td>string</td>
        <td>
          initiatorName is the custom iSCSI Initiator Name.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>iscsiInterface</b></td>
        <td>string</td>
        <td>
          iscsiInterface is the interface Name that uses an iSCSI transport.
Defaults to 'default' (tcp).<br/>
          <br/>
            <i>Default</i>: default<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>portals</b></td>
        <td>[]string</td>
        <td>
          portals is the iSCSI Target Portal List.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the ReadOnly setting in VolumeMounts.
Defaults to false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexiscsisecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef is the CHAP Secret for iSCSI target and initiator authentication<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].iscsi.secretRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexiscsi)



secretRef is the CHAP Secret for iSCSI target and initiator authentication

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].nfs
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



nfs represents an NFS mount on the host that shares a pod's lifetime
More info: https://kubernetes.

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
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path that is exported by the NFS server.
More info: https://kubernetes.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>server</b></td>
        <td>string</td>
        <td>
          server is the hostname or IP address of the NFS server.
More info: https://kubernetes.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the NFS export to be mounted with read-only permissions.
Defaults to false.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].persistentVolumeClaim
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



persistentVolumeClaimVolumeSource represents a reference to a
PersistentVolumeClaim in the same...

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
        <td><b>claimName</b></td>
        <td>string</td>
        <td>
          claimName is the name of a PersistentVolumeClaim in the same namespace as the pod using this volume.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly Will force the ReadOnly setting in VolumeMounts.
Default false.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].photonPersistentDisk
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets...

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
        <td><b>pdID</b></td>
        <td>string</td>
        <td>
          pdID is the ID that identifies Photon Controller persistent disk<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].portworxVolume
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



portworxVolume represents a portworx volume attached and mounted on kubelets host machine.

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
        <td><b>volumeID</b></td>
        <td>string</td>
        <td>
          volumeID uniquely identifies a Portworx volume<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fSType represents the filesystem type to mount
Must be a filesystem type supported by the host...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly defaults to false (read/write).<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



projected items for all in one resources secrets, configmaps, and downward API

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
        <td><b>defaultMode</b></td>
        <td>integer</td>
        <td>
          defaultMode are the mode bits used to set permissions on created files by default.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindex">sources</a></b></td>
        <td>[]object</td>
        <td>
          sources is the list of volume projections. Each entry in this list
handles one source.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojected)



Projection that may be projected along with other supported volume types.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexclustertrustbundle">clusterTrustBundle</a></b></td>
        <td>object</td>
        <td>
          ClusterTrustBundle allows a pod to access the `.spec.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexconfigmap">configMap</a></b></td>
        <td>object</td>
        <td>
          configMap information about the configMap data to project<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexdownwardapi">downwardAPI</a></b></td>
        <td>object</td>
        <td>
          downwardAPI information about the downwardAPI data to project<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexsecret">secret</a></b></td>
        <td>object</td>
        <td>
          secret information about the secret data to project<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexserviceaccounttoken">serviceAccountToken</a></b></td>
        <td>object</td>
        <td>
          serviceAccountToken is information about the serviceAccountToken data to project<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].clusterTrustBundle
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindex)



ClusterTrustBundle allows a pod to access the `.spec.

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
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Relative path from the volume root to write the bundle.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexclustertrustbundlelabelselector">labelSelector</a></b></td>
        <td>object</td>
        <td>
          Select all ClusterTrustBundles that match this label selector.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Select a single ClusterTrustBundle by object name.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          If true, don't block pod startup if the referenced ClusterTrustBundle(s)
aren't available.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>signerName</b></td>
        <td>string</td>
        <td>
          Select all ClusterTrustBundles that match this signer name.
Mutually-exclusive with name.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].clusterTrustBundle.labelSelector
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexclustertrustbundle)



Select all ClusterTrustBundles that match this label selector.

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexclustertrustbundlelabelselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].clusterTrustBundle.labelSelector.matchExpressions[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexclustertrustbundlelabelselector)



A label selector requirement is a selector that contains values, a key, and an operator that...

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].configMap
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindex)



configMap information about the configMap data to project

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexconfigmapitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          items if unspecified, each key-value pair in the Data field of the referenced
ConfigMap will be...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          optional specify whether the ConfigMap or its keys must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].configMap.items[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexconfigmap)



Maps a string key to a path within a volume.

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the key to project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          
path is the relative path of the file to map the key to.
May not be an absolute path.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          mode is Optional: mode bits used to set permissions on this file.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].downwardAPI
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindex)



downwardAPI information about the downwardAPI data to project

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexdownwardapiitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          Items is a list of DownwardAPIVolume file<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].downwardAPI.items[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexdownwardapi)



DownwardAPIVolumeFile represents information to create the file containing the pod field

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
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Required: Path is  the relative path name of the file to be created.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexdownwardapiitemsindexfieldref">fieldRef</a></b></td>
        <td>object</td>
        <td>
          Required: Selects a field of the pod: only annotations, labels, name, namespace and uid are...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          
Optional: mode bits used to set permissions on this file, must be an octal value
between 0000 and...<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexdownwardapiitemsindexresourcefieldref">resourceFieldRef</a></b></td>
        <td>object</td>
        <td>
          
Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].downwardAPI.items[index].fieldRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexdownwardapiitemsindex)



Required: Selects a field of the pod: only annotations, labels, name, namespace and uid are...

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
        <td><b>fieldPath</b></td>
        <td>string</td>
        <td>
          Path of the field to select in the specified API version.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          Version of the schema the FieldPath is written in terms of, defaults to "v1".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].downwardAPI.items[index].resourceFieldRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexdownwardapiitemsindex)




Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.

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
        <td><b>resource</b></td>
        <td>string</td>
        <td>
          Required: resource to select<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>containerName</b></td>
        <td>string</td>
        <td>
          Container name: required for volumes, optional for env vars<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>divisor</b></td>
        <td>int or string</td>
        <td>
          Specifies the output format of the exposed resources, defaults to "1"<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].secret
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindex)



secret information about the secret data to project

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
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexsecretitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          items if unspecified, each key-value pair in the Data field of the referenced
Secret will be...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          optional field specify whether the Secret or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].secret.items[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindexsecret)



Maps a string key to a path within a volume.

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the key to project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          
path is the relative path of the file to map the key to.
May not be an absolute path.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          mode is Optional: mode bits used to set permissions on this file.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].projected.sources[index].serviceAccountToken
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexprojectedsourcesindex)



serviceAccountToken is information about the serviceAccountToken data to project

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
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is the path relative to the mount point of the file to project the
token into.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>audience</b></td>
        <td>string</td>
        <td>
          audience is the intended audience of the token.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>expirationSeconds</b></td>
        <td>integer</td>
        <td>
          expirationSeconds is the requested duration of validity of the service
account token.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].quobyte
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



quobyte represents a Quobyte mount on the host that shares a pod's lifetime.

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
        <td><b>registry</b></td>
        <td>string</td>
        <td>
          registry represents a single or multiple Quobyte Registry services
specified as a string as...<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>volume</b></td>
        <td>string</td>
        <td>
          volume is a string that references an already created Quobyte volume by name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>group</b></td>
        <td>string</td>
        <td>
          group to map volume access to
Default is no group<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the Quobyte volume to be mounted with read-only permissions.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tenant</b></td>
        <td>string</td>
        <td>
          tenant owning the given Quobyte volume in the Backend
Used with dynamically provisioned Quobyte...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          user to map volume access to
Defaults to serivceaccount user<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].rbd
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



rbd represents a Rados Block Device mount on the host that shares a pod's lifetime.

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
        <td><b>image</b></td>
        <td>string</td>
        <td>
          image is the rados image name.
More info: https://examples.k8s.io/volumes/rbd/README.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>monitors</b></td>
        <td>[]string</td>
        <td>
          monitors is a collection of Ceph monitors.
More info: https://examples.k8s.io/volumes/rbd/README.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type of the volume that you want to mount.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keyring</b></td>
        <td>string</td>
        <td>
          keyring is the path to key ring for RBDUser.
Default is /etc/ceph/keyring.<br/>
          <br/>
            <i>Default</i>: /etc/ceph/keyring<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>pool</b></td>
        <td>string</td>
        <td>
          pool is the rados pool name.
Default is rbd.
More info: https://examples.k8s.io/volumes/rbd/README.<br/>
          <br/>
            <i>Default</i>: rbd<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the ReadOnly setting in VolumeMounts.
Defaults to false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexrbdsecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef is name of the authentication secret for RBDUser. If provided
overrides keyring.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          user is the rados user name.
Default is admin.
More info: https://examples.k8s.<br/>
          <br/>
            <i>Default</i>: admin<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].rbd.secretRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexrbd)



secretRef is name of the authentication secret for RBDUser. If provided
overrides keyring.

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].scaleIO
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes.

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
        <td><b>gateway</b></td>
        <td>string</td>
        <td>
          gateway is the host address of the ScaleIO API Gateway.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexscaleiosecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef references to the secret for ScaleIO user and other
sensitive information.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>system</b></td>
        <td>string</td>
        <td>
          system is the name of the storage system as configured in ScaleIO.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.<br/>
          <br/>
            <i>Default</i>: xfs<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>protectionDomain</b></td>
        <td>string</td>
        <td>
          protectionDomain is the name of the ScaleIO Protection Domain for the configured storage.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly Defaults to false (read/write).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sslEnabled</b></td>
        <td>boolean</td>
        <td>
          sslEnabled Flag enable/disable SSL communication with Gateway, default false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storageMode</b></td>
        <td>string</td>
        <td>
          storageMode indicates whether the storage for a volume should be ThickProvisioned or...<br/>
          <br/>
            <i>Default</i>: ThinProvisioned<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storagePool</b></td>
        <td>string</td>
        <td>
          storagePool is the ScaleIO Storage Pool associated with the protection domain.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeName</b></td>
        <td>string</td>
        <td>
          volumeName is the name of a volume already created in the ScaleIO system
that is associated with...<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].scaleIO.secretRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexscaleio)



secretRef references to the secret for ScaleIO user and other
sensitive information.

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].secret
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



secret represents a secret that should populate this volume.
More info: https://kubernetes.

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
        <td><b>defaultMode</b></td>
        <td>integer</td>
        <td>
          defaultMode is Optional: mode bits used to set permissions on created files by default.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexsecretitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          items If unspecified, each key-value pair in the Data field of the referenced
Secret will be...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          optional field specify whether the Secret or its keys must be defined<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>secretName</b></td>
        <td>string</td>
        <td>
          secretName is the name of the secret in the pod's namespace to use.
More info: https://kubernetes.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].secret.items[index]
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexsecret)



Maps a string key to a path within a volume.

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
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the key to project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is the relative path of the file to map the key to.
May not be an absolute path.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          mode is Optional: mode bits used to set permissions on this file.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].storageos
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



storageOS represents a StorageOS volume attached and mounted on Kubernetes nodes.

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
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly defaults to false (read/write).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexstorageossecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef specifies the secret to use for obtaining the StorageOS API
credentials.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeName</b></td>
        <td>string</td>
        <td>
          volumeName is the human-readable name of the StorageOS volume.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeNamespace</b></td>
        <td>string</td>
        <td>
          volumeNamespace specifies the scope of the volume within StorageOS.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].storageos.secretRef
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindexstorageos)



secretRef specifies the secret to use for obtaining the StorageOS API
credentials.

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
          Name of the referent.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.deployment.template.spec.volumes[index].vsphereVolume
[Go to parent definition](#gatewayclassparametersspeckubernetesdeploymenttemplatespecvolumesindex)



vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine.

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
        <td><b>volumePath</b></td>
        <td>string</td>
        <td>
          volumePath is the path that identifies vSphere volume vmdk<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is filesystem type to mount.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storagePolicyID</b></td>
        <td>string</td>
        <td>
          storagePolicyID is the storage Policy Based Management (SPBM) profile ID associated with the...<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storagePolicyName</b></td>
        <td>string</td>
        <td>
          storagePolicyName is the storage Policy Based Management (SPBM) profile name.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.spec.kubernetes.service
[Go to parent definition](#gatewayclassparametersspeckubernetes)





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
        <td><b>annotations</b></td>
        <td>map[string]string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>externalTrafficPolicy</b></td>
        <td>string</td>
        <td>
          ServiceExternalTrafficPolicy describes how nodes distribute service traffic they
receive on one of...<br/>
          <br/>
            <i>Default</i>: Cluster<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>labels</b></td>
        <td>map[string]string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>loadBalancerClass</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Service Type string describes ingress methods for a service<br/>
          <br/>
            <i>Default</i>: LoadBalancer<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GatewayClassParameters.status
[Go to parent definition](#gatewayclassparameters)





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
        <td><b><a href="#gatewayclassparametersstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### GatewayClassParameters.status.conditions[index]
[Go to parent definition](#gatewayclassparametersstatus)



Condition contains details for one aspect of the current state of this API Resource.

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
          lastTransitionTime is the last time the condition transitioned from one status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          message is a human readable message indicating details about the transition.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          reason contains a programmatic identifier indicating the reason for the condition's last transition.<br/>
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
          type of condition in CamelCase or in foo.example.com/CamelCase.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          observedGeneration represents the .metadata.generation that the condition was set based upon.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Minimum</i>: 0<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
