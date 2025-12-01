# Go E-Commerce API

A comprehensive e-commerce backend system built with Go, featuring secure authentication, payment processing, and scalable architecture. This project showcases modern API design patterns and real-world integration with databases, caching systems, and external payment services.



## Features

- **User Authentication:** Signup, login, JWT-based authentication, and role-based access control.
- **Product Management:** CRUD operations for products, categories, and product images.
- **Shopping Cart:** Add, update, and remove items from the cart.
- **Order Processing:** Place orders, manage order items, and track order status.
- **Payment Integration:** Connects to payment gateways (e.g., Zarinpal), handles payment verification, and updates order/payment status.
- **Product Reviews:** Users can leave reviews and ratings for products; average rating is auto-updated.
- **Admin Panel:** Admin endpoints for managing products, categories, and users.
- **Rate Limiting & Security:** Middleware for rate limiting, CORS, and authentication.
- **Clean Architecture:** Separation of concerns with domain, usecase, repository, and delivery layers.
- **Dockerized:** Ready for deployment with Docker and docker-compose.

## Project Structure

```
├── src/
│   ├── api/            # HTTP handlers, DTOs, routers
│   ├── cmd/            # Application entrypoint
│   ├── config/         # Configuration files
│   ├── constants/      # Constant values
│   ├── dependency/     # Dependency injection
│   ├── domain/         # Models, business entities
│   ├── events/         # Domain events
│   ├── infra/          # Database, cache, payment gateway
│   ├── pkg/            # Shared packages (errors, etc.)
│   ├── usecase/        # Business logic
│   └── ...
├── docker/             # Docker and compose files
├── README.md           # Project documentation
└── ...
```
