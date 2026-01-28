-- +goose Up
-- +goose StatementBegin
ALTER TABLE services DROP CONSTRAINT services_point_code_cat_subcat_unique;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE services ADD CONSTRAINT services_point_code_cat_subcat_unique
    UNIQUE (point_code, category_id, subcategory_id);
-- +goose StatementEnd
