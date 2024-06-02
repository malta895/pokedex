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
  - [Design considerations](#design-considerations)
    - [API Gateway](#api-gateway)
    - [Containerization and Containers Orchestration](#containerization-and-containers-orchestration)
    - [Caching](#caching)
    - [Help from External Libraries](#help-from-external-libraries)


## How to run the project

You can either build a Docker image, or build the source code and run the project locally.

Both require you download a local copy of the source code.

### Download the Source Code

This project is hosted on a Git repository.

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

## Design considerations 

This program has been developed for a simple local demonstration, it is not production-ready.

To make it so, several aspects must be addressed, detailed in the next paragraphs.

### API Gateway

The service should not be directly exposed to the internet, an API Gateway reverse proxy should be used to redirect the traffic to it. This ensures better security, among other things

Leveraging the features of modern API Gateways like Nginx or Envoy can help with dealing with problems typical of services with high traffic rates.

### Containerization and Containers Orchestration

The project is ready to be shipped as a Docker container, in a production environment it can be deployed with a Container Orchestrator software; for example, it can be  on a Kubernetes cluster. 
By leveraging the features offered by this container orchestrator, with a proper setup, the service can be scaled to serve a very high number of end-users with an high availability, suitable for a production-ready scenario. Features like Horizontal Pod Autoscaling can be used to help in handling peaks of requests.

### Caching

Since we don't expect the external resources responses to change often in time, we can cache the responses in a in-memory database, such as Redis. This offers several advantages: cache hits guarantee a more rapid response, and also they do not count towards the rate limit of the external services APIs. 
A proper TTL should be setup in order to achieve a reasonable trade-off taking into account factors such as a rough mean frequency of updates of the external resources, the rate limit imposed by the external APIs, the number of daily users of our APIs, the amount of used space on the Redis machine.

### Help from External Libraries

The Go testing library has been used. Although it can be a perfect choice even for big projects, the usage of external testing libraries such as testify can ease the developer work, especially when working with big structs, where the diffs of failed assertion help in finding errors.


