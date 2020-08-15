CREATE EXTENSION "uuid-ossp";

CREATE TABLE posts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    userid uuid,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    cheer INTEGER,
    created_date TIMESTAMP DEFAULT current_timestamp,
    update_date TIMESTAMP DEFAULT current_timestamp
);

create table users (
    id uuid primary key default uuid_generate_v4(),
    username varchar(20) not null,
    email text not null,
    password text not null,
    created_date timestamp default current_timestamp
);
