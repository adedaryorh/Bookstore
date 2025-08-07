
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(20) UNIQUE,
    publication_year VARCHAR(4),
    genre VARCHAR(100),
    price DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO books (title, author, isbn, publication_year, genre, price) VALUES
('The Go Programming Language', 'Alan Donovan, Brian Kernighan', '9780134190440', '2015', 'Programming', 45.99),
('Clean Code', 'Robert C. Martin', '9780132350884', '2008', 'Programming', 42.99),
('The Pragmatic Programmer', 'David Thomas, Andrew Hunt', '9780135957059', '2019', 'Programming', 39.99)
ON CONFLICT (isbn) DO NOTHING;
