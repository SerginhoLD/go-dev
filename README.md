## Запуск

* `docker-compose up -d --build`
* `docker exec -i -t go-dev-app sh`
* `make migrate`

## Создание миграций

`make create-migration` or `make create-migration name=name`

## Api

* `http://127.0.0.90/` Список продуктов
* `http://127.0.0.90/product/1` Один продукт
* `http://127.0.0.90/metrics` Метрики

## Тесты

`make test`

or

`make coverage-html`