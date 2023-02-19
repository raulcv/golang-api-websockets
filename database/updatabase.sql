drop table if exists users;
create table users (
  id varchar(32) primary key,
  password varchar(225) not null,
  email varchar(225) not null,
  created_at timestapm not null default now()
  modified_at timestapm null
  deleted_at timestapm null
)

drop table if exists posts;
create table posts (
  id varchar(32) primary key,
  content varchar(32) not null,
  user_id varchar(32) not null
  created_at timestapm not null default now()
  modified_at timestapm null
  deleted_at timestapm null
  foreign key(user_id) references users(id)
)