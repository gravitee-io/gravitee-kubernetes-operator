{{/*
Expand the name of the chart.
*/}}
{{- define "helm.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "helm.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "helm.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "helm.labels" -}}
helm.sh/chart: {{ include "helm.chart" . }}
{{ include "helm.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "helm.selectorLabels" -}}
app.kubernetes.io/name: {{ include "helm.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "helm.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "helm.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}


{{/*
 Create the name of the service account to use for the manager
 */}}
{{- define "rbac.serviceAccountName" -}}
{{- default "gko-controller-manager" .Values.serviceAccount.name }}
{{- end }}

{{/*
 Create the name of the manager cluster role
 */}}
{{- define "rbac.ClusterRoleName" -}}
   {{ template "rbac.serviceAccountName" . }}-cluster-role
{{- end }}

{{/*
 Create the name of the manager cluster role binding
 */}}
{{- define "rbac.ClusterRoleBindingName" -}}
   {{ template "rbac.ClusterRoleName" . }}-binding
{{- end }}

{{/*
 Create the name of the manager role
 */}}
{{- define "rbac.RoleName" -}}
   {{ template "rbac.serviceAccountName" . }}-role
{{- end }}


{{/*
 Create the name of the manager role binding
 */}}
{{- define "rbac.RoleBindingName" -}}
   {{ template "rbac.RoleName" . }}-binding
{{- end }}

{{/*
 Create the name of the manager role for leader election
 */}}
{{- define "rbac.LeaderElectionRoleName" -}}
   {{ template "rbac.serviceAccountName" . }}-leader-election-role
{{- end }}

{{/*
 Create the name of the manager role binding for leader election
 */}}
{{- define "rbac.LeaderElectionRoleBindingName" -}}
   {{ template "rbac.LeaderElectionRoleName" . }}-binding
{{- end }}


{{/*
 Create the name of the manager cluster role for CRD patch
 */}}
{{- define "rbac.ResourcePatchClusterRoleName" -}}
   {{ template "rbac.serviceAccountName" . }}-crd-patch-cluster-role
{{- end }}

{{/*
 Create the name of the manager cluster role binding for CRD patch
 */}}
{{- define "rbac.ResourcePatchClusterRoleBindingName" -}}
   {{ template "rbac.ResourcePatchClusterRoleName" . }}-binding
{{- end }}

{{/*
 Create the name of the manager cluster role for metrics
 */}}
{{- define "rbac.MetricsClusterRoleName" -}}
   {{ template "rbac.serviceAccountName" . }}-metrics-cluster-role
{{- end }}

{{/*
 Create the name of the manager cluster role binding for metrics
 */}}
{{- define "rbac.MetricsClusterRoleBindingName" -}}
   {{ template "rbac.MetricsClusterRoleName" . }}-binding
{{- end }}

{{/*
 Create the name of the kube rbac provy cluster role for metrics
 */}}
{{- define "rbac.ProxyClusterRoleName" -}}
   {{ template "rbac.serviceAccountName" . }}-proxy-cluster-role
{{- end }}

{{/*
 Create the name of the kube rbac procy role binding for metrics
 */}}
{{- define "rbac.ProxyClusterRoleBindingName" -}}
   {{ template "rbac.ProxyClusterRoleName" . }}-binding
{{- end }}

{{/*
 merge list of ingress classes into a single string that will be parsed later in the code
 */}}
{{- define "ingress.Classes" -}}
{{- join "," .Values.ingress.ingressClasses }}
{{- end -}}