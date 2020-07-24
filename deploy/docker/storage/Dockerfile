FROM golang:1.14-buster as builder

WORKDIR /go/src/colossus
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o storage cmd/storage/main.go

FROM gcr.io/distroless/base-debian10
COPY --from=builder /go/src/colossus/storage /
CMD ["/storage"]
