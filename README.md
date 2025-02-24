# transcation-processor
Multi-Tenant POS Transaction Processor

# Instructions

## Run
 - Install docker
 - Run docker
 - Go the project directory
 - type "docker compose up --build"
(it will take around 5-10 minutes downloading, building, running images)

## Shutdown
 - type "docker compose down"
 - if you want to erase the data(volume) type "docker compose down -v"

## API

# System architecture

# Project structure

# Database modeling

# Optimizations

# Concurrency techniques used

# Trade-offs made between caching, transaction processing, and API design

# Todos
 - Refactor
    - encapsulate workers
    - encapsulate logger
    - add flags for configs
    - replace environment variables with flags
    - Move /metrics to a different/separate port e.g. 8081 (internal usage)
 - Send email
 - Loom demo

# Q&A

Is it safe to commit the .env file?
no, this is to make the running steps easier. in production you shouldn't commit .env files.


# Next
- Kafka/RabbitMQ
- gRPC
- Kubernetes
- Load Testing
- Unit & Integration Testing
- Rate Limit
- Authentication & Authorization
- Physical & Logical Isolation
- Frontend (React, Hotwire)
