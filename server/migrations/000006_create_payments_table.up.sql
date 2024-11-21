CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    payment_method VARCHAR(50) NOT NULL, -- e.g., 'credit_card', 'paypal'
    payment_status VARCHAR(50) NOT NULL DEFAULT 'pending', -- e.g., 'pending', 'completed'
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
