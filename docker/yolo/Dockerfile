FROM debian:buster

RUN apt-get update && apt-get install --no-install-recommends -y \
    apt-transport-https \
    build-essential \
    ca-certificates \
    git \
    wget \
    && rm -rf /var/lib/apt/lists/*

RUN git clone --depth=1 https://github.com/pjreddie/darknet.git
WORKDIR /darknet
RUN make
RUN wget https://pjreddie.com/media/files/yolov3.weights
