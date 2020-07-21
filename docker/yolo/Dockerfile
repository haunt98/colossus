FROM debian:buster

RUN apt-get update && apt-get install --no-install-recommends -y \
    apt-transport-https \
    build-essential \
    ca-certificates \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*

RUN git clone --depth=1 https://github.com/AlexeyAB/darknet.git
WORKDIR /darknet
RUN make
RUN curl -O -L -C - https://github.com/AlexeyAB/darknet/releases/download/darknet_yolo_v4_pre/yolov4-tiny.weights
