CREATE TABLE IF NOT EXISTS sentences (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sentence_hash TEXT NOT NULL,
    sentence TEXT NOT NULL,
    correction TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS sentence_index ON sentences(sentence_hash);

-- DROP TABLE IF EXISTS sentences;
