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
- REST
    - [ ] gin
    - [ ] beego
    - [ ] fiber
    - [ ] fasthttp
- [ ] GraphQL

### Controllers
- [ ] REST
- [ ] gRPC

### Adapters
- [x] MariaDB - [github.com/go-sql-driver/mysql (v1.7.0)](https://github.com/go-sql-driver/mysql)
- [x] MongoDB - [go.mongodb.org/mongo-driver (v1.11.3)](https://github.com/mongodb/mongo-go-driver)
- [ ] MySQL
- [x] PostgreSQL - [github.com/jackc/pgx/v5 (v5.3.1)](https://github.com/jackc/pgx)
- [x] Redis - [github.com/go-redis/redis/v8 (v8.11.5)](https://github.com/redis/go-redis)
- [ ] SQLite


### Templates
There is no single way to structure things in Go. You may not like the way this is exactly structured, or you may have different needs than what this tool can satisfiy. With templates, you can swiftly create modules from external template repositories, and then modify them to your liking.

I have implemented a few templates that I have used in the past, all of which have inspired this tool, but I am sure there are more out there. If you have a template you would like to share, please open a pull request.

- [evrone](https://github.com/evrone) / [go-clean-template](https://github.com/evrone/go-clean-template/tree/master)
    - Clean Architecture template for Golang services
- [thangchung](https://github.com/thangchung) / [go-coffeeshop](https://github.com/thangchung/go-coffeeshop) 
    - A practical event-driven microservices demo built with Golang. Nomad, Consul Connect, Vault, and Terraform for deployment

*For more information about project layouts, I recommend: [golang-standards/project-layout](https://github.com/golang-standards/project-layout)*


### Structure
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
 ┃ ┃ ┗ zap.go
 ┣ README.md
 ┣ debug.log
 ┣ go.mod
 ┣ go.sum
 ┗ main.go
```

## Development
Rename `Makefile.local` to `Makefile`, change the variables at the top and run any of the commands to get started.

## Contributing
Pull requests are welcome. For major or breaking changes, please open an issue first to discuss what you would like to change. 