CREATE TABLE IF NOT EXISTS clients(
    client_id VARCHAR(50) PRIMARY KEY,
    api_key TEXT NOT NULL,
    capacity NUMERIC NOT NULL,
    rps NUMERIC NOT NULL
);
