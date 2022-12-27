USE test_snippetbox;

create table snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

create index idx_snippets_created on snippets(created);

create table users (
    id integer not null primary key auto_increment,
    name varchar(255) not null,
    email varchar(255) not null,
    hashed_password char(60) not null,
    created DATETIME not null
);

alter table users add constraint  users_uc_email unique (email);

insert into users (name, email, hashed_password, created) values (
    'Alice Jones',
    'alice@exmaple.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 10:00:00'
)