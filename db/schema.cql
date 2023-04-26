CREATE TABLE IF NOT EXISTS workspace (
                                     id uuid PRIMARY KEY,
                                     owner text,
                                     editor text,
                                     size int,
                                     name text,
                                     facet_ids set<uuid>
);

CREATE TABLE IF NOT EXISTS facet (
                                     id uuid PRIMARY KEY,
                                     name text,
                                     value text,
                                     workspace_ids set<uuid>
);