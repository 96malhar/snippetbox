-- migrate -path ./internal/store/migration -database "postgresql://postgres@localhost:5432/snippetbox?sslmode=disable" -verbose down

DROP TABLE IF EXISTS snippets;

DROP USER IF EXISTS web;
