# API по созданию сокращённых ссылок

Сервис  принимает запросы по http:
1. Метод Post, который сохраняет оригинальный URL в базе и возвращать сокращённый. Метод ожидает тело в формате json 
Пример запроса:
```
{
  "url": "https://example.com"
}
```
Пример ответа:
```
{
    "result": "http://localhost:8080/o8DOU7s683"
}
```
Адрес сервиса:
```
http://localhost:8080/
```
2. Метод Get, который принимает сокращённый URL как PATH параметр и редиректит на оригинальный.

```
http://localhost:8080/{id}
```

Работа сервиса так же доступна через GRPC


## Хранилище
По умолчанию приложение использует PostgreSQL. Чтобы сменить тип хранилища нужно отредактировать флаг d из Dockerfile

Вариант для PostgreSQL:
```
CMD ["-d="]
```

Вариант для in-memory решения:
```
CMD ["-d"]
```