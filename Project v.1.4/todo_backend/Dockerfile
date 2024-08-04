FROM golang:latest

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /todo

EXPOSE 8081
CMD ["/todo"]