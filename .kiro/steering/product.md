# Product Overview

**dfood** is a food-related web application built with Go. The application provides user authentication and event management functionality through a REST API.

## Core Features
- User registration and authentication with JWT tokens
- Password management and updates
- Event creation and management
- Multi-environment support (dev, staging, production)

## Architecture
The application follows a clean architecture pattern with clear separation of concerns:
- RESTful API built with Gin framework
- SQLite database for data persistence
- JWT-based authentication system
- Repository pattern for data access
- Service layer for business logic