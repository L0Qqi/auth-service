CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY, -- JTI токена
    user_id INTEGER NOT NULL, -- ID пользователя, которому принадлежит токен
    token_hash TEXT NOT NULL, -- Хеш токена (bcrypt)
    revoked BOOLEAN NOT NULL DEFAULT FALSE, -- Статус отзыва
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Дата создания
    finish_at TIMESTAMPTZ NOT NULL -- Когда токен истекает
);