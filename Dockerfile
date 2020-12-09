FROM frinx/wasm-worker as wasmer
#FROM facebookconnectivity-hub-docker.jfrog.io/wasm-worker as wasmer

FROM node:12 as node
WORKDIR /resMgr
# Copy RM
COPY pools/allocating_strategies/strategies pools/allocating_strategies/strategies
COPY build_strategies.sh build_strategies.sh
# Build allocating strats
RUN ./build_strategies.sh

FROM golang:1.14.6-stretch
ARG RM_LOG_FILE=rm.log
ARG GITHUB_TOKEN_EXTERNAL_DOCKERFILE

# Add log rotation
RUN apt-get update && apt-get --yes install logrotate
RUN echo "${RM_LOG_FILE} { \n rotate 5 \n weekly \n copytruncate \n compress \n missingok \n notifempty \n } \n " > /etc/logrotate.d/rm

WORKDIR /resMgr

COPY --from=wasmer /app/wasm-worker/.wasmer ./.wasmer
COPY --from=wasmer /app/wasm-worker/wasm ./wasm
# COPY built allocation strategies
COPY --from=node /resMgr/pools/allocating_strategies/strategies ./pools/allocating_strategies/strategies

# Test wasmer
RUN ./.wasmer/bin/wasmer ./wasm/quickjs/quickjs.wasm -- --std -e 'console.log("Wasmer works!")'

# Copy RM
COPY . .

ENV GITHUB_TOKEN_EXTERNAL=$GITHUB_TOKEN_EXTERNAL_DOCKERFILE
RUN ./build.sh
RUN go get github.com/go-delve/delve/cmd/dlv
ENV SKIP_BUILD=true
CMD ./run.sh
