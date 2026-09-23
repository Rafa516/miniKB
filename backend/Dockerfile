# --- Build stage ---
FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# --- Run stage ---
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/server ./server

# A porta real é lida de $PORT em tempo de execução (padrão 8080 se não
# definida) — ver backend/cmd/server/main.go. Sem volume: os dados agora
# ficam no PostgreSQL (DATABASE_URL), não em disco local.
EXPOSE 8080

ENTRYPOINT ["/app/server"]
