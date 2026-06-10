<div align="center">

# Sundance

##### A centralized platform for managing and distributing forms.

<img src="./imgs/logos/Forms Hub Logo.png" style="width:175px;"/>

</div>

## Getting Started

### Install Tools 🛠️

#### Go

This project is built using Go. You can download and install Go from the official website: [https://golang.org/dl/](https://golang.org/dl/).

#### Docker

Docker is a platform that allows you to develop, ship, and run applications in containers. You can download Docker from the official website: [https://www.docker.com/get-started](https://www.docker.com/get-started). This project uses Docker to for local development by running the application and its dependencies in containers.

#### IDE

You can use any code editor or IDE that supports Go development. Some popular options include:

- [Visual Studio Code](https://code.visualstudio.com/) with the Go extension
- [GoLand](https://www.jetbrains.com/go/)
- [Neovim](https://neovim.io/) with Go plugins

#### Insomnia

Insomnia is a popular API client that allows you to test and interact with APIs. You can download Insomnia from the official website: [https://insomnia.rest/download](https://insomnia.rest/download).

### Pre-Requisites ⚙️

#### Cosmtrek Air (Optional)

[Cosmtrek Air](https://github.com/air-verse/air) is a live-reloading tool for Go applications. It allows you to automatically restart your application whenever you make changes to the code, which can speed up your development workflow.

1. Install Cosmtrek Air by following the instructions on their [GitHub repository](https://github.com/air-verse/air).

2. Once installed, you can run your Go application with live-reloading by using the following command in your terminal:

   ```bash
   $ air
   ```

> [!NOTE]
> You may need to confirm your Go bin directory is in your system's PATH environment variable to use the `air` command globally.

### Download Dependencies 📦

To download the necessary dependencies for this project, you can use the following command in your terminal:

```bash
# Download dependencies for all projects
$ go mod download all

# Download dependencies for a specific project
$ go -C [project-directory] mod download

# Example:
$ go -C services/tenants mod download
```

### Running the Services 💻

To run the services in this project, open the `services/[service-name]` directory in your IDE and create a new `settings.json` or `settings.yaml` file based on the `settings-example.json` or `settings-example.yaml` file provided in the same directory. This file should contain the necessary configuration settings for the service, such as database connection strings, port numbers, and other environment variables.

Once you have created the settings file, you can use Docker Compose to run the services. Docker Compose allows you to define and manage multi-container Docker applications. You can use the following command in your terminal to start the services:

```bash
$ docker-compose up
```

Alternatively, you can run the services individually using the terminal.

```bash
# Without Cosmtrek Air
$ go -C services/[service-name] run ./cmd/server

# Example:
$ go -C services/tenants run ./cmd/server

# With Cosmtrek Air
$ air -c services/[service-name]/.air.toml
```

> [!NOTE]
> If you are running the services with Cosmtrek Air, you will need to ensure that the `.air.toml` configuration file is properly set-up with the build arguments required to locate the settings file. By default, the `.air.toml` file is configured to look for a `settings.json` file in the same directory as the service's `cmd` directory. If you choose to name your settings file differently or place it in a different location, you will need to update the `.air.toml` file accordingly

Verify that the service is running by checking the terminal output for any errors and ensuring that the service is listening on the expected port. You can also use Insomnia to send test requests to the service's API endpoints to confirm that it is functioning correctly or visit `http://localhost:[port]/swagger/index.html` to view Swagger documentation for the service's API.

### Building the Services 🏗️

To build the services in this project, run the following command from the root of the repository:

```bash
$ go -C [project-directory] build -o [output-binary-name] ./cmd/[command-name]

# Example:
$ go -C services/tenants build -o tenants ./cmd/server
```

By default, Go builds the application for the operating system and architecture of the machine you are building on. If you want to build for a different target platform, you can set the `GOOS` and `GOARCH` environment variables before running the build command. For example, to build for Linux on an AMD64 architecture, you can use the following command:

```bash
$ GOOS=linux GOARCH=amd64 go -C [project-directory] build -o [output-binary-name] ./cmd/[command-name]
```

### Testing the Projects 🧪

To run the unit tests for the projects in this repository, use the following command:

```bash
$ go -C [project-directory] test -v ./...

# Example:
$ go -C services/tenants test -v ./...
```

This command will run all the unit tests in the specified project directory and its subdirectories. To generate code coverage report, you can use the following command:

```bash
$ go -C [project-directory] test -coverprofile=coverage.out ./...
$ go tool cover -html=coverage.out -o coverage.html
```

You can then open the `coverage.html` file in your web browser to view the code coverage report. Alternatively, Visual Studio Code has built-in support for runnig Go tests and viewing code coverage directly within the editor when using the Go extension.

## Project Structure

The repository is organized into several directories, each serving a specific purpose:

### `services/`

Contains the individual microservices that make up the Sundance platform, organized by domain (e.g. `tenants`, `forms`). Each service directory contains the source code for the service, including the main application code, configuration files, and tests. Each service follows a similar structure with a `cmd` directory for the main application entry point(s) and an `internal` directory for the service's internal logic and implementation.

Example structure for Tenants Service:

```
tenants/
├── cmd/
│   └── server/
│       └── main.go             # The main entry point for the application, responsible for initializing the service and starting the server.
├── internal/
│   ├── adapters/
│       ├── persistence/        # Contains the implementation(s) for data storage and retrieval, such as MongoDB repositories.
│       ├── rest/               # REST API routes and handlers for processing incoming HTTP requests.
│       └── workers/            # Background workers for processing asynchronous tasks, such as message queue consumers or scheduled jobs.
│   ├── core/
│       ├── domain/             # Contains the core business logic and domain models for the service.
│       ├── ports/              # The interfaces that define the expected behavior of the service's dependeencies, (e.g. repositories) and the public API of the service.
│       ├── services/           # The implementation of the sevice's business logic, coordinating between teh domain models and the adapters.
│       └── strategies/         # Contains any strategy pattern implementations for handling different algorithms or behaviors within the service.
```

### `pkg/`

Contains shared packages and utilities that are used across multiple services in the project, such as database connection logic, common data models, and helper functions.

### `docs/`

Contains documentation related to Sundance, including API documentation, architectural diagrams, and other relevant information.

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Go by Example](https://gobyexample.com/)
- [Go: The Complete Guide](https://www.udemy.com/course/go-the-complete-guide/)
- [Docker Documentation](https://docs.docker.com/)
- [Docker & Kubernetes: The Complete Guide](https://www.udemy.com/course/docker-and-kubernetes-the-complete-guide/)

