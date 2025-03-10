CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    number VARCHAR(20) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('stylist', 'customer')),
    password TEXT NOT NULL,
    location TEXT, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);