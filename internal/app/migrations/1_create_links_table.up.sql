BEGIN;

CREATE TABLE IF NOT EXISTS links(
    id serial primary key, 
    shorturl char(10) UNIQUE , 
    originalurl text);
    
CREATE INDEX originalurlIndex ON links(originalurl);

COMMIT;


