# Orijinplus - Backend

## This projects is based upon Golang v1.16

## Description
This is an of implementation of Clean Architecture in Go (Golang) projects.

### Project Structure
This project has  4 Domain layer :
 * Models Layer --> models
 * Repository Layer --> store
 * Master Layer --> master
 * Usecase Layer   --> services
 * Delivery Layer --> api

_Dependencies_

- [go](https://golang.org/)
- [postgres](https://www.postgresql.org/)
- [golang-migrate](https://github.com/golang-migrate/migrate/releases)
- [gqlgen](https://gqlgen.com/)
- [docker](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
- [docker-compose](https://docs.docker.com/compose/install/)
- [solidity](https://github.com/ethereum/solidity)
- [geth](https://geth.ethereum.org)

_Included dependent binaries_

- [migrate](https://github.com/golang-migrate/migrate)
- [go-ethereum](https://github.com/ethereum/go-ethereum)


## Prerequisites

### Setup docker and docker compose
- Docker installation guide - https://docs.docker.com/engine/install/ubuntu/
- Docker Compose installation guide - https://docs.docker.com/compose/install/


### Setup Golang: go version go1.17.5 linux/amd64
- Golang setup guide - https://golang.org/doc/install
- `wget -c https://golang.org/dl/go1.17.5.linux-amd64.tar.gz`
- `sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.17.5.linux-amd64.tar.gz`
- `export PATH=$PATH:/usr/local/go/bin`
- Setup go path: in root directory, go to file .profile and paste the following line `export PATH=$PATH:/usr/local/go/bin`

#### Setup Golang with gobrew: alternative
- `curl -sLk https://git.io/gobrew | sh -`
- add following command in .bashrc file `export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"`



### Environment
- Environment variables are saved inside `/env` directory.
- Environment variables are accesssed by package `config`.


### Database

#### Setup Postgres
```bash
make postgres
make createdb
make migratedown
make migrateup
```

#### Create and Migrate DB
```bash
make createdb
make migratedown-local
make migrateup-local
```

#### Seed DB
```bash
make dbseed
```

#### Drop DB
```bash
make migratedown
make dropdb
```

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
\q
```

#### Generate GraphQL
```bash
make gqlgen
```
### Server

#### Test server

```bash
make test
```

#### Run server

```bash
make run
```

#### Build server

```bash
make build
```

## Docker

#### Build image

```bash
docker build -t orijinplus:latest .
```
