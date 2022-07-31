#
# Set of tasks related to API building and running.
#
setup-api:
	go mod download

audit-api:
	go mod tidy

format-api:
	gofmt -w .; \
	golint ./...

# The '-mod=mod' instruction tells the go command to update go.mod and go.sum 
# if anything is missing or inconsistent related to dependency management in GoLang.
generate-gql:
	go run -mod=mod github.com/99designs/gqlgen generate --config internal/transport/presentation/handler/graphql/gqlgen/gqlgen.yml

checkversion-api:
	go run cmd/api/main.go version

run-api:
	. ./scripts/setup_env_vars.sh; \
	go run cmd/api/main.go run

#
# Set of tasks related to API testing.
#
build-mocks:
	. ./scripts/build_mocks.sh

test-api:
	. ./scripts/setup_env_vars.test.sh; \
	go test ./... -v -coverprofile=./docs/api/tests/unit/coverage.out

analyze-api:
	go tool cover -func=./docs/api/tests/unit/coverage.out > ./docs/api/tests/unit/coverage_report.out

#
# Set of tasks related to APP container
#
startup-app:
	docker-compose up -d --build api

shutdown-app:
	docker-compose down -v --rmi all

#
# Set of tasks related to APP container testing
#
start-deps:
	docker-compose up -d --build postgrestestdb

finish-deps:
	docker-compose rm --force --stop -v postgrestestdb

test-app:
	docker exec --env-file ./.env.test api_container go test ./... -v -coverprofile=./docs/api/tests/unit/coverage.out

analyze-app:
	docker-compose exec api go tool cover -func=./docs/api/tests/unit/coverage.out

#
# Set of tasks related to APP deployment.
#
init-deploy:
	cd deployments/heroku/terraform; \
	terraform init

plan-deploy:
	. ./deployments/heroku/scripts/setup_env_vars.sh; \
	cd deployments/heroku/terraform; \
	terraform plan

apply-deploy:
	. ./deployments/heroku/scripts/build_app.sh; \
	. ./deployments/heroku/scripts/setup_env_vars.sh; \
	cd deployments/heroku/terraform; \
	terraform apply

destroy-deploy:
	. ./deployments/heroku/scripts/destroy_app.sh; \
	. ./deployments/heroku/scripts/setup_env_vars.sh; \
	cd deployments/heroku/terraform; \
	terraform destroy
