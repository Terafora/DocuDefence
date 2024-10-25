FROM golang:1.20


WORKDIR /app

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/ ./

EXPOSE 8000

RUN go build -o main .

CMD ["./main"]
