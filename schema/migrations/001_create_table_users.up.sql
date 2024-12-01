CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE songs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- UUID с автоматической генерацией
    group_name VARCHAR(255) NOT NULL,               -- Группа, исполняющая песню
    name VARCHAR(255) NOT NULL,                     -- Название песни
    release_date DATE NOT NULL,                     -- Дата выхода песни
    text TEXT NOT NULL,                             -- Текст песни
    link TEXT NOT NULL                              -- Ссылка на песню
);
