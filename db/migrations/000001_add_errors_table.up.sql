CREATE TABLE errors (
    id         SERIAL PRIMARY KEY,
    expression text NOT NULL,
    endpoint   text NOT NULL,
    frequency  integer NOT NULL default 0,
    kind       text NOT NULL
);

ALTER TABLE errors
    ADD CONSTRAINT errors_unique_constraint
    UNIQUE(expression, endpoint, kind);