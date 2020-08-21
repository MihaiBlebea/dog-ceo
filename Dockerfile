FROM golang:1.14-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/dog_ceo

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o ./out/dog_ceo .

# Start fresh from a smaller image
FROM alpine:3.9

RUN apk add ca-certificates

WORKDIR /app

COPY --from=build_base /tmp/dog_ceo/out/dog_ceo /app/dog_ceo

# Copy the static folder and get all the static files inside it
COPY ./static /app/static

EXPOSE ${HTTP_PORT}

# Run the binary program produced by `go install`
CMD ["/app/dog_ceo"]