FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG CGO_ENABLED=0
RUN go build -o /server

EXPOSE 8080

CMD [ "/server" ]
