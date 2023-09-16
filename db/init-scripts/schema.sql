\c lil;

CREATE TABLE IF NOT EXISTS transaction (
                                           id serial PRIMARY KEY,
                                           name text,
                                           value int,
                                           user_id int,
                                           creation_date date
);

CREATE TABLE IF NOT EXISTS teleuser (
                                      id serial PRIMARY KEY,
                                      first_name text,
                                      last_name text
);