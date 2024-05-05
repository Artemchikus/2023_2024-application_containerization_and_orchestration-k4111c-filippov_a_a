{{- define "app-cm" }}
data:
  config.yaml: |
    ### CONFIG START ###
    app:
      port: {{ .Values.app.port }}
    db:
      host: {{ .Values.db.host }}
      port: {{ .Values.db.port }}
      user: {{ .Values.db.user }}
      name: {{ .Values.db.name }}
      search_path: {{ .Values.db.search_path }}
      password: {{ .Values.db.password }}
    ### CONFIG END ###
{{- end }}