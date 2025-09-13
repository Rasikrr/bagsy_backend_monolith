### Bugsy Backend ###

Для запуска нужен .env файл с переменными окружения:
```
CONFIG_PATH=./config/твой_файл.yaml

POSTGRES_DSN=postgres://postgres:rasik1234@localhost:5432/test?sslmode=disable

REDIS_PORT=6379
REDIS_USER=
REDIS_HOST=localhost
REDIS_PASSWORD=
REDIS_DB=0

NATS_DSN=nats://localhost:4222
```