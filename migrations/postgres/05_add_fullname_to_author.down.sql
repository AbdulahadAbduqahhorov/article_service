ALTER TABLE author ADD COLUMN middlename VARCHAR(55);
ALTER TABLE author ADD COLUMN lastname VARCHAR(55);
ALTER TABLE author ADD COLUMN firstname VARCHAR(55);



UPDATE author SET firstname = (SELECT split_part(fullname, ' ', 1));
UPDATE author SET lastname = (SELECT split_part(fullname, ' ', 2));
UPDATE author SET middlename = (SELECT split_part(fullname, ' ', 3));

ALTER TABLE author DROP COLUMN fullname;