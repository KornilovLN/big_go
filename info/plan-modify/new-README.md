# README.md - проект big_go (Версия 2.0)

## Обзор проекта
Проект представляет собой систему сбора и отображения данных с датчиков (температура, давление, влажность) для разных пользователей. Датчики организованы в группы на определенных постах. Посты расположены по определенным адресам.

## Архитектура системы

### Компоненты системы:
* **Генератор данных** - создает данные о температуре, давлении и влажности с разной периодичностью
* **Коллектор** - получает данные через RabbitMQ, сохраняет в Redis и отправляет уведомления
* **Пользовательские приложения** (User1 и User2) - получают уведомления и данные из Redis
* **Инфраструктура** - PostgreSQL, Redis, RabbitMQ в отдельных контейнерах

### Поток данных:
1. Генератор создает данные и отправляет их в очередь RabbitMQ
2. Коллектор получает данные из RabbitMQ
3. Коллектор сохраняет данные в Redis с TTL (время жизни)
4. Коллектор отправляет уведомление через Redis Pub/Sub
5. Пользовательские приложения получают уведомление и запрашивают данные из Redis
6. Пользовательские приложения отображают данные в веб-интерфейсе

## Состав проектируемой системы и как происходит работа

* **Инфраструктура коммуникаций проекта опирается на сеть big_go_network**
  * Общая сеть сервисов проекта - big_go_network создается в docker-compose.yml  

* **Сервис: rabbitmq**
  * Создается как сервис (big_go_rabbitmq) со всеми настройками в docker-compose.yml
  * Работает в сети big_go_network 
  * Используется для асинхронной передачи данных от генератора к коллектору

* **Сервис: redis** 
  * Создается как сервис (big_go_redis) со всеми настройками в docker-compose.yml 
  * Работает в сети big_go_network 
  * Используется для:
    * Хранения данных датчиков с TTL
    * Отправки уведомлений через механизм Pub/Sub
    * Обеспечения отказоустойчивости и масштабируемости системы

* **Сервис: postgresql** 
  * Создается как сервис (big_go_postgres) со всеми настройками в docker-compose.yml  
  * Работает в сети big_go_network
  * Подготовлен для будущего использования (долгосрочное хранение данных)

* **Сервис: generator**
  * Создается как сервис (big_go_generator) со всеми настройками в docker-compose.yml 
  * Работает в сети big_go_network 
  * Подписывается на RabbitMQ как публикатор данных (producer)
  * Создает канал с rabbitmq
  * Cоздает очередь с именем "sensor_data"
  * Подключение к RabbitMQ происходит с использованием данных из config_rabbitmq.json

* **Сервис: collector**
  * Создается как сервис (big_go_collector) со всеми настройками в docker-compose.yml
  * Работает в сети big_go_network 
  * Подписывается на RabbitMQ как получатель данных (consumer)
  * Создает канал с rabbitmq
  * Использует очередь с именем "sensor_data"
  * Сохраняет полученные данные в Redis с TTL
  * Отправляет уведомления через Redis Pub/Sub
  * Логика работы коллектора:
    * Получение сообщения из RabbitMQ
    * Десериализация данных
    * Сохранение данных в Redis
    * Отправка уведомления через Redis Pub/Sub

* **Сервис: Приложение user1**
  * Создается как сервис (big_go_user1) со всеми настройками в docker-compose.yml
  * Работает в сети big_go_network 
  * Подписывается на канал уведомлений Redis Pub/Sub
  * При получении уведомления запрашивает данные из Redis
  * Отображает данные в веб-интерфейсе
  * Предоставляет API для получения данных в формате JSON
  * User1 сервис запущен на порту 8082 и доступен по адресу http://localhost:8082

* **Сервис: Приложение user2**
  * Аналогично приложению user1, но для другого пользователя
  * Создается как сервис (big_go_user2) со всеми настройками в docker-compose.yml
  * Работает в сети big_go_network 
  * User2 сервис запущен на порту 8083 и доступен по адресу http://localhost:8083

## Структура проекта
```
big_go/
├── cmd/
│   ├── collector/
│   │   └── main.go
│   ├── generator/
│   │   └── main.go
│   ├── user1/
│   │   └── main.go
│   └── user2/
│       └── main.go
├── config/
│   ├── config.go
│   ├── opentsdb.go
│   ├── postgresql.go
│   ├── rabbitmq.go
│   └── redis.go
├── docker/
│   ├── collector/
│   │   └── Dockerfile
│   ├── generator/
│   │   └── Dockerfile
│   ├── user1/
│   │   └── Dockerfile
│   └── user2/
│       └── Dockerfile
├── internal/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   │   └── redis/
│   │       └── redis.go
│   ├── routes/
│   ├── services/
│   │   ├── generator/
│   │   │   └── generator.go
│   │   ├── collector/
│   │   │   └── collector.go
│   │   └── user/
│   │       └── user.go
│   └── templates/
│       └── index.html
├── config_go.json
├── config_postgresql.json
├── config_rabbitmq.json
├── config_redis.json
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Порядок запуска проекта
* **Клонирование репозитория и переход в директорию проекта:**
```bash
git clone <ваш-репозиторий>
cd big_go
```
* **Создание необходимых конфигурационных файлов:**
  * config_go.json - основная конфигурация приложения
  * config_postgresql.json - конфигурация PostgreSQL
  * config_redis.json - конфигурация Redis
  * config_rabbitmq.json - конфигурация RabbitMQ
* **Запуск проекта с помощью Docker Compose:**
```bash
docker-compose up -d
```

### Эти команды запустят:
* PostgreSQL
* Redis
* RabbitMQ
* Генератор данных
* Коллектор
* User1
* User2

### Проверка работы сервисов:
* **После запуска вы можете проверить работу сервисов:**
  * Веб-интерфейс RabbitMQ: http://localhost:15672 (логин: guest, пароль: guest)
  * User1 Dashboard: http://localhost:8082
  * User2 Dashboard: http://localhost:8083
* **Мониторинг логов:**
```bash
docker-compose logs -f
```
* **Для просмотра логов конкретного сервиса:**
```bash
docker-compose logs -f generator
docker-compose logs -f collector
docker-compose logs -f user1
docker-compose logs -f user2
```
* **Остановка проекта:**
```bash
docker-compose down
```
* **Для полной очистки (включая удаление томов):**
```bash
docker-compose down -v
```

## Преимущества новой архитектуры

1. **Отказоустойчивость**
   * Если пользовательские приложения временно недоступны, данные сохраняются в Redis
   * Данные могут быть получены позже, когда приложения снова станут доступны

2. **Масштабируемость**
   * Можно легко добавлять новых пользователей без изменения логики коллектора
   * Redis обеспечивает высокую производительность даже при большом количестве клиентов

3. **Разделение ответственности**
   * Коллектор отвечает только за сбор и сохранение данных
   * Пользовательские приложения отвечают за получение и отображение данных

4. **Производительность**
   * Redis обеспечивает высокую скорость чтения/записи данных
   * Механизм Pub/Sub позволяет эффективно уведомлять клиентов о новых данных

5. **Удобство использования**
   * Улучшенный пользовательский интерфейс
   * Возможность ручного обновления данных
   * API для интеграции с другими системами


## Шаг 10: Создадим скрипт для запуска проекта

```bash:start.sh
#!/bin/bash

# Проверка наличия конфигурационных файлов
if [ ! -f "config_rabbitmq.json" ]; then
    echo "Создание config_rabbitmq.json..."
    cat > config_rabbitmq.json << EOF
{
  "host": "rabbitmq",
  "port": 5672,
  "user": "guest",
  "password": "guest",
  "vhost": "/"
}
EOF
fi

if [ ! -f "config_redis.json" ]; then
    echo "Создание config_redis.json..."
    cat > config_redis.json << EOF
{
  "host": "redis",
  "port": 6379,
  "password": "",
  "db": 0
}
EOF
fi

if [ ! -f "config_postgresql.json" ]; then
    echo "Создание config_postgresql.json..."
    cat > config_postgresql.json << EOF
{
  "host": "postgres",
  "port": 5432,
  "user": "postgres",
  "password": "postgres",
  "dbname": "big_go",
  "sslmode": "disable"
}
EOF
fi

if [ ! -f "config_go.json" ]; then
    echo "Создание config_go.json..."
    cat > config_go.json << EOF
{
  "log_level": "info",
  "environment": "development"
}
EOF
fi

# Запуск проекта с помощью Docker Compose
echo "Запуск проекта..."
docker-compose up -d

# Проверка статуса контейнеров
echo "Проверка статуса контейнеров..."
docker-compose ps

echo "Проект успешно запущен!"
echo "Веб-интерфейс RabbitMQ: http://localhost:15672 (логин: guest, пароль: guest)"
echo "User1 Dashboard: http://localhost:8082"
echo "User2 Dashboard: http://localhost:8083"
```

Сделать скрипт исполняемым:

```bash
chmod +x start.sh
```

### Заключение
* Добавили Redis в рабочий процесс. Cистема стала более отказоустойчивой, масштабируемой и производительной. Коллектор сохраняет данные в Redis и отправляет уведомления через Redis Pub/Sub, а пользовательские приложения получают уведомления и запрашивают данные из Redis.
* Эта архитектура соответствует рекомендациям из документа big_go-var2.md и обеспечивает все указанные преимущества: отказоустойчивость, масштабируемость и производительность.