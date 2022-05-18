# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# It starts from the postgres base image.
FROM postgres

# Add Maintainer Info.
LABEL maintainer="Ícaro Ribeiro <icaroribeiro@hotmail.com>"

# Add the content of the scripts directory to the image.
# All files in docker-entrypoint-initdb.d are automatically executed during container startup.
COPY ./scripts ./docker-entrypoint-initdb.d