############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
ARG DOCKER_TAG=0.0.0
# checkout the project 
WORKDIR /builder
COPY . .
# Fetch dependencies.
RUN go get -d 
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /defi-portal-scanner -ldflags="-s -w -extldflags \"-static\" -X main.Version=$DOCKER_TAG"
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable + data
COPY --from=builder /defi-portal-scanner /
# Copy config file 
COPY ./private/config.yaml /
COPY ./private/protocols.json /
# Run the whole shebang.
ENTRYPOINT [ "/defi-portal-scanner" ]
CMD [ "listen", "--config", "/config.yaml", "--protocols", "/protocols.json", "--scan", "--http"]
