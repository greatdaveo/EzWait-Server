CREATE TABLE stylists (
    id SERIAL PRIMARY KEY,
    stylist_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    active_status BOOLEAN DEFAULT TRUE,
    profile_picture TEXT,
    ratings FLOAT,
    services JSONB,
    service_img JSONB,
    available_time_slots JSONB,
    no_of_customer_bookings INTEGER DEFAULT 0,
    no_of_current_customers INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
