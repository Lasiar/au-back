-- создании схемы для игр
CREATE SCHEMA IF NOT EXISTS game

-- создание таблицы для сессии игр
CREATE TABLE game.sessions (
    id SERIAL PRIMARY KEY,
    id_user INT REFERENCES auth.users (id) on UPDATE CASCADE,
    secret text not null
)

-- создания таблицы для роундов (попытка отгадывания)
CREATE TABLE game.lap (
    id_session INT REFERENCES game.sessions (id) ON UPDATE CASCADE,
    input text not null,
    dt TIMESTAMP DEFAULT NOW()
}

-- создание представляения для отображения завершенных сессий
CREATE OR REPLACE VIEW game.v_sessions_completed AS
SELECT s.id,
       s.id_user,
       s.secret,
       lap.dt IS NOT NULL AS completed
FROM game.sessions s
LEFT JOIN game.lap lap ON lap.id_session = s.id
AND lap.input = s.secret;


-- создание представляения для отображения доски лидеров
CREATE OR REPLACE VIEW game.v_session_count_laps
AS SELECT s.id,
    s.id_user,
    s.secret,
    count(lap.*) AS laps
   FROM game.sessions s
     LEFT JOIN game.lap lap ON s.id = lap.id_session
  GROUP BY s.id, s.id_user, s.secret;
