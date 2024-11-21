CREATE TABLE stock_history (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    change INTEGER NOT NULL, -- Positive or negative value indicating the stock change
    reason TEXT, -- E.g., 'restock', 'order purchase'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
