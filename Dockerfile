# Start from golang base image
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /go/src/consumer-payment-service

# Copy go mod and go sum files to determine container is to download
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Copy the source codes from the current directory to the working directory inside the container
COPY . .

# Build the Go app for linux for
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

EXPOSE 8800

# Run go build to compile the binary executable of golang program
CMD [ "./main" ]
