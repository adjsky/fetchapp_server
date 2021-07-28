FROM alpine:3.14
RUN apk add --no-cache \
    go \
    python3 \
    py3-pip
COPY . /source_code
WORKDIR /source_code
RUN pip3 install -r ./requirements.txt
CMD go run ./cmd
