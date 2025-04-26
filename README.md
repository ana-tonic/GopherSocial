# GopherSocial

GopherSocial is a modern social network backend written in Go. The project is designed with a focus on scalability, performance, and good architecture.

## 🚀 Features

- RESTful API with Swagger documentation
- JWT authentication
- Rate limiting for DoS protection
- Redis caching for improved performance
- PostgreSQL database
- Email notifications through SendGrid
- Logging with Zap logger
- Docker support for easier deployment

## 🛠️ Technologies

- Go 1.23.3
- Chi router for HTTP routing
- JWT for authentication
- PostgreSQL for database
- Redis for caching
- SendGrid for email notifications
- Swagger for API documentation
- Zap for logging

## 📋 Prerequisites

- Go 1.23.3 or newer
- PostgreSQL
- Redis (optional)
- Docker and Docker Compose (optional)

## 🚀 Getting Started

1. Clone the repository:
```bash
git clone https://github.com/ana-tonic/GopherSocial.git
cd GopherSocial
```

2. Set up environment variables:
```bash
cp .envrc.example .envrc
# Edit .envrc file with your configurations
```

3. Run the application:
```bash
# With air for development
air

# Or directly
go run cmd/api/main.go
```

4. For Docker:
```bash
docker-compose up
```

## 📚 API Documentation

Swagger documentation is available at:
```
http://localhost:8080/v1/swagger/doc.json
```

## 🧪 Testing

```bash
go test ./...
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📫 Contact

Ana Tonic - [@ana-tonic](https://github.com/ana-tonic) 