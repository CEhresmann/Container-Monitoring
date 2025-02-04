CREATE TABLE IF NOT EXISTS ips (
    id SERIAL PRIMARY KEY,
    ip VARCHAR(100),
    ping_time INTEGER,
    last_ok TIMESTAMP
);