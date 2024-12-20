# Song Library API

Приложение реализует CRUD операции с песнями.

## Запуск приложения

Перед запуском приложения необходимо создать `.env` файл в корне проекта.
Можно воспользоваться `.env.example`:

```
DB_PASSWORD=password_for_db
EXTERNAL_API_URL=api_url
APP_ENV=dev
```

Назначение полей:
- `DB_PASSWORD`: пароль к базе данных
- `EXTERNAL_API_URL`: необходимо указать url внешнего api, куда будет отправлять запрос на `/info`
- Можно указать `APP_ENV=prod` для запуска в режиме продакшена, при любом другом значении будет запуск в режиме разработки

Также следует изменить файл конфигурации `config/config.yml`, при необходимости.
В нём содержатся поля, отвечающие за адрес, на котором будет запущено api, подключение к базе и выбор реализации репозитория:

```yaml
http:
  host: "localhost"
  port: "8080"

db:
  username: "postgres"
  host: "localhost"
  port: "5432"
  dbname: "songlibrary"
  sslmode: "disable"

repo_implement:
  engine: "postgresql"
  sqldriver: "pgx/v5"
```

После выполнения настроек можно воспользоваться Makefile для запуска:

1) `make build` - создаёт исполняемый файл проекта `build/bin`
2) `make run` - запускает исполняемый файл

Или выполнить команду:
```shell
go run main.go
```

Если всё указано корректно и подключение к базе пройдёт успешно, то приложение начнёт работать на хосте и порте, указанных в `config.yml`.

Также при запуске приложения будет создана база данных в Postgres, если её нет.
Название возьмётся из `config.yml`. А после будут выполнены миграции.

Документация `swagger` будет доступна по адресу `/docs/index.html`, где можно найти полную документацию api.

Во время работы приложение создаст `logs/all.log`, куда будет писать логи (также дублируются в консоль).

## Используемые технологии
### База данных
В качестве базы данных для хранения песен используется `PostgreSQL`.
### Пакеты/Библиотеки
- `pgx/v5` подключение к базе данных
- `squirrel` создание sql-запросов
- `gin` маршрутизация запросов
- `logrus` логирование
- `swag` создание документации swagger
- `viper` чтение конфигов
- `godotenv` чтение окружения

## Дополнительно

При разработке старался учесть все принципы чистого кода, такие как SOLID, KISS, DRY.

Отделял реализацию чего-либо с помощью интерфейсов.

Старался сделать проект не зависимым от выбора сторонних технологий.
