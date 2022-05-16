# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# It starts from the golang base image.
FROM golang

# Add maintainer info.
LABEL maintainer="Ícaro Ribeiro <icaroribeiro@hotmail.com>"

# Set the environment variable that deals with the application deployment in Heroku Platform.
ENV DEPLOY=NO

# Set the working directory inside the container.
WORKDIR /app

# Copy the source from the current directory to the working directory inside the container.
COPY . .

# Download all dependencies.
RUN go mod download

# Build the Go application.
RUN cd cmd/api && go build -o api .

# Command to run the application.
# SIM
CMD ./cmd/api/api run
# NÃO
# CMD ["./cmd/api/api run"]
# SIM
#CMD ["./cmd/api/api", "run"]
# SIM
#CMD ["sh", "-c", "./cmd/api/api run"]
# NÃO
# CMD ["echo $DB_NAME"]
# NÃO
#CMD ["sh", "-c", "echo $DB_NAME"]
# NÃO
# CMD [ "sh", /
#     "-c", /
#     "./cmd/api/api", \
#     "run" \
# ]
# SIM
# CMD [ "./cmd/api/api", \
#     "run" \
# ]
