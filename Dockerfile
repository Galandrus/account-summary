# Build stage
FROM golang:1.25-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Crear usuario no-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Establecer directorio de trabajo
WORKDIR /app

# Copiar el binario desde el stage de build
COPY --from=builder /app/main .

# Copiar archivos estáticos
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/frontend.html ./
COPY --from=builder /app/src/templates ./src/templates

# Cambiar propiedad de los archivos
RUN chown -R appuser:appgroup /app

# Cambiar al usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
