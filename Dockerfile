FROM integration_wasm-worker as wasmer

FROM node:12 as node
WORKDIR /resMgr
# Copy RM
COPY . .
# Build allocating strats
RUN ./build_strategies.sh

FROM golang:1.14.6-stretch
WORKDIR /resMgr

# Copy RM
COPY . .

COPY --from=wasmer /app/wasm-worker/.wasmer ./.wasmer
COPY --from=wasmer /app/wasm-worker/wasm ./wasm
# COPY built allocation strategies
COPY --from=node /resMgr/pools/allocating_strategies/strategies ./pools/allocating_strategies/strategies

# Test wasmer
RUN ./.wasmer/bin/wasmer ./wasm/quickjs/quickjs.wasm -- --std -e 'console.log("Wasmer works!")'

RUN ./build.sh
RUN go get github.com/go-delve/delve/cmd/dlv
CMD ./run.sh
