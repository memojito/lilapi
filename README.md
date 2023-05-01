# lilapi
Just trying out things with a simple api.
### How to run
To run server & cassandra
```sh
docker-compose up --build
```
To run cassandra only
```sh
docker-compose -f docker-compose-cassandra.yml up
```
Enter cassandra container
```sh
docker exec -it lilcassandra bash
```
In container terminal execute the **db-init.sh** which creates keyspace and schema
```sh
sh /docker-entrypoint-initdb.d/db-init.sh
```
Then run the app
```sh
go run server/main.go
```
