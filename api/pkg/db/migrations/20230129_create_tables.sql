
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  nickname VARCHAR(255) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL UNIQUE,
  token VARCHAR(255) NOT NULL UNIQUE,
  is_active BOOLEAN DEFAULT false,
  password VARCHAR(100) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE todos (
  id SERIAL PRIMARY KEY,
  status VARCHAR(100) NOT NULL,
  description VARCHAR(255) NOT NULL,
  author_id INTEGER NOT NULL REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE activation_links (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    is_used BOOLEAN DEFAULT false,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE reset_passwords (
	id SERIAL PRIMARY KEY,
	email VARCHAR(100) NOT NULL,
	token VARCHAR(255) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NOW()
);
