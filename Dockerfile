FROM frinx/wasm-worker as wasmer

FROM node:12 as node
WORKDIR /resMgr
# Copy RM
COPY pools/allocating_strategies/strategies pools/allocating_strategies/strategies
COPY build_strategies.sh build_strategies.sh
# Build allocating strats
RUN ./build_strategies.sh

FROM golang:1.21.7 as build
ARG GITHUB_TOKEN_EXTERNAL_DOCKERFILE
WORKDIR /resMgr

# Copy RM
COPY . .

ENV GITHUB_TOKEN_EXTERNAL=$GITHUB_TOKEN_EXTERNAL_DOCKERFILE
RUN ./build.sh

# final image:
FROM golang:1.21.7

ARG RM_LOG_FILE=rm.log
ARG git_commit=unspecified
LABEL git_commit="${git_commit}"
LABEL org.opencontainers.image.source="https://github.com/FRINXio/resource-manager"

WORKDIR /resMgr

# Add log rotation
RUN apt-get update && apt-get --yes install logrotate
RUN echo "${RM_LOG_FILE} { \n rotate 5 \n weekly \n copytruncate \n compress \n missingok \n notifempty \n } \n " > /etc/logrotate.d/rm

# wasmer
COPY --from=wasmer /app/wasm-worker/.wasmer ./.wasmer
COPY --from=wasmer /app/wasm-worker/wasm ./wasm
# COPY built allocation strategies
RUN ./.wasmer/bin/wasmer ./wasm/quickjs/quickjs.wasm -- --std -e 'console.log("Wasmer works!")'

COPY run.sh run.sh
COPY go.mod go.mod
COPY --from=build /resMgr/resourceManager resourceManager

RUN groupadd -r resManager && \
    useradd --no-log-init -r -g resManager resManager && \ 
    chown -R resManager:resManager /resMgr && \
    mkdir -p /home/resManager && \
    chown -R resManager:resManager /home/resManager && \
    chown -R resManager:resManager /var/log 

USER resManager

CMD ./run.sh
