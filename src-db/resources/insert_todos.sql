-- insert testdata into tables in database demo_todo.db

use demo_todo;

DELETE FROM todos;

INSERT INTO todos (title, completed, created_at, completed_at) 
  VALUES ('Bake a cake', TRUE, STR_TO_DATE('2025-02-18 15:44:04', '%Y-%m-%d %H:%i:%s'), STR_TO_DATE('2025-02-18 16:44:04', '%Y-%m-%d %H:%i:%s'));
INSERT INTO todos (title, completed, created_at, completed_at) 
  VALUES ('Feed the cat', FALSE, STR_TO_DATE('2025-02-18 15:55:04', '%Y-%m-%d %H:%i:%s'), NULL);
INSERT INTO todos (title, completed, created_at, completed_at) 
  VALUES ('Take out the trash', FALSE, STR_TO_DATE('2025-02-18 15:57:10', '%Y-%m-%d %H:%i:%s'), NULL);

SELECT * FROM todos;
