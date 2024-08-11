-- +migrate Up
CREATE TABLE IF NOT EXISTS services (
    service_id          TEXT PRIMARY KEY,
    state               TEXT,

    occurred_at         TIMESTAMP WITHOUT TIME ZONE,
    updated_at          TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS services_id_idx ON services(service_id);

-- +migrate Down
DROP TABLE IF EXISTS services;
