FROM golang:latest

RUN apt-get update && apt-get install -y \
    git \
    make \

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main .

CMD ["go", "test", "./...", "-v"]
