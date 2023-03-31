FROM golang:latest
WORKDIR /app
COPY . .
RUN go get -v 
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]