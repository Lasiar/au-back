-- удаление старых таблиц
DROP TABLE IF EXISTS auth.sessions;
DROP TABLE IF EXISTS auth.users;
DROP TABLE IF EXISTS auth.permissions;
DROP SCHEMA IF EXISTS auth;

CREATE SCHEMA IF NOT EXISTS auth;


-- создание таблицы уровней доступа
CREATE TABLE auth.permissions (
    id SERIAL PRIMARY KEY,
    code TEXT not null,
    name TEXT
)

-- создание таблицы пользователей
CREATE TABLE auth.users (
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    pass TEXT NOT NULL,
    name TEXT,
    id_permission INT REFERENCES auth.permissions (id) on UPDATE CASCADE
);


-- создание таблицы сессий
CREATE TABLE auth.sessions (
    token TEXT PRIMARY KEY NOT NULL,
    user_id INT REFERENCES auth.users (id) ON UPDATE CASCADE,
    last_update TIMESTAMP DEFAULT NULL
);



-- заполнение стандартными значениями
INSERT INTO auth.permissions (code, name, mask) VALUES ('user', 'обычный пользователь', 1);
INSERT INTO auth.permissions (code, name, mask) VALUES ('admin_full', 'администратор', 2);
