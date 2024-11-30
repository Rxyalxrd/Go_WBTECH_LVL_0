# Демонстрационный Сервис Обработки Заказов

Сервис обработки заказов - это демонстрационный сервис для управления и отображения данных о заказах. Сервис подключается к базе данных PostgreSQL, подписывается на Kafka для получения данных о заказах, кэширует данные в памяти и предоставляет HTTP API для получения информации о заказах по их идентификаторам.

---

## Что сделано

- Интеграция с PostgreSQL: Хранение данных о заказах в базе данных PostgreSQL.
- Apache Kafka: Подписка на брокер сообщений Kafka для получения данных о заказах.
- Кэширование в памяти: Кэширование данных о заказах в памяти для быстрого доступа.
- HTTP API: Предоставление API для получения информации о заказах по их идентификаторам.
- Устойчивость: Восстановление кэша из базы данных в случае перезапуска сервиса.
- Docker-контейнеризация: Удобство развертывания сервисов для тестирования через docker-compose.
---

## Стэк

- Golang
- Kafka
- Zookeper
- PostgreSQL
- Makefile
- Docker
- Docker-compose
- Postman
---

## Примеры запросов
- POST-запрос

  ```bash
  {
     "order_uid": "b563feg7bdg3286ttjr67siej",
     "track_number": "WBILMTESTTRACK",
     "entry": "WBIL",
     "delivery": {
        "name": "Test Testov",
        "phone": "+9720000000",
        "zip": "2639809",
        "city": "Kiryat Mozkin",
        "address": "Ploshad Mira 15",
        "region": "Kraiot",
        "email": "test@gmail.com"
     },
     "payment": {
        "transaction": "b563feb7b2b84b6test",
        "request_id": "",
        "currency": "USD",
        "provider": "wbpay",
        "amount": 1817,
        "payment_dt": 1637907727,
        "bank": "alpha",
        "delivery_cost": 1500,
        "goods_total": 317,
        "custom_fee": 0
     },
     "items": [
        {
           "chrt_id": 9934930,
           "track_number": "WBILMTESTTRACK",
           "price": 453,
           "rid": "ab4219087a764ae0btest",
           "name": "Mascaras",
           "sale": 30,
           "size": "0",
           "total_price": 317,
           "nm_id": 2389212,
           "brand": "Vivienne Sabo",
           "status": 202
        }
     ],
     "locale": "en",
     "internal_signature": "",
     "customer_id": "test",
     "delivery_service": "meest",
     "shardkey": "9",
     "sm_id": 99,
     "date_created": "2021-11-26T06:22:19Z",
     "oof_shard": "1"
  }
  ```

  Ответ
  ```bash
  {"message": "Order processed and saved successfully"}
  ```

- GET-запрос
  ```bash
  http://127.0.0.1:8080/orders/b563feg7b25sdg3284b6ttjrfhsiej
  ```

  Ответ
  ```bash
    "order": {
        "id": 1,
        "order_uid": "b563feg7b2b84b6test",
        "track_number": "WBILMTESTTRACK",
        "entry": "WBIL",
        "locale": "en",
        "internal_signature": "",
        "customer_id": "test",
        "delivery_service": "meest",
        "shardkey": "9",
        "sm_id": 99,
        "date_created": "2021-11-26T06:22:19Z",
        "oof_shard": "1",
        "delivery": {
            "ID": 1,
            "OrderID": 0,
            "Name": "Test Testov",
            "Phone": "+9720000000",
            "Zip": "2639809",
            "City": "Kiryat Mozkin",
            "Address": "Ploshad Mira 15",
            "Region": "Kraiot",
            "Email": "test@gmail.com"
        },
        "payment": {
            "ID": 1,
            "OrderID": 0,
            "Transaction": "b563feb7b2b84b6test",
            "RequestID": "",
            "Currency": "USD",
            "Provider": "wbpay",
            "Amount": 1817,
            "PaymentDT": 0,
            "Bank": "alpha",
            "DeliveryCost": 0,
            "GoodsTotal": 0,
            "CustomFee": 0
        },
        "items": [
            {
                "ID": 1,
                "OrderID": 1,
                "ChrtID": 0,
                "TrackNumber": "",
                "Price": 453,
                "Rid": "ab4219087a764ae0btest",
                "Name": "Mascaras",
                "Sale": 30,
                "Size": "0",
                "TotalPrice": 0,
                "NmID": 0,
                "Brand": "Vivienne Sabo",
                "Status": 202
            },
            {
                "ID": 2,
                "OrderID": 1,
                "ChrtID": 0,
                "TrackNumber": "",
                "Price": 453,
                "Rid": "ab4219087a764ae0btest",
                "Name": "Mascaras",
                "Sale": 30,
                "Size": "0",
                "TotalPrice": 0,
                "NmID": 0,
                "Brand": "Vivienne Sabo",
                "Status": 202
            }
        ]
    }
  ```
