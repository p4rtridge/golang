CREATE TABLE IF NOT EXISTS tasks (
  id        serial  PRIMARY KEY,
  title     text,
  completed boolean,
)
