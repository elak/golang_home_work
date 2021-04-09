-- +goose Up
CREATE TABLE Events
(
    ID varchar(36) NOT NULL,
    Title text NOT NULL,
    Date date NOT NULL,
    DueDate date NOT NULL,
    Description text NOT NULL,
    Owner varchar(36) NOT NULL,
    NotifyDate date NOT NULL,

    PRIMARY KEY (ID)
);

-- +goose Down
DROP TABLE Events;