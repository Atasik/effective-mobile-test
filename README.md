# Тех. задание отбора на стажировку effective mobile
## Для работы приложения нужно задать .env файл:
```
POSTGRES_PASSWORD=qwerty
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres

HTTP_HOST=localhost

REDIS_ADDRESS=redis:6379
REDIS_PASSWORD=qwerty

ZOOKEEPER_CLIENT_PORT=2181
ZOOKEEPER_TICK_TIME=2000

KAFKA_BROKER_ID=1
KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181
KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092,PLAINTEXT_HOST://localhost:29092
KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT
KAFKA_AUTO_CREATE_TOPICS_ENABLE=true
KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1
KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=100

KAFKA_CLUSTERS_0_NAME=local
KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS=kafka:9092
DYNAMIC_CONFIG_ENABLED=true
```
### Запуск приложения:
```
make run
```
### Миграции:
```
make migrate
```
### Запуск kafka-ui:
```
make kafka-ui
```
