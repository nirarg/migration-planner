-- +goose Up
-- +goose StatementBegin
ALTER TABLE sources ADD COLUMN share_token VARCHAR(64) UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sources DROP COLUMN share_token;
-- +goose StatementEnd 