FROM golang:1.14.6-buster as build

WORKDIR /go/src/colossus
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o storage cmd/storage/main.go

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/src/colossus/storage /
CMD ["/storage"]
