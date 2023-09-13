# Dockerfile definition for Backend application service.

# From which image we want to build. This is basically our environment.
# FROM golang:1.19-alpine as Build

# This will copy all the files in our repo to the inside the container at root location.
# COPY . .

# Build our binary at root location.
# RUN GOPATH= go build -o /main cmd/main.go
# RUN go build -o /main cmd/main.go

####################################################################
# This is the actual image that we will be using in production.
# FROM alpine:latest

# We need to copy the binary from the build image to the production image.
#  COPY --from=Build /main .

# This is the port that our application will be listening on.
# EXPOSE 1323

# This is the command that will be executed when the container is started.
# ENTRYPOINT ["./main"]

FROM golang:1.19-alpine  as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy file
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod vendor

# Build the Go app
RUN go build -o /main cmd/main.go
# RUN GOPATH= go build -o /main cmd/main.go


######## Start a new stage from scratch #######
FROM alpine:latest  


WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /main .

# Expose port 1323 to the outside world
EXPOSE 1323

# Command to run the executable
CMD ["./main"]
