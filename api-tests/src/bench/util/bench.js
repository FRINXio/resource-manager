import Benchmark from 'benchmark';
import hdr from 'hdr-histogram-js';
import microtime from 'microtime';
import tap from 'tap';

// configuration
const useWasmForHistograms = false
const outputHistogramsToPlotter = false

// global state
const globalCtx = {} // holds benchmark -> histograms mapping

export async function executeBenchmark(t, histograms, fn) {
    const name = t.name;
    console.log(`Starting '${name}'`)
    const suite = new Benchmark.Suite
    return new Promise(function (resolve, reject) {
        suite
            .add(name, async (deferred) => {
                await record(histograms, 'total', async () => await fn(histograms));
                deferred.resolve();
            }, { defer: true })
            .on('cycle', (event) => {
                finishBenchmark(name, event, histograms)
                resolve();
                t.end();
            })
            .run();
    })
}

function finishBenchmark(name, event, histograms) {
    console.log(String(event.target))
    globalCtx[name] = outputHistograms(histograms)
}

function outputHistograms(histograms) {
    const summaries = {}
    const encodedSummaries = {}
    const output = { summaries, encodedSummaries }
    for (const key in histograms) {
        const h = histograms[key]
        summaries[key] = {
            mean: h.mean,
            ...h.summary
        }
        if (outputHistogramsToPlotter) {
            encodedSummaries[key] = hdr.encodeIntoCompressedBase64(h)
        }
        h.destroy()
    }
    console.log(summaries)
    if (outputHistogramsToPlotter) {
        console.log(encodedSummaries) // view using https://hdrhistogram.github.io/HdrHistogramJSDemo/plotFiles.html
    }
    return output
}

function newHist() {
    return hdr.build({ useWebAssembly: useWasmForHistograms })
}

export function newFixedHist(value) {
    const result = newHist()
    result.recordValue(value)
    return result
}

export function recordFixedValue(histograms, key, value) {
    if (!histograms[key]) {
        histograms[key] = newHist()
    }
    histograms[key].recordValue(value)
    return value
}


export async function record(histograms, key, fn) {
    if (!histograms[key]) {
        histograms[key] = newHist()
    }
    const start = microtime.now()
    const result = await fn()
    histograms[key].recordValue((microtime.now() - start) / 1000)
    return result
}

export function bench(name, fnWithHistograms) {
    tap.test(name, async (t) =>
        await executeBenchmark(t, {}, fnWithHistograms)
    )
}
