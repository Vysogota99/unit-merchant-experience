CREATE TABLE IF NOT EXISTS offers (
    id int PRIMARY KEY,
    saler_id int NOT NULL,
    name varchar(255) NOT NULL,
    price numeric NOT NULL,
    quantity int NOT NULL
);

CREATE INDEX offers_saler_id_index ON offers(saler_id);
CREATE INDEX offers_name_index ON offers(name);
CREATE INDEX offers_ns_index ON offers(saler_id, name);