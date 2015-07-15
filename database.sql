create table events (ts timestamp with time zone not null, proto varchar(8) not null, src inet not null, dst inet not null, sport int not null, dport int not null);

CREATE EXTENSION dblink;

CREATE SCHEMA jobmon;
CREATE EXTENSION pg_jobmon SCHEMA jobmon;
INSERT INTO jobmon.dblink_mapping_jobmon (username, pwd) VALUES ('conntrack', 'conntrack');
grant usage on schema jobmon to conntrack;
grant usage on schema public to conntrack;
grant select, insert, update, delete on all tables in schema jobmon to conntrack;
grant execute on all functions in schema jobmon to conntrack;
grant all on all sequences in schema jobmon to conntrack;

CREATE SCHEMA partman;
CREATE EXTENSION pg_partman SCHEMA partman;

grant all on all tables in schema public to conntrack;
SELECT partman.create_parent('public.events', 'ts', 'time', 'daily');
UPDATE partman.part_config SET retention = '30 days', retention_keep_table = false WHERE parent_table = 'public.events';




