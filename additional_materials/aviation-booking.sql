-- 1. Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL, -- Захэшированный пароль
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Справочник аэропортов (IATA коды)
CREATE TABLE airports (
    id SERIAL PRIMARY KEY,
    iata_code VARCHAR(3) UNIQUE NOT NULL, -- Например, 'MOW', 'LED', 'SVO'
    name VARCHAR(255) NOT NULL, -- 'Шереметьево'
    city_name VARCHAR(255) NOT NULL, -- 'Москва'
    country_code VARCHAR(2) NOT NULL, -- 'RU'
    latitude DECIMAL(9,6),
    longitude DECIMAL(9,6),
    timezone VARCHAR(50)
);
-- Надо наполнить эту таблицу данными заранее

-- 3. Справочник авиакомпаний
CREATE TABLE airlines (
    id SERIAL PRIMARY KEY,
    iata_code VARCHAR(2) UNIQUE NOT NULL, -- Например, 'SU', 'S7'
    name VARCHAR(255) NOT NULL, -- 'Аэрофлот'
    logo_url VARCHAR(255)
);

-- 4. История поисковых запросов (ЗОЛОТАЯ жила для ML)
CREATE TABLE searches (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL, -- Поиск может быть и анонимным
    origin_airport_id INTEGER NOT NULL REFERENCES airports(id),
    destination_airport_id INTEGER NOT NULL REFERENCES airports(id),
    departure_date DATE NOT NULL,
    return_date DATE, -- NULL для one-way
    adults INTEGER DEFAULT 1,
    children INTEGER DEFAULT 0,
    infants INTEGER DEFAULT 0,
    class VARCHAR(20) DEFAULT 'economy',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- Индексы для ускорения аналитических запросов
CREATE INDEX idx_searches_created_at ON searches(created_at);
CREATE INDEX idx_searches_route ON searches(origin_airport_id, destination_airport_id);

-- 5. Таблица с найденными вариантами перелетов (сырые данные из парсинга/API)
CREATE TABLE flights (
    id SERIAL PRIMARY KEY,
    search_id INTEGER NOT NULL REFERENCES searches(id) ON DELETE CASCADE, -- К какому запросу относится этот результат
    airline_id INTEGER NOT NULL REFERENCES airlines(id),
    flight_number VARCHAR(10) NOT NULL, -- 'SU 123'
    departure_airport_id INTEGER NOT NULL REFERENCES airports(id),
    arrival_airport_id INTEGER NOT NULL REFERENCES airports(id),
    departure_utc TIMESTAMPTZ NOT NULL,
    arrival_utc TIMESTAMPTZ NOT NULL,
    duration INTEGER, -- В минутах
    transfers INTEGER DEFAULT 0, -- Количество пересадок
    currency VARCHAR(3) DEFAULT 'RUB',
    price DECIMAL(10, 2) NOT NULL, -- Актуальная цена на момент поиска
    created_at TIMESTAMPTZ DEFAULT NOW(),
    -- Уникальный constraint, чтобы избежать дубликатов в рамках одного поиска
    UNIQUE(search_id, airline_id, flight_number, departure_utc)
);

-- 6. ИСТОРИЯ ЦЕН (Самая важная таблица для ML и прогноза)
-- Здесь будем хранить снимок цены каждого рейса в момент парсинга
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    flight_id INTEGER NOT NULL REFERENCES flights(id) ON DELETE CASCADE,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    recorded_at TIMESTAMPTZ DEFAULT NOW() -- Момент, когда мы увидели эту цену
);
-- Индекс для быстрого построения графиков цен по одному рейсу
CREATE INDEX idx_price_history_flight_id ON price_history(flight_id);
CREATE INDEX idx_price_history_recorded_at ON price_history(recorded_at);

-- 7. Избранные направления пользователей (для персонализации)
CREATE TABLE user_favorites (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    origin_airport_id INTEGER NOT NULL REFERENCES airports(id),
    destination_airport_id INTEGER NOT NULL REFERENCES airports(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, origin_airport_id, destination_airport_id) -- Чтобы не дублировать
);