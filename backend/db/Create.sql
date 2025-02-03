
CREATE TABLE IF NOT EXISTS employees (
                                         id SERIAL PRIMARY KEY,
                                         name VARCHAR(100),
                                         position VARCHAR(100),
                                         hire_date DATE
);