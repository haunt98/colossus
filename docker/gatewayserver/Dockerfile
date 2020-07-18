FROM cyndie7216/go

WORKDIR /go/src/colossus
COPY ./vendor ./vendor
COPY ./configs ./configs
COPY ./internal ./internal
COPY ./cmd/gatewayserver ./cmd/gatewayserver

RUN go build -tags static_all -o gatewayserver ./cmd/gatewayserver/*.go

CMD ["./gatewayserver"]
