FROM cyndie7216/go

WORKDIR /go/src/colossus
COPY ./vendor ./vendor
COPY ./configs ./configs
COPY ./internal ./internal
COPY ./cmd/tesseractprocessor ./cmd/tesseractprocessor

RUN go build -tags static_all -o tesseractprocessor ./cmd/tesseractprocessor/*.go

# Install tesseract
# https://github.com/tesseract-ocr/tesseract/wiki
RUN apt-get update && apt-get install --no-install-recommends -y \
    tesseract-ocr \
    libtesseract-dev \
    && rm -rf /var/lib/apt/lists/*

CMD ["./tesseractprocessor"]
