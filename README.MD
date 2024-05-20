## Запуск
```shell
  go mod download -
  docker-compose up -d

  API_TOKEN="token" go run ./
```

#### Для проверки можно использовать curl

```shell
curl --location 'http://127.0.0.1:3000/api/v1/fact' \
--header 'Content-Type: application/json' \
--data '[
    {
        "period_start": "2024-05-01",
        "period_end": "2024-05-31",
        "period_key": "month",
        "indicator_to_mo_id": 227373,
        "value": 852,
        "fact_time": "2024-05-31",
        "is_plan": false,
        "auth_user_id": 40,
        "comment": "buffer Last_name"
    },
    {
        "period_start": "2024-05-01",
        "period_end": "2024-05-31",
        "period_key": "month",
        "indicator_to_mo_id": 227373,
        "value": 321,
        "fact_time": "2024-05-31",
        "is_plan": false,
        "auth_user_id": 40,
        "comment": "buffer Last_name"
    }
]'
```

Есть возможность добавлять множество фактов одновременно