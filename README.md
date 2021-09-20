# auth-api
Gokit Basic API with JWT Authentication

The technologies implemented in this API are:
- Go
- Net/HTTP
- Gorilla/mux
- Gokit
- JWT
- MySQL

# System requirements

- GO version > 1.16
- Docker

# Get going

To start this app you will need to complete the following steps:

- On the base project directory, execute `go mod tidy` and then `go mod vendor` to donwload all the dependencies of the project

- You can use the pre-define SQL schema located in the db/schema.sql file, if you wanna use your own database, make sure to change the database URI variable (DB_URI) on the .env file

- Run the command `docker compose build` to build the docker images, then run `docker compose up` to start the application

And thats it, you are ready to GO :)



