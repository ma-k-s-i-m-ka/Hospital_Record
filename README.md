# Hospital_Record
  REST API application on GO for hospital operation with SOLID principles and clean architecture.
  
Сontains the following functionality:
 - CRUD systems for working with basic objects (users, doctors, appointment ...)
 - JWT-based authentication
 - Error Handling
 - Logging of the system operation
 - The ability to change the configuration

The application uses the following auxiliary and replaceable packages at your discretion:
 - Routing: [httprouter](https://github.com/julienschmidt/httprouter)
 - Database access: [pgx/v5](https://github.com/jackc/pgx)
 - Logging: [logger](https://github.com/google/logger/blob/master/logger.go)
 - JWT: [jwt-go](https://github.com/dgrijalva/jwt-go)

## Getting Started
  The server works at http://localhost:3000. Optionally, you can change the connection settings of both the server and the database in the config.yml file
## Testing
  Tested the application using POSTMAN. Folder with requests [Postman](https://drive.google.com/drive/folders/1Vmrq3W1DxLjh2Qcuo3HNCxI5Ll-u01pM?usp=sharing)
## Project Layout
```sh
.
├── app                  
│   ├── cmd                         main applications of the project
│   ├── internal
│   │    ├── config                 application configuration
│   │    ├── domain
│   │    │    ├── apperror          application-side error handler
│   │    │    ├── auth              authentication feature
│   │    │    ├── disease           working with disease
│   │    │    ├── doctor            working with doctor
│   │    │    ├── handler           route registration
│   │    │    ├── portfolio         working with portfolio
│   │    │    ├── record            working with record
│   │    │    ├── response          error handler from the client side
│   │    │    ├── specialization    working with specialization
│   │    │    └── user              working with user
│   │    ├── http/db                postgresql database
│   │    └── server                 the API server application
│   ├── pkg/logger                  application logging system
│   └── doctorimages                image storage
└── logs                            application files
```
