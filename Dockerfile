# Builder image, where we build the example.
FROM golang:1.18 AS builder
WORKDIR /go/src/maersk.com/prometheus/client_golang
COPY . .
#
# Download all dependencies
#
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
RUN go mod init
RUN go mod download && go mod verify
RUN go get

#
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w'

# Final image.
FROM scratch
LABEL maintainer "Prometheus Example <martin.bower@maersk.com>"
COPY --from=builder /go/src/maersk.com/prometheus/client_golang/client_golang .
EXPOSE 8080
ENTRYPOINT ["/client_golang"]
