# Chirpy: A CRUD Web Server That Simulates a Twitter-Like Experience
![License: MIT](https://img.shields.io/badge/License-MIT-red.svg)
## Overview

The **Chirpy** app is a Twitter-like server that connects to a **PostgreSQL** database. 
The web server is written in **Go**. It is a lighweight and efficient server that uses Go's standard
`net/http` library to handle requests with ease. It was designed for handling RESTful API requests, and
managing user data. It leverages **PostgreSQL** as the backend database to ensure robust and reliable data storage. The interface 
between the two is generated using [SQLC](https://github.com/sqlc-dev/sqlc).
## How to Install and Run 

### Prerequisites
- Go must be installed (download it from [golang.org](https://go.dev/dl/).
- Ensure your `$GOPATH/bin` is in your sytem's `PATH`.
- PostgreSQL must be installed (download it using your preferred package manager) 
- Run all the migrations in your database ([Goose](https://github.com/pressly/goose) makes this extremely easy).

### **Option 1: Install the Binary (Recommended)**
Run the following command to install the binary globally: 
```sh
go install github.com/tsyrdev/chirpy
```
Once installed, you can run the server using `chirpy`

### **Option 2: Build from Source**
1. Clone the repo: 
```sh
git clone https://github.com/tsyrdev/chirpy.git
cd chirpy
```
2. Build the project: 
```sh
go build -o chirpy 
```
3. Run the tool locally: 
```sh
./chirpy
```
4. If you want to run it globally, move the binary to your `$GOPATH/bin`:
```sh
mv chirpy ~/go/bin
```

## How to use it 

The `chirpy` server exposes the following endpoints for users to connect to:
- `POST /admin/reset` - Resets the users in the database.
- `GET /api/healthz` - Returns the status of the server.
- `POST /api/users` - Creates a new user. 
- `POST /api/chirps` - Creates a new chirp.
- `GET /api/chirps` - Gets all the chirps in the database.
- `GET /api/chirps/{chirpID}` - Gets the specified chirp.
- `DELETE /api/chirps/{chirpID}` - Deletes the specified chirp.
- `POST /api/login` - Logs into a user account. 
- `POST /api/revoke` - Revokes a user's access token.
- `PUT /api/users` - Updates a user's credentials.
- `POST /api/polka/webhooks` - Third-party connection for users to upgrade their membership.

## Credits

This project is from a Go Servers tutorial on [boot.dev](https://www.boot.dev/tracks/backend)

## License

This project is licensed under the [MIT License](LICENSE)
