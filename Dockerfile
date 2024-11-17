# Stage 1: Build stage
FROM golang:1.22.1 AS builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# Stage 2: Final stage
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .
COPY --from=builder /app/helper/format.html /app/helper/format.html


# Optionally copy the .env file if it's needed
COPY --from=builder /app/.env .


#COPY --from=builder /app/api/model.conf ./api/model.conf
#COPY --from=builder /app/api/policy.csv ./api/policy.csv

# Expose port 8080
EXPOSE 8080

CMD ["./myapp"]
