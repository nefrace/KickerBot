# KickerBot

Бот, генерирующий капчу в виде картинки с набором логотипов игровых движков. Проверяемому необходимо выбрать правильный и написать ответ цифрой. Если ответ неверный: бот банит пользователя.

## Как запустить:
Создать два файла: `bot.env` и `mongo.env` для настройки базы и самого бота.

**bot.env**
```
TOKEN=<токен бота Telegram>
DB_USER=<логин для базы>
DB_PASSWORD=<пароль для базы>
MONGO_URI=mongodb://mongo:27017
```

**mongo.env**
```
MONGO_INITDB_ROOT_USERNAME=<логин для базы>
MONGO_INITDB_ROOT_PASSWORD=<пароль для базы>

ME_CONFIG_MONGODB_ADMINUSERNAME=<логин для базы>
ME_CONFIG_MONGODB_ADMINPASSWORD=<пароль для базы>
ME_CONFIG_MONGODB_URL=mongodb://<логин>:<пароль>@mongo:27017/
```

Затем запустить команду `docker-compose up -d --build`, чтобы собрать образ бота и запустить контейнеры в стэке.
