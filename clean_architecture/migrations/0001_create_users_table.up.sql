create table Users
(
    id            int generated always as identity primary key,
    email         varchar(255) unique                   not null,
    password_hash VARCHAR(255)                          NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
create index idx_users_email on Users(lower(email));