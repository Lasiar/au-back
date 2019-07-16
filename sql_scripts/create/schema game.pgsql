// создании схемы для игр
CREATE SCHEMA IF NOT EXISTS game

// создание таблицы для сессии игр
CREATE TABLE game.sessions (
    id SERIAL PRIMARY KEY,
    id_user INT REFERENCES auth.users (id) on UPDATE CASCADE,
    secret text not null
)

// создания таблицы для роундов (попытка отгадывания)
CREATE TABLE game.lap (
    id_session INT REFERENCES game.sessions (id) ON UPDATE CASCADE,
    input text not null,
    dt TIMESTAMP DEFAULT NOW()
)