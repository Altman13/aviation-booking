-- 1. Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    email_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Справочник аэропортов (IATA коды)
CREATE TABLE airports (
    id SERIAL PRIMARY KEY,
    iata_code VARCHAR(3) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    city_name VARCHAR(255) NOT NULL,
    country_code VARCHAR(2) NOT NULL,
    latitude DECIMAL(9,6),
    longitude DECIMAL(9,6),
    timezone VARCHAR(50)
);
-- Индексы для аэропортов
CREATE INDEX idx_airports_iata ON airports(iata_code);
CREATE INDEX idx_airports_country ON airports(country_code);

-- 3. Справочник авиакомпаний
CREATE TABLE airlines (
    id SERIAL PRIMARY KEY,
    iata_code VARCHAR(2) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    logo_url VARCHAR(255)
);
CREATE INDEX idx_airlines_iata ON airlines(iata_code);

-- 4. Таблица рейсов (основная информация о рейсах)
CREATE TABLE flights (
    id SERIAL PRIMARY KEY,
    airline_id INTEGER NOT NULL REFERENCES airlines(id),
    flight_number VARCHAR(10) NOT NULL,
    departure_airport_id INTEGER NOT NULL REFERENCES airports(id),
    arrival_airport_id INTEGER NOT NULL REFERENCES airports(id),
    departure_date DATE NOT NULL,
    departure_time TIME NOT NULL,
    arrival_date DATE NOT NULL,
    arrival_time TIME NOT NULL,
    duration_minutes INTEGER NOT NULL,
    aircraft_type VARCHAR(50),
    base_price_economy DECIMAL(10, 2) NOT NULL,
    base_price_business DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    available_seats_economy INTEGER DEFAULT 0,
    available_seats_business INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(flight_number, departure_date, departure_time)
);
-- Индексы для рейсов
CREATE INDEX idx_flights_departure_date ON flights(departure_date);
CREATE INDEX idx_flights_route ON flights(departure_airport_id, arrival_airport_id);
CREATE INDEX idx_flights_airline_departure ON flights(airline_id, departure_date);
CREATE INDEX idx_flights_active ON flights(is_active) WHERE is_active = TRUE;

-- 5. История поисковых запросов (для ML и аналитики)
CREATE TABLE searches (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    origin_airport_id INTEGER NOT NULL REFERENCES airports(id),
    destination_airport_id INTEGER NOT NULL REFERENCES airports(id),
    departure_date DATE NOT NULL,
    return_date DATE,
    adults INTEGER DEFAULT 1,
    children INTEGER DEFAULT 0,
    infants INTEGER DEFAULT 0,
    travel_class VARCHAR(20) DEFAULT 'ECONOMY',
    currency VARCHAR(3) DEFAULT 'RUB',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- Индексы для поисков
CREATE INDEX idx_searches_created_at ON searches(created_at);
CREATE INDEX idx_searches_route ON searches(origin_airport_id, destination_airport_id);
CREATE INDEX idx_searches_dates ON searches(departure_date, return_date);
CREATE INDEX idx_searches_user_date ON searches(user_id, departure_date) WHERE user_id IS NOT NULL;

-- 6. Таблица бронирований
CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    booking_reference VARCHAR(20) UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    flight_id INTEGER NOT NULL REFERENCES flights(id),
    travel_class VARCHAR(20) NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    status VARCHAR(20) DEFAULT 'CONFIRMED',
    payment_method VARCHAR(50),
    payment_status VARCHAR(20) DEFAULT 'PENDING',
    contact_email VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(20) NOT NULL,
    special_requests TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    -- Индексы через CREATE INDEX для лучшей производительности
    UNIQUE(booking_reference)
);
-- Индексы для бронирований
CREATE INDEX idx_bookings_user ON bookings(user_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_payment_status ON bookings(payment_status);
CREATE INDEX idx_bookings_created_at ON bookings(created_at);
CREATE INDEX idx_bookings_flight ON bookings(flight_id, travel_class);

-- 7. Таблица пассажиров
CREATE TABLE passengers (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    birth_date DATE NOT NULL,
    document_type VARCHAR(20) DEFAULT 'PASSPORT',
    document_number VARCHAR(50) NOT NULL,
    nationality VARCHAR(3),
    seat_preference VARCHAR(20),
    special_requests TEXT,
    seat_assignment VARCHAR(10),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- Индексы для пассажиров
CREATE INDEX idx_passengers_booking ON passengers(booking_id);
CREATE INDEX idx_passengers_document ON passengers(document_number);
CREATE INDEX idx_passengers_name ON passengers(first_name, last_name);

-- 8. ИСТОРИЯ ЦЕН (для ML и прогноза)
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    flight_id INTEGER NOT NULL REFERENCES flights(id) ON DELETE CASCADE,
    travel_class VARCHAR(20) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    available_seats INTEGER NOT NULL,
    recorded_at TIMESTAMPTZ DEFAULT NOW()
);
-- Индексы для истории цен
CREATE INDEX idx_price_history_flight_id ON price_history(flight_id);
CREATE INDEX idx_price_history_recorded_at ON price_history(recorded_at);
CREATE INDEX idx_price_history_flight_class ON price_history(flight_id, travel_class);
CREATE INDEX idx_price_history_timestamp ON price_history(flight_id, travel_class, recorded_at DESC);

-- 9. Токены пользователей (для аутентификации)
CREATE TABLE user_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_token VARCHAR(512) NOT NULL,
    refresh_token VARCHAR(512) NOT NULL,
    token_type VARCHAR(20) DEFAULT 'Bearer',
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, access_token)
);
-- Индексы для токенов
CREATE INDEX idx_user_tokens_user ON user_tokens(user_id);
CREATE INDEX idx_user_tokens_expires ON user_tokens(expires_at);
CREATE INDEX idx_user_tokens_refresh ON user_tokens(refresh_token);

-- 10. Избранные направления пользователей
CREATE TABLE user_favorites (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    origin_airport_id INTEGER NOT NULL REFERENCES airports(id),
    destination_airport_id INTEGER NOT NULL REFERENCES airports(id),
    note VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, origin_airport_id, destination_airport_id)
);
-- Индексы для избранного
CREATE INDEX idx_user_favorites_user ON user_favorites(user_id);
CREATE INDEX idx_user_favorites_route ON user_favorites(origin_airport_id, destination_airport_id);

-- 11. Таблица сессий пользователей
CREATE TABLE user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    device_info TEXT,
    ip_address INET,
    user_agent TEXT,
    last_activity TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- Индексы для сессий
CREATE INDEX idx_user_sessions_user ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_expires ON user_sessions(expires_at);
CREATE INDEX idx_user_sessions_token ON user_sessions(session_token);

-- 12. Таблица уведомлений (новое)
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    read_at TIMESTAMPTZ
);
-- Индексы для уведомлений
CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_read ON notifications(is_read) WHERE is_read = FALSE;
CREATE INDEX idx_notifications_created ON notifications(created_at DESC);

-- 13. Таблица для отслеживания изменений мест (новое)
CREATE TABLE seat_availability_log (
    id SERIAL PRIMARY KEY,
    flight_id INTEGER NOT NULL REFERENCES flights(id) ON DELETE CASCADE,
    travel_class VARCHAR(20) NOT NULL,
    seats_before INTEGER NOT NULL,
    seats_after INTEGER NOT NULL,
    change_reason VARCHAR(50), -- 'BOOKING', 'CANCELLATION', 'ADMIN_UPDATE'
    booking_id INTEGER REFERENCES bookings(id) ON DELETE SET NULL,
    changed_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);
-- Индексы для лога мест
CREATE INDEX idx_seat_log_flight ON seat_availability_log(flight_id, travel_class);
CREATE INDEX idx_seat_log_date ON seat_availability_log(changed_at DESC);

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггеры для обновления updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_flights_updated_at BEFORE UPDATE ON flights
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bookings_updated_at BEFORE UPDATE ON bookings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_favorites_updated_at BEFORE UPDATE ON user_favorites
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Функция для генерации booking_reference
CREATE OR REPLACE FUNCTION generate_booking_reference()
RETURNS TRIGGER AS $$
BEGIN
    NEW.booking_reference := 'BOOK-' || LPAD(NEW.id::TEXT, 6, '0');
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггер для автоматической генерации booking_reference
CREATE TRIGGER set_booking_reference BEFORE INSERT ON bookings
    FOR EACH ROW EXECUTE FUNCTION generate_booking_reference();