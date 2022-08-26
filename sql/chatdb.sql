drop table if exists users;
create table users(
    id int unsigned not null auto_increment,
    created_at datetime(3) null ,
    updated_at datetime(3) null on update current_timestamp(3),
    deleted_at datetime(3) null ,
    username varchar(20) not null,
    password varchar(20) not null,
    primary key (id),
    unique index idx_users_username(username)
);
drop table if exists rooms;
create table rooms(
    id int unsigned not null auto_increment,
    created_at datetime(3) null ,
    updated_at datetime(3) null on update current_timestamp(3),
    deleted_at datetime(3) null ,
    user_id int unsigned not null,
    room_name varchar(20) not null,
    primary key (id),
    index idx_rooms_room_name(room_name)
);
drop table if exists messages;
create table messages(
    id int unsigned not null auto_increment,
    created_at datetime(3) null ,
    updated_at datetime(3) null on update current_timestamp(3),
    deleted_at datetime(3) null ,
    user_id int unsigned not null,
    room_id int unsigned not null,
    content longtext null,
    primary key (id),
    index idx_messages_user_id(user_id)
);