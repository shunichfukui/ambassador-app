FROM golang:1.16

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

# 全てコピー
COPY . .

CMD ["go", "run", "main.go"]