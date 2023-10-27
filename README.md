# Запуск

~~~zsh
git clone https://github.com/realPointer/EnrichInfo
cd EnrichInfo
make compose-up
~~~

Механизм миграции БД происходит автоматически при запуске 

# Swagger

После запуска приложения доступна Swagger-документация по адресу [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

![swagger](https://github.com/realPointer/EnrichInfo/assets/50529632/2b717bd2-da42-4fbd-8575-90480e533ef5)

# Запросы

### Добавление персоны

patronymic является необязательным

~~~zsh
curl -X POST "http://localhost:8080/v1/people" \
  -H 'Content-Type: application/json' \
  -d '{
        "name": "name",
        "surname": "surname",
        "patronymic": "patronymic"
    }'
~~~

---

### Изменение персоны

При изменении имени перезапишутся дополнительные показатели. Возможно изменить name, surname, patronymic

~~~zsh
curl -X PUT "http://localhost:8080/v1/people/{id}" \
  -H 'Content-Type: application/json' \
  -d '{
        "name": "name",
        "surname": "surname"
    }'
~~~

---

### Удаление персоны

~~~zsh
curl -X DELETE "http://localhost:8080/v1/people/{id}"
~~~

---

### Получение данных

Фильтр по name, surname, patronymic, age, gender, nationality. page - номер страницы, perPage - количество записей на странице

Комбинирование параметров происходит через &. Например: name=Andrew&surname=Forest

~~~zsh
curl "http://localhost:8080/v1/people?"
~~~


## Что явно стоило бы сделать тут
- Невозможность записи дубликатов
- Идемпотентность

Надеюсь, что не сильно задержался из-за поломки ноута 😅 Буду рад фидбеку!
