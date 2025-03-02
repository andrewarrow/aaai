CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT UNIQUE NOT NULL,
  bio TEXT,
  linked_in_url TEXT,
  github_url TEXT,
  photo_url TEXT,
  password TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  token TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT OR IGNORE INTO users (username, bio, linked_in_url, github_url, photo_url, password) VALUES 
('Andrew Arrow', '30 years in software development, I use Claude Code to vibe.', 'https://www.linkedin.com/in/andrewarrow/', 'https://github.com/andrewarrow', 'https://avatars.githubusercontent.com/u/127054?v=4', 'testing');

INSERT OR IGNORE INTO users (username, bio, linked_in_url, github_url, photo_url, password) VALUES 
('Jane Smith', NULL, NULL, NULL, 'https://i.imgur.com/HQs9d76.png', 'testing');

INSERT OR IGNORE INTO users (username, bio, linked_in_url, github_url, photo_url, password) VALUES 
('John Smith', NULL, NULL, NULL, 'https://i.imgur.com/ANgLsc6.jpeg', 'testing');