-- +goose Up
-- +goose StatementBegin
-- Основная таблица медиа-файлов
CREATE TABLE IF NOT EXISTS media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- S3 информация
    file_key TEXT NOT NULL UNIQUE,
    bucket_name TEXT NOT NULL,

    -- Метаданные файла
    original_filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,

    -- Для изображений
    width INTEGER,
    height INTEGER,

    -- Статус обработки
    status TEXT NOT NULL DEFAULT 'pending',

    -- Audit
    uploaded_by TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Индексы для media
CREATE INDEX idx_media_file_key ON media(file_key) WHERE deleted_at IS NULL;
CREATE INDEX idx_media_status ON media(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_media_uploaded_by ON media(uploaded_by) WHERE deleted_at IS NULL;
CREATE INDEX idx_media_created_at ON media(created_at DESC);

COMMENT ON TABLE media IS 'Central storage for all media files (images, photos)';
COMMENT ON COLUMN media.file_key IS 'S3 object key: YYYY/MM/{uuid}.{ext}';
COMMENT ON COLUMN media.status IS 'pending: uploaded to S3, active: in use, processing: being processed, failed: error';

-- ============================================
-- TRIGGERS
-- ============================================

-- Триггеры updated_at
CREATE TRIGGER trigger_update_media_updated_at
    BEFORE UPDATE ON media
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd






-- +goose Down
-- +goose StatementBegin
-- Удаляем триггеры
DROP TRIGGER IF EXISTS trigger_update_media_updated_at ON media;

-- Удаляем таблицы
DROP TABLE IF EXISTS media;
-- +goose StatementEnd
