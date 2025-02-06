FROM golang:1.23.6-alpine3.21

RUN apk add --no-cache make \
    && go install github.com/google/wire/cmd/wire@latest \
    && go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#COPY go.mod go.sum ./
#RUN go mod download && go mod verify

COPY . .
RUN wire && go build -o main

CMD ["./main"]
EXPOSE 8080

# Download Go modules
#COPY go.mod go.sum ./
#RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
#COPY *.go ./

# Build
#RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can (optionally) document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
#EXPOSE 8080

# Run
#CMD [ "/docker-gs-ping" ]
