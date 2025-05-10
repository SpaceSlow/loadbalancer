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
- `rate_limiter.default_bucket.rps` - скорость пополнения токенов в бакете по умолчанию

**Примечание**: пример конфигурационного файла доступен по пути `config/config.example.yaml`

### Сборка и запуск
```bash 
go build -o loadbalancer cmd/loadbalancer/main.go
./loadbalancer -config <путь до файла конфигурации>
```
