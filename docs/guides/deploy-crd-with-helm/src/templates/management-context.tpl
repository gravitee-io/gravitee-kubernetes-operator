apiVersion: gravitee.io/v1alpha1
kind: ManagementContext
metadata:
  name: {{ regexReplaceAll "\\W+" .Values.context.name "-" | lower }}
spec:
  baseUrl: {{ .Values.context.baseUrl }}
  environmentId: {{ .Values.context.environmentId }}
  organizationId: {{ .Values.context.organizationId }}
  auth:
    bearerToken: {{ .Values.context.token | quote }}

