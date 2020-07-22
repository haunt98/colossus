FROM debian:buster

RUN apt-get update && apt-get install --no-install-recommends -y \
    curl \
    openjdk-11-jdk \
    unzip \
    && rm -rf /var/lib/apt/lists/*

RUN curl -O -L -C - https://www.languagetool.org/download/LanguageTool-5.0.zip \
    && unzip LanguageTool-5.0.zip \
    && rm -rf LanguageTool-5.0.zip
WORKDIR /LanguageTool-5.0
