
CREATE TABLE IF NOT EXISTS ips (
    id SERIAL PRIMARY KEY,
    ip VARCHAR(100),
    ping_time int,
    last_ok DATE
);