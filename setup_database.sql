
create database weather;
use weather;

create table user (
  name varchar(255) primary key,
  salt varbinary(20) not null,
  bcrypt varchar(255) not null);
    
create table session (
  token varchar(255) primary key,
  username varchar(255) not null,
  deleted datetime,
  expires datetime not null,
  constraint foreign key (username) references user(name));

create table stream_token (
  token varchar(255) primary key,
  username varchar(255) not null,
  created datetime not null,
  deleted datetime,
  constraint foreign key (username) references user(name));

