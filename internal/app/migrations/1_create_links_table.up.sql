BEGIN;
CREATE TABLE IF NOT EXISTS links(
    id serial primary key, 
    shorturl char(10) UNIQUE , 
    originalurl text);
COMMIT;


