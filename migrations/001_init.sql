CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS books (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    category TEXT NOT NULL,
    description TEXT NOT NULL,
    price INTEGER NOT NULL,
    cover_url TEXT,
    file_path TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO users (
    name,
    email,
    password,
    role
)
VALUES (
    'Admin Bukuin',
    'admin@bukuin.test',
    'admin123',
    'admin'
);

INSERT INTO books (
    title,
    author,
    category,
    description,
    price,
    cover_url,
    file_path
)
VALUES
(
    'Laskar Pelangi',
    'Andrea Hirata',
    'Novel',
    'Novel tentang perjuangan sekelompok anak dalam mengejar pendidikan.',
    50000,
    '',
    ''
),
(
    'Bumi Manusia',
    'Pramoedya Ananta Toer',
    'Novel',
    'Novel sejarah yang menggambarkan kehidupan masyarakat pada masa kolonial.',
    60000,
    '',
    ''
),
(
    'Madilog',
    'Tan Malaka',
    'Filsafat',
    'Buku yang membahas materialisme, dialektika, dan logika.',
    45000,
    '',
    ''
),
(
    'Filosofi Teras',
    'Henry Manampiring',
    'Pengembangan Diri',
    'Pengenalan konsep filsafat Stoisisme dalam kehidupan modern.',
    70000,
    '',
    ''
),
(
    'Atomic Habits',
    'James Clear',
    'Pengembangan Diri',
    'Buku tentang membangun kebiasaan kecil secara konsisten.',
    80000,
    '',
    ''
);