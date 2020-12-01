export WASMER_BIN=/resMgr/.wasmer/bin/wasmer
export WASMER_JS=/resMgr/wasm/quickjs/quickjs.wasm
export WASMER_PY=/resMgr/wasm/python/bin/python.wasm
export WASMER_PY_LIB=/resMgr/wasm/python/lib/
docker run --rm -it \
    -e WASMER_BIN=$WASMER_BIN \
    -e WASMER_JS=$WASMER_JS \
    -e WASMER_PY=$WASMER_PY \
    -e WASMER_PY_LIB=$WASMER_PY_LIB \
    -e WASMER_MAX_TIMEOUT_MILLIS=15000 \
    frinx/resource-manager go test -run Integration ./pools/...
