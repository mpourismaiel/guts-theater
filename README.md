# Guts Theater

## Project description

This project is an api created using go-chi and couchdb, to provide a way to
seat multiple customer groups inside a theater hall. The project is written in
Go and uses Docker for deployment.

## Usage

As the project doesn't require any explicit configuration, you can start a
couchdb instance on port `5984` and build and run the project. For client access
[another repository](https://github.com/mpourismaiel/guts-theater-fe) has been
provided which is written in React.

In order to ease the deployment process, a `docker-compose.yml` configuration
has been provided which can be run using:

```
docker-compose up --build
```

This docker-compose will run the api project (current repository), the client (
granted that the client project is stored in the same directory as this repo),
a couchdb instance, and Prometheus and Grafana to provide instrumentation.

After starting the docker-compose, you can visit the [client](http://localhost:4005)
: `http://localhost:4005` or visit the [api](http://localhost:4000): `http://localhost:4000`.

To make sure the API is online, please visit [healthz](http://localhost:4000/healthz).

In order to test the functionality of the application, please:

1. Create at least one section
2. Create at least one row
3. Create at least one seat
4. Create at least one group
5. Call the trigger seating API

## API

The API is written using go-chi router and is located in the `/api` directory.
The API follows the rest pattern, providing CRUD operations for all the models (
some operations have been omitted to prevent breaking the flow of the application).

```
Get     /seats
Get     /seats/{section}

Get     /section
Post    /section
Put     /section/{section}
Delete  /section/{section}

Get     /section/{section}/rows
Post    /section/{section}/row
Delete  /section/{section}/row/{row}

Get     /section/{section}/seats
Post    /section/{section}/row/{row}/seat
Put     /section/{section}/row/{row}/seat/{seat}
Delete  /section/{section}/row/{row}/seat/{seat}

Get     /groups
Post    /groups

Get     /ticket
Get     /ticket/{groupId}

Post    /trigger-seating
```
