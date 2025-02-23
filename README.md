# transcation-processor
Multi-Tenant POS Transaction Processor


volume create
 docker volume create db-vol

network create
 docker network create -d bridge net
 
container run
 docker run --detach --name db-server --hostname db --network net -p 3306:3306 -p 8080:8080 -e MYSQL_ROOT_PASSWORD=pass -v db-vol:/var/lib/mysql mysql:8

container enter
 docker exec -it db-server mysql -u root -p
 password

    database setup
    CREATE DATABASE transaction_processor_db;
    CREATE USER 'dbuser'@'%' IDENTIFIED BY 'pass';
    GRANT ALL PRIVILEGES ON transaction_processor_db.* TO 'dbuser'@'%';
    FLUSH PRIVILEGES;

image build
 docker build --tag transaction-processor-rest .

container run
 docker run -it --detach --network net --name rest-server --publish 80:8080 --env DSN="root:pass@(db:3306)/transaction_processor_db?parseTime=true" transaction-processor-rest

 ctrl+c or cmd+c


 docker rm
 docker stop
 docker restart

 # app app-db-net db

test - 1gb
docker build -f Dockerfile -t transaction-processor-test --progress plain --no-cache --target run-test-stage .



find todos
encapsulate workers

sales limit 10? top-selling
encapsulate logger
encapsulate handlers
flags to read DSN & REDIS ADDR instead of environment variables
flags to read configs
using root

Q&A
is it safe to commit the .env file?
no, this is to make the running steps easier. in production you shouldn't commit .env files.

- retry
- document
- run

next
- kafka
- grpc
- tests
- rate limit
- authentication
- authorization
- physical isolation
- logical isolation