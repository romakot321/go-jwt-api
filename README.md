# REST API in Golang with JWT authentication

## Frameworks
- Fiber (web-framework)
- gorm (sql ORM)
- golang-jwt
- viper (for environment)

## Run
- By Docker: `docker compose up -d --build`
- By native golang: `go mod download && go run cmd/app/app.go`

## Features
- Email warning on IP change while refresh requested(only mock now).
- JWT tokens with scopes(for access and refresh operations)
- Postgresql database
- Project easy to expand

## Endpoints
V1 endpoints is not recommended to use
- POST `/api/auth/v1/login` - Login by GUID. Parameters passed by query!
- POST `/api/auth/v1/refresh` - Refresh access token by old access token and refresh token. Refresh token is one-time. Parameters passed by query!

- POST `/api/auth/register` - Register an account
- POST `/api/auth/login` - Login to account for tokens
- POST `/api/auth/refresh` - Refresh access token by refresh token
- GET `/api/user/me` - Get user info(GUID and username) by access token(passed as Bearer token)

## Tests
Tests written in Python3 and using _requests_ library. On error simply print error, on full pass print all tests.
Before running: `pip install requests`
