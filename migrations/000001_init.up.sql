CREATE TABLE IF NOT EXISTS cities (
    id serial primary key,
    name TEXT NOT NULL unique ,
    country TEXT NOT NULL,
    longitude TEXT NOT NULL,
    latitude TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS forecast (
    cityId INTEGER NOT NULL,
    temp DOUBLE PRECISION NOT NULL,
    date DATE NOT NULL,
    misc JSONB NOT NULL
);