swagger:
	swag init -g internal/ports/http/server.go -o docs/swagger

swagger-serve:
	swagger serve docs/swagger.yaml