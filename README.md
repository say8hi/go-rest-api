# Go API Project

This repository contains a Go-based API project that demonstrates a structured approach to building a Go application using various technologies such as Docker, PostgreSQL, RabbitMQ, and testing strategies. The project is structured to support scalability and maintainability.

## Project Structure

The project is organized as follows:

- `cmd/go-api-test`: Contains the entry point of the application.
- `internal`: Houses the core logic of the application, including database interactions, handlers, middlewares, and models.
- `services/datacollector`: A separate service for collecting data and interacting with RabbitMQ.
- `tests`: Contains integration tests for the application.
- `Dockerfile` and `docker-compose.yml`: For containerization and orchestration.
- `go.mod` and `go.sum`: Go module files for managing dependencies.

## Technologies Used

- **Go**: The primary programming language for building the application.
- **Docker and Docker Compose**: Used for containerizing the application and its services.
- **PostgreSQL**: The database system for storing application data.
- **RabbitMQ**: For asynchronous message queuing.
- **Gorilla Mux**: A powerful URL router and dispatcher for matching incoming requests to their respective handler.

## Setup and Running

To run this project, you need to have Docker and Docker Compose installed on your system. Follow these steps:

1. Clone the repository to your local machine.
```bash
git clone https://github.com/say8hi/go-api-test.git
```
2. Go into **go-api-test** folder:
```bash
cd go-api-test
```
3. Rename `.env.example` to `.env` and adjust the configuration according to your environment.
```bash
mv .env.example .env
```
4. Build and start the services with Docker Compose:
```bash
docker-compose up -d --build
```
This command will start all the required services, including the Go application, PostgreSQL, and RabbitMQ.

## Testing

To run the integration tests:
```bash
./run_tests.sh
```
This script will set up the test environment, run the tests, and tear down the environment afterwards.
