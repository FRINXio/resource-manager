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
#### Create ipv4 prefix pool 100x serially x 0.59 ops/sec ±2.04% (7 runs sampled)
Measure 100x insert a row into DB serially.
```
{
  total: {
    mean: 1698.857142857143,
    p50: 1700,
    p75: 1730,
    p90: 1740,
    p97_5: 1740,
    p99: 1740,
    p99_9: 1740,
    p99_99: 1740,
    p99_999: 1740,
    max: 1742.669,
    totalCount: 7
  },
  setup: {
    mean: 98.14285714285714,
    p50: 93,
    p75: 100,
    p90: 130,
    p97_5: 130,
    p99: 130,
    p99_9: 130,
    p99_99: 130,
    p99_999: 130,
    max: 130.475,
    totalCount: 7
  },
  createAllocationPool: {
    mean: 1600.142857142857,
    p50: 1600,
    p75: 1630,
    p90: 1650,
    p97_5: 1650,
    p99: 1650,
    p99_9: 1650,
    p99_99: 1650,
    p99_999: 1650,
    max: 1655.364,
    totalCount: 7
  }
}
```
#### Create ipv4 prefix pool 100x parallelly x 0.33 ops/sec ±3.02% (6 runs sampled)
Measure 100x insert a row into DB parallelly. Currently 1,8x slower than when allocating serially.
```
{
  total: {
    mean: 3015.6666666666665,
    p50: 2990,
    p75: 3040,
    p90: 3160,
    p97_5: 3160,
    p99: 3160,
    p99_9: 3160,
    p99_99: 3160,
    p99_999: 3160,
    max: 3164.317,
    totalCount: 6
  },
  setup: {
    mean: 101.16666666666667,
    p50: 99,
    p75: 102,
    p90: 116,
    p97_5: 116,
    p99: 116,
    p99_9: 116,
    p99_99: 116,
    p99_999: 116,
    max: 116.284,
    totalCount: 6
  }
}
```
#### allocate 100 ipv4_prefix_pool resources serially x 0.12 ops/sec ±2.18% (5 runs sampled)
Measure e2e performance of ipv4 prefix strategy.
```
{
  total: {
    mean: 8549.6,
    p50: 8550,
    p75: 8660,
    p90: 8730,
    p97_5: 8730,
    p99: 8730,
    p99_9: 8730,
    p99_99: 8730,
    p99_999: 8730,
    max: 8728.536,
    totalCount: 5
  },
  setup: {
    mean: 113.4,
    p50: 118,
    p75: 120,
    p90: 124,
    p97_5: 124,
    p99: 124,
    p99_9: 124,
    p99_99: 124,
    p99_999: 124,
    max: 124.453,
    totalCount: 5
  },
  allocate: {
    mean: 8434.4,
    p50: 8440,
    p75: 8540,
    p90: 8600,
    p97_5: 8600,
    p99: 8600,
    p99_9: 8600,
    p99_99: 8600,
    p99_999: 8600,
    max: 8604.046,
    totalCount: 5
  }
}
```
#### allocate 100 ipv4_prefix_pool resources parallelly x 0.05 ops/sec ±3.80% (5 runs sampled)
Measure e2e performance of ipv4 prefix strategy. Currently 2,3x slower than when allocating serially.
```
{
  total: {
    mean: 19502.4,
    p50: 19500,
    p75: 19800,
    p90: 20200,
    p97_5: 20200,
    p99: 20200,
    p99_9: 20200,
    p99_99: 20200,
    p99_999: 20200,
    max: 20248.69,
    totalCount: 5
  },
  setup: {
    mean: 121.2,
    p50: 119,
    p75: 132,
    p90: 135,
    p97_5: 135,
    p99: 135,
    p99_9: 135,
    p99_99: 135,
    p99_999: 135,
    max: 135.701,
    totalCount: 5
  },
  awaitPromisses: {
    mean: 19384,
    p50: 19400,
    p75: 19700,
    p90: 20100,
    p97_5: 20100,
    p99: 20100,
    p99_9: 20100,
    p99_99: 20100,
    p99_999: 20100,
    max: 20115.599,
    totalCount: 5
  }
}
```
#### allocate 100 ipv4_pool resources serially x 0.14 ops/sec ±1.18% (5 runs sampled)
Measure e2e performance of ipv4 strategy.
```
{
  total: {
    mean: 7227.6,
    p50: 7220,
    p75: 7230,
    p90: 7330,
    p97_5: 7330,
    p99: 7330,
    p99_9: 7330,
    p99_99: 7330,
    p99_999: 7330,
    max: 7332.764,
    totalCount: 5
  },
  setup: {
    mean: 322.8,
    p50: 308,
    p75: 335,
    p90: 373,
    p97_5: 373,
    p99: 373,
    p99_9: 373,
    p99_99: 373,
    p99_999: 373,
    max: 373.2,
    totalCount: 5
  },
  allocate: {
    mean: 6902.8,
    p50: 6890,
    p75: 6920,
    p90: 6950,
    p97_5: 6950,
    p99: 6950,
    p99_9: 6950,
    p99_99: 6950,
    p99_999: 6950,
    max: 6959.528,
    totalCount: 5
  }
}
```
#### allocate 100 ipv4_pool resources parallelly x 0.05 ops/sec ±20.05% (5 runs sampled)
Measure e2e performance of ipv4 strategy. Currently 2,6x slower than when allocating serially.
```
{
  total: {
    mean: 18724.8,
    p50: 17400,
    p75: 18400,
    p90: 23900,
    p97_5: 23900,
    p99: 23900,
    p99_9: 23900,
    p99_99: 23900,
    p99_999: 23900,
    max: 23997.025,
    totalCount: 5
  },
  setup: {
    mean: 340.6,
    p50: 322,
    p75: 349,
    p90: 440,
    p97_5: 440,
    p99: 440,
    p99_9: 440,
    p99_99: 440,
    p99_999: 440,
    max: 440.361,
    totalCount: 5
  },
  awaitPromisses: {
    mean: 18383.2,
    p50: 17100,
    p75: 18100,
    p90: 23500,
    p97_5: 23500,
    p99: 23500,
    p99_9: 23500,
    p99_99: 23500,
    p99_999: 23500,
    max: 23556.596,
    totalCount: 5
  }
}
```
