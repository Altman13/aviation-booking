-- 1. Таблица пользователей
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email_verified_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 2. Справочник аэропортов (IATA коды)
CREATE TABLE airports (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iata_code VARCHAR(3) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    city_name VARCHAR(255) NOT NULL,
    country_code VARCHAR(2) NOT NULL,
    latitude DECIMAL(9,6),
    longitude DECIMAL(9,6),
    timezone VARCHAR(50)
);

-- 3. Справочник авиакомпаний
CREATE TABLE airlines (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iata_code VARCHAR(2) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    logo_url VARCHAR(255)
);

-- 4. История поисковых запросов
CREATE TABLE searches (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NULL,
    origin_airport_id INT NOT NULL,
    destination_airport_id INT NOT NULL,
    departure_date DATE NOT NULL,
    return_date DATE NULL,
    adults INT DEFAULT 1,
    children INT DEFAULT 0,
    infants INT DEFAULT 0,
    class VARCHAR(20) DEFAULT 'economy',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (origin_airport_id) REFERENCES airports(id),
    FOREIGN KEY (destination_airport_id) REFERENCES airports(id)
);

-- 5. Таблица с найденными вариантами перелетов
CREATE TABLE flights (
    id INT AUTO_INCREMENT PRIMARY KEY,
    search_id INT NOT NULL,
    airline_id INT NOT NULL,
    flight_number VARCHAR(10) NOT NULL,
    departure_airport_id INT NOT NULL,
    arrival_airport_id INT NOT NULL,
    departure_utc DATETIME NOT NULL,
    arrival_utc DATETIME NOT NULL,
    duration INT,
    transfers INT DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'RUB',
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (search_id) REFERENCES searches(id) ON DELETE CASCADE,
    FOREIGN KEY (airline_id) REFERENCES airlines(id),
    FOREIGN KEY (departure_airport_id) REFERENCES airports(id),
    FOREIGN KEY (arrival_airport_id) REFERENCES airports(id),
    UNIQUE KEY unique_flight_per_search (search_id, airline_id, flight_number, departure_utc)
);

-- 6. История цен
CREATE TABLE price_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    flight_id INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RUB',
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (flight_id) REFERENCES flights(id) ON DELETE CASCADE
);

-- 7. Избранные направления пользователей
CREATE TABLE user_favorites (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    origin_airport_id INT NOT NULL,
    destination_airport_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (origin_airport_id) REFERENCES airports(id),
    FOREIGN KEY (destination_airport_id) REFERENCES airports(id),
    UNIQUE KEY unique_user_favorite (user_id, origin_airport_id, destination_airport_id)
);