-- +migrate Up
CREATE TABLE IF NOT EXISTS services (
    service_id          TEXT PRIMARY KEY,
    state               TEXT,
    occurred_at         TIMESTAMP WITHOUT TIME ZONE,
    updated_at          TIMESTAMP WITHOUT TIME ZONE
);

-- +migrate Down
DROP TABLE IF EXISTS services;
