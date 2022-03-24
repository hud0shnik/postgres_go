CREATE TABLE names
(
    id               serial       not null unique,
    name             varchar(255) not null,
    meaning          varchar(1500) not null,
    gender           varchar(255) not null,
    origin           varchar(255) not null,
    PeoplesCount     int not null,
    WhenPeoplesCount varchar(255) not null
);