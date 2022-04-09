FROM golang:1.17-alpine
RUN mkdir /app
ADD . /app
WORKDIR /app
EXPOSE 27604:27604
RUN go build -o router .
CMD ["./router"]
