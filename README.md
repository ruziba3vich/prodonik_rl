# prodonik_rl

> **Redis-based Token Bucket Rate Limiter for Go**

[![Go Reference](https://pkg.go.dev/badge/github.com/ruziba3vich/prodonik_rl.svg)](https://pkg.go.dev/github.com/ruziba3vich/prodonik_rl)

---

## Overview

`prodonik_rl` is a lightweight and efficient **rate limiter** built with the **Token Bucket algorithm** and **Redis** as the storage backend.  
It is perfect for distributed environments where you need a fast, scalable way to throttle user requests based on IP addresses or custom keys.

✅ Distributed & scalable (uses Redis)  
✅ Implements **Token Bucket** strategy (allows bursts and refills smoothly)  
✅ Simple and minimalistic API  
✅ Versioned and ready to use via `go get`

---

## Installation

```bash
go get github.com/ruziba3vich/prodonik_rl
```
## Benchmarks

The `TokenBucketLimiter` is designed to be **lightweight** and **high-performance**.

In typical scenarios with a local or well-optimized Redis setup:

- **Request Check Time**: ~0.3ms - 0.8ms
- **Memory Usage**: Very low per bucket (just a few Redis hash fields)

>  Suitable for high-traffic APIs, microservices, and public endpoints.

_Real-world performance may vary based on network latency and Redis configuration._

---

## Common Pitfalls

- **Clock Drift**: Make sure all your services use synchronized time (e.g., NTP) since token calculation depends on timestamps.
- **Redis Persistence**: If Redis restarts and persistence is not configured properly, bucket states will be lost.
- **Window Size**: Setting a very small `window` might cause Redis to evict buckets before they are fully utilized.

### Best Practices
- Use a **replicated Redis** (like a Redis cluster) for production environments.
- Tune `maxTokens` and `refillRate` carefully based on your API's needs.
- Monitor Redis memory usage if you have a very large number of keys (buckets).

---

## How It Works

```text
+-----------------+
|  Incoming       |
|  Request        |
+--------+--------+
         |
         v
+--------+--------+
| Fetch current   |
| token count and |
| last update time|
+--------+--------+
         |
         v
+--------+--------+
| Calculate how   |
| many new tokens |
| should be added |
+--------+--------+
         |
         v
+--------+--------+
| Enough tokens?  |--No--> [ Deny Request ]
| (tokens > 0)    |
+--------+--------+
         |
        Yes
         |
         v
+--------+--------+
| Consume 1 token |
| Save back to    |
| Redis           |
+--------+--------+
         |
         v
 [ Allow Request ]
```

Time Axis ---->
```
+-------+       +-------+       +-------+
| Bucket|       | Bucket|       | Bucket|
| 80%   | ----> | 100%  | ----> | 100%  |
| Full  |       | Full  |       | Full  |
+-------+       +-------+       +-------+
(Refill happens over time based on refillRate)



