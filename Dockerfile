FROM golang:1.16 AS builder

LABEL maintainer="yashvardhan-kukreja"

RUN mkdir -p /export/kube-bench

RUN mkdir -p /kube-bench-exporter
WORKDIR /kube-bench-exporter

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go get -d ./...

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/kube-bench-exporter -v .

# Packaging stage
FROM alpine

LABEL maintainer="yashvardhan-kukreja"

COPY --from=builder /bin/kube-bench-exporter /bin/kube-bench-exporter
COPY --from=builder /etc/passwd /etc/passwd

USER 10001

ENTRYPOINT ["/bin/kube-bench-exporter"]
