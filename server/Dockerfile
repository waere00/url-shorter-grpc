FROM golang:1.16-alpine
ENV GOPATH=/
COPY ./ ./
RUN go mod download
RUN go build -o server server.go
CMD ["./server"]
EXPOSE 9080
