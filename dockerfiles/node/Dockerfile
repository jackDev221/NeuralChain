FROM golang:1.12-alpine as builder
RUN echo "------ Builder ------"
RUN apk add --no-cache make gcc musl-dev linux-headers ca-certificates git

# Create the user and group files that will be used in the running
# container to run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Copy the predefined token file into the location that git depends on
COPY dockerfiles/node/token /root/.netrc
RUN chmod 600 /root/.netrc

COPY . /evrynet-node
WORKDIR "/evrynet-node"
RUN go build ./cmd/gev
RUN pwd
RUN ls

FROM alpine:latest
RUN apk add --no-cache ca-certificates bash perl

RUN echo "------ Node ------"

COPY --from=builder /evrynet-node/gev /
COPY ./deploy/devnet/node/genesis.json /
COPY ./dockerfiles/node/entrypoint.sh /

RUN chmod +x entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

CMD ["/gev", "--help"]
