
### TODO APP

We perform CRUD and to-do list operations using Go-Postgres-Gin-Gorm-JWT.




## Directory structure

<a href="https://ibb.co/R3JrKFm"><img src="https://i.ibb.co/hBTQN6j/directory.png" alt="directory" border="0"></a>

.
 ├── main.go
 ├── api
 |    ├── server.go
 |    |
 |    ├── auth    
 |    |    └── token.go
 |    |
 |    ├── controllers
 |    |     └── base.go
 |    |     └── home_controller.go
 |    |     └── login_controller.go
 |    |     └── routes.go
 |    |     └── todo_controller.go
 |    |     └── users_controller.go
 |    |
 |    ├── middlewares
 |    |     └── middlewares.go
 |    |
 |    |
 |    ├── models
 |    |     └── todo.go
 |    |     └── user.go
 |    |
 |    |
 |    ├── responses 
 |    |     └── json.go
 |    |
 |    ├── seed
 |    |     └── seeder.go
 |    |
 |    ├── utils
 |    |     └── formaterror.go
 |    ├── ...
...
