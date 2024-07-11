# stage 1
FROM golang:1.22.5-alpine3.20 AS builder

LABEL author="Dwi Prasetiyo"
LABEL project="prasorganic-auth-service"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./main ./src/main.go

# stage 2
FROM alpine:3.20  

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3300

CMD [ "./main" ]