FROM golang:1.19

# Set the working directory to /app
WORKDIR /app

# Copy the Go module files and download dependencies
# COPY go.mod go.sum ./

# Copy the rest of the application source code
COPY . .
COPY ./srv/main.go ./srv/main.go
RUN go mod download

WORKDIR /app/srv

# Build the Go application
ENTRYPOINT [ "go", "run", "main.go" ]


# RUN mkdir app
# COPY ./go.mod /app/go.mod
# COPY ./go.sum /app/go.sum
# ENV GOPATH=/app
# RUN go mod download
# COPY ./srv/main.go /app/main.go
# CMD go run /app/main.go
