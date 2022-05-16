#!/bin/bash

#
# Internal
#
# Generate a mock object related to auth's datastore repository.
AUTH_DATASTORE_REPOSITORY_PATH="internal/core/ports/infrastructure/storage/datastore/repository/auth"
MOCK_AUTH_DATASTORE_REPOSITORY_PATH="internal/core/ports/infrastructure/storage/datastore/mockrepository/auth"
mockery --dir "$AUTH_DATASTORE_REPOSITORY_PATH" --name IRepository --outpkg auth --structname Repository --output "$MOCK_AUTH_DATASTORE_REPOSITORY_PATH" --filename mock_repository.go

# Generate a mock object related to login's datastore repository.
LOGIN_DATASTORE_REPOSITORY_PATH="internal/core/ports/infrastructure/storage/datastore/repository/login"
MOCK_LOGIN_DATASTORE_REPOSITORY_PATH="internal/core/ports/infrastructure/storage/datastore/mockrepository/login"
mockery --dir "$LOGIN_DATASTORE_REPOSITORY_PATH" --name IRepository --outpkg login --structname Repository --output "$MOCK_LOGIN_DATASTORE_REPOSITORY_PATH" --filename mock_repository.go

# Generate a mock object related to user's datastore repository.
USER_DATASTORE_REPOSITORY_PATH="internal/core/ports/infrastructure/storage/datastore/repository/user"
MOCK_USER_DATASTORE_REPOSITORY_PATH="internal/core/ports/infrastructure/storage/datastore/mockrepository/user"
mockery --dir "$USER_DATASTORE_REPOSITORY_PATH" --name IRepository --outpkg user --structname Repository --output "$MOCK_USER_DATASTORE_REPOSITORY_PATH" --filename mock_repository.go

# Generate a mock object related to auth's service.
AUTH_SERVICE_PATH="internal/core/ports/application/service/auth"
MOCK_AUTH_SERVICE_PATH="internal/core/ports/application/mockservice/auth"
mockery --dir "$AUTH_SERVICE_PATH" --name IService --outpkg auth --structname Service --output "$MOCK_AUTH_SERVICE_PATH" --filename mock_service.go

# Generate a mock object related to healthcheck's service.
HEALTHCHECK_SERVICE_PATH="internal/core/ports/application/service/healthcheck"
MOCK_HEALTHCHECK_SERVICE_PATH="internal/core/ports/application/mockservice/healthcheck"
mockery --dir "$HEALTHCHECK_SERVICE_PATH" --name IService --outpkg healthcheck --structname Service --output "$MOCK_HEALTHCHECK_SERVICE_PATH" --filename mock_service.go

# Generate a mock object related to user service.
USER_SERVICE_PATH="internal/core/ports/application/service/user"
MOCK_USER_SERVICE_PATH="internal/core/ports/application/mockservice/user"
mockery --dir "$USER_SERVICE_PATH" --name IService --outpkg user --structname Service --output "$MOCK_USER_SERVICE_PATH" --filename mock_service.go

#
# PKG
#
# Generate a mock object related to auth.
AUTH_PATH="pkg/auth"
MOCK_AUTH_PATH="tests/mocks/pkg/mockauth"
mockery --dir "$AUTH_PATH" --name IAuth --outpkg mockauth --structname Auth --output "$MOCK_AUTH_PATH" --filename mock_auth.go

# Generate a mock object related to security.
SECURITY_PATH="pkg/security"
MOCK_SECURITY_PATH="tests/mocks/pkg/mocksecurity"
mockery --dir "$SECURITY_PATH" --name ISecurity --outpkg mocksecurity --structname Security --output "$MOCK_SECURITY_PATH" --filename mock_security.go

# Generate a mock object related to validator.
VALIDATOR_PATH="pkg/validator"
MOCK_VALIDATOR_PATH="tests/mocks/pkg/mockvalidator"
mockery --dir "$VALIDATOR_PATH" --name IValidator --outpkg mockvalidator --structname Validator --output "$MOCK_VALIDATOR_PATH" --filename mock_validator.go