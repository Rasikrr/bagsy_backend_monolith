-- +goose Up
-- +goose StatementBegin
--------------------------------------
-- Media Assets --
--------------------------------------

CREATE TABLE media_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Данные для S3 / Хранилища
    bucket VARCHAR(255) NOT NULL,
    object_key VARCHAR(1024) NOT NULL, -- Путь файла: 'organizations/123/logo.png'

    -- Метаданные для фронтенда и валидации
    filename VARCHAR(255) NOT NULL,    -- Исходное имя файла (например, 'my_photo.jpg')
    mime_type VARCHAR(100) NOT NULL,   -- 'image/jpeg', 'image/png', 'video/mp4'
    size_bytes BIGINT NOT NULL,        -- Важно для лимитов

    -- Статус для паттерна Direct-to-S3 Upload
    status VARCHAR(50) NOT NULL, -- 'pending', 'uploaded', 'failed'

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    -- В одном бакете не может быть двух файлов с одинаковым путем
    CONSTRAINT unique_media_assets_bucket_key
    UNIQUE (bucket, object_key)
);

-- Индекс для быстрого поиска "зависших" загрузок воркером (Garbage Collector)
CREATE INDEX idx_media_assets_status ON media_assets(status);




--------------------------------------
-- Media Junctions --
--------------------------------------

ALTER TABLE employees
    ADD COLUMN avatar_id UUID REFERENCES media_assets(id) ON DELETE SET NULL;

-- Фотографии организации (логотип, обложка, фото интерьера)
CREATE TABLE organization_media (
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES media_assets(id) ON DELETE CASCADE,

    type VARCHAR(50) NOT NULL, -- 'logo', 'cover', 'gallery'
    sort_order INT DEFAULT 0,

    PRIMARY KEY (organization_id, media_id)
);

-- Фотографии локации (фасад, интерьер, рабочие места)
CREATE TABLE location_media (
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES media_assets(id) ON DELETE CASCADE,

    type VARCHAR(50) NOT NULL, -- 'gallery', maybe add other
    sort_order INT DEFAULT 0,

    PRIMARY KEY (location_id, media_id)
);

-- Фотографии категорий услуг (иконки или обложки категорий)
CREATE TABLE service_category_media (
    category_id UUID NOT NULL REFERENCES service_categories(id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES media_assets(id) ON DELETE CASCADE,

    type VARCHAR(50) NOT NULL, -- 'icon', 'cover'
    sort_order INT DEFAULT 0,

    PRIMARY KEY (category_id, media_id)
);

-- Фотографии конкретных услуг (результаты работ "до/после", примеры)
CREATE TABLE service_media (
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES media_assets(id) ON DELETE CASCADE,

    type VARCHAR(50) NOT NULL, -- 'gallery'
    sort_order INT DEFAULT 0,

    PRIMARY KEY (service_id, media_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
