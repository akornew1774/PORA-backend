ALTER TABLE users
ADD CONSTRAINT users_phone_key UNIQUE (phone),
ADD CONSTRAINT users_email_key UNIQUE (email);