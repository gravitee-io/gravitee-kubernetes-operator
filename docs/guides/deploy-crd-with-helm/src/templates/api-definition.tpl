

{{- $common := .Values }}
{{- $ns := .Release.Namespace}}
{{- range $base := $common.apis }}
{{- $api := deepCopy $common.api | merge $base}}
---
apiVersion: "gravitee.io/v1alpha1"
kind: "ApiDefinition"
metadata:
    name: {{ regexReplaceAll "\\W+" $api.name "-" | lower }}
spec:
    name: {{ $api.name }}
    version: "{{ $api.version }}"
{{- if $api.description }}
    description: {{ $api.description }}
{{- else }}
    description: "API {{ $api.name }}"
{{- end }}
    visibility: {{ $api.visibility }}
{{- if $api.flows }}
    flows: 
{{ toYaml $api.flows | indent 4 }}
{{ else }}
    flows: []
{{- end }}
    gravitee: "2.0.0"
    flow_mode: {{ $api.flow_mode }}
{{- if $common.resources }}
    resources:
{{ toYaml $common.resources | indent 4 }}
{{ else }}
    resources: []
{{- end }}
{{- if $api.properties }}
    properties:
{{ toYaml $api.properties | indent 4 }}
{{ else }}
    properties: []
{{- end }}
{{- if $api.plans }}
    plans:
{{ toYaml $api.plans | indent 4 }}
{{ else }}
    plans: []
{{- end }}
{{- if $api.path_mappings }}
    path_mappings:
{{ toYaml $api.path_mappings | indent 4 }}
{{ else }}
    path_mappings: []
{{- end }}
    proxy:
{{- if $api.proxy.virtual_hosts }}
        virtual_hosts:
{{ toYaml $api.proxy.virtual_hosts | indent 6 }}
{{ else }}
        virtual_hosts:
            - path: "/{{ regexReplaceAll "\\W+" $api.name "-" | lower }}"
{{- end }}
        strip_context_path: {{ $api.proxy.strip_context_path }}
        preserve_host: {{ $api.proxy.preserve_host }}
{{- range $api.proxy.groups }}
        groups:
            - name: {{ .name | default "default-groups" }}
              endpoints:
{{- range .endpoints }}
                  - backup: {{ .backup | default false }}
                    inherit: {{ .inherit | default true }}
                    name: {{ .name | default "default" }}
                    weight: {{ .weight | default 1 }}
                    type: {{ .type | default "http" }}
                    target: {{ .target }}
{{- end }}
{{- if .load_balancing }}
              load_balancing:
                  type: {{ .load_balancing.type }}
{{- else }}
              load_balancing:
                  type: "ROUND_ROBIN"
{{- end }}
{{- if .http }}
              http:
{{ toYaml .http | indent 16 }}
{{- else }}
              http:
                connectTimeout: 5000
                idleTimeout: 60000
                keepAlive: true
                readTimeout: 10000
                pipelining: false
                maxConcurrentConnections: 100
                useCompression: true
                followRedirects: false
{{- end }}
{{- end }}
{{- if $api.response_templates }}
{{ toYaml $api.response_templates | indent 4 }}
{{- else }}
    response_templates: {}
{{- end }}
{{- if $common.context }}
    contextRef:
        name: {{ $common.context.name }}
        namespace: {{ $ns }}
{{- end }}
{{- end }}
