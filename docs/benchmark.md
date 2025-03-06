# Apache Benchmark Output

### Command

```bash
wallet_id=$(curl -s http://localhost:8080/api/v1/wallets -X POST)
ab -n 10000 -c 1000 http://localhost:8080/api/v1/wallets/$wallet_id
```

### Server Information

- **Server Hostname:** `localhost`
- **Server Port:** `8080`

### Request Details

- **Document Path:** `/api/v1/wallets/2adeb9bf-a962-46d1-aaae-6f7751c02154`
- **Document Length:** `1 bytes`

### Performance Metrics

- **Concurrency Level:** `1000`
- **Time taken for tests:** `3.014 seconds`
- **Complete requests:** `10000`
- **Failed requests:** `0`
- **Total transferred:** `1020000 bytes`
- **HTML transferred:** `10000 bytes`
- **Requests per second:** `3318.07 [#/sec] (mean)`
- **Time per request:** `301.380 [ms] (mean)`
- **Time per request:** `0.301 [ms] (mean, across all concurrent requests)`
- **Transfer rate:** `330.51 [Kbytes/sec] received`

### Connection Times (ms)

| Type       | Min | Mean | +/-sd | Median | Max |
| ---------- | --- | ---- | ----- | ------ | --- |
| Connect    | 0   | 6    | 14.0  | 0      | 70  |
| Processing | 69  | 284  | 61.8  | 268    | 532 |
| Waiting    | 4   | 283  | 61.7  | 267    | 532 |
| **Total**  | 74  | 289  | 68.9  | 271    | 552 |

### Request Time Percentiles (ms)

| Percentile | Time                  |
| ---------- | --------------------- |
| 50%        | 271                   |
| 66%        | 302                   |
| 75%        | 316                   |
| 80%        | 331                   |
| 90%        | 374                   |
| 95%        | 443                   |
| 98%        | 500                   |
| 99%        | 519                   |
| 100%       | 552 (longest request) |
