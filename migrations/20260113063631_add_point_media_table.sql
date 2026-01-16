-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS point_media (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   point_code TEXT NOT NULL,
   media_id UUID NOT NULL,
-- Порядок отображения
   display_order INTEGER NOT NULL DEFAULT 0,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   deleted_at TIMESTAMPTZ
);

-- Индексы для point_media
CREATE INDEX idx_point_media_point_code ON point_media(point_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_point_media_media_id ON point_media(media_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_point_media_display_order ON point_media(point_code, display_order) WHERE deleted_at IS NULL;

-- Уникальный порядок для каждого point
CREATE UNIQUE INDEX idx_point_media_unique_order
    ON point_media(point_code, display_order)
    WHERE deleted_at IS NULL;


CREATE TRIGGER trigger_update_point_media_updated_at
    BEFORE UPDATE ON point_media
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


COMMENT ON TABLE point_media IS 'Point location photos';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS point_media;
DROP TRIGGER IF EXISTS trigger_update_point_media_updated_at ON point_media;
-- +goose StatementEnd
