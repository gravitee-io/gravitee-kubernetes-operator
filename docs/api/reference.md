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
Most of the configuration properties defined here are already documented
in the APIM Console API Reference.
See https://docs.gravitee.io/apim/3.x/apim_installguide_rest_apis_documentation.html<br/>
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
Most of the configuration properties defined here are already documented
in the APIM Console API Reference.
See https://docs.gravitee.io/apim/3.x/apim_installguide_rest_apis_documentation.html

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#apidefinitionspecdefinition_context">definition_context</a></b></td>
        <td>object</td>
        <td>
          The definition context is used to inform a management API instance that this API definition
is managed using a kubernetes operator<br/>
        </td>
        <td>true</td>
      </tr><tr>
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
          The list of categories the API belongs to.
Categories are reflected in APIM portal so that consumers can easily find the APIs they need.<br/>
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
this ID identifies the API across those different environments.
Setting this ID also allows to take control over an existing API on an APIM instance
(by setting the same value as defined in APIM).
If empty, a UUID will be generated based on the namespace and name of the resource.<br/>
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
        <td><b>groups</b></td>
        <td>[]string</td>
        <td>
          List of groups associated with the API<br/>
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
          local defines if the api is local or not.

If true, the Operator will create the ConfigMaps for the Gateway and pushes the API to the Management API
but without setting the update flag in the datastore.

If false, the Operator will not create the ConfigMaps for the Gateway.
Instead, it pushes the API to the Management API and forces it to update the event in the datastore.
This will cause Gateways to fetch the APIs from the datastore<br/>
          <br/>
            <i>Default</i>: true<br/>
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
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecpropertiesindex">properties</a></b></td>
        <td>[]object</td>
        <td>
          List of Properties for the API<br/>
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
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
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
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecflowsindexpreindex">pre</a></b></td>
        <td>[]object</td>
        <td>
          Flow pre step<br/>
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
        <td><b>displayName</b></td>
        <td>string</td>
        <td>
          Member display name<br/>
        </td>
        <td>false</td>
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
          This is the display name of the page in APIM and on the portal.
This field can be edited safely if you want to rename a page.<br/>
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
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecplansindexflowsindex">flows</a></b></td>
        <td>[]object</td>
        <td>
          List of different flows for this Plan<br/>
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
            <i>Enum</i>: PUBLISHED<br/>
            <i>Default</i>: PUBLISHED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tags</b></td>
        <td>[]string</td>
        <td>
          List of plan tags<br/>
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
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apidefinitionspecplansindexflowsindexpreindex">pre</a></b></td>
        <td>[]object</td>
        <td>
          Flow pre step<br/>
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
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>allowMethods</b></td>
        <td>[]string</td>
        <td>
          Access Control - List of allowed methods<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>allowOrigin</b></td>
        <td>[]string</td>
        <td>
          Access Control -  List of Allowed origins<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>exposeHeaders</b></td>
        <td>[]string</td>
        <td>
          Access Control - List of Exposed Headers<br/>
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
        <td><b>inherit</b></td>
        <td>boolean</td>
        <td>
          Is service inherited or not?<br/>
        </td>
        <td>true</td>
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
        <td><b>type</b></td>
        <td>string</td>
        <td>
          The KeyStore type to use (possible values are PEM, PKCS12, JKS)<br/>
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
        <td><b>type</b></td>
        <td>string</td>
        <td>
          The TrustStore type to use (possible values are PEM, PKCS12, JKS)<br/>
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
        <td><b>type</b></td>
        <td>string</td>
        <td>
          The KeyStore type to use (possible values are PEM, PKCS12, JKS)<br/>
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
        <td><b>type</b></td>
        <td>string</td>
        <td>
          The TrustStore type to use (possible values are PEM, PKCS12, JKS)<br/>
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
            <i>Enum</i>: PROXY, MESSAGE<br/>
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
          API Analytics<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>categories</b></td>
        <td>[]string</td>
        <td>
          The list of categories the API belongs to.
Categories are reflected in APIM portal so that consumers can easily find the APIs they need.<br/>
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
        <td><b><a href="#apiv4definitionspecflowexecution">flowExecution</a></b></td>
        <td>object</td>
        <td>
          API Flow Execution<br/>
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
        <td><b>groups</b></td>
        <td>[]string</td>
        <td>
          List of groups associated with the API<br/>
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
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecresourcesindex">resources</a></b></td>
        <td>[]object</td>
        <td>
          Resources can be either inlined or reference the namespace and name
of an <a href="#apiresource">existing API resource definition</a>.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecresponsetemplateskeykey">responseTemplates</a></b></td>
        <td>map[string]map[string]object</td>
        <td>
          A list of Response Templates for the API<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecservices">services</a></b></td>
        <td>object</td>
        <td>
          API Services<br/>
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
        <td><b>type</b></td>
        <td>string</td>
        <td>
          The KeyStore type to use (possible values are PEM, PKCS12, JKS)<br/>
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
        <td><b>type</b></td>
        <td>string</td>
        <td>
          The TrustStore type to use (possible values are PEM, PKCS12, JKS)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.analytics
[Go to parent definition](#apiv4definitionspec)



API Analytics

<table>
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
            <i>Default</i>: KUBERNETES<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### ApiV4Definition.spec.flowExecution
[Go to parent definition](#apiv4definitionspec)



API Flow Execution

<table>
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
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the flow this field is mainly used for compatibility with
APIM exports and can be safely ignored.<br/>
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
          List of Request flow steps<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecflowsindexresponseindex">response</a></b></td>
        <td>[]object</td>
        <td>
          List of Response flow steps<br/>
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
        <td><b>displayName</b></td>
        <td>string</td>
        <td>
          Member display name<br/>
        </td>
        <td>false</td>
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
          This is the display name of the page in APIM and on the portal.
This field can be edited safely if you want to rename a page.<br/>
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
            <i>Enum</i>: PUBLISHED<br/>
            <i>Default</i>: PUBLISHED<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tags</b></td>
        <td>[]string</td>
        <td>
          List of plan tags<br/>
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
        <td><b>id</b></td>
        <td>string</td>
        <td>
          The ID of the flow this field is mainly used for compatibility with
APIM exports and can be safely ignored.<br/>
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
          List of Request flow steps<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#apiv4definitionspecplanskeyflowsindexresponseindex">response</a></b></td>
        <td>[]object</td>
        <td>
          List of Response flow steps<br/>
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



API Services

<table>
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
          The base64 encoded picture to use for this application when displaying it on the portal (if not relying on an URL)<br/>
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
        <td>string</td>
        <td>
          Oauth client application type<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>grantTypes</b></td>
        <td>[]string</td>
        <td>
          List of Oauth client grant types<br/>
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
        <td><b>displayName</b></td>
        <td>string</td>
        <td>
          Member display name<br/>
        </td>
        <td>false</td>
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
          When API has been created regardless of errors, this field is
used to persist the error message encountered during admission<br/>
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
          The processing status of the Application.
The value is `Completed` if the sync with APIM succeeded, Failed otherwise.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Application.status.errors
[Go to parent definition](#applicationstatus)



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
