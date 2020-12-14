# Build Geth in a stock Go builder container
FROM golang:1.12-alpine as builder
RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /neuralChain
WORKDIR /neuralChain
RUN go mod download && go build ./cmd/gev

# -------------------------------------------------------------
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR "/root"
COPY --from=builder /neuralChain/gev /usr/local/bin/
COPY --from=builder /neuralChain/tests/onenode/nodedata ./nodedata
COPY --from=builder /neuralChain/tests/onenode/init_node.sh ./init_node.sh

RUN chmod +x nodedata/one_node_genesis.json ./init_node.sh

CMD ./init_node.sh

EXPOSE 8545
EXPOSE 30303