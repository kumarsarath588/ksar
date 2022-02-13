{{/*
Expand the name of the chart.
*/}}
{{- define "tabsquare-app.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "tabsquare-app.fullname" -}}
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
{{- define "tabsquare-app.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "tabsquare-app.labels" -}}
helm.sh/chart: {{ include "tabsquare-app.chart" . }}
{{ include "tabsquare-app.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "tabsquare-app.selectorLabels" -}}
app.kubernetes.io/name: {{ include "tabsquare-app.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "tabsquare-app.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "tabsquare-app.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "mysql.fullname" -}}
{{- printf "%s-%s" .Release.Name "mysql" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Return the Mysql Host
*/}}
{{- define "tabsquare.databaseHost" -}}
{{- if .Values.mysql.enabled }}
    {{- printf "%s" (include "mysql.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Mysql Port
*/}}
{{- define "tabsquare.databasePort" -}}
{{- if .Values.mysql.enabled }}
    {{- printf "3306" -}}
{{- end -}}
{{- end -}}

{{/*
Return the Mysql Database Name
*/}}
{{- define "tabsquare.databaseName" -}}
{{- if .Values.mysql.enabled }}
    {{- printf "%s" .Values.mysql.mysqlDatabase -}}
{{- end -}}
{{- end -}}

{{/*
Return the Mysql User
*/}}
{{- define "tabsquare.databaseUser" -}}
{{- if .Values.mysql.enabled }}
    {{- printf "%s" .Values.mysql.mysqlUser -}}
{{- end -}}
{{- end -}}

{{/*
Return the Mysql User Password
*/}}
{{- define "tabsquare.databaseSecretName" -}}
{{- if .Values.mysql.enabled }}
    {{- printf "%s" (include "mysql.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Init container name
*/}}
{{- define "tabsquare-app.InitContainerName" -}}
{{- if .Values.initContainers.enabled }}
    {{- printf "%s-%s" .Release.Name "init" -}}
{{- end -}}
{{- end }}