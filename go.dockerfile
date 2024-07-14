FROM golang:1.16.3-alpine3.13

WORKDIR /api

COPY . .

# Download and install the dependencies:
RUN go get -d -v ./...

# Build the go app
RUN go build -o api cmd/goapp/main.go 


EXPOSE 8000

CMD ["./api"]