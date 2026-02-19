{{ with secret "kv/mongo" }}{{ .Data.MONGO_INITDB_ROOT_PASSWORD }}{{ end }}
