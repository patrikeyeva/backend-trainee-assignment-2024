## Инструкция по запуску

make run-all собирает и запускает два контейнера - для постгреса и для сервера, так же делает миграцию для БД

Для Postman
#### 1. Пример get_user_banner запроса 
```bash
GET localhost:8080/user_banner?tag_id=1&feature_id=1&use_last_revision=true

headers:
token: user_token
```

#### 2. Пример get_banner запроса 
```bash
GET localhost:8080/banner?feature_id=1&limit=5

headers:
token: admin_token
```

#### 3. Пример post_banner запроса 
```bash
POST localhost:8080/banner

headers:
token: admin_token
```
```json
json:
{
  "tag_ids": [1, 2, 3],
  "feature_id": 123,
  "content": {
    "title": "Example Title",
    "text": "This is an example text.",
    "url": "https://example.com"
  },
  "is_active": true
}
```

#### 4. Пример patch_banner запроса 
```bash
PATCH localhost:8080/banner/1

headers:
token: admin_token
```
```json
json:
{
  "tag_ids": [1, 2, 3],
  "feature_id": 123456,
  "content": {
    "title": "patched Title",
    "text": "This is an example text.",
    "url": "https://example.com"
  },
  "is_active": false
}
```

#### 5. Пример delete_banner запроса 
```bash
DELETE localhost:8080/banner/1

headers:
token: admin_token
```
