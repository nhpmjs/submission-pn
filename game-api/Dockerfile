FROM golang:1.24 as build
WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 go build -a -o http-service ./cmd/http

FROM gcr.io/distroless/cc-debian11:nonroot
WORKDIR /app
COPY --from=build --chown=nonroot:nonroot /go/src/http-service .
CMD ["./http-service"]
