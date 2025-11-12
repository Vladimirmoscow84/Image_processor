BEGIN;

DROP INDEX IF EXISTS idx_images_original_path;
DROP INDEX IF EXISTS idx_images_status;

DROP TABLE IF EXISTS images;

COMMIT;