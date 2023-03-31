# gowizard
Let a wizard guide you through your Go module set up quickly.
![Gopher Wizard](wizard.png)

> This project is still in development so use at your own risk!

## Installation
```bash
go get -u github.com/mahcks/gowizard
```

## Usage
Using the standard wizard:
```bash
gowizard
```

Quickly generate a new project with the needed services, controllers and adapters while bypassing the wizard:
```bash
gowizard generate --module github.com/username/module --path /path/to/module --mariadb --redis
```

## Services
- [ ] REST - Standard HTTP
- [ ] REST - fasthttp
- [ ] GraphQL

## Controllers
- [ ] REST
- [ ] gRPC

## Adapters
- [x] MariaDB - [github.com/go-sql-driver/mysql (v1.7.0)](https://github.com/go-sql-driver/mysql)
- [x] MongoDB - [go.mongodb.org/mongo-driver (v1.11.3)](https://github.com/mongodb/mongo-go-driver)
- [ ] MySQL
- [ ] PostgreSQL
- [x] Redis - [github.com/go-redis/redis/v8 (v8.11.5)](github.com/go-redis/redis/v8)
- [ ] SQLite

