FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy


CMD ["go", "test", "-tags=integration", "./integration-tests"]