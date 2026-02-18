package test

const (
	configContentJSON = `
{
  "config_name": "gostrap",
  "GIN": {
    "port": "8088",
    "enable": true,
    "mode": "test",
    "middleware": {
      "cors": {
        "enabled": true,
        "allow_origins": [
          "http://localhost:5173"
        ],
        "allow_methods": [
          "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"
        ],
        "allow_headers": [
          "Origin", "Content-Type", "Authorization"
        ],
        "expose_headers": [
          "Content-Length"
        ],
        "allow_credentials": true,
        "max_age": 43200
      }
    }
  }
}
`
	configContentYAML = `
configName: "gostrap"
gin:
  port: "8088"
  enable: true
  mode: "test"
  middleware:
    cors:
      enabled: true
      allowOrigins:
        - "http://localhost:5173"
      allowMethods:
        - "GET"
        - "POST"
        - "PUT"
        - "DELETE"
        - "OPTIONS"
        - "PATCH"
      allowHeaders:
        - "Origin"
        - "Content-Type"
        - "Authorization"
      exposeHeaders:
        - "Content-Length"
      allowCredentials: true
      maxAge: 43200
`
)
