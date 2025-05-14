### Конфигурационный файл
Конфигурационный файл должен содержать следующие ключи:
- `load_balancer.port` - порт, по которому будет доступен балансировщик нагрузки
- `load_balancer.strategy` - стратегия балансировщика нагрузки. *Доступные значения:*
  - `round-robin`
- `load_balancer.backends` - список бэкенд серверов для балансировки нагрузки
- `load_balancer.backends[].url` - адрес бэкенд сервера для балансировки нагрузки
- `load_balancer.backends[].healthcheck.path` - path опроса доступности бэкенд сервера
- `load_balancer.backends[].healthcheck.interval` - интервал опроса доступности бэкенд сервера
- `rate_limiter.default_bucket` - настройки бакета по умолчанию 
- `rate_limiter.default_bucket.capacity` - количество токенов в бакете по умолчанию
- `rate_limiter.default_bucket.rps` - скорость пополнения токенов в бакете по умолчанию (поддерживает десятичные числа)
- `db` - настройки БД
- `db.dbms` - название СУБД
- `db.dsn` - DSN строка подключения к БД 

**Примечания**: 
- Пример конфигурационного файла доступен по пути `config/config.example.yaml`.
- В `docker-compose.yaml` указывается `CONFIG_PATH` (по умолчанию используется: `./config/config.example.yaml`)
- Документация к API доступна по `api/openapi.yaml`
- API-key указывается в query параметрах для проксируемых запросов. Пример: http://localhost:\<port>/\<path>/?api_key=<api_key>

### Сборка и запуск балансировщика нагрузки
```bash 
docker-compose up -d
make migrate_up
```
