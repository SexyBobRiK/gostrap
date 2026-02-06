![Logo](./assets/gostrap.png)
# Gostrap 🚀

**Gostrap** is a lightweight, configuration-driven Go framework designed to eliminate boilerplate code when initializing microservices. It handles the heavy lifting of setting up HTTP servers, database connections, and cache providers through a single JSON configuration.

## ✨ Key Features

* 🍸 **Gin Provider**: Ready-to-go HTTP server with CORS and static files support.
* 🗄️ **GORM Provider**: Automatic PostgreSQL/MySQL connection management.
* ⚡ **Redis Provider**: Native support for Redis (RESP2/RESP3) out of the box.
* 🛠️ **Clean Architecture**: Decouples infrastructure setup from your business logic.
* 🚦 **Graceful Control**: Built-in lifecycle management with `Pulse()` and `Wait()`.

## 📦 Installation

```bash
go get [github.com/SexyBobRiK/gostrap](https://github.com/SexyBobRiK/gostrap)

🛠 Usage
1. Configuration (gostrap.json)
Define your infrastructure in a simple JSON file:

```JSON
{
  "gin": {
    "port": "8080",
    "enable": true,
    "mode": "release",
    "middleware": {
      "cors": {
        "enabled": true,
        "allow_origins": ["*"]
      }
    }
  },
  "gorm": {
    "enable": true,
    "host": "localhost",
    "port": "5432",
    "user": "postgres",
    "password": "password",
    "dbname": "my_service",
    "sslmode": "disable"
  },
  "redis": {
    "enable": true,
    "host": "localhost",
    "port": "6379",
    "db": 0
  }
}
2. Implementation
Bootstrap your entire application in just a few lines of code:

```Go
package main

import (
    "[github.com/SexyBobRiK/gostrap](https://github.com/SexyBobRiK/gostrap)"
    "your_project/internal/routes"
    "your_project/internal/repository"
)

func main() {
    // 🚀 Initialize all providers via config
    boot, err := gostrap.LetsGo("./config/gostrap.json")
    if err != nil {
        panic(err)
    }

    // 🏗️ Inject dependencies into your layers
    // Example: taking the first available DB connection
    var dbConn *gorm.DB
    for _, db := range boot.Database {
        dbConn = db
        break
    }

    routes.RouterInit(boot.Gin, repository.NewRepoContainer(dbConn))

    // ❤️ Start the application heartbeat
    if err := boot.Pulse(); err != nil {
        panic(err)
    }

    // 🛑 Block until interrupt signal (SIGTERM/SIGINT)
    boot.Wait()
}
🏗 Lifecycle Management
LetsGo(path): Parses config and establishes all active connections.

Pulse(): Fires up the HTTP server.

Wait(): Keeps the application alive and listens for termination signals for a clean exit.

🛡 License
This project is licensed under the MIT License.

Developed with ❤️ by SexyBobRiK