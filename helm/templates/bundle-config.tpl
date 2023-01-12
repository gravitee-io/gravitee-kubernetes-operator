kind: ConfigMap 
apiVersion: v1 
metadata:
  name: {{ .Values.manager.configMap.name }}
data:
{{- if not .Values.manager.scope.cluster }}
  NAMESPACE: {{ .Release.Namespace }}
{{- end }}
{{- if not .Values.manager.logs.json }}
  DEV_MODE: "false"
{{- end }}
