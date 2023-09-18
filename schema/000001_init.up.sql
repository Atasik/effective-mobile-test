CREATE TABLE persons
(
    id serial not null unique,
    name varchar(255) not null,
    surname varchar(255) not null,
    patronymic varchar(255),
    age int not null,
    gender varchar(255) not null,
    nationality varchar(255) not null
);