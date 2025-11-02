FROM golang:1.25.3

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api

CMD [ "/app/api" ]

