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
#### Create ipv4 prefix pool 100x serially x 0.76 ops/sec ±2.96% (8 runs sampled)
Measure insert a row into DB serially. One iteration is 100 mutations.
```
{
  total: {
    mean: 1310.375,
    p50: 1300,
    p75: 1320,
    p90: 1370,
    p97_5: 1370,
    p99: 1370,
    p99_9: 1370,
    p99_99: 1370,
    p99_999: 1370,
    max: 1379.577,
    totalCount: 8
  },
  setup: {
    mean: 9.625,
    p50: 7,
    p75: 7,
    p90: 31,
    p97_5: 31,
    p99: 31,
    p99_9: 31,
    p99_99: 31,
    p99_999: 31,
    max: 31.38,
    totalCount: 8
  },
  createAllocationPool: {
    mean: 1300,
    p50: 1290,
    p75: 1300,
    p90: 1370,
    p97_5: 1370,
    p99: 1370,
    p99_9: 1370,
    p99_99: 1370,
    p99_999: 1370,
    max: 1372.01,
    totalCount: 8
  }
}
```
#### Create ipv4 prefix pool 100x parallelly x 0.36 ops/sec ±4.24% (6 runs sampled)
Measure insert a row into DB parallelly. One iteration is 100 mutations.
```
{
  total: {
    mean: 2781.3333333333335,
    p50: 2730,
    p75: 2880,
    p90: 2940,
    p97_5: 2940,
    p99: 2940,
    p99_9: 2940,
    p99_99: 2940,
    p99_999: 2940,
    max: 2944.977,
    totalCount: 6
  },
  setup: {
    mean: 7.166666666666667,
    p50: 7,
    p75: 8,
    p90: 8,
    p97_5: 8,
    p99: 8,
    p99_9: 8,
    p99_99: 8,
    p99_999: 8,
    max: 8.741,
    totalCount: 6
  }
}
```
#### allocate 100 ipv4_prefix_pool resources serially x 0.16 ops/sec ±1.46% (5 runs sampled)
Measure e2e performance of ipv4 prefix strategy. One iteration is 100 mutations.
```
{
  total: {
    mean: 6271.6,
    p50: 6240,
    p75: 6300,
    p90: 6380,
    p97_5: 6380,
    p99: 6380,
    p99_9: 6380,
    p99_99: 6380,
    p99_999: 6380,
    max: 6386.231,
    totalCount: 5
  },
  setup: {
    mean: 21,
    p50: 21,
    p75: 22,
    p90: 23,
    p97_5: 23,
    p99: 23,
    p99_9: 23,
    p99_99: 23,
    p99_999: 23,
    max: 23.478,
    totalCount: 5
  },
  allocate: {
    mean: 6250,
    p50: 6210,
    p75: 6280,
    p90: 6360,
    p97_5: 6360,
    p99: 6360,
    p99_9: 6360,
    p99_99: 6360,
    p99_999: 6360,
    max: 6367.2,
    totalCount: 5
  }
}
```
#### allocate 100 ipv4_prefix_pool resources parallelly x 0.07 ops/sec ±11.85% (5 runs sampled)
Measure e2e performance of ipv4 prefix strategy. One iteration is 100 mutations.
```
{
  total: {
    mean: 14260,
    p50: 14000,
    p75: 14900,
    p90: 16100,
    p97_5: 16100,
    p99: 16100,
    p99_9: 16100,
    p99_99: 16100,
    p99_999: 16100,
    max: 16167.741,
    totalCount: 5
  },
  setup: {
    mean: 22.8,
    p50: 22,
    p75: 23,
    p90: 27,
    p97_5: 27,
    p99: 27,
    p99_9: 27,
    p99_99: 27,
    p99_999: 27,
    max: 27.565,
    totalCount: 5
  },
  awaitPromisses: {
    mean: 14237.6,
    p50: 14000,
    p75: 14800,
    p90: 16100,
    p97_5: 16100,
    p99: 16100,
    p99_9: 16100,
    p99_99: 16100,
    p99_999: 16100,
    max: 16145.263,
    totalCount: 5
  }
}
```
#### allocate 100 ipv4_pool resources serially x 0.20 ops/sec ±1.27% (5 runs sampled)
Measure e2e performance of ipv4 strategy. One iteration is 100 mutations.
```
{
  total: {
    mean: 5046,
    p50: 5040,
    p75: 5060,
    p90: 5120,
    p97_5: 5120,
    p99: 5120,
    p99_9: 5120,
    p99_99: 5120,
    p99_999: 5120,
    max: 5122.371,
    totalCount: 5
  },
  setup: {
    mean: 83.6,
    p50: 83,
    p75: 84,
    p90: 87,
    p97_5: 87,
    p99: 87,
    p99_9: 87,
    p99_99: 87,
    p99_999: 87,
    max: 87.74,
    totalCount: 5
  },
  allocate: {
    mean: 4961.2,
    p50: 4950,
    p75: 4980,
    p90: 5030,
    p97_5: 5030,
    p99: 5030,
    p99_9: 5030,
    p99_99: 5030,
    p99_999: 5030,
    max: 5034.225,
    totalCount: 5
  }
}
```
#### allocate 100 ipv4_pool resources parallelly x 0.08 ops/sec ±3.27% (5 runs sampled)
Measure e2e performance of ipv4 strategy. One iteration is 100 mutations.
```
{
  total: {
    mean: 12258.4,
    p50: 12100,
    p75: 12500,
    p90: 12500,
    p97_5: 12500,
    p99: 12500,
    p99_9: 12500,
    p99_99: 12500,
    p99_999: 12500,
    max: 12594.177,
    totalCount: 5
  },
  setup: {
    mean: 85.2,
    p50: 83,
    p75: 88,
    p90: 90,
    p97_5: 90,
    p99: 90,
    p99_9: 90,
    p99_99: 90,
    p99_999: 90,
    max: 90.055,
    totalCount: 5
  },
  awaitPromisses: {
    mean: 12170.4,
    p50: 12100,
    p75: 12400,
    p90: 12500,
    p97_5: 12500,
    p99: 12500,
    p99_9: 12500,
    p99_99: 12500,
    p99_999: 12500,
    max: 12510.955,
    totalCount: 5
  }
}
```
