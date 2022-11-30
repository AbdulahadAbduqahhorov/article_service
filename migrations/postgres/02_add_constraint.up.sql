-- ALTER TABLE article ADD CONSTRAINT u_title_deleted_at UNIQUE NULLS NOT DISTINCT (title,deleted_at);

-- ALTER TABLE author  ADD CONSTRAINT u_firstname_lastname_deleted_at UNIQUE NULLS NOT DISTINCT (firstname,lastname,deleted_at);
