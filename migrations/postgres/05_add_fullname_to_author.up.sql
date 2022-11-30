ALTER TABLE author ADD COLUMN fullname VARCHAR(610);

UPDATE author SET fullname = firstname || ' ' || lastname ||
(SELECT CASE WHEN middlename IS NULL THEN '' ELSE  ' ' || middlename END AS middlename);

ALTER TABLE author DROP COLUMN firstname;
ALTER TABLE author DROP COLUMN lastname;
ALTER TABLE author DROP COLUMN middlename;