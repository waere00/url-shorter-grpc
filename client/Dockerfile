FROM golang:1.16-alpine
ENV GOPATH=/
COPY ./ ./
RUN go mod download
RUN go build -o client client.go
CMD ["./client"]
EXPOSE 1080
