# User Management

## API - REST
Instancia básica de API REST con Gin

### Inicialización del api
```bash
go run cmd/api/main.go
```

### Prueba de rutas
- Rutas Públicas
```bash
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/metrics
curl http://localhost:8080/docs
```

- Rutas Protegidas
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer valid-token" \
  -d '{"name": "Juan", "email": "juan@test.com", "age": 30}'

curl http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer valid-token"
```