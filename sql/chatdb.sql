drop table if exists chat_history;
create table chat_history(
   id int unsigned not null auto_increment,
   created_at datetime(3) null ,
   updated_at datetime(3) null on update current_timestamp(3),
   deleted_at datetime(3) null ,
   user_id varchar(32) not null,
   room_id varchar(64) not null,
   to_user_id varchar(32) not null,
   content longtext,
   image_url varchar(64),
   primary key (id),
   index idx_chat_history_room_id(room_id)
);
drop table if exists user_info;
create table user_info(
  id int unsigned not null auto_increment,
  created_at datetime(3) null ,
  updated_at datetime(3) null on update current_timestamp(3),
  deleted_at datetime(3) null ,
  username varchar(20) not null,
  password varchar(20) not null,
  primary key (id),
  unique index idx_user_info_username(username)
);
drop table if exists room_info;
create table room_info(
  id int unsigned not null auto_increment,
  created_at datetime(3) null ,
  updated_at datetime(3) null on update current_timestamp(3),
  deleted_at datetime(3) null ,
  room_name varchar(20) not null,
  user_id varchar(32) not null,
  username varchar(20) not null,
  primary key (id),
  index idx_room_info_room_name(room_name)
);