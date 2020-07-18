FROM cyndie7216/go

WORKDIR /go/src/colossus
COPY ./vendor ./vendor
COPY ./configs ./configs
COPY ./internal ./internal
COPY ./cmd/tesseractserver ./cmd/tesseractserver

RUN go build -tags static_all -o tesseractserver ./cmd/tesseractserver/*.go

CMD ["./tesseractserver"]
