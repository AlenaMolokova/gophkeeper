# Настройка пула соединений базы данных

GophKeeper использует настраиваемый пул соединений PostgreSQL для оптимизации производительности и управления ресурсами.

> **Примечание**: Эта документация описывает настройку пула соединений для сервера GophKeeper.
> Для получения дополнительной информации о архитектуре см. [README.md](../README.md).

## Переменные окружения

Следующие переменные окружения позволяют настроить поведение пула соединений:

### `DB_MAX_CONNS`
Максимальное количество соединений в пуле.
- **По умолчанию**: 20
- **Пример**: `DB_MAX_CONNS=50`

### `DB_MIN_CONNS`
Минимальное количество соединений в пуле.
- **По умолчанию**: 5
- **Пример**: `DB_MIN_CONNS=10`

### `DB_MAX_CONN_LIFETIME`
Максимальное время жизни соединения.
- **По умолчанию**: 1 час
- **Пример**: `DB_MAX_CONN_LIFETIME=2h`

### `DB_MAX_CONN_IDLE_TIME`
Максимальное время простоя соединения.
- **По умолчанию**: 30 минут
- **Пример**: `DB_MAX_CONN_IDLE_TIME=15m`

### `DB_HEALTH_CHECK_PERIOD`
Период проверки здоровья соединений.
- **По умолчанию**: 1 минута
- **Пример**: `DB_HEALTH_CHECK_PERIOD=30s`

## Рекомендации по настройке

### Для разработки
```bash
export DB_MAX_CONNS=10
export DB_MIN_CONNS=2
export DB_MAX_CONN_LIFETIME=30m
export DB_MAX_CONN_IDLE_TIME=10m
export DB_HEALTH_CHECK_PERIOD=1m
```

### Для продакшена
```bash
export DB_MAX_CONNS=50
export DB_MIN_CONNS=10
export DB_MAX_CONN_LIFETIME=1h
export DB_MAX_CONN_IDLE_TIME=30m
export DB_HEALTH_CHECK_PERIOD=30s
```

### Для высоконагруженных систем
```bash
export DB_MAX_CONNS=100
export DB_MIN_CONNS=20
export DB_MAX_CONN_LIFETIME=2h
export DB_MAX_CONN_IDLE_TIME=1h
export DB_HEALTH_CHECK_PERIOD=15s
```

## Мониторинг

Сервер автоматически логирует статистику пула каждые 5 минут:

```
Pool stats - Total: 25, Idle: 15, Acquired: 10, Constructing: 0
```

### Интерпретация статистики

- **Total**: Общее количество соединений в пуле
- **Idle**: Количество свободных соединений
- **Acquired**: Количество активных соединений
- **Constructing**: Количество соединений в процессе создания

## Примечания

1. **MaxConns** не должно превышать лимиты PostgreSQL (`max_connections`)
2. **MinConns** должно быть достаточно для базовой нагрузки
3. **MaxConnLifetime** помогает избежать проблем с устаревшими соединениями
4. **HealthCheckPeriod** влияет на скорость обнаружения проблем с соединениями

## Связанные файлы

- `internal/server/storage/storage.go` - основная реализация пула соединений
- `cmd/server/main.go` - настройка и мониторинг пула
- `internal/server/storage/storage_mock.go` - моки для тестирования
- `start-server.ps1` - скрипт запуска с настройками пула
- `docker-compose.yml` - пример конфигурации для Docker 