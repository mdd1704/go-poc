# Database Transaction

![Alt text](db-transaction.png?raw=true "Database Transaction")

## Objectives

There are 2 main reasons:

1. We want our unit of work to be reliable and consistent, even in case of system failure.
2. We want to provide isolation between programs that access the database concurrently.

## Summary

In order to achieve these 2 objectives, a database transaction must satisfy the ACID properties, where:

* `A` is Atomicity, which means either all operations of the transaction complete successfully, or the whole transaction fails, and everything is rolled back, the database is unchanged.
* `C` is Consistency, which means the database state should remains valid after the transaction is executed. More precisely, all data written to the database must be valid according to predefined rules, including constraints, cascades, and triggers.
* `I` is Isolation, meaning all transactions that run concurrently should not affect each other. There are several levels of isolation that defines when the changes made by 1 transaction can be visible to others. We will learn more about it in another lecture.
* The last property is `D`, which stands for Durability. It basically means that all data written by a successful transaction must stay in a persistent storage and cannot be lost, even in case of system failure.

## Scenario

| No | Scenario | Goals |
| ------------- | ------------- | ------------- |
| 1  | Failure in the middle proccess of usecase function with & without database transaction | Benchmark |

## Success Criteria

### 1. Failure in the middle proccess of usecase function with & without database transaction

#### Started state

MySQL Database:
```
mysql> select * from channels;
+--------------------------------------+------+---------------------+---------------------+
| id                                   | code | created_at          | updated_at          |
+--------------------------------------+------+---------------------+---------------------+
| 1bfe1431-50e1-4038-b6a6-a03c40e41e4b | A    | 2023-08-29 15:00:55 | 2023-08-29 15:00:55 |
+--------------------------------------+------+---------------------+---------------------+
1 row in set (0.00 sec)
```

#### Upsert channel without database transaction

Request:
```
curl --location 'localhost:8000/api/channel/upsert-batch-fetching' \
--header 'Content-Type: application/json' \
--data '[
    {
        "code": "B"
    },
    {
        "code": "A"
    }
]'
```

Response:
```
{
    "transaction_id": "30d4002d-3822-4fdc-acf5-1ab721e9713d",
    "success": false,
    "data": [
        {
            "id": "fa7a8949-9909-431c-90d4-6081ab565081",
            "code": "A",
            "message": "Error 1062 (23000): Duplicate entry 'A' for key 'channels.code'"
        }
    ],
    "error": null
}
```

MySQL Database:
```
mysql> select * from channels;
+--------------------------------------+------+---------------------+---------------------+
| id                                   | code | created_at          | updated_at          |
+--------------------------------------+------+---------------------+---------------------+
| 1bfe1431-50e1-4038-b6a6-a03c40e41e4b | A    | 2023-08-29 15:00:55 | 2023-08-29 15:00:55 |
| 715975b4-111d-4b06-a320-9b98fbc47839 | B    | 2023-08-29 15:04:47 | 2023-08-29 15:04:47 |
+--------------------------------------+------+---------------------+---------------------+
2 rows in set (0.00 sec)
```

#### Upsert channel with database transaction

Request:
```
curl --location 'localhost:8000/api/channel/upsert-with-transaction' \
--header 'Content-Type: application/json' \
--data '[
    {
        "code": "B"
    },
    {
        "code": "A"
    }
]'
```

Response:
```
{
    "transaction_id": "921a0ba9-36f4-42eb-bd23-bd62a9295923",
    "success": false,
    "data": [
        {
            "id": "51e504a0-d48f-4351-8686-1b9c40965965",
            "code": "A",
            "message": "Error 1062 (23000): Duplicate entry 'A' for key 'channels.code'"
        }
    ],
    "error": null
}
```

MySQL Database:
```
mysql> select * from channels;
+--------------------------------------+------+---------------------+---------------------+
| id                                   | code | created_at          | updated_at          |
+--------------------------------------+------+---------------------+---------------------+
| 1bfe1431-50e1-4038-b6a6-a03c40e41e4b | A    | 2023-08-29 15:00:55 | 2023-08-29 15:00:55 |
+--------------------------------------+------+---------------------+---------------------+
1 row in set (0.00 sec)
```

#### Benchmark

| Without DB transaction | With DB transaction | Summary |
|  ------------- | ------------- | ------------- |
| No rollback | Rollback  | With DB transaction is more reliable |