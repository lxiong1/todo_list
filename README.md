# TODO List
Web application that stores a list of things to do

## Installation Pre-requisites
- Golang
- Glide
- Docker

## Running Locally
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
