create table events (ts timestamp with time zone not null, proto varchar(8) not null, src inet not null, dst inet not null, sport int not null, dport int not null);