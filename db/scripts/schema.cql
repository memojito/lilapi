USE lilapi;

CREATE TABLE IF NOT EXISTS transaction (
                                     id uuid PRIMARY KEY,
                                     name text,
                                     value int
);
