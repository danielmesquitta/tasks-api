-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `tasks` (
  id VARCHAR(36) NOT NULL PRIMARY KEY DEFAULT (UUID()),
  summary TEXT NOT NULL,
  assigned_to_user_id VARCHAR(36),
  created_by_user_id VARCHAR(36) NOT NULL,
  finished_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT fk_assigned_to_user FOREIGN KEY (assigned_to_user_id) REFERENCES users(id) ON DELETE
  SET NULL,
    CONSTRAINT fk_created_by_user FOREIGN KEY (created_by_user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE `tasks`;
-- +goose StatementEnd