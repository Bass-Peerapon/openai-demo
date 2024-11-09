# Stage 1: Build
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Specify Linux as the target OS to ensure the executable will work in the alpine container
RUN GOOS=linux GOARCH=amd64 go build -o /app/build/app /app/cmd

# Stage 2: Run
FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /usr/app

# Copy the built binary from the builder
COPY --from=builder /app/build/app ./app
COPY --from=builder /app/migrate ./migrate

# Ensure the binary is executable
RUN chmod +x ./app

EXPOSE 3000

ENTRYPOINT [ "./app"]
