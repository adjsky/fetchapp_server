FROM ubuntu:21.04
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN apt-get install -y golang
RUN apt-get install -y python3
RUN apt-get install -y python3-pip
RUN apt-get install -y npm
COPY . /source_code
WORKDIR /source_code
RUN bash ./init.sh
CMD go run ./cmd