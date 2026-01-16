-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_media (
  user_phone TEXT PRIMARY KEY,
  media_id UUID NOT NULL, -- FK media(id)
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индексы для user_media
CREATE INDEX idx_user_media_media_id ON user_media(media_id);
COMMENT ON TABLE user_media IS 'User avatars (one avatar per user)';


-- Тригеры для updated_at
CREATE TRIGGER trigger_update_user_media_updated_at
    BEFORE UPDATE ON user_media
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_media;
DROP TRIGGER IF EXISTS trigger_update_user_media_updated_at ON user_media;
-- +goose StatementEnd
