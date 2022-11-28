# new-go-code-challenge-template-2-2

# Hi there! 👋

Be very welcome to my solution to X's code challenge.

- [Introduction](#introduction)
- [Architecture](#architecture)
- [Database](#database)
- [How to run the project?](#how-to-run-the-project)
- [API documentation](#api-documentation)
- [Test cases](#test-cases)
- [How to run the tests?](#how-to-run-the-tests)
- [Deployment](#deployment)
- [How to deploy the project?](#how-to-deploy-the-project)
- [References](#references)

## Introduction

This project consists of the development of a **REST API** using **Go** programming language, **Json Web Token** and **Postgres** database for managing authentication operations and accessing users data.

## Architecture

The architecture of the project was designed according to my understanding and my code structuring decisions based on my research of the concepts of **Domain Driven Design** and **Hexagonal Architecture**.

### Domain Driven Design

This approach is intended to simplify the complexity developers face by connecting the implementation to an evolving model.

To do it, the implementation is basically divided up into the following essential layers in order to have a separation of interests by arranging responsibilities:

#### Application

This layer is responsible for serving the application purposes. It contains services (or use cases) that are used to implement the business logic acting as intermediaries for communication between the repositories and handlers.

In this way, the services represent the implementation of business logic, regardless of the type of database used, or how the service will be exposed externally (http or grpc, for example).

Also, they include the validation of the input parameter values from the API requests payloads.

#### Core/Domain

This layer is resposible for holding the schema of entities and ports used for the communication between the handlers and services, as well as between the services and repositories.

#### Infrastructure

This layer is responsible for serving as a supporting layer for other layers.

It contains the procedures to establish connection to the database and the implementation of repositories that interact with the database by retrieving and/or modifing records.

#### Interfaces

This layer is responsible for the interaction with user by accepting API requests, calling out the relevant services and then delivering the response.

It contains the handling of requests by exposing the routes associated with each API endpoints, applying authentication actions when needed that mediate the access to them, as well as the elaboration of API responses.

### Hexagonal Architecture

This approach (also known as Ports and Adapters pattern) allows creating an application where the business logic is in a core (*core*) and there is no dependence on external systems, thus facilitating the development of regression tests.

It was designed in such a way that adapters (*adapters*) can be "plugged" (*dependency injection*) into the system from ports (*ports*), not affecting the business logic that was defined in the system's core.

Dependency injection is a technique where adapters are plugged in with their respective ports and that can be used to inject the dependencies of a class into the class. It helped to keep the code simple and easy to understand. Also, it facilitates the development of tests by mocking dependencies.

In this context, it was enabled the use of Ports represented as interfaces that contain the signatures of the methods that are used by the adapters, in order to perform the desired operations.

Essentially, the interfaces are implemented by services and repositories placed in application and infrastructure layers, respectively, that belong to the nucleus and define how the communication between the nucleus and actors that want to interact with it are carried out; and adapters that were responsible for translating the information between the core and these actors.

Adapters are implemented in the infrastructure (known as repositories) and interface layers (known as handlers) and are responsible for http communication and database communication, respectively.

Such structuring of the code makes it possible to focus on the implementation of business logic, since it can be developed completely independently from the rest of the system, as well as on the separation of dependencies, the ease of changing the infrastructure (such as a change of a database), and even allows tests in isolation to be carried out in a simple way.

## Database

Two Postgres dabases need to be configured to use the project. One of them is intended to common (or usual) use and the other is directed to test execution. However, both of them contain the same tables named auths, logins and users defined in the **database/scripts/1-create_tables.sql** file.

### Tables

**Auths**

The **auths** table contains the authentication data.

| Fields     | Data type | Extra                       |
|:-----------|:----------|:----------------------------|
| id         | UUID      | NOT NULL PRIMARY KEY        |
| user_id    | UUID      | NOT NULL UNIQUE FOREIGN KEY |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP   |

**Note**:

A record is created in this table whenever a user performs login and this same record is deleted as soon as the related user performs logout.

**Logins**

The **logins** table contains the users credentials.

| Fields     | Data type | Extra                       |
|:-----------|:----------|:----------------------------|
| id         | UUID      | NOT NULL PRIMARY KEY        |
| user_id    | UUID      | NOT NULL UNIQUE FOREIGN KEY |
| username   | TEXT      | NOT NULL                    |
| password   | TEXT      | NOT NULL                    |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP   |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP   |

**Users**

The **users** table contains the users data.

| Fields          | Data type | Extra                     |
|:----------------|:----------|:--------------------------|
| id              | UUID      | NOT NULL PRIMARY KEY      |
| username        | TEXT      | NOT NULL UNIQUE           |
| created_at      | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| updated_at      | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |

## How to run the project?

The project can be run either **locally** or using a [**Docker**](https://www.docker.com/) container. However, in order to facilitate explanations, this documentation will focus on running using a Docker container.

### Makefile file

A **Makefile** file was created as a single entry point containing a set of instructions to run the project in these two different ways via commands in the terminal.

Furthermore, this file also contains a series of routines used throughout the development of the project, such as reformatting the **.go** file and printing style errors, generating API documentation, creating *mocks* used in tests of the solution, among others.

To run the project with a Docker container, run the command:

```
make startup-app
```

Note:

- The **.env** file contains the environment variables used by the Docker container. However, it is not necessary to make changes to this file before running the project, so the variables can be kept as they are defined.

To close the application, run the command:

```
make shutdown-app
```

## API documentation

### API endpoints

The API *endpoints* were documented using the Github repository called [swaggo/swag](https://github.com/swaggo/swag) which converts code annotations in **Go** into **Swagger 2.0** documentation based on **Swagger** files located in the **docs/api/swagger** directory.

After running the project, access the following URL through your web browser to view an HTML page that illustrates the information of the API *endpoints*:

```
http://{host}:8080/swagger/index.html
```

### Postman Collection

To support the use of the API, it was created the file **new-go-code-challenge-template.postman_collection.json** in the directory **docs/api/postman_collection** which contains a group of requests that can be imported into the **Postman** tool (an API client used to facilitate the creation, sharing, testing and documentation of APIs by developers.).

## Test cases

The test cases were designed as [**Table Driven Tests**](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests) so that each test case was built by declaring a structure that contains actions that can be performed before and after executing them, as well as expected inputs and outputs, following the **unit** and **integration** tests approaches.

### Unit Tests

The unit tests are located inside the **internal** and **pkg** directories at the project root.

They are evaluated using the **Black-Box** testing strategy, where the test code is not in the same package as the code under evaluation.

For this, files were created with the suffix **_test** added to their names and also to the names of their test packages. For example, the code from the package (*pkg*) called **validator** is tested by a file called **validator_test.go**, which is defined in another package, called **validator_test**.

The separation of codes into distinct packages aims to ensure that only the identifiers exported from the packages under evaluation are tested. By doing this, the test code is compiled as a separate package and then linked and run with the main test binary.

#### Mocks

Some of the tests were written using mock objects in order to simulate dependencies so that the layers could interact with each other through **interfaces** rather than concrete implementations, made possible by the *design pattern* of **Dependency Injection**.

Basically, the purpose of mocking is to isolate and focus on the code being tested and not on the behavior or state of external dependencies. In simulation, dependencies are replaced with well-controlled replacement objects that simulate the behavior of real ones. Thus, each layer is tested independently, without relying on other layers. Also, you don't have to worry about the accuracy of the dependencies (the other layers).

For the mocking purpose, the Github repositories called [DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) e [vektra/mockery](https://github.com/vektra/mockery) were used for mocking the SQL driver behavior without needing to actually connect to a database and for generating the mock objects from interface, respectively.

### Integration Tests

The integration tests are located inside the **tests/api** directory at the project root.

They were written by combining and testing the project layers together to simulate the production environment.

Note:

- The unit and integration tests check a large and relevant part of the different components of the solution, but not all of them. In addition, not all tests written have **100%** coverage of the tested code.

## How to run the tests?

Before running the project tests, it is needed to start up the Docker containers named **api_container** and **postgrestestdb_container** successfully.

The **postgrestestdb_container** container is necessary to execute the integration tests and it can be initialized by running the command:

```
make start-deps
```

After all these containers are successfully initialized, to execute the tests of the project, run the command:

```
make test-app
```

After running any of the tests, it is possible to check the percentage of code coverage that is met by each test case displayed in the test execution output.

The statistics collected from the run are saved in the **docs/api/tests/unit/coverage.out** file for coverage analysis. To check the **unit** test coverage report informed in the **coverage.out** file, run the command:

```
make analyze-app
```

Notes:

- The **coverage.out** file contains only **unit** test execution statistics. (There are no statistics on the execution of the **integration** tests.)

## Deployment

The project was deployed as a Docker container on the **Heroku** hosting service using the infrastructure tool as code **Terraform**.

The API endpoints can be accessed from the hosted project using the following base URL:

```
https://icaroribeiro-<something>.api.herokuapp.com
```

For example, to check API documentation through your web browser, go to the following URL:

```
Method: HTTP GET
URL: https://icaroribeiro-<something>.api.herokuapp.com/swagger/index.html
```

To validate the application status, check the result of the request to the following API endpoint:

```
Method: HTTP GET
URL: https://icaroribeiro-<something>.api.herokuapp.com/status
```

## How to deploy the project?

Below there are the procedures used to deploy the project.

Note:

- To proceed with the deployment process it is necessary to have a Heroku account and the [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli) and [Terraform](https://www.terraform.io/downloads.html) softwares installed on the local machine.

First, it were created a manifest file called **heroku.yml** and a file named **Dockerfile.multistage** in the **deployments/heroku/app** directory that was used to build the project as a Docker container.

The Docker file was designed using a Docker stages feature that allows creating multiple images in the same Dockerfile (Docker's multi-stage image):

In summary, the first **FROM** statement of this file is related to an image that uses an alternative name "as constructor" to be referred to later in the file and the respective code contains the generation of the middle layer where the compilation of Go takes place; and the second **FROM** statement is directed to an image of [**alpine**](https://hub.docker.com/_/alpine) where simply the instruction "--from=builder" is defined to get the executable from the intermediate layer.

The objective of this approach is to build a final image that is as lean as possible, that is, with reduced size, containing only the binary application and the base operating system necessary to run it. This way, the application could be deployed quickly, even under slow network conditions.

Next, it was needed to create a file called **deployments/heroku/scripts/setup_env.sh** that have the following environment variables:

```sh
#!/bin/bash

#
# Heroku platform settings
#
export TF_VAR_heroku_email=<heroku_email>
export TF_VAR_heroku_api_key=<heroku_api_key>

#
# Heroku application settings
#
export TF_VAR_heroku_app_name=<heroku_app_name>
```

The first and second variables above are related to the Heroku Platform API settings and refer to the email address of Heroku account and a Heroku API key, respectively. The third variable refers to the name of the application that will be hosted on Heroku.

After installing the Heroku CLI software, to get a Heroku API key, run the Heroku CLI command:

```
heroku login
```

This way, you will be redirected to the web browser so that you can login to the Heroku website. After that, run the command to get the Heroku API key:

```
heroku auth:token
```

After obtaining the Heroku API key, configure the values ​​of the environment variables defined in the **setup_env.sh** file.

Then, run the below commands located in the Makefile file.

To initialize everything Terraform requires to provision the infrastructure, run the command:

```
make init-deploy
```

The previous command downloads the plugin from the Heroku provider and stores it in a hidden .terraform folder.

The infrastructure resources referring to the API were defined in the **deployments/heroku/terraform/resources.tf** file.

Then, for details on what will happen to the infrastructure without making any changes to it, run the command:

```
make plan-deploy
```

To make the necessary changes to reach the desired state of the configuration, run the command:

```
make apply-deploy
```

After applying the changes, it is possible to set up the database tables by means of CLI Heroku commands.

Firstly, to identify what is the identifier of the Heroku Postgres database, execute the command:

```
heroku pg:info -a=<HEROKU_APP_NAME>
```

The output should look something like this:

```
=== DATABASE_URL
...
Add-on: <HEROKU_POSTGRES>
```

Then, to create the database tables, run the command:

```
heroku pg:psql -a=<HEROKU_APP_NAME> <HEROKU_POSTGRES> < database/postgres/scripts/1-create_tables.sql
```

Finally, in order to terminate all the provisioned infrastructure components, run the command:

```
make destroy-deploy
```

### Accessing remote Postgres database locally

It is possible to verify the data generated in Heroku using [pgAdmin](https://www.pgadmin.org/) tool.

To achieve this, firstly access the Heroku website in order to check the datastore settings.

In the Settings tab, click the View Credentials... button and take note of the following data: **Host**, **Database**, **User**, **Port** and **Password**.

In what follows, there are the steps to configure a remote server in pgAdmin and to establish access to the Postgres database using the previous data:

In pgAdmin, right click Server(s) icon, and then navigate to Create and Server options.

After that, it is necessary to fill out the following parameters:

In the General tab, name the server whatever you want.

Under the Connection tab, inform the Host name and port. The Host name is the one configured like ...amazonaws.com and the port is 5432. In Maintenance database, informe the Database name from the previous data and do the same procedure to fill out the Username and Password fields.

In the SSL tab, mark SSL mode as Require.

Before finalising, it is necessary to apply one more configuration:

The database name needs to be informed in a "desired database list" in order to avoid parsing many other databases that are not cared about. (This has to do with how Heroku configures their servers.)

In this regard, go to the Advanced tab and under DB restriction copy the Database name (It's the same value filled in the Maintenance database field earlier), and then click Save button.

Finally, navigate through the options structure: Databases, database name, Schemas, public and inside Tables, check the tables. (In case of the tables are not displayed, try right click the related Server created and then click Refresh option.)

**Note**

This project was configured with a Heroku Postgres database resource in a **Free plan** (Hobby Dev - Free). Because of that, the database will only support a limited number of records (10.000 rows). Therefore, please evaluate the operations to be carried out before using the application in this way.

## References

Project layout:

- https://github.com/golang-standards/project-layout

Domain Driven Design

- https://dev.to/stevensunflash/using-domain-driven-design-ddd-in-golang-3ee5

Hexagonal Architecture

- https://medium.com/avenue-tech/arquitetura-hexagonal-com-golang-c344411aa940

Dependency Injection

- https://medium.com/avenue-tech/dependency-injection-in-go-35293ef7b6

Database Transaction

- https://medium.com/wesionary-team/implement-database-transactions-with-repository-pattern-golang-gin-and-gorm-application-907517fd0743