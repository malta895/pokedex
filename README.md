# Pokedex

A simple web API interface to retrieve information about your Pokemons!

- [Pokedex](#pokedex)
  - [How to run the project](#how-to-run-the-project)
    - [Download the Source Code](#download-the-source-code)
      - [Clone with Git](#clone-with-git)
      - [Download the zip archive](#download-the-zip-archive)
    - [Build and run with Docker](#build-and-run-with-docker)
    - [Build from source](#build-from-source)
      - [Build](#build)
      - [Run](#run)
      - [Testing](#testing)
  - [Usage](#usage)
    - [Basic Pokemon Information](#basic-pokemon-information)
    - [Translated Pokemon Information](#translated-pokemon-information)
  - [Project Design and Architecture](#project-design-and-architecture)
  - [Production-Ready Considerations](#production-ready-considerations)
    - [Containerization and Containers Orchestration](#containerization-and-containers-orchestration)
    - [API Gateway](#api-gateway)
    - [Caching](#caching)
    - [Logging and Monitoring](#logging-and-monitoring)
    - [Authentication and Authorization](#authentication-and-authorization)


## How to run the project

You can either build a Docker image, or build the source code and run the project locally.

Both require you download a local copy of the source code.

### Download the Source Code

This project source code is hosted on this git repository.

Binaries are not provided, so you have to download the source code and build it yourself.

You can also build the project with Docker.

Further instructions are provided below.

#### Clone with Git

If you want to clone the project with git, first make sure [Git](https://www.git-scm.com/) is installed on your operating system;

Then, run in a terminal:

```bash
git clone git@github.com:malta895/pokedex.git
```

The project files will be available in the `pokedex/` folder.

#### Download the zip archive

You can also download the [source zip archive](https://github.com/malta895/pokedex/archive/refs/heads/main.zip).

Using a zip compatible archive manager, extract it to a folder, naming it, for example, `pokedex/`.

### Build and run with Docker

Make sure [Docker](https://www.docker.com/get-started/) is correctly installed on your operating system.

Once you have downloaded the project, open up a terminal in the project folder.

The Docker image is not available on any registry, and it must be built locally; a Dockerfile is available in the project root for the purpose.

Build the project:

```bash
docker build . -t my-pokedex
```
This will build an executable and create a local image in your machine.

To run it:
```bash
docker run --name pokedex -p 8080:8080 my-pokedex
```

This command will expose the HTTP server built within the project on your local network, on the port 8080.

You can reach it at the address:
```
http://localhost:3000
```

You can now contact the exposed APIs; follow the instructions at the [Usage section](#usage).

### Build from source

To build from source, the [Go programming language](https://go.dev/doc/install) has to be installed on your system. 

This project requires at least Go 1.22.

#### Build

```bash
go build .
```
An executable file will be created in the project folder, named `pokedex` or `pokedex.exe` depending on the platform.

#### Run

On unix-like systems (e.g. Linux, MacOS, WSL):

```bash
./pokedex
```

On Windows:

```batch
pokedex.exe
```

Once spinned up, the executable will expose the API service at `http://localhost:3000`.

To expose the project on a different port, you can set the `HTTP_PORT` to a value of your choice.

For example, to use the port 8080, run on a unix-like OS:

```bash
HTTP_PORT=8080 ./pokedex
```

You can achieve the same in Windows with:

```batch
set HTTP_PORT=8080
pokedex.exe
```

#### Testing

To run the tests, you can use the `go test` command.

```bash
go test ./...
```

Integration tests are not run by default for performance reasons, as they require an internet connection to run, and they can concur to the rate limiting of the external APIs.
To run them, you can use the `-tags=integration` flag.

```bash
go test -tags=integration ./...
```

The above command will run the integration tests, along with the unit tests.

## Usage

The project exposes two endpoints, described below in the dedicated paragraphs.

### Basic Pokemon Information

Endpoint signature: `GET /pokemon/{pokemonName}`

Retrieves information about a Pokemon, given its name, by leveraging the [PokeAPI](https://pokeapi.co/).

Example usage:
  
  ```bash
  curl http://localhost:3000/pokemon/pikachu
  ```

Example response:

```json
{
  "name": "pikachu",
  "description": "When several of these POKéMON gather, their electricity could build and cause lightning storms.",
  "habitat": "forest",
  "isLegendary": false
}
```
  

### Translated Pokemon Information

Endpoint signature: `GET /pokemon/translated/{pokemonName}`

Retrieves information about a Pokemon, given its name, and translates its description, by leveraging the [PokeAPI](https://pokeapi.co/) and the [Fun Translations API](https://funtranslations.com/).

Example usage:
  
  ```bash
  curl http://localhost:3000/pokemon/translated/pikachu
  ```

Example response:

```json
{
  "name": "pikachu",
  "description": "At which hour several of these pokémon gather,  their electricity couldst buildeth and cause lightning storms.",
  "habitat": "forest",
  "isLegendary": false
}
```  

## Project Design and Architecture

The project is a simple web API service, written in Go.

Given the trivial nature of the project, the architecture is very simple, yet attention has been given to the separation of concerns, and to the testability of the code.

Here is the package structure of the project:

```
.
├── apiclients
│   ├── funtranslations
│   │   ├── client.go
│   │   └── client_test.go
│   └── pokeapi
│       ├── client.go
│       ├── client_test.go
│       ├── doc.go
│       └── pokeapi.go
├── main.go
├── main_test.go
├── pokemonmux
│   ├── mux.go
│   └── mux_test.go
├── testutils
│   └── testutils.go
└── types
    └── types.go
```

For separation of concerns, the external API clients have been placed in the `apiclients` package, and the HTTP server has been placed in the `pokemonmux` package.

The API clients expose simple interfaces, that has been mocked in the tests, to allow for easy testing of the server.

The `pokemonmux` package contains the HTTP server, that uses the Go standard library `net/http` `ServeMux` to handle the incoming requests.

For simplicity, handlers contains some business logic, such as the code to decide which translation type should be used. 
This could be extracted in a separate package, but given the simplicity of the project, it has been simply left in the handlers package, taking care however to separate it in a different function.

The project has been tested with unit tests, that mock the external API clients, and with integration tests, that test the server with the real external API clients.

## Production-Ready Considerations

This program has been developed for a simple local demonstration, it is not production-ready.

To make it so, several aspects must be addressed, explained in details in the following paragraphs.

### Containerization and Containers Orchestration

The project is ready to be shipped as a Docker container, in a production environment it can be deployed with a Container Orchestrator software; for example, it can be deployed on Kubernetes.

To be deployed on Kubernetes, the service should be modified to expose a health check endpoint, to allow Kubernetes to check the health of the service, and to restart it in case of a failure.
This ensures that the service is always up and running, and that it can recover from failures.

By leveraging the features offered by this container orchestrator, with a proper setup, the service can be scaled to serve a very high number of end-users with an high availability and reliability, suitable for a production-ready scenario.

Since the service is stateless, it can be easily scaled horizontally, by running multiple instances of the service behind a load balancer.
Having multiple instances of the service running also increases the availability of the service, since if one instance fails, the load balancer can redirect the traffic to the other instances.

Kubernetes can also help with rolling updates, by deploying new versions of the service without downtime, and with self-healing, by restarting the service in case of a failure.

### API Gateway

Used in combination with a container orchestrator, an API Gateway can be used to manage the incoming requests to the service.

This offers several advantages, such as:

- Load balancing: the API Gateway can distribute the incoming requests to multiple instances of the service, to handle a higher number of users.
- Rate limiting: the API Gateway can be configured to limit the number of requests per second, to avoid abuse of the service. A 429 Too Many Requests response can be sent to the client when the limit is reached.
- Security features: the API Gateway can be configured to filter out malicious requests, or to enforce HTTPS.
- Logging and monitoring: the API Gateway can log the incoming requests, and provide metrics about the usage of the service.
- API versioning: the API Gateway can be configured to handle different versions of the service, to avoid breaking changes for the clients.
- Caching: the API Gateway can cache the responses of the service, to reduce the load on the backend.

### Caching

In a production environment, caching can be used to reduce the load on the backend service, and to improve the response time for the clients. Caching can be implemented at different levels:

- In-memory caching: the responses of the service can be cached in memory, to avoid recomputing them for each request. Several libraries offer this feature. In memory caching is fast, but it is limited by the amount of memory available on the machine, and the cache is lost when the service is restarted. However, given the stateless nature of the service, this is not a big issue, unless the number of users significantly increases.
- Caching on a in-memory database: the responses can be cached in a in-memory database, such as Redis, to share the cache between multiple instances of the service. This fixes the issues of in-memory caching, also allowing for a greater amount of data to be cached, but it introduces a bit of latency, since Redis has to be reached over the network.
- Caching on a distributed cache: the responses can be cached in a distributed cache, such as Memcached or Hazelcast, to share the cache between multiple instances of the service, and to scale the cache horizontally. This is the most scalable solution, but it introduces more complexity, and it requires more resources.

### Logging and Monitoring

Besides the request logging provided by the API Gateway, the service can be instrumented with a logging and monitoring system, to keep track of the health of the service, and to debug issues.

This project uses the standard Go logging package, which has limited features. In a production environment, a more advanced logging library can be used, such as Logrus, which allows to log to different outputs, and to format the logs in different ways.

Better error handling can be implemented, to log the errors in a structured way, and to include more information in the logs, such as the request ID, the user ID, and the error message returned by the external APIs.

For monitoring, a monitoring system can be used, such as Prometheus, which can scrape metrics from the service, and Grafana, which can visualize the metrics in dashboards. This allows to keep track of the performance of the service, and to detect issues before they become critical.

Alerts can be setup in the monitoring system, to notify maintainers when the service is not performing as expected, or when it is down.

### Authentication and Authorization

If the project requirements include the need for authentication and authorization, the service can be secured with a proper authentication and authorization system.

For authentication, the service can leverage the OAuth2 standard, which allows clients to authenticate with the service using access tokens. The service can be configured to accept only requests with a valid access token, and to reject requests without a token.

For authorization, the service can use role-based access control, to restrict the access to certain endpoints to certain users. The service can be configured to allow only users with a certain role to access certain endpoints, and to reject users without the required role.