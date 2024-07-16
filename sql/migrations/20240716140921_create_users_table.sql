-- +goose Up
-- +goose StatementBegin
CREATE TABLE `users` (
  id VARCHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
  role TINYINT UNSIGNED NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT chk_role CHECK (role IN (1, 2))
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE `users`;
-- +goose StatementEnd