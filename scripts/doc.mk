swagger:
	swag init -g internal/ports/http/server.go --parseDependency --parseInternal --parseDepth 2 -o docs/swagger

swagger-serve:
	swagger serve docs/swagger.yaml