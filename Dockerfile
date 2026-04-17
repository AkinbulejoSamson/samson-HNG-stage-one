FROM golang:1.25-alpine

WORKDIR /app

# Install git in case your go.mod has private repos or git-based dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Explicitly disable CGO for the Pure Go driver
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

EXPOSE 8080

CMD ["./main"]