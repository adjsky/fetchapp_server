FROM ubuntu:21.04
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN apt-get install -y golang
RUN apt-get install -y python3
RUN apt-get install -y python3-pip
COPY . /source_code
WORKDIR /source_code
RUN pip3 install -r ./requirements.txt
CMD go run ./cmd