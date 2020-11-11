# TODO List
Web service that stores a list of things to do

## Installation Pre-requisites
- Golang
- Docker

## Running Locally
### Install Build Dependencies
```
go mod download
```

### Create MySQL Database
Create MySQL instance:
```
docker run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=root mysql
```
Create database called `todolist`:
```
docker exec -it mysql mysql -uroot -proot -e 'CREATE DATABASE todolist'
```
Verify `todolist` database creation:
```
docker exec -it mysql mysql -uroot -proot -e 'SHOW DATABASES'
```
This should result in similar output with `todolist` present in database list:
```
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
| todolist           |
+--------------------+
```

### Stand Up Service
Build service:
```
go build cmd/main.go
```
Run service:
```
./main
```

### Create requests
Available endpoints:
```
GET /todos
GET /todos/complete
GET /todos/incomplete
POST /todo
POST /todo/{id}
DELETE /todo/{id}
```
Example of creating todo item:
```
curl -X POST localhost:8000/todo -d "description=Wash dishes"
```
