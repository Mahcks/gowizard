# gowizard
Let a wizard guide you through your Go module set up quickly.
![Gopher Wizard](wizard.png)

> This project is still in development so use at your own risk!

## Services
- [ ] REST - Standard HTTP
- [ ] REST - fasthttp
- [ ] GraphQL

## Controllers
- [ ] REST
- [ ] gRPC

## Adapters
- [x] MariaDB
- [x] Redis
- [ ] MongoDB
- [ ] PostgreSQL
- [ ] MySQL
- [ ] SQLite

Start the wizard:
```bash
gowizard
```

Quickly generate a new project with a MariaDB and Redis database:
```bash
gowizard generate --module github.com/username/module --path /path/to/module --mariadb --redis
```