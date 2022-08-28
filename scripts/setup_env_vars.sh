#!/bin/bash

#
# HTTP server settings
#
export HTTP_PORT="8080"

#
# RSA Keys
#
export RSA_PUBLIC_KEY_PATH="./configs/auth/rsa_keys/rsa.public"
export RSA_PRIVATE_KEY_PATH="./configs/auth/rsa_keys/rsa.private"

#
# JWT settings
#
export TOKEN_EXP_TIME_IN_SEC="120"
export TIME_BEFORE_TOKEN_EXP_TIME_IN_SEC="30"

#
# Datastore settings
#
export DB_DRIVER="postgres"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_NAME="db"