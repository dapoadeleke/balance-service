CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   name VARCHAR(50) NOT NULL,
   balance DECIMAL(10, 2) DEFAULT 0.00,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);