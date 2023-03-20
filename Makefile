#
# API building and running locally
# Set of tasks related to API building and running.
#
audit-api:
	go mod tidy

format-api:
	gofmt -w .; \
	golint ./...

#
# The '-mod=mod' instruction tells the go command to update go.mod and go.sum 
# if anything is missing or inconsistent related to dependency management in GoLang.
#
generate-gql:
	go run -mod=mod github.com/99designs/gqlgen generate --config internal/presentation/api/gqlgen/gqlgen.yml

run-api:
	. ./scripts/setup_env_vars.sh; \
	go run cmd/api/main.go

#
# API test
# Set of tasks related to API testing locally.
#
build-mocks:
	. ./scripts/build_mocks.sh

test-api:
	. ./scripts/setup_env_vars.test.sh; \
	go test ./... -v -coverprofile=./docs/api/tests/unit/coverage.out; \
	go tool cover -func=./docs/api/tests/unit/coverage.out > ./docs/api/tests/unit/coverage_report.out; \
	go test ./tests/api/... -v

#
# APP test container
# Set of tasks related to APP testing container.
#
start-deps:
	docker network create testapp_network; \
	docker build -t postgrestestdb --no-cache -f ./database/postgres/Dockerfile .; \
	docker run --name postgrestestdb_container --env-file ./database/postgres/.env.test -d -p 5434:5432 -v postgrestestdb-data:/var/lib/postgresql/data --restart on-failure postgrestestdb; \
	docker network connect testapp_network postgrestestdb_container

init-app:
	docker build -t apitest -f Dockerfile.test .; \
	docker run --name apitest_container --env-file ./.env.test -d -p 8080:8080 --restart on-failure apitest; \
	docker network connect testapp_network apitest_container

test-app:
	docker exec --env-file ./.env.test apitest_container go test ./... -v

destroy-app:
	docker network disconnect testapp_network apitest_container; \
	docker stop apitest_container; \
	docker rm apitest_container; \
	docker rmi apitest

finish-deps:
	docker network disconnect testapp_network postgrestestdb_container; \
	docker stop postgrestestdb_container; \
	docker rm postgrestestdb_container; \
	docker volume rm postgrestestdb-data; \
	docker rmi postgrestestdb; \
	docker network rm testapp_network

#
# APP production container
# Set of tasks related to APP production container starting up and shutting down.
#
startup-app:
	docker-compose up -d --build api

shutdown-app:
	docker-compose down -v --rmi all
