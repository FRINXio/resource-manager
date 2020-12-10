# Benchmarks

## Running benchmarks
Execute a benchmark directly using node:
```sh
node src/bench/strategy.bench.js
```

## Results:
### strategy.bench.js
#### create and test simple JS strategy once x 24.44 ops/sec ±2.39% (60 runs sampled)
One JS execution takes around 30-50ms.
```
{
  total: {
    mean: 40.166666666666664,
    p50: 40,
    p75: 42,
    p90: 43,
    p97_5: 45,
    p99: 48,
    p99_9: 68,
    p99_99: 68,
    p99_999: 68,
    max: 68.066,
    totalCount: 120
  },
  setup: {
    mean: 4.925,
    p50: 5,
    p75: 5,
    p90: 6,
    p97_5: 7,
    p99: 10,
    p99_9: 30,
    p99_99: 30,
    p99_999: 30,
    max: 30.319,
    totalCount: 120
  },
  testStrategy: {
    mean: 34.775,
    p50: 35,
    p75: 36,
    p90: 37,
    p97_5: 38,
    p99: 39,
    p99_9: 39,
    p99_99: 39,
    p99_999: 39,
    max: 39.265,
    totalCount: 120
  }
}
```

#### create and test simple PY strategy once x 1.90 ops/sec ±0.85% (14 runs sampled)
One PY execution takes arund 500-800ms, so more than 10x slower than JS and with
less predictable duration.
```
{
  total: {
    mean: 526.7857142857143,
    p50: 523,
    p75: 534,
    p90: 538,
    p97_5: 539,
    p99: 539,
    p99_9: 539,
    p99_99: 539,
    p99_999: 539,
    max: 539.409,
    totalCount: 14
  },
  setup: {
    mean: 5.928571428571429,
    p50: 6,
    p75: 6,
    p90: 6,
    p97_5: 7,
    p99: 7,
    p99_9: 7,
    p99_99: 7,
    p99_999: 7,
    max: 7.311,
    totalCount: 14
  },
  testStrategy: {
    mean: 520.2857142857143,
    p50: 517,
    p75: 528,
    p90: 531,
    p97_5: 532,
    p99: 532,
    p99_9: 532,
    p99_99: 532,
    p99_999: 532,
    max: 532.979,
    totalCount: 14
  }
}
```


#### create and test simple JS strategy 100x x 1.32 ops/sec ±1.27% (11 runs sampled)
Running 100 JS executions in parallel. On averge, one execution takes 7,5ms (awaitPromisses.mean/100).

```
{
  total: {
    mean: 756.9090909090909,
    p50: 752,
    p75: 765,
    p90: 766,
    p97_5: 789,
    p99: 789,
    p99_9: 789,
    p99_99: 789,
    p99_999: 789,
    max: 789.721,
    totalCount: 11
  },
  setup: {
    mean: 6.090909090909091,
    p50: 6,
    p75: 7,
    p90: 7,
    p97_5: 7,
    p99: 7,
    p99_9: 7,
    p99_99: 7,
    p99_999: 7,
    max: 7.337,
    totalCount: 11
  },
  testStrategy: {
    mean: 425.2218181818182,
    p50: 424,
    p75: 593,
    p90: 696,
    p97_5: 738,
    p99: 754,
    p99_9: 778,
    p99_99: 779,
    p99_999: 779,
    max: 779.609,
    totalCount: 1100
  },
  awaitPromisses: {
    mean: 745.4545454545455,
    p50: 740,
    p75: 753,
    p90: 753,
    p97_5: 778,
    p99: 778,
    p99_9: 778,
    p99_99: 778,
    p99_999: 778,
    max: 778.69,
    totalCount: 11
  }
}
```

#### create and test simple PY strategy 100x x 0.08 ops/sec ±9.21% (5 runs sampled)
Running 100 PY executions in parallel. On averge, one execution takes 129,72ms (awaitPromisses.mean/100).
Compared to JS, this is 17x slower.

```
{
  total: {
    mean: 12984.8,
    p50: 13100,
    p75: 13400,
    p90: 14200,
    p97_5: 14200,
    p99: 14200,
    p99_9: 14200,
    p99_99: 14200,
    p99_999: 14200,
    max: 14233.261,
    totalCount: 5
  },
  setup: {
    mean: 7.6,
    p50: 8,
    p75: 8,
    p90: 9,
    p97_5: 9,
    p99: 9,
    p99_9: 9,
    p99_99: 9,
    p99_999: 9,
    max: 9.772,
    totalCount: 5
  },
  testStrategy: {
    mean: 7067.562,
    p50: 7030,
    p75: 10200,
    p90: 12000,
    p97_5: 13300,
    p99: 13700,
    p99_9: 14200,
    p99_99: 14200,
    p99_999: 14200,
    max: 14220.835,
    totalCount: 500
  },
  awaitPromisses: {
    mean: 12972,
    p50: 13100,
    p75: 13400,
    p90: 14200,
    p97_5: 14200,
    p99: 14200,
    p99_9: 14200,
    p99_99: 14200,
    p99_999: 14200,
    max: 14216.477,
    totalCount: 5
  }
}
```

### ipv4.bench.js
#### Create ipv4 prefix pool 100x serially x 0.73 ops/sec ±2.04% (8 runs sampled)
Measure insert a row into DB serially. One iteration is 100 mutations.
```
{
  total: {
    mean: 1368,
    p50: 1360,
    p75: 1380,
    p90: 1430,
    p97_5: 1430,
    p99: 1430,
    p99_9: 1430,
    p99_99: 1430,
    p99_999: 1430,
    max: 1433.152,
    totalCount: 8
  },
  setup: {
    mean: 16.375,
    p50: 13,
    p75: 14,
    p90: 40,
    p97_5: 40,
    p99: 40,
    p99_9: 40,
    p99_99: 40,
    p99_999: 40,
    max: 40.281,
    totalCount: 8
  },
  createAllocationPool: {
    mean: 1350.75,
    p50: 1340,
    p75: 1360,
    p90: 1390,
    p97_5: 1390,
    p99: 1390,
    p99_9: 1390,
    p99_99: 1390,
    p99_999: 1390,
    max: 1392.296,
    totalCount: 8
  }
}
```
#### Create ipv4 prefix pool 100x parallelly x 2.42 ops/sec ±7.06% (16 runs sampled)
Measure insert a row into DB parallelly. One iteration is 100 mutations.
```
{
  total: {
    mean: 413.0625,
    p50: 403,
    p75: 421,
    p90: 514,
    p97_5: 516,
    p99: 516,
    p99_9: 516,
    p99_99: 516,
    p99_999: 516,
    max: 516.215,
    totalCount: 16
  },
  setup: {
    mean: 14.9375,
    p50: 14,
    p75: 16,
    p90: 17,
    p97_5: 18,
    p99: 18,
    p99_9: 18,
    p99_99: 18,
    p99_999: 18,
    max: 18.253,
    totalCount: 16
  }
}
```
#### allocate 100 ipv4_prefix_pool resources serially x 0.15 ops/sec ±1.74% (5 runs sampled)
Measure e2e performance of ipv4 prefix strategy. One iteration is 100 mutations.
```
{
  total: {
    mean: 6877.2,
    p50: 6830,
    p75: 6930,
    p90: 7010,
    p97_5: 7010,
    p99: 7010,
    p99_9: 7010,
    p99_99: 7010,
    p99_999: 7010,
    max: 7019.732,
    totalCount: 5
  },
  setup: {
    mean: 30.2,
    p50: 31,
    p75: 31,
    p90: 33,
    p97_5: 33,
    p99: 33,
    p99_9: 33,
    p99_99: 33,
    p99_999: 33,
    max: 33.93,
    totalCount: 5
  },
  allocate: {
    mean: 6846.8,
    p50: 6800,
    p75: 6900,
    p90: 6980,
    p97_5: 6980,
    p99: 6980,
    p99_9: 6980,
    p99_99: 6980,
    p99_999: 6980,
    max: 6985.26,
    totalCount: 5
  }
}
```
#### allocate 100 ipv4_prefix_pool resources parallelly x 0.06 ops/sec ±6.36% (5 runs sampled)
Measure e2e performance of ipv4 prefix strategy. One iteration is 100 mutations.
```
{
  total: {
    mean: 16483.2,
    p50: 16400,
    p75: 17200,
    p90: 17400,
    p97_5: 17400,
    p99: 17400,
    p99_9: 17400,
    p99_99: 17400,
    p99_999: 17400,
    max: 17399.918,
    totalCount: 5
  },
  setup: {
    mean: 32.4,
    p50: 32,
    p75: 32,
    p90: 37,
    p97_5: 37,
    p99: 37,
    p99_9: 37,
    p99_99: 37,
    p99_999: 37,
    max: 37.646,
    totalCount: 5
  },
  awaitPromisses: {
    mean: 16452.8,
    p50: 16400,
    p75: 17200,
    p90: 17300,
    p97_5: 17300,
    p99: 17300,
    p99_9: 17300,
    p99_99: 17300,
    p99_999: 17300,
    max: 17367.717,
    totalCount: 5
  }
}
```
#### allocate 100 ipv4_pool resources serially x 0.18 ops/sec ±2.10% (5 runs sampled)
Measure e2e performance of ipv4 strategy. One iteration is 100 mutations.
```
{
  total: {
    mean: 5510,
    p50: 5470,
    p75: 5530,
    p90: 5660,
    p97_5: 5660,
    p99: 5660,
    p99_9: 5660,
    p99_99: 5660,
    p99_999: 5660,
    max: 5660.634,
    totalCount: 5
  },
  setup: {
    mean: 113.2,
    p50: 112,
    p75: 112,
    p90: 126,
    p97_5: 126,
    p99: 126,
    p99_9: 126,
    p99_99: 126,
    p99_999: 126,
    max: 126.921,
    totalCount: 5
  },
  allocate: {
    mean: 5395.6,
    p50: 5350,
    p75: 5420,
    p90: 5550,
    p97_5: 5550,
    p99: 5550,
    p99_9: 5550,
    p99_99: 5550,
    p99_999: 5550,
    max: 5548.543,
    totalCount: 5
  }
}
```
#### allocate 100 ipv4_pool resources parallelly x 0.07 ops/sec ±8.66% (5 runs sampled)
Measure e2e performance of ipv4 strategy. One iteration is 100 mutations.
```
{
  total: {
    mean: 13533.6,
    p50: 12900,
    p75: 14200,
    p90: 14800,
    p97_5: 14800,
    p99: 14800,
    p99_9: 14800,
    p99_99: 14800,
    p99_999: 14800,
    max: 14820.072,
    totalCount: 5
  },
  setup: {
    mean: 108.4,
    p50: 106,
    p75: 110,
    p90: 117,
    p97_5: 117,
    p99: 117,
    p99_9: 117,
    p99_99: 117,
    p99_999: 117,
    max: 117.434,
    totalCount: 5
  },
  awaitPromisses: {
    mean: 13424.8,
    p50: 12800,
    p75: 14100,
    p90: 14700,
    p97_5: 14700,
    p99: 14700,
    p99_9: 14700,
    p99_99: 14700,
    p99_999: 14700,
    max: 14702.587,
    totalCount: 5
  }
}
```
