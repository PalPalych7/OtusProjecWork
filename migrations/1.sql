create database otusProj;
create user otus1 with encrypted password '123456';
grant all privileges on database otusProj to otus1;
create table banner (
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    descr text
);
create table slot (
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    descr text
);

create table soc_group (
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    descr text
);

insert into banner (descr) values ('banner 1') ;
insert into banner (descr) values ('banner 2') ;
insert into banner (descr) values ('banner 3') ;
insert into banner (descr) values ('banner 4') ;
insert into banner (descr) values ('banner 5') ;
insert into banner (descr) values ('banner 6') ;
insert into banner (descr) values ('banner 7') ;
insert into banner (descr) values ('banner 8') ;
insert into banner (descr) values ('banner 9') ;
insert into banner (descr) values ('banner 10') ;
insert into banner (descr) values ('banner 11') ;
insert into banner (descr) values ('banner 12') ;
insert into banner (descr) values ('banner 13') ;
insert into banner (descr) values ('banner 14') ;
insert into banner (descr) values ('banner 15') ;
insert into banner (descr) values ('banner 16') ;
insert into banner (descr) values ('banner 17') ;
insert into banner (descr) values ('banner 18') ;
insert into banner (descr) values ('banner 18') ;
insert into banner (descr) values ('banner 20') ;
