# Simple Banking Backend

Backend transaction system deployed on AWS for testing.

### Concepts
* Go
* Postgres
* Redis
* Gin
* gRPC
* Docker
* Kubernetes
* AWS
* CI/CD

## Getting Started

### Prerequisites

Install Tools
- [Visual Studio Code](https://code.visualstudio.com/download)
- [Golang](https://golang.org/)
- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [DB Docs](https://dbdocs.io/docs)
- [Gomock](https://github.com/golang/mock)
- [Postman](https://www.postman.com/downloads/)
- Sign up for [Amazon Web Services](https://aws.amazon.com/)

### Setup

1. Clone the repo
   ```
   git clone https://github.com/xiao-yangg/simple-bank.git
   ```
2. Open using Visual Studio Code
3. Create network
   ```
   make network
   ```
4. Make Postgres container in Docker
   ```
   make postgres
   ```
5. Create database with MySQL
   ```
   make createdb
   ```
6. Run db migration up to latest version
   ```
   make migrateup
   ```
7. Run unit tests:
    ```
    make test
    ```
8. Run server:
    ```
    make server
    ```
9. Test with Postman / any frontend

## Acknowledgement
https://www.udemy.com/course/backend-master-class-golang-postgresql-kubernetes/
