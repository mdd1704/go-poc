# Locking Row

## Objectives

Ensuring data integrity and concurrency in multi-user applications.

## Summary

Locking row is a locking strategy that allows multiple transactions to access and modify different rows of the same table at the same time, without interfering with each other. Locking row prevents data inconsistencies and conflicts by locking only the rows that are affected by a transaction, while leaving the rest of the table available for other transactions. Locking row is often used in scenarios where there is a high degree of concurrency and contention for the same data, such as online shopping, banking, or reservation systems.

## Scenario

| No | Scenario | Goals |
| ------------- | ------------- | ------------- |
| 1  | Concurrent process of usecase function with & without locking row | Benchmark |

## Success Criteria

### 1. Concurrent process of usecase function with & without locking row

#### Upsert channel without locking row

#### Upsert channel with locking row

#### Benchmark

| Without DB transaction | With DB transaction | Summary |
|  ------------- | ------------- | ------------- |
| No rollback | Rollback  | With DB transaction is more reliable |