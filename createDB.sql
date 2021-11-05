create table if not exists Archive (
    ID integer primary key autoincrement,
    FilePath text unique,
    FileSize integer,
    IndexPath text unique,
    Processed boolean default false
);

create table if not exists Page (
    ID integer primary key,
    StreamID integer not null,
    foreign key(StreamID) references Stream(id)
);

create table if not exists Stream (
    ID integer primary key autoincrement,
    ByteBegin integer,
    ByteEnd integer,
    ArchiveID integer not null,
    foreign key(ArchiveID) references Archive(ID),
    unique(ByteBegin,ArchiveID)
);
