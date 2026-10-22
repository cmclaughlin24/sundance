.PHONY: openapi

openapi:
	(cd ./services/forms && swag init -g ./internal/adapters/rest/routes.go -o ./internal/adapters/rest/docs --parseDependency --useStructName)
	(cd ./services/tenants && swag init -g ./internal/adapters/rest/routes.go -o ./internal/adapters/rest/docs --parseDependency --useStructName)
