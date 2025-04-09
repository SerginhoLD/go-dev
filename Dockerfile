FROM golang:1.24.2-alpine3.21 AS builder

RUN apk add make \
    && go install github.com/google/wire/cmd/wire@v0.6.0 \
    && go install github.com/pressly/goose/v3/cmd/goose@v3.24.2

# Build and cache the dependencies
WORKDIR /srv
COPY go.mod go.sum ./

RUN go mod download

# Copy the actual code files and build the application
COPY . .

RUN make build app=web \
    && make build app=scheduler

# Build the final image
FROM alpine:3.21

RUN apk add make

WORKDIR /srv
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations migrations
COPY go.mod Makefile .env ./
COPY --from=builder /srv/build ./build

#EXPOSE 8080
CMD ["./build/web"]
