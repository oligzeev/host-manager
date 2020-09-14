# Tested on docker 19.03
# docker build -t host-manager -f build/docker/Dockerfile .
# docker run --rm -p9090:8080 host-manager

FROM golang:latest AS builder
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main /app/cmd/host-manager/

FROM alpine:latest
RUN mkdir -p /app/config
COPY --from=builder /app/config/host-manager.yaml /app/config
COPY --from=builder /app/config/kube-config.yaml /app/config
COPY --from=builder /app/main /app
WORKDIR /app
CMD ["/app/main"]