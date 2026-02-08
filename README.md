# Movie Platform

Платформа для просмотра и обсуждения фильмов с использованием Go (Gin) и React.

## Требования проекта

- ✅ **Usage of middleware** - `AuthMiddleware`, `AdminMiddleware`, CORS middleware
- ✅ **Front-End (ready html template is ok)** - Go HTML templates в `templates/`
- ✅ **Usage of golang templating tags** - Используются `{{.TotalMovies}}`, `{{.TotalUsers}}`, `{{.TotalReviews}}`
- ✅ **Clean code and solid project structure** - Четкое разделение на пакеты
- ✅ **Bonus: usage of JS framework** - React компоненты встроены в Go templates

## Структура проекта

```
GO_movie_platform/
├── main.go                 # Точка входа
├── controllers/            # Обработчики запросов
│   ├── home.go            # Главная страница (Go template)
│   ├── movies.go          # Страницы с React
│   ├── movie.go           # API для фильмов
│   ├── review.go          # API для отзывов
│   └── user.go            # API для пользователей
├── middleware/            # Middleware
│   ├── auth.go           # Аутентификация
│   └── admin.go          # Проверка админ прав
├── models/                # Модели данных
├── routes/                # Маршруты
├── database/              # Подключение к БД
├── utils/                 # Утилиты
├── templates/             # Go HTML templates с React
│   ├── index.html        # Главная (Go template)
│   ├── movies.html       # Список фильмов (React)
│   └── movie-details.html # Детали фильма (React)
└── .env                   # Переменные окружения
```

## Установка

### Требования
- Go 1.21+
- MongoDB

### Шаги

1. **Установите зависимости Go**
```bash
go mod download
```

2. **Создайте файл `.env`**
```env
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=movie_platform
SECRET_KEY=your-secret-key-here
SECRET_REFRESH_KEY=your-refresh-secret-key-here
PORT=8080
```

## Запуск

```bash
go run main.go
```

Сервер запустится на `http://localhost:8080`

**URLs:**
- Главная (Go Template): http://localhost:8080/
- Фильмы (React): http://localhost:8080/movies
- API: http://localhost:8080/api

## API Endpoints

### Публичные
- `GET /api/movies` - Список фильмов
- `GET /api/movies/:id` - Детали фильма
- `GET /api/movies/:id/reviews` - Отзывы фильма
- `POST /api/register` - Регистрация
- `POST /api/login` - Вход

### Требуют авторизации
- `POST /api/movies/:id/reviews` - Добавить отзыв
- `DELETE /api/reviews/:id` - Удалить отзыв

### Требуют админ прав
- `POST /api/movies` - Добавить фильм
- `PUT /api/movies/:id` - Обновить фильм
- `DELETE /api/movies/:id` - Удалить фильм

## Технологии

**Backend:**
- Go 1.21+
- Gin (веб-фреймворк)
- MongoDB
- JWT (аутентификация)

**Frontend:**
- Go HTML Templates (требование)
- React 18 (CDN, встроен в templates)
- Vanilla CSS

## Особенности

- **Простота**: React загружается через CDN, не нужен npm/node
- **Go Templates**: Главная страница использует Go templating
- **React**: Компоненты встроены прямо в HTML через Babel
- **Один сервер**: Все работает на одном Go сервере
