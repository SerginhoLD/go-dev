# О чём

Пример парсинга объявлений.

Очередь на KeyDB Stream, база на Manticore. Dependency Injection через Wire.

## Запуск

* `docker compose up -d` (create `compose.override.yaml`)
* `docker exec -i -t vep-app sh`
* `make migrate`

## Routes

* `http://127.0.0.13/` Список
* `http://127.0.0.13/metrics` Метрики

## Тесты

`make test` or `make coverage-html`

## Слои

`./internal/app` -> `./internal/infrastructure` -> `./internal/domain`

* `./internal/app/server(|scheduler|consumer)` Ввод/Вывод
* `./internal/infrastructure` Реализация интерфейсов и вспомогательные компоненты
* `./internal/domain` Доменная область (Entities + UseCases)

```mermaid
flowchart LR
    App ==> Domain
    subgraph App
        Controller_1 -- Not allowed x--x Controller_2
    end
    subgraph Domain
        UseCases ==> Entities
    end
    subgraph UseCases
        UseCase_1 -- Not allowed x--x UseCase_2
    end
    Controller_1 ---> UseCase_1
    Controller_2 ---> UseCase_2
    subgraph UseCase_1
        Command_1 --> Handler_1 --> ViewModel_1
    end
    subgraph UseCase_2
        Command_2 --> Handler_2 --> ViewModel_2
    end
    subgraph Entities
        Entity
        Message
        RepositoryInterface
        MsgBusInterface
    end
    subgraph Infrastructure
        SqlRepository --> RepositoryInterface
        HttpClient --> RepositoryInterface
        Kafka --> MsgBusInterface
        JsonEncoder -- Serialize response --> Controller_1
        JsonEncoder -- Serialize response --> Controller_2
    end
```
