create table if not exists datafiles (
    id integer primary key autoincrement,
    path text unique,
    size integer,
    indexpath text unique,
    indexed boolean default false
);

create table if not exists pages (
    id integer primary key,
    streamid integer not null,
    foreign key(streamid) references streams(id)
);

create table if not exists streams (
    id integer primary key autoincrement,
    bytebegin integer,
    byteend integer,
    fileid integer not null,
    foreign key(fileid) references datafiles(id),
    unique(bytebegin,fileid)
);
