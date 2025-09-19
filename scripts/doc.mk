swagger:
	swag init -g cmd/app/main.go -o docs/swagger

swagger-serve:
	swagger serve docs/swagger.yaml