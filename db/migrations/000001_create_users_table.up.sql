BEGIN;

CREATE TABLE todos (
  "id" PRIMARY KEY, 
  "status" varchar(32) NOT NULL DEFAULT 'ACTIVE',
  "description" varchar(128) NOT NULL,   
  "Author"
  "AuthorID"  PRIMARY KEY
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()), 
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "deleted_at"   TIMESTAMPTZ
  
);
COMMIT;

CREATE TABLE users (
  "id" PRIMARY KEY, 
  "name" varchar(32) NOT NULL DEFAULT 'ACTIVE',
  "email" varchar(32) NOT NULL,   
  "password"  NOT NULL , 
  "created_at"  TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "deleted_at"  TIMESTAMPTZ NOT NULL DEFAULT (now())
);
COMMIT;