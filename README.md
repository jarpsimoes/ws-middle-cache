# Web Services Middle Cache Layer - WS Middle Cache Layer

## Overview
**WS Middle Cache Layer** is a Go-based web service that implements a middle caching layer using the Gin framework. It provides health check and liveness endpoints to monitor the status of the application and supports caching using both in-memory storage and Azure Table Storage.

The service is designed to improve performance by reducing redundant calls to backend services and providing a caching mechanism for frequently accessed data.

## Features
- **In-memory caching**: Fast and efficient caching for frequently accessed data.
- **Persistent caching**: Data is stored in Azure Table Storage for long-term access.
- **Dynamic expiration**: Cache entries can have configurable expiration times.
- **Health and liveness checks**: Endpoints to monitor the service's health and availability.
- **Environment-based configuration**: Easily configurable via environment variables.

## Project Structure
```
ws-middle-cache
├── cmd
│   └── server
│       └── main.go          # Entry point of the application
├── internal
│   ├── handlers
│   │   ├── cache.go         # Cache handler
│   │   ├── health.go        # Health check endpoint
│   │   └── liveness.go      # Liveness check endpoint
│   ├── middleware
│   │   └── logging.go       # Logging middleware
│   ├── routes
│   │   └── router.go        # Route setup
│   └── services
│       ├── aztable.go       # Azure Table Storage service
│       └── cache.go         # In-memory cache service
├── pkg
│   └── utils
│       ├── environment.go   # Environment variable utility
│       └── response.go      # Utility functions for responses
├── go.mod                    # Go module definition
├── go.sum                    # Module dependency checksums
└── README.md                 # Project documentation
```

## Setup Instructions
1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd ws-middle-cache
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set environment variables:**
   Create a `.env` file or export the required variables:
   ```bash
   export PORT=8080
   export BACKEND_ENDPOINT=https://example.com
   export CACHE_EXPIRATION_SECONDS=600
   export AZURE_STORAGE_ACCOUNT_NAME=your_account_name
   export AZURE_STORAGE_ACCOUNT_KEY=your_account_key
   export AZURE_STORAGE_ACCOUNT_TABLE_NAME=cachetable
   export LOG_LEVEL=INFO
   ```

4. **Run the application:**
   ```bash
   go run cmd/server/main.go
   ```

5. **Access the endpoints:**
   - Cache API: `GET /api/*any`
   - Health Check: `GET /health`
   - Liveness Check: `GET /liveness`

## Configuration

The service can be configured using the following environment variables:

| Variable Name                   | Default Value       | Description                                                                 |
|---------------------------------|---------------------|-----------------------------------------------------------------------------|
| `PORT`                          | `8080`              | The port on which the service runs.                                         |
| `BACKEND_ENDPOINT`              | `""` (empty)        | The base URL of the backend service to fetch data from.                     |
| `CACHE_EXPIRATION_SECONDS`      | `600` (10 minutes)  | The default expiration time for cache entries (in seconds).                |
| `AZURE_STORAGE_ACCOUNT_NAME`    | `""` (empty)        | The Azure Storage account name.                                             |
| `AZURE_STORAGE_ACCOUNT_KEY`     | `""` (empty)        | The Azure Storage account key.                                              |
| `AZURE_STORAGE_ACCOUNT_TABLE_NAME` | `cachetable`     | The Azure Table Storage table name.                                         |
| `LOG_LEVEL`                     | `INFO`              | The logging level (`DEBUG`, `INFO`, `WARN`, `ERROR`).                       |

## Usage
The service provides the following endpoints:

### Cache API
- **GET `/api/*any`**: Fetches data from the cache or backend service.
- **DELETE `/api/cache`**: Clears the in-memory and Azure Table Storage caches.

### Health and Liveness
- **GET `/health`**: Returns `200 OK` if the service is healthy.
- **GET `/liveness`**: Returns `200 OK` if the service is alive.

## Logging

The service uses a custom logging utility. The log level can be configured using the `LOG_LEVEL` environment variable. Supported levels are:
- `DEBUG`
- `INFO`
- `WARN`
- `ERROR`

## Contributing
Contributions are welcome! Please feel free to submit a pull request or open an issue for any enhancements or bug fixes.

## License
This project is licensed under the MIT License. See the `LICENSE` file for more details.