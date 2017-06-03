# optimization_test
In this project i want to create a server in golang that connects to mySql and client, sends result over tcp to the server which saves it in mySql.

docker mysql:
docker run --name mysqldocker -e MYSQL_ROOT_PASSWORD="pass" -e MYSQL_USER="user" -e MYSQL_PASSWORD="pass" -e MYSQL_DATABASE="black_hole_test" -e MYSQL_ROOT_HOST="172.17.0.1" -p 3306:3306 -d mysql/mysql-server

