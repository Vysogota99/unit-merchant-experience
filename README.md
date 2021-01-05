# unit-merchant-experience
<h2>Инструкция по запуску</h2>
<h3>Структура проекта</h3>
<a href="https://github.com/golang-standards/project-layout">https://github.com/golang-standards/project-layout</a>
<br>
<ul>
    <li>build - содержит все необходимое для запуска и работы сервер: .env, Dockerfile для golang, миграции и файлы базы данных. После сборки сервера здесь появится скомпилированный файл для запуска;</li>
    <li>cmd - содержит main package;</li>
    <li>deployments - содержит файл для docker-compose.yml. Контейнеры запускаются здесь;</li>
    <li>internal - содержит пакеты для работы сервера.
        <ul>
            <li>models - модели сущностей;</li>
            <li>server - конфиг, сервер, роутер и обработчики;</li>
            <li>store - включает в себя две реализации интерфейса для работы с б.д: mock для тестирования http и prod для работы сервера;</li>
        </ul>
    </li>
</ul>
<h3>Запуск контейнеров</h3>
<p>В директории ./deployments прописать:
    <br>
    <code>docker-compose up -d</code>
    <br>
    После этого соберутся и запустятся контейнеры для работы с go и postgres
</p>
<p>
 Посмотреть их названия можно командой:
    <br>
    <code>
        docker ps
    </code>
</p>

<h3>Создание пользователя для работы с базой данных </h3>
<ul>
    <li>Войти в контейнере c postgres в учетную запись postgres (пароль - "qwerty"):
    <br>
    <code>docker exec -it deployments_psql_1 /bin/bash</code>
    <br>
        <code>psql -U postgres -p 5432 -h store</code>
    <br>
    </li>
    <li>Создать нового пользователя:
    <br>
        <code>CREATE ROLE user1 WITH PASSWORD 'password' LOGIN CREATEDB;</code>
    <br>
    </li>
    <li>Создать базы данных:
    <br>
        <code>CREATE DATABASE user1;
        <br>
        CREATE DATABASE app
        </code>
    <br>
    </li>
</ul>

<h3>Запуск миграций </h3>
<ul>
    <li>
        Зайти в контейнер с golang:
        <br>
        <code>docker exec -it deployments_go_1 /bin/bash</code>
    </li>
    <li>В директории ./build запустить миграции командой:
    <br>
    <code>migrate -database ${POSTGRESQL_URL} -path ./migrations up</code>
</ul>


<h3>Запуск сервера</h3>
<ul>
    <li>
        В директории ./build необходимо создать .env файл и заполнить его по примеру .env_example(скопировать все из .env_example в .env)
        <br>
        <code>
            cp .env_example .env
        </code>
    </li>
    <li>
        В директории ./build необходимо прописать команду для сборки сервера:
        <br>
            <code>go build ../cmd/app/main.go</code>
        <br>
    </li>
    <li>
        Теперь его можно запустить командой:
        <br>
            <code>./main</code>
        <br>
    </li>
</ul>
<h3>API</h3>
<ul>
    </li>
    <li>POST http://127.0.0.1:3000/offer:
    <br>
    Добавление offer-а в базу данных.
    <br>
    Тело запроса включает в себя id - продавец, url - ссылка на файл с xlsx:
    <br>
        <code>
        {
            "id": 1,
            "url": "http://nginx:80/files/1.xlsx"
        }
        </code>
    <br>
    В ответ возвращает ошибку или id задания, при помощи которого можно узнать статус его выполнения и результат работы.
    <br>
    <pre>
    {
        "error": "",
        "result": {
            "ID задачи": 91
        }
    }
    </pre>
    <br>
    </li>
    <li>GET http://127.0.0.1:3000/status/:id:
    <br>
    Получение статуса/результата работы задачи
    <br>
    :id - номер задачи
    <br>
    В ответ возвращает ошибку, статус или краткую статистику.
    <br>
    <pre>
{
    "error": "",
    "result": {
        "status": "Готово",
        "статистика": {
            "количество созданных строк": 0,
            "количество обновленных строк": 5,
            "количество удаленных строк": 0,
            "количество строк с ошибками": 0
        }
    }
}
    </pre>
    <br>
    </li>
        </li>
    <li>GET http://127.0.0.1:3000/offer
    <br>
    Получение списка загруженных товаров.
    <br>
    Метод принимает на вход параметры GET запроса - id продавца, offer_id, подстрока названия товара (по тексту "теле" находились и "телефоны", и "телевизоры"). Ни один параметр не является обязательным, все указанные параметры применяются через логический оператор "AND"
    <br>
    <code>
    http://127.0.0.1:3000/offer?saler_id=1&offer=_X
    </code>
    <br>
    <pre>
{
    "error": "",
    "result": [
        {
            "OfferID": 1,
            "SalerID": 1,
            "Name": "iphone_X",
            "Price": 40000,
            "Quantity": 10
        },
        {
            "OfferID": 5,
            "SalerID": 1,
            "Name": "iphone_XS_MAX",
            "Price": 60000,
            "Quantity": 15
        },
        {
            "OfferID": 2,
            "SalerID": 1,
            "Name": "iphone_XR",
            "Price": 42000,
            "Quantity": 100
        }
    ]
}
    </pre>
    <br>
    </li>
</ul>

<h3>Нагрузочное тестирование</h3>
<p>
Программа для проведения нагрузочного тестирования находится в папке /test/loadtester. Для запуска необходимо войти в эту дирректорию и прописать ./main.
</p>
<h4>Результаты тестирования для 1 горутины</h4>
<pre>

</pre>

