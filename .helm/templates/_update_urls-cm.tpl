{{- define "update_urls-cm" }}
data:
  config.yaml: |
    ### CONFIG START ###
    service_url: {{ .Chart.Name }}
    service_api: api/{{ .Values.app.api_version }}
    service_port: {{ .Values.app.port }}
    ### CONFIG END ###
{{- end }}