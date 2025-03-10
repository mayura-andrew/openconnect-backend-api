# Build stage
FROM golang:1.24-alpine AS builder

# Install migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install required build dependencies
RUN apk --no-cache add gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata postgresql-client

# Copy migrate binary from builder
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Copy compiled application and migrations
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/internal/mailer/templates ./internal/mailer/templates
COPY --from=builder /app/uploads ./uploads

# Environment variables
ENV DB_DSN=${DB_DSN}
ENV SMTPPORT=${SMTPPORT}
ENV SMTPSENDER=${SMTPSENDER}
ENV SMTPHOST=${SMTPHOST}
ENV SMTPUSERNAME=${SMTPUSERNAME}
ENV SMTPPASS=${SMTPPASS}
ENV GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID}
ENV GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET}
ENV GOOGLE_REDIRECT_URI=${GOOGLE_REDIRECT_URI}
ENV FRONTEND_URL=${FRONTEND_URL}
# Add entrypoint script
COPY scripts/entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
