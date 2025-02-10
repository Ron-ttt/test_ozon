BEGIN;
CREATE TABLE IF NOT EXISTS links(
    id serial primary key, 
    shorturl text UNIQUE, 
    originalurl text);
COMMIT;


