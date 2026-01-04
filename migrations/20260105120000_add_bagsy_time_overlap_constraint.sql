-- +goose Up
-- +goose StatementBegin

-- Включаем расширение для GiST индексов (если еще не включено)
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Создаем EXCLUDE constraint для предотвращения пересечения временных интервалов
-- Для одного мастера в одно время может быть только одна активная бронь
ALTER TABLE bagsies
ADD CONSTRAINT bagsy_no_time_overlap
EXCLUDE USING gist (
    master_phone WITH =,
    tstzrange(start_at, end_at) WITH &&
)
WHERE (
    deleted_at IS NULL
    AND status NOT IN ('canceled',  'completed')
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Удаляем constraint
ALTER TABLE bagsies
DROP CONSTRAINT IF EXISTS bagsy_no_time_overlap;

-- Расширение не удаляем, т.к. может использоваться в других местах
-- DROP EXTENSION IF EXISTS btree_gist;

-- +goose StatementEnd
