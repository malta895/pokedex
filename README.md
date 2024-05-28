# Pokedex

A simple web API interface to retrieve information about your Pokemons!

## How to run the project

To compile and run this project, you need a working [Go installation](https://go.dev/doc/install).
This project requires at least Go 1.22.

<!-- TODO: guide for at least unix like and Windows OSes -->

Once you have Go installed, to build the project run: 
```bash
go build .
```
An executable file will be created in the project folder, named `pokedex` or `pokedex.exe` depending on the platform.

``

Once spinned up, the executable wil expose the API service at `http://localhost:8080`.

You can now contact the endpoints as explained in the usage section

## Usage

<!-- TODO create me! -->


## Production-ready API

This program has been developed for a simple local demonstration.

To make it production ready, several aspects must be addressed:


- The service should not be exposed to the internet, an API Gateway reverse proxy should be used to redirect the traffic to it. This ensures better security and separation of concerns.
- The project is ready to be shipped as a Docker container, in a production environment it can be deployed on a Kubernetes cluster. By leveraging the features offered by this container orchestrator, with a proper setup, the service can be scaled to serve a very high number of end-users with an high availability, suitable for a production-ready scenario.
- Since we don't expect the external resources responses to change often in time, we can cache the responses in a in-memory database, such as Redis. This offers several advantages: cache hits guarantee a more rapid response, and also they do not count towards the rate limit of the external services APIs. A proper TTL should be setup in order to achieve a reasonable trade-off taking into account factors such as a rough mean frequency of updates of the external resources, the rate limit imposed by the external APIs, the number of daily users of our APIs, the amount of used space on the Redis machine.


Additionally, given the simplicity of the task, the following design decisions were taken:

- The GraphQL queries are manually built; a proper GraphQL library should be used for a better maintainable and solid solution in case more APIs and functionalities need to be implemented
- The Go testing library has been used. Although it can be a perfect choice even for big projects, the usage of external testing libraries such as testify can ease the developer work, especially when working with big structs, where the diffs of failed assertion help in finding errors.