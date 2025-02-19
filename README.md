## Запуск

* `docker-compose up -d --build`
* `docker exec -i -t go-dev-app sh`
* `make migrate`

## Создание миграций

`make create-migration` or `make create-migration name=name`

## Api

* `http://127.0.0.90/` Список продуктов
* `http://127.0.0.90/products/1` Один продукт
* `POST http://127.0.0.90/products` Создать продукт `{"name": "Test", "price": 12.35}`
* `http://127.0.0.90/metrics` Метрики

## Тесты

`make test`

or

`make coverage-html`

## Слои

`./cmd/web` -> `./internal/infrastructure` -> `./internal/domain`

* `./cmd/web` Ввод/Вывод
* `./internal/infrastructure` Реализация интерфейсов и вспомогательные компоненты
* `./internal/domain` Доменная область (Entities + UseCases)