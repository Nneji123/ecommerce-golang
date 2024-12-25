# LeadzAura API

This project is a comprehensive API service that includes various functionalities such as email management, lead generation, and task scheduling. It is designed to be scalable, maintainable, and easy to deploy using Docker and Docker Compose.

## **TODO**

### Documentation
- [ ] Create extensive project documentation for new developers.

### Scraping
- [ ] Implement Google Maps Scraper
- [ ] Implement Twitter Scraper
- [x] Test ProxyCurl Linkedin Scraper and write scraped data to database

### Server Settings
- [X] Add CORS Handling (read cors list from .env or env variables)
- [x] Add air for hot reload
- [x] Add Graceful shutdown
- [X] Add Rate Limiter Middleware for specific routes
- [x] Add Caching
- [x] Add API Documentation with swag
- [x] Change banner
- [x] Add server health check
- [x] Create custom structs for errors in swag (instead of using map[string]string
- [ ] create bash script that uses air when $ENV is set to development
- [ ] restructure API modules 

### Email Errors
- [x] Modify Routes handlers to send requests to task queue (send mail, scrape emails etc)
- [x] Fix Attachment upload for email sending
- [ ] Create routes for email sequence tasks (for testing purposes
- [ ] Refactor code and make it cleaner.




## Table of Contents

- [LeadzAura API](#leadzaura-api)
  - [**TODO**](#todo)
    - [Documentation](#documentation)
    - [Scraping](#scraping)
    - [Server Settings](#server-settings)
    - [Email Errors](#email-errors)
  - [Table of Contents](#table-of-contents)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
  - [Usage](#usage)
    - [API Endpoints](#api-endpoints)
  - [API Documentation](#api-documentation)
  - [Contributing](#contributing)
  - [License](#license)
  - [Docker Compose](#docker-compose)
  - [Environment Variables](#environment-variables)
  - [Make Commands](#make-commands)


## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Docker
- Docker Compose
- Go (for building from source)

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/leadzaura/go-scraper.git
   ```
2. Navigate to the project directory:
   ```
   cd go-scraper
   ```

## Usage

Once the Docker container is running, you can access the API service at `http://localhost:8080`.

### API Endpoints

- **Email Management**: `/api/email`
- **Lead Generation**: `/api/leads`
- **Task Scheduling**: `/api/tasks`

## API Documentation

API documentation is generated using Swagger and can be accessed at `http://localhost:8080/swagger/index.html`.

## Contributing

Contributions are welcome. Please read the [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Docker Compose

To run the project using Docker Compose, follow these steps:

1. Ensure Docker and Docker Compose are installed on your machine.
2. Navigate to the project directory.
3. Run the following command to start the services defined in `docker-compose.yml`:
   ```
   docker-compose up -d
   ```
4. To stop the services, run:
   ```
   docker-compose down
   ```

## Environment Variables

Environment variables can be set in the `.env` file located in the project root. These variables are used by the application and can be customized as needed.

## Make Commands

This project includes a Makefile with several commands for common tasks. Here are some of the most useful ones:

- **Build the project**:
 ```
 make build
 ```
- **Run the project**:
 ```
 make run
 ```
- **Clean the build artifacts**:
 ```
 make clean
 ```
- **Run tests**:
 ```
 make test
 ```
- **Generate API documentation**:
 ```
 make generate-docs
 ```

For a full list of available commands, run:
```
make help
```
