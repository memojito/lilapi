USE lilapi;

CREATE TABLE IF NOT EXISTS transaction (
                                     id uuid PRIMARY KEY,
                                     name text,
                                     value int,
                                     user_id int
);

CREATE TABLE IF NOT EXISTS user (
                                           id int PRIMARY KEY,
                                           first_name text,
                                           last_name text
);
