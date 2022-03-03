FROM golang:latest

RUN mkdir -p cats-go/
WORKDIR /cats-go/
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /cats-go/
RUN go build -o cats-go-docker ./main.go

EXPOSE 8000

CMD ["./cats-go-docker"]