# optimization_test
In this project i want to create a server in golang that connects to mySql, client that is sending result over tcp with gRPC to the server and reader for viewing the results.


docker mysql:
docker run --name mysqldocker -e MYSQL_ROOT_PASSWORD="pass" -e MYSQL_USER="user" -e MYSQL_PASSWORD="pass" -e MYSQL_DATABASE="black_hole_test" -e MYSQL_ROOT_HOST="172.17.0.1" -p 3306:3306 -d mysql/mysql-server
