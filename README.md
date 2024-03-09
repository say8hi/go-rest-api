# Go API Project

This repository contains a Go-based API project that demonstrates a structured approach to building a Go application using various technologies such as Docker, PostgreSQL, RabbitMQ, and testing strategies. The project is structured to support scalability and maintainability.

- [Go API Project](#go-api-project)
  - [Project Structure](#project-structure)
  - [Technologies Used](#technologies-used)
  - [Authentication Method](#authentication-method)
    - [Hashing Utility](#hashing-utility)
      - [Usage](#usage)
  - [Setup and Running](#setup-and-running)
  - [API Endpoints](#api-endpoints)
      - [Unauthorized Endpoints](#unauthorized-endpoints)
      - [Authorized Endpoints](#authorized-endpoints) 
  - [Application Usage Examples](#application-usage-examples)
      - [Creating a User](#creating-a-user) 
      - [Creating a Category](#creating-a-category)
      - [Creating a Product](#creating-a-product)
  - [Testing](#testing)

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

## Authentication Method

The application utilizes an authentication mechanism that employs a combination of the username and password, concatenated in the format "passwordusername", and then hashed using the SHA-256 algorithm. This approach ensures enhanced security by storing only hashed versions of the credentials, thereby protecting sensitive user information.

### Hashing Utility

For convenience, the application includes a hashing utility script located in the `utils/` directory, named `hashcli.go`. This script allows for easy generation of SHA-256 hashes of arbitrary strings, adhering to the application's authentication format.

#### Usage

To hash a string, such as a concatenated password and username, use the following command with the compiled binary:

```bash
bin/hashcli -s "passwordusername"
```

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

## API Endpoints

The application exposes several RESTful endpoints divided into unauthorized (public) and authorized (secured) categories. Authorization is managed through a bearer token provided in the `Authorization` header.

### Unauthorized Endpoints

- **Users**
  - `POST /users/create`: Create a new user.

- **Categories**
  - `GET /category/{id}`: Get a category by ID.
  - `GET /category/`: Get all categories.

- **Products**
  - `GET /product/{id}`: Get a product by ID.
  - `GET /category/{id}/products`: Get all products in a category.

### Authorized Endpoints

- **Categories**
  - `POST /category/create`: Create a new category.
  - `PATCH /category/{id}`: Update a category.
  - `DELETE /category/{id}`: Delete a category.

- **Products**
  - `POST /product/create`: Create a new product.
  - `PATCH /product/{id}`: Update a product.
  - `DELETE /product/{id}`: Delete a product.

Use the provided `curl` examples in the [Application Usage Examples](#application-usage-examples) section to interact with these endpoints.

## Application Usage Examples

Below are examples of how to interact with the application using `curl`, a command-line tool for transferring data with URLs. These examples demonstrate user creation, category creation, and product creation through the application's API.

### Creating a User

To create a new user, send a POST request with the username and password in JSON format:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"username": "username", "password": "pass"}' http://0.0.0.0:8080/users/create
```
### Creating a Category

To create a new category, you must first authenticate using a valid token. Then, send a POST request with the category name and description:
```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_TOKEN_HERE" -d '{"name": "new_category", "description": "desc"}' http://0.0.0.0:8080/category/create
```
Replace YOUR_TOKEN_HERE with your actual authentication token.

### Creating a Product

Similarly, to create a new product, use an authenticated POST request. Provide the product name, description, price, and associated categories in JSON format:
```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer YOUR_TOKEN_HERE" -d '{"name": "new_product", "description": "desc", "price": 2.5, "categories": ["new_category"]}' http://0.0.0.0:8080/product/create
```
Again, replace YOUR_TOKEN_HERE with your actual authentication token. This example assumes you have already created a category named "new_category" to which the product is being associated.

## Testing

To run the integration tests:
```bash
./run_tests.sh
```
This script will set up the test environment, run the tests, and tear down the environment afterwards.
