FROM golang:latest

ARG DATABASE_URL
ENV DATABASE_URL=$DATABASE_URL

ARG SECRET_KEY
ENV SECRET_KEY=$SECRET_KEY

WORKDIR /app

ENV APP_ENV=production

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod vendor

RUN go build -o gogym .

EXPOSE 8080

CMD ["./gogym"]