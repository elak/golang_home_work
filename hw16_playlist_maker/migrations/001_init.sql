-- +goose Up
CREATE TABLE `CATEGORIES`
(
    ID    varchar(36) NOT NULL PRIMARY KEY,
    Title text        NOT NULL
);

CREATE TABLE `GROUPS`
(
    ID         varchar(36) NOT NULL PRIMARY KEY,
    Title      text        NOT NULL,
    ParentID   varchar(36) NOT NULL,
    `Order`    int         NOT NULL,
    CategoryID varchar(36) NOT NULL
);

CREATE TABLE `VIDEOS`
(
    ID         varchar(36) NOT NULL PRIMARY KEY,
    Title      text        NOT NULL,
    ParentID   varchar(36) NOT NULL,
    `Order`    int         NOT NULL,
    CategoryID varchar(36) NOT NULL,
    Duration   int         NOT NULL
);

CREATE TABLE `TEMPLATES`
(
    ID    varchar(36) NOT NULL PRIMARY KEY,
    Title text        NOT NULL
);

CREATE TABLE `TEMPLATES_ITEMS`
(
    ID            varchar(36) NOT NULL PRIMARY KEY,
    Title         text        NOT NULL,
    `Order`       int         NOT NULL,
    Duration      int         NOT NULL,
    TemplateID    varchar(36) NOT NULL,
    TemplateBlock int         NOT NULL
);

CREATE TABLE `TEMPLATES_FILLERS`
(
    ID             varchar(36) NOT NULL PRIMARY KEY,
    `Order`        int         NOT NULL,
    CategoryID     varchar(36) NOT NULL,
    AllowRepeat    bool        NOT NULL,
    GroupsPriority int         NOT NULL,
    VideosPriority int         NOT NULL,
    OwnerID        varchar(36) NOT NULL
);

CREATE TABLE `RESTRICTIONS`
(
    ID         varchar(36) NOT NULL PRIMARY KEY,
    Title      text        NOT NULL,
    Scope      int         NOT NULL,
    CategoryID varchar(36) NOT NULL,
    GroupID    varchar(36) NOT NULL,
    Duration   int         NOT NULL,
    Amount     int         NOT NULL,
    OwnerID    varchar(36) NOT NULL
);

CREATE TABLE `HISTORY`
(
    VideoID  varchar(36) NOT NULL PRIMARY KEY,
    LastSeen datetime    NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS `CATEGORIES`;
DROP TABLE IF EXISTS `GROUPS`;
DROP TABLE IF EXISTS `VIDEOS`;
DROP TABLE IF EXISTS `TEMPLATES`;
DROP TABLE IF EXISTS `TEMPLATES_ITEMS`;
DROP TABLE IF EXISTS `TEMPLATES_FILLERS`;
DROP TABLE IF EXISTS `RESTRICTIONS`;
DROP TABLE IF EXISTS `HISTORY`;