-- +goose Up
CREATE TABLE mappings(
    path VARCHAR(255) PRIMARY KEY,
    url VARCHAR(255) NOT NULL
);

INSERT INTO mappings
VALUES
('/test', 'https://github.com/kevinjimenez96/url-shortener#url-shortener' );

-- +goose Down
DROP TABLE mappings;
