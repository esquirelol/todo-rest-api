CREATE TABLE IF NOT EXISTS tasks(
  id SERIAL PRIMARY KEY,
  author TEXT  NOT NULL,
  title TEXT  NOT NULL,
  description TEXT,
  status BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT now(),
  completed_at TIMESTAMP
);