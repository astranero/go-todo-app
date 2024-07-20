FROM golang:latest

WORKDIR /app
COPY * .

RUN go mod download
RUN go build -o /todo

EXPOSE 8080
CMD (/todo)