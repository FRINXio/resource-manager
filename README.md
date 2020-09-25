# Resource manager

## Build & Run & Test

Make sure you have access to symphony on github.
Generate an access token in github and place it into .env.

Execute build.sh and then run.sh

Note: See build.sh for git configuration required for using symphony as a dependency
Note: run.sh runs a mysql container

Subsequently, you can use test-queries.graphql to exercise the APIs

Note: .graphqlconfig can be used to point graphql client tools to the right direction

### Unit tests
```sh
go test -v -short ./pools/...
```

### Integration tests
Requirements: wasmer, quickjs and python packages available.
Example execution:
```sh
export WASMER_BIN=~/.wasmer/wasmer
export WASMER_JS=~/.wasmer/globals/wapm_packages/_/quickjs@0.0.3/build/qjs.wasm
export WASMER_PY=~/.wasmer/globals/wapm_packages/_/python@0.1.0/bin/python.wasm
export WASMER_PY_LIB=~/.wasmer/globals/wapm_packages/_/python@0.1.0/lib/
go test -run Integration ./pools/...
```

## TODOs

* !! Unify property (de)serialisation across the pool package
* !! graphql error handling ... no msg is returned to caller e.g. when creating a duplicate pool
* graphql.schema - add support for directives to better control entity field values
* !! pool - support labels
* ent - research and implement policies/privacy for entities based on user role (RBAC)
* docker-compose - create docker compose for graphql server + mysql
* security ?
* logging
* actions triggers ?

## Glossary

* Telemetry - tracing data streaming. Not used in RM
* Privacy/Policy - RBAC control over entity CRUD operations. Defined as part of ent.go schema
* Hook - Custom code invoked when interacting ent.go entity. Not used in RM
* Features - Special tags coming as HTTP headers that can grant additional permissions to users ? Not used in RM
* Directive - Custom extension to graphql schema for graphqlgen framework. Needs to be defined in the schema and implemented as go code. Example: IntRange restriction directive. Not used yet
* Actions -
* Triggers -
* Events -
* Jobs -
* Exporter - batch export/import of data (CSV). Not used in RM

## Features
List of important features of resource manager

### Building on FBC inventory
Lots of components and parts of the DB schema are reused from the inventory project.

### Model driven DB
Database schema is derived/generated from ent.go schema definition. Ent.go hides/handles all DB interactions. Ent schema can be found at ent/schema

**IMPORTANT** the underlying database needs to be set to the serializable isolation level. In the current database - MySQL - this is
done by adding the `--transaction-isolation=SERIALIZABLE ` starting parameter setting the isolation level on all sessions
 from clients. The reason for this is the currently implemented resource claiming algorithm running in a transaction.
 The algorithm has these steps:
 
 1. Prepare all necessary data (pool, strategy etc.) and prerequisite checks to claim a particular resource 
 2. Call a JS or Python script to generate a new resource
 3. Update the database with this particular resource
 
 The problem is if multiple clients claim the resource at the same time. This causes the same resources to be computed
 in step 2 and both clients update the database with the same resource (e.g. the same IP address). We need to isolate 
 transactions from each other (i.e. run them serially not in parallel because we might get the same resources otherwise) 
 thus increasing the isolation level.  


### Model driven graphql server
GraphQL server is derived/generated from graphql.schema. Code which ties graphql and ent together is written manually.
Schema needs to be kept in sync with ent.go DB schema, they are not connected in any automated way.

### APIs
Northbound APIs:

#### HTTP
Exposes grahpql API

#### webSockets
???

### Multitenancy
Multitenancy is supported throughout the stack.
In DB, each tenant has their own database. The database is created whenever a new tenant is detected.
GraphQL server switches to appropriate tenant context using TenantHandler baked into the HTTP API.

### RBAC
Privacy rules attached to ent.go schema definitions define the permissions. They can be anything from alwaysAllow, alwaysDeny, but usually they are tied to user role e.g. only a superuser can CUD and entity.

There are also additional optional features coming in as HTTP headers that can alter the permissions granted from user's role.

### Logging
zap logging framework is used. Main parameters allow control over logging level and format.

??? connection with events ? are the logs streamed ?

### Telementry
Support for tracing (distributed tracing). Streams data into a collector such as Jaeger.
Default is Nop.
See main parameters or telementry/config.go for further details to enable jaeger tracing

### Health
Basic health info of the app (also checks if mysql connection is healthy)

```
# server can serve requests
http://localhost:8884/healthz/liveness
# server works fine
http://localhost:8884/healthz/readiness
```

#### Metrics
Prometheus style metrics are exposed at:

```
http://localhost:8884/metrics
```

### Security
???

## Development

Prerequisites:
* go

### Wire

Install:
```
go get github.com/google/wire/cmd/wire
```

Generate wiring code:
```
wire ./graph/...
```

### ent.go

Generate ent.go entities:
```
go generate ./eng/...
```

### graphqlgen

Generate graphql resolvers:
```
go generate ./graph/graphql/...
```
