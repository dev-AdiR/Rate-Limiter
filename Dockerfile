FROM golang:1.25-alpine
RUN go install github.com/go-task/task/v3/cmd/task@latest
WORKDIR /rate_limiter
COPY . .
RUN go mod download
EXPOSE 8081
CMD ["task", "run"]
