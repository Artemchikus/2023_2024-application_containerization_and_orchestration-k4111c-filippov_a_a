-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vk_channels (
 id SERIAL PRIMARY KEY,
 channel_name VARCHAR(255) NOT NULL,
 channel_url VARCHAR(255) NOT NULL,
 channel_type VARCHAR(255) NOT NULL,
 site_url VARCHAR(255) NOT NULL,
 created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE FUNCTION upd_trig() RETURNS trigger
   LANGUAGE plpgsql AS
$$BEGIN
   NEW.updated_at := CURRENT_TIMESTAMP;
   RETURN NEW;
END;$$;

CREATE TRIGGER upd_trig BEFORE UPDATE ON vk_channels
   FOR EACH ROW EXECUTE PROCEDURE upd_trig();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vk_channels;
-- +goose StatementEnd
