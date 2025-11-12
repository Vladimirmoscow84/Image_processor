BEGIN;

CREATE TABLE IF NOT EXISTS images(
    id SERIAL PRIMARY KEY,
    original_path TEXT NOT NULL,
    processed_path TEXT DEFAULT NULL,
    thumbnail_path TEXT DEFAULT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_images_original_path ON images(original_path);
CREATE INDEX IF NOT EXISTS idx_images_status ON images(status);

COMMIT;