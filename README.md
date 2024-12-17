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

## How to run
Prerequisites:


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