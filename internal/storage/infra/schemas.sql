CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(72) NOT NULL,
    email VARCHAR(80) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO
    users (username, password, email)
VALUES
    (
        'luisnquin',
        'dumm-password',
        'lpaandres2020@gmail.com'
    );