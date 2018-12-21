CREATE TABLE if not exists users(
    id serial PRIMARY KEY,
    email varchar(255) unique,
    nickname varchar(255) unique,
    age int,
    password text,
    token text,
    department text,
    super boolean default false,
    mod boolean default false,
    confirmed boolean default false,    
    created_at TIMESTAMP not null DEFAULT current_timestamp

);

CREATE TABLE if not exists messages(
    id serial PRIMARY key,
    nickname varchar(255) References users(nickname) on delete cascade,
    content text,
    created_at TIMESTAMP not null DEFAULT current_timestamp
)