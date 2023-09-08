# Hexagonal Architecture

![Alt text](hexagonal.png?raw=true "Hexagonal Architecture")

## Objectives

Hexagonal architecture is a general-purpose architectural style that aims to create decoupled software.

## Summary

The Hexagonal Architecture, also referred to as Ports and Adapters, is an architectural pattern that allows input by users or external systems to arrive into the Application at a Port via an Adapter, and allows output to be sent out from the Application through a Port to an Adapter. This creates an abstraction layer that protects the core of an application and isolates it from external — and somehow irrelevant — tools and technologies.

## Scenario

| No | Scenario | Goals |
| ------------- | ------------- | ------------- |
| 1  | Create 2 domains with different database  | Decoupling business needs |
| 2  | Change database in 1 domain  | Change the concrete implementations without having to touch anything inside the business rules |

## Success Criteria

### 1. Create 2 domains with different database

#### Upsert channel in sales channel's domain (database: MySQL)

<span id="before"></span>
Environment:
```
SALES_CHANNEL_MAIN=mysql
```

Request:
```
curl --location 'localhost:8000/api/channel/upsert' \
--header 'Content-Type: application/json' \
--data '[
    {
        "code": "TKPD"
    }
]'
```

Response:
```
{
    "transaction_id": "ebf5b326-8f97-4309-92a7-1601a0b86f09",
    "success": true,
    "data": null,
    "error": null
}
```

MySQL Database:
```
mysql> select * from channels;
+--------------------------------------+------+---------------------+---------------------+
| id                                   | code | created_at          | updated_at          |
+--------------------------------------+------+---------------------+---------------------+
| 4d0481d5-178f-4edc-8f14-59fc210a7995 | TKPD | 2023-08-23 15:34:35 | 2023-08-23 15:34:35 |
+--------------------------------------+------+---------------------+---------------------+
1 row in set (0.01 sec)
```

#### Upsert location in inventory's domain (database: Postgres)

Environment:
```
INVENTORY_MAIN=postgres
```

Request:
```
curl --location 'localhost:8000/api/location/upsert' \
--header 'Content-Type: application/json' \
--data '[
    {
        "code": "LGK"
    }
]'
```

Response:
```
{
    "transaction_id": "703a15a7-0659-46c5-9c3c-22139ca3f744",
    "success": true,
    "data": null,
    "error": null
}
```

Postgres Database:
```
poc=# select * from locations;
                  id                  | code |         created_at         |         updated_at         
--------------------------------------+------+----------------------------+----------------------------
 73f48821-295f-4cd1-bc10-d77857b13fa4 | LGK  | 2023-08-23 15:29:12.300906 | 2023-08-23 15:29:12.300913
(1 row)
```

### 2. Change database in 1 domain

#### Change sales channel's domain database to Postgres

Before [here](#before)

Environment:
```
SALES_CHANNEL_MAIN=postgres
```

Request:
```
curl --location 'localhost:8000/api/channel/upsert' \
--header 'Content-Type: application/json' \
--data '[
    {
        "code": "SHPE"
    }
]'
```

Response:
```
{
    "transaction_id": "811f37ba-364a-4a17-91f9-2969b2b7d4a3",
    "success": true,
    "data": null,
    "error": null
}
```

Postgres Database:
```
poc=# select * from channels;
                  id                  | code |        created_at        |         updated_at         
--------------------------------------+------+--------------------------+----------------------------
 38014949-5f40-429d-bed8-559240973c07 | SHPE | 2023-08-23 15:50:46.3491 | 2023-08-23 15:50:46.349105
(1 row)
```