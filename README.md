# gowizard
Let a wizard guide you through your Go module set up quickly.

> This project is still in development so use at your own risk!

![Gopher Wizard](wizard.png)


## Installation
```bash
go install github.com/mahcks/gowizard@latest
```

## Usage
Using the standard wizard:
```bash
gowizard
```

Quickly generate a new project with the needed services, controllers and adapters while bypassing the wizard:
```bash
gowizard generate --module github.com/username/module --path /path/to/module --adapter mariadb,redis,mongodb
```

Using a template:
```bash
gowizard generate --module github.com/username/module --template github.com/valid/template --path /path/to/module
```

### Services
Each service has multiple "flavors" that can be used to generate the service. The following are the available flavors for each service.

#### REST
- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [beego/beego](https://github.com/beego/beego)
- [gofiber/fiber](https://github.com/gofiber/fiber)
- [valyala/fasthttp](https://github.com/valyala/fasthttp)

#### GraphQL
- [github.com/99designs/gqlgen](https://github.com/99designs/gqlgen)

### Controllers

#### REST
- [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- [github.com/beego/beego](https://github.com/beego/beego)
- [github.com/gofiber/fiber](https://github.com/gofiber/fiber)
- [github.com/valyala/fasthttp](https://github.com/valyala/fasthttp)

#### gRPC
- [github.com/grpc/grpc-go](https://github.com/grpc/grpc-go)

### Adapters
- MariaDB - [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- MongoDB - [go.mongodb.org/mongo-driver](https://github.com/mongodb/mongo-go-driver)
- MySQL - [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- PostgreSQL - [github.com/jackc/pgx/v5](https://github.com/jackc/pgx)
- Redis - [github.com/go-redis/redis/v8](https://github.com/redis/go-redis)

### Structure
This structure is by no means the official structure for Go projects; however, it is a set of [common historical and emerging project layout patterns in the Go ecosystem](https://github.com/golang-standards/project-layou).


The following project was generated using the following command:
```bash
gowizard generate --module github.com/user/module --adapter mariadb,redis
```

```md
generated-by-gowizard
 ┣ cmd
 ┃ ┗ cli
 ┃ ┃ ┗ main.go
 ┣ domain
 ┃ ┗ module.go
 ┣ templates
 ┃ ┣ adapters
 ┃ ┃ ┣ mariadb
 ┃ ┃ ┃ ┣ service.go
 ┃ ┃ ┃ ┗ template.go
 ┃ ┃ ┗ redis
 ┃ ┃ ┃ ┣ config.yaml
 ┃ ┃ ┃ ┣ service.go
 ┃ ┃ ┃ ┗ template.go
 ┃ ┗ logger
 ┣ README.md
 ┣ debug.log
 ┣ go.mod
 ┣ go.sum
 ┗ main.go
```

### Templates
There are excellent repositories already set up with specific project layouts that may be more to your liking. Templates allow you to choose from a list of external repositories with pre-configured project layouts.

- [evrone](https://github.com/evrone) / [go-clean-template](https://github.com/evrone/go-clean-template/tree/master)
    - Clean Architecture template for Golang services
- [thangchung](https://github.com/thangchung) / [go-coffeeshop](https://github.com/thangchung/go-coffeeshop) 
    - A practical event-driven microservices demo built with Golang. Nomad, Consul Connect, Vault, and Terraform for deployment

If the repository isn't listed, you may use `gowizard template --custom` to use a custom template. If you'd like the template to be added to the list, please open an issue.

> What makes this different from just cloning the repository? 

The wizard will ask you a few questions to help you get started with your project. It will also rename the module, use the optional path, and run a setup function if it's a pre-defined template that needs additional setup.

## Development
Rename `Makefile.local` to `Makefile`, change the variables at the top, and run any of the commands to get started.

## Contributing
Pull requests are welcome. For major or breaking changes, please open an issue first to discuss what you would like to change. 


todo
- [ ] generate makefile
- [ ] generate readme with just the first line being # module
add support for controllers
    - [ ] REST
    - [ ] gRPC
