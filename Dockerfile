## We specify the base image we need for our
## go application
FROM golang:1.20-alpine AS builder
## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app
## We copy everything in the root directory
## into our /app directory
ADD . /app
## We specify that we now wish to execute
## any further commands inside our /app
## directory
WORKDIR /app
## we run go build to compile the binary
## executable of our Go program
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o poc .

FROM alpine:3.14  

COPY --from=builder /app/poc .
COPY --from=builder /app/service/saleschannel/migration service/saleschannel/migration
COPY --from=builder /app/service/inventory/migration service/inventory/migration

ENTRYPOINT ./poc
