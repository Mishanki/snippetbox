# Latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Mishanki <Mishanki@yandex.ru>"

# Workdir inside the container
WORKDIR /app

# Copy all the app sources
COPY . .

# Download dependency
RUN go mod download

# Install the package
RUN go install -v ./...

# Set http port
ENV PORT 4000

# Build the app
RUN go build -work /app/cmd/*

# Remove source file
RUN rm -rf /app/cmd
RUN rm -rf /app/handlers
RUN rm -rf /app/internal
RUN rm -rf /app/tmp
RUN rm -rf /app/.DS_Store
RUN rm -rf /app/.fleet
RUN rm -rf /app/_docker
RUN rm -rf /app/docker-compose.yml
RUN rm -rf /app/go.mod
RUN rm -rf /app/go.sum

# Make port 4000 available to the world outside this container
EXPOSE $PORT

# Run the app
CMD ["./web"]