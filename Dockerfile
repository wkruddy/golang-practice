FROM golang:alpine

WORKDIR /app

COPY . .

CMD ["go", "run", "poker.go"]
