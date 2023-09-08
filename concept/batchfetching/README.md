# Batch Fetching

## Objectives

By leveraging batch-fetching, developers can optimize their applications for improved performance, reduced network overhead, and efficient data retrieval.

## Summary

Batch-fetching offers several benefits for data scientists and developers working with Hibernate:

Improved Performance: By reducing the number of database queries and network round-trips, batch-fetching can greatly improve the performance of applications. This is particularly beneficial when dealing with large datasets or complex object graphs that involve multiple associations.

Reduced Network Overhead: Network latency can be a significant bottleneck in distributed applications. Batch-fetching minimizes the amount of data transferred between the application and the database, resulting in reduced network overhead and improved response times.

Optimized Resource Utilization: Since batch-fetching fetches data in bulk, it allows the database to optimize its query execution plans and utilize resources more efficiently. This can result in better utilization of CPU, memory, and disk I/O, leading to overall improved database performance.

Avoidance of N+1 Select Problem: By fetching associated entities in batches, batch-fetching mitigates the N+1 select problem, which occurs when lazy loading leads to excessive database queries. This ensures that the application efficiently retrieves all the required data without unnecessary round-trips.


## Scenario

| No | Scenario | Goals |
| ------------- | ------------- | ------------- |
| 1  | Create usecase function with & without batch fetching | Benchmark |

## Success Criteria

### 1. Create usecase function with & without batch fetching

#### Upsert channel without batch fetching

Request upsert with 26 datas

Jaeger:
```
{
    "timestamp": 1693312392232594,
    "fields": [
        {
            "key": "event",
            "type": "string",
            "value": "channel upsert success"
        },
        {
            "key": "type",
            "type": "string",
            "value": "Success"
        }
    ]
}
```

![Alt text](jaeger-upsert.png?raw=true "Jaeger upsert")

Postman:

![Alt text](postman-upsert.png?raw=true "Postman upsert")

#### Upsert channel with batch fetching

Request upsert with 26 datas

Jaeger:
```
{
    "timestamp": 1693312517313033,
    "fields": [
        {
            "key": "event",
            "type": "string",
            "value": "channel upsert batch fetching success"
        },
        {
            "key": "type",
            "type": "string",
            "value": "Success"
        }
    ]
}
```

![Alt text](jaeger-upsert-batch-fetching.png?raw=true "Jaeger upsert Batch Fetching")

Postman:

![Alt text](postman-upsert-batch-fetching.png?raw=true "Postman upsert Batch Fetching")

#### Benchmark

| Result by | Without batch fetching | With batch fetching | Summary |
| ------------- | ------------- | ------------- | ------------- |
| Jaeger | 117.14ms | 60.6ms | Batch fetching 93.3% faster |
| Postman | 179ms | 102ms | Batch fetching 75.5% faster |