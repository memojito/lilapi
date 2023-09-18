\c lil;

CREATE TABLE IF NOT EXISTS transaction (
                                           id serial PRIMARY KEY,
                                           name text,
                                           value int,
                                           user_id int,
                                           creation_date date,
                                           category_id int
);

CREATE TABLE IF NOT EXISTS teleuser (
                                      id serial PRIMARY KEY,
                                      first_name text,
                                      last_name text
);

CREATE TABLE IF NOT EXISTS category (
                                        id serial PRIMARY KEY,
                                        name text,
                                        user_id int
);
