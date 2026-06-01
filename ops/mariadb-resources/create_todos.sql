-- create tables in database demo_todo.db
use demo_todo;

CREATE TABLE IF NOT EXISTS todos (
    id int unsigned NOT NULL AUTO_INCREMENT,
    title varchar(255) CHARACTER SET utf8 NOT NULL,
    completed boolean NOT NULL DEFAULT FALSE, 
    created_at datetime NOT NULL,
    completed_at datetime DEFAULT NULL,
    PRIMARY KEY (id)
);
