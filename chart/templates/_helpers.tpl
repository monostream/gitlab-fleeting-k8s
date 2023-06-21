{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "gitlab-runner-docker-autoscaler.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "gitlab-runner-docker-autoscaler.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "gitlab-runner-docker-autoscaler.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "gitlab-runner-docker-autoscaler.labels" -}}
helm.sh/chart: {{ include "gitlab-runner-docker-autoscaler.chart" . }}
{{ include "gitlab-runner-docker-autoscaler.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "gitlab-runner-docker-autoscaler.selectorLabels" -}}
app.kubernetes.io/name: {{ include "gitlab-runner-docker-autoscaler.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "gitlab-runner-docker-autoscaler.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "gitlab-runner-docker-autoscaler.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Unregister runners on pod stop
*/}}
{{- define "gitlab-runner-docker-autoscaler.unregisterRunners" -}}
{{- if or (and (hasKey .Values "unregisterRunners") .Values.unregisterRunners) (and (not (hasKey .Values "unregisterRunners")) .Values.runnerRegistrationToken) -}}
lifecycle:
  preStop:
    exec:
      command: ["/entrypoint", "unregister", "--all-runners"]
{{- end -}}
{{- end -}}

{{/*
Define the name of the secret containing the tokens
*/}}
{{- define "gitlab-runner-docker-autoscaler.secret" -}}
{{- default (include "gitlab-runner-docker-autoscaler.fullname" .) .Values.runners.secret | quote -}}
{{- end -}}