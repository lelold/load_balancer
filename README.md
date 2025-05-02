
# load_balancer

Балансировщик нагрузки

  

## Запуск

1) Клонируем:

> git clone https://github.com/lelold/load_balancer.git

2) Меняем директорию:

> cd load_balancer

3) Собираем и запускаем контейнер:

> docker-compose up --build

4) Для запуска тестов:
> go test ./... -bench=. -race
  

Сервис доступен по адресу: http://localhost:8080

  
  

## Описание проекта

Балансировщик нагрузки, распределяющий HTTP-запросы по пулу бэкенд-серверов.

Сервера проходят health-чеки каждые 10 секунд.

В конфигурационном файле assets/config.json можно настроить порт (port); алгоритм балансировщика (strategy), по дефолту round-robin, запушен random, есть ещё least_connections; список серверов (backends); дефолтную вместимость токенов для клиентов (capacity); дефолтную скорость набора токенов (refill_rate).
Конфигурационный файл независим от кода.

Все пакеты задокументированы.

В assets/clients.json хранятся клиенты.

Весь проект написан с использованием стандартной библиотеки Go.

  

## Структура

assets/ - json-файлы (clients - хранение пользователей с их конфигурацией, config - конфиг проекта)

backend/ - запускает 2 бэкенд сервера с endpoint`ом health/ (в докер компоузе они тоже описаны)

cmd/ - инициализация и запуск проекта

integration/ - интеграционный тест

internal/app/ - сервис для инкапсуляции лимитера и балансировщика

internal/config/ - импортирует конфиг файл в проект

internal/domain/lb/ - логика балансировщика

internal/domain/ratelimiter/ - логика лимитера

internal/handlers/ - хендлеры и роуты

internal/logger/ - логгер

testutils/ - вспомогательные функции для интеграционных тестов

По структуре старался придерживаться чистой архитектуры и DDD.

## Реализация лимитера и балансировщика

**Балансировщик** может работать по 1 из 3 алгоритмов:

1. round-robin

2. least_connections

3. random

  

Для смены алгоритма надо в assets/config.json вписать один из этих трех вариантов

**Лимитер** работает по алгоритму Token Bucket, дефолтные значения для лимита токенов всех пользователей и кол-во обновляемых токенов можно так же изменить в конфиг файле

## Логгер

В проекте реализован логгер на базе log из стандартного пакета

## Ручки

Реализован endpoint /clients (GET, POST), а так же /clients/{client_id} (DELETE).

При неверных запросах возвращаются структурированные json-ответы.

## Сервера

Для проверки работы в директории backend/ реализован запуск двух бэкенд серверов на портах 9001 и 9002, в докер-компоуз файле все сделано для их запуска вместе с балансировщиком.

## Скриншоты

**Запуск сервера**

![image](https://github.com/user-attachments/assets/0ce328a9-18c1-450c-b448-eb496ca4e1ab)

**Health-чеки**

![image](https://github.com/user-attachments/assets/4f03390f-c841-4601-a4be-e1b41a283426)

**Пример нагрузки с недостатком токенов через ab**

![image](https://github.com/user-attachments/assets/4bb9aaec-db05-4c1c-84c5-67dc7e672274)

![image](https://github.com/user-attachments/assets/b0f752db-c1c1-4932-a255-fcd5421839a0)


**Работа лимитера**

![image](https://github.com/user-attachments/assets/92d4da9b-b6b9-4a8e-ae6e-5dd9e74c4087)

**Пример нагрузки с достаточным количеством токенов через ab**

![image](https://github.com/user-attachments/assets/6dd9e24a-efd2-484d-a569-7ab15d390c0c)

**Ручки и ответы**

![image](https://github.com/user-attachments/assets/be63bf7f-2a4a-4d51-8bda-0c226aae44d7)

![image](https://github.com/user-attachments/assets/b6879b9a-cd4c-4afa-8b41-3db818f54f39)

![image](https://github.com/user-attachments/assets/7b5d2937-76c5-431f-9529-d31dcacfd8dd)

![image](https://github.com/user-attachments/assets/443afbec-b137-40de-80ef-fff565beab7e)

![image](https://github.com/user-attachments/assets/d91d4c37-f858-4ff3-bcf0-5854bf442ec7)

![image](https://github.com/user-attachments/assets/0aa1d19e-4d7b-49cc-a23d-bd5b90049346)

![image](https://github.com/user-attachments/assets/ee1aeeaa-1c5c-4237-b325-a41ced452227)

![image](https://github.com/user-attachments/assets/d0b1fb18-4eaf-4044-9de3-7611f9b98101)


