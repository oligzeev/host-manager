# docker build -t host-manager -f build/docker1/Dockerfile .
# docker run --rm -p9090:8080 host-manager

FROM golang:latest

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o main /app/cmd/host-manager/

CMD ["/app/main"]