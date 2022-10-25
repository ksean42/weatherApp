CREATE schema postgres;
CREATE TABLE IF NOT EXISTS cities (
    city_id serial NOT NULL unique ,
    name TEXT NOT NULL unique ,
    country TEXT NOT NULL,
    longitude TEXT NOT NULL,
    latitude TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS forecast (
    city_id serial NOT NULL unique ,
    temp DOUBLE PRECISION NOT NULL,
    date DATE NOT NULL,
    misc JSONB NOT NULL
);