FROM golang:1.23.6-alpine

WORKDIR /app

#We need GCC compiler
RUN apk add gcc musl-dev

#Copy go modules
COPY go.mod .
COPY go.sum .

# Download go modules
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -o /pastytext

EXPOSE 8080
CMD ["/pastytext"]
