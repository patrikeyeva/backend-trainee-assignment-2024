-- +goose Up
-- +goose StatementBegin
CREATE TABLE Banner (
    banner_id BIGSERIAL NOT NULL PRIMARY KEY,
    feature_id INT NOT NULL,
    "text" TEXT,
    title TEXT,
    "url" TEXT,
    is_active BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE
    BannerTag (
        banner_id INT,
        tag_id INT,
        PRIMARY KEY (banner_id, tag_id),
        FOREIGN KEY (banner_id) REFERENCES Banner (banner_id)
    );

ALTER TABLE IF EXISTS Banner OWNER to "postgres";

ALTER TABLE IF EXISTS BannerTag OWNER to "postgres";

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Banner CASCADE;

DROP TABLE IF EXISTS BannerTag CASCADE;


-- +goose StatementEnd