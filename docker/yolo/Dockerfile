FROM debian:buster

RUN apt-get update && apt-get install --no-install-recommends -y \
    git \
    && rm -rf /var/lib/apt/lists/*

RUN git clone --depth=1 https://github.com/pjreddie/darknet.git
WORKDIR /darknet
RUN make
