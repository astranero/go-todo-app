FROM golang:latest

WORKDIR /usr/src/app/writer

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /logoutput-writer

EXPOSE 8080
CMD ["/logoutput-writer"]