FROM golang:1.17-alpine
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o server .
CMD ["./server"]
