FROM golang:latest
WORKDIR /app

ENV APP_ENV=production

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o gogym .

EXPOSE 8080

CMD ["./gogym"]