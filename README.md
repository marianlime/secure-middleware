# Secure Middleware

Secure Middleware is a robust, high-performance security and middleware framework developed in Go. Designed for modern web applications, it seamlessly integrates advanced security, caching, and proxy functionalities—all built using a test-driven development approach.

### Overview
Secure Middleware empowers developers with a suite of tools to enhance the security and efficiency of their applications. The project is organized into several key packages:

    Authentication:
    Implements custom authentication strategies including OAuth and JWT, ensuring secure and flexible user validation.

    Encryption:
    Provides secure data handling by leveraging the AES-256-GCM algorithm, offering state-of-the-art encryption for data integrity and confidentiality.

    Logging Monitor:
    Offers a comprehensive logging system with configurable levels (debug, info, error) to facilitate real-time monitoring and troubleshooting.

    Custom Cache & Proxy:
    Features a Redis-inspired caching system combined with a proxy component that clones requests and expertly manages hop-by-hop headers.

### Features

    Go-Powered Performance:
    Harnesses the speed and efficiency of Go to deliver a lightweight yet powerful middleware solution.

    Advanced Authentication:
    OAuth integration for secure third-party authentication.
    JWT-based custom authentication for token management and validation.

    Robust Encryption:
    Utilizes the AES-256-GCM algorithm to protect sensitive data with industry-standard encryption practices.

    Dynamic Logging:
    Supports multiple log levels (debug, info, error) to assist in rapid diagnosis and enhanced observability.

    Efficient Caching & Proxying:
    Includes a Redis-like cache for high-speed data retrieval.
    Features a proxy that clones and sanitizes HTTP requests by managing hop-by-hop headers.

    Test-Driven Development:
        Developed with TDD, ensuring a high standard of code quality, maintainability, and reliability.

### Installation

Ensure you have Go installed (version 1.XX or later recommended). Then, clone the repository and build the project:

git clone https://github.com/marianlime/secure-middleware.git
cd secure-middleware
go build ./...

