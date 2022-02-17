CREATE TABLE cats (
    ID UUID,
    Name varchar(120)
);

CREATE TABLE users (
    ID SERIAL PRIMARY KEY,
    Name varchar(120) NOT NULL,
    Username varchar(120),
    Password varchar(120)
);

CREATE SEQUENCE users_sequence;
