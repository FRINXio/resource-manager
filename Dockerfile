FROM integration_wasm-worker as wasmer
FROM golang:1.14.6-stretch
WORKDIR /resMgr

# Copy RM
COPY . .

# COPY wasmer from wasm workflows worker
COPY --from=wasmer /app/wasm-worker/.wasmer ./.wasmer
COPY --from=wasmer /app/wasm-worker/wasm ./wasm

# Test wasmer
RUN ./.wasmer/bin/wasmer ./wasm/quickjs/quickjs.wasm -- --std -e 'console.log("Wasmer works!")'

RUN ./build.sh
CMD ./run.sh
