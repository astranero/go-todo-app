FROM golang:latest

WORKDIR /app
COPY /static/* ./static/
COPY go.mod .
COPY go.sum .
COPY main.go .

RUN go mod download
RUN go build -o /todo

EXPOSE 8080
CMD ["/todo"]