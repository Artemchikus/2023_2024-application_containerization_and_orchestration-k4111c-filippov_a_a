# 2023_2024-application_containerization_and_orchestration-k4111c-filippov_a_a

Labs for "Application containerization and orchestration" course

## Документация проекта

### Описание

Сервис представляет собой серверное приложение, предоставляющее информацию об сайте, находящемся в описании группы VK. Сам сервис реализует функционал доступа и модификации инофрмации в БД PostgreSQL, а также скраппинг информации со страницы паблика VK (так как доступ к API VK ограничен только для сайтов с купленными доменами). Также в проекте реализован CronJob срабатывающий каждый день в 8 утра и обновляющий информацию о сайте из описания паблика, если URL был изменен.

### API

Подробнуй документацию по API можно посмотреть по URI `/api/v1/docs/`, или в [swagger.json](cmd/app/api//docs/swagger.json) файле.

### Билд и развертывание

Для работы сервиса необходим GO 1.21.6 версии, а также docker, если планируется локальное развертывание.

Помимо [основного образа](Dockerfile) приложения для сервиса также неоходимо билдить образ для прогонки [миграций](storage/Dockerfile) и образ для [CronJob](cmd/cron_lobs/update_urls/Dockerfile). Билд и развертывание сервиса можно проводить с помощью Makefile, в нем реализованы билд, запуск приложения, создание и пуш образа (только введите свой профиль в dockerHub), а также локальное развертывание с помощью docker-compose и развертывние в кластере Kubernetes с помощью helm.

### Конфигурация

Примеры кофигурации можно найти в папке [examples](examples), миграции же в папке [storage/migrations](storage/migrations). Пример Helm конфигурации можно найти [тут](.helm/values.yaml).
