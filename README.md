# Docker container monitoring system as a test case for VK
## Структура проекта

```.
├── backend
│   ├── config
│   │   ├── config.go       # Загрузка конфигурации через Viper
│   │   └── config.yaml     # Конфигурационный файл (порт, БД, очередь)
│   ├── db
│   │   ├── connect.go      # Подключение к PostgreSQL
│   │   ├── Create.sql      # Скрипт создания таблиц
│   │   └── init.sql        # Скрипт инициализации БД
│   ├── Dockerfile          # Dockerfile для backend
│   ├── go.mod
│   ├── go.sum
│   ├── handlers
│   │   ├── IpToBDApi.go    # HTTP-обработчики для API (/api/ip)
│   │   └── IpToBDApi_test.go
│   ├── main.go             # Запуск API-сервера, метрики и Kafka-потребитель
│   └── queue
│       └── kafka.go        # Чтение сообщений из Kafka и запись в БД
├── db
│   └── init.sql            # Инициализация базы данных PostgreSQL
├── docker-compose.yml      # Описание всех сервисов (backend, frontend, Kafka, exporters, Prometheus)
├── frontend                # Каталог с фронтендом (React + Nginx)
│   ├── Dockerfile          # Dockerfile для сборки фронтенда
│   ├── default.conf        # Конфигурация Nginx для проксирования API-запросов
│   ├── nginx.conf          # Основная конфигурация Nginx
│   ├── package.json        # Конфигурация npm
│   ├── public              # Публичные файлы (index.html и др.)
│   ├── src                 # Исходный код React-приложения
│   │   ├── App.css
│   │   ├── App.test.js
│   │   ├── components
│   │   │   ├── App.js
│   │   │   └── IPStatusTable.js
│   │   ├── index.css
│   │   ├── index.js
│   │   ├── logo.svg
│   │   ├── nginx.conf
│   │   └── setupTests.js
│   └── README.md           # Документация по фронтенду (не собственного сочинения)
├── init_kafka.sh           # Скрипт инициализации Kafka (создание топика pingTopic)
├── jmx_exporter            # Директория с конфигурацией для JMX Exporter (jmx_prometheus_javaagent.jar, kafka.yml и др.)
├── kafka-logs              # Том для хранения данных Kafka
├── Makefile                # Набор команд для сборки, запуска, логов и чистки
└── prometheus.yml          # Конфигурация Prometheus для сбора метрик
```

**_Этот демонстративный проект для системы мониторинга контейнеров. Бэкенд отвечает за управление базой данных, обработку HTTP-запросов и взаимодействие с брокером сообщений Kafka. сервис Пингер пингует доступные контейнеры и передаёт данные во фронтэнд, которые становятся доступны на странице http://localhost Также реализован сбор и экспорт метрик Prometheus для мониторинга системы._**

### Предварительные требования

    Docker и Docker Compose установлены.
    Проверьте, что порты, указанные в docker-compose.yml, свободны на вашем хосте.

### Сервисы
- db - PostgreSQL база данных для хранения информации о контейнерах.
- Frontend – веб-интерфейс, написанный на React с использованием Bootstrap, который:
  - Отображает данные (таблицу IP-адресов с их статусом).
  - Доступен через веб-сервер Nginx, который также проксирует API-запросы к Backend.
- Kafka – брокер сообщений, работающий в режиме KRaft (без Zookeeper) с Bitnami‑образом, который:
  - Передаёт сообщения от сервиса pinger к Backend.
  - Позволяет другим сервисам собирать метрики.
- kafka-init - создание топиков
- backend - Основной API-сервер, обрабатывающий HTTP-запросы и взаимодействующий с базой данных.
- pinger - Сервис, который получает список IP-адресов контейнеров, пингует их и отправляет метрики в Kafka.
- prometheus - Система мониторинга для сбора метрик из backend, kafka-exporter и postgres-exporter.

### Запуск проекта

- **скопируйте репозиторий себе на машину**
```git clone ```

- **Запустите через Docker Compose**

```make up```

#### Остановка сервисов и очистка данных

```make down```
```make clean```

#### Линтинг кода

```make lint```

### Метрики

Бэкенд и связанные сервисы собирают метрики для мониторинга. Основные метрики доступны по следующим эндпоинтам:

#### *[Backend service](http://localhost:8080/metrics):* `http://localhost:8080/metrics`

#### *[Kafka Exporter](http://localhost:9308/metrics):* `http://localhost:9308/metrics`

#### *[Prometheus UI](http://localhost:9090):* `http://localhost:9090`

### Работающие контейнеры можно посмотреть на http://localhost после запуска сервисов в докер



