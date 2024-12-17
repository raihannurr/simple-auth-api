# Simple Auth REST API
### Author: [Raihannur Reztaputra](https://www.linkedin.com/in/raihannurr)

## API List & Functionalities
- [x] CSRF Token Validation
- [x] Register New User
- [x] Login
- [x] Get user profile
- [x] Update user profile

## Prerequisites
- [Go](https://golang.org/dl/) version 1.23.4 or higher
- [Make](https://www.gnu.org/software/make/)
- [MySQL](https://www.mysql.com/) 8.0 or higher, running server (optional, you can also run this service by using `DB_ADAPTER=memory` to run simple in-memory database)
- [Docker](https://www.docker.com/), to build the docker image
## Project Structure

```
.
├── api
├── bin/
│   └── rest-api -> executable 
├── cmd/
│   └── rest-api/
│       └── server.go -> main golang function
├── deployments
│   └── rest-api/
│       └── Dockerfile -> dockerfile for building the rest-api
├── internal/
│   ├── config/ ->
│   │   ├── config.go -> defining the config for the application
│   │   └── config_test.go
│   ├── db/ -> storing database migration files
│   │   ├── {version}_{name}.up.sql -> for applying database migration
│   │   └── {version}_{name}.down.sql -> for reverting/rollback database migration
│   ├── entity/ -> storing entity objects
│   │   ├── user.go -> definition of user entity object
│   │   └── user_test.go
│   ├── errors/ ->
│   │   └── errors.go -> for registering errors
│   ├── repository/
│   │   ├── mocks/ -> storing mock functions
│   │   │   └── mock_repository.go -> mock for repository interface, generated using gomock mockgen
│   │   ├── repository.go -> define the interface for repository
│   │   ├── mysql.go -> implementation of mysql repository
│   │   └── simple_struct.go -> implementation of in-memory simple struct repository
│   ├── server/
│   │   ├── handler/ -> storing handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── user_handler.go
│   │   │   ├── csrf_handler.go
│   │   │   └── ...
│   │   ├── middleware/ -> storing middlewares
│   │   │   ├── authenticate.go
│   │   │   ├── validate_csrf.go
│   │   │   └── ...
│   │   ├── router.go -> file to register routes in the application
│   │   └── router_test.go -> unit test
│   └── utils/ -> storing utility functions
│       └── ...
├── .env.sample
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## How to run

```bash
# create .env file from the example
cp .env.sample .env
# ensure SESSION_SECRET_KEY is filled with the correct size for AES encryption

# run the server,
#ensure running this command on the same directory as the Makefile and the .env file
make start
```

## How to test
```bash
make test
```

## How to generate coverage report
```bash
make coverage
```

## How to migrate database
Prerequisites:
- Ensure you have the correct .env file
- Ensure you have the correct database credentials
- Ensure you have the correct database is ready to be migrated
- Ensure you have already created the database in the database server
- Install the migration tool: [golang-migrate](https://github.com/golang-migrate/migrate), or by running `make install-dependencies`

```bash
make db-migrate
# this will run all the up.sql files in the internal/db/migrations/ directory
```

## How to rollback database
```bash
make db-rollback
# this will rollback the last migration
```

## Build Docker Image
```bash
make docker-build
```