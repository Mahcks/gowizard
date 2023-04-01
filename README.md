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

### Services
- [ ] REST - Standard HTTP
- [ ] REST - fasthttp
- [ ] GraphQL

### Controllers
- [ ] REST
- [ ] gRPC

### Adapters
- [x] MariaDB - [github.com/go-sql-driver/mysql (v1.7.0)](https://github.com/go-sql-driver/mysql)
- [x] MongoDB - [go.mongodb.org/mongo-driver (v1.11.3)](https://github.com/mongodb/mongo-go-driver)
- [ ] MySQL
- [ ] PostgreSQL
- [x] Redis - [github.com/go-redis/redis/v8 (v8.11.5)](https://github.com/redis/go-redis)
- [ ] SQLite

## Development
Rename `Makefile.local` to `Makefile` and run `make` to build the project.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. 
*Please refrain from making PRs for services/controllers since I don't have those systems in yet. It's a top priority however: <issue>*

For more details read the [CONTRIBUTING.md](CONTRIBUTING.md) file.