
![Logo](https://ik.imagekit.io/pj3r6oe9k/prasorganic-high-resolution-logo-transparent.svg?updatedAt=1726835541390)
# Prasorganic Auth Service

Prasorganic Auth Service is one of the components in the Prasorganic Microservice architecture
built with Go (Golang). This service supports operations for authentication and authorization
users via RESTful API, gRPC, and Message Broker.

## Tech Stack

[![My Skills](https://skillicons.dev/icons?i=go,docker,redis,rabbitmq,bash,git&theme=light)](https://skillicons.dev)

## Features

- **Authentication and Authorization:** Supports operations such as registration, login, refresh token, and OAuth2 authentication.

- **RESTful API:** Provides a RESTful API using Fiber with various middleware for managing requests and responses.

- **gRPC:** Utilizes gRPC for inter-service communication, equipped with interceptors for handling requests and responses.

- **Message Broker:** This service acts as a producer for the RabbitMQ Email Service.

- **Caching:** Redis is used for data caching.

- **Logging:** Logs are recorded using Logrus.

- **Error Handling:** Error handling is implemented to ensure proper detection and handling of errors, minimizing impact on both the client and server.

- **System Resilience:** Uses a Circuit Breaker to improve application resilience and fault tolerance, protecting the system from cascading failures.

- **Configuration and Security:** Employs Viper and HashiCorp Vault for integrated configuration and security management.

- **Testing:** Implements unit testing using Testify.


## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

This project makes use of third-party packages and tools. The licenses for these
dependencies can be found in the `LICENSES` directory.

## Dependencies and Their Licenses

- `Docker:` Licensed under the Apache License 2.0. For more information, see the [Docker License](https://github.com/docker/docs/blob/main/LICENSE).

- `Docker Compose:` Licensed under the Apache License 2.0. For more information, see the [Docker Compose License](https://github.com/docker/compose/blob/main/LICENSE).

- `Go:` Licensed under the BSD 3-Clause "New" or "Revised" License. For more information, see the [Go License](https://github.com/golang/go/blob/master/LICENSE).

- `Redis:` Follows a dual-licensing model with RSALv2 and SSPLv1. For more information, see the [Redis License](https://redis.io/legal/licenses/).

- `RedisInsight:` Licensed under the RedisInsight License. For more information, see the [RedisInsight License](https://github.com/RedisInsight/RedisInsight/blob/main/LICENSE).

- `Htpasswd:` Licensed under the Apache License, Version 2.0. For more information, see the [Htpasswd License](https://www.apache.org/licenses/LICENSE-2.0).

- `GNU Make:` Licensed under the GNU General Public License v3.0. For more information, see the [GNU Make License](https://www.gnu.org/licenses/gpl.html).

- `GNU Bash:` Licensed under the GNU General Public License v3.0. For more information, see the [Bash License](https://www.gnu.org/licenses/gpl-3.0.html).

- `Git:` Licensed under the GNU General Public License version 2.0. For more information, see the [Git License](https://opensource.org/license/GPL-2.0).

## Thanks 👍
Thank you for viewing my project.