-- this query will be executed before gorm table migration


-------- enable uuid extension -----------
--CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-------- end uuid extension -----------

------ create schemas ---------------
CREATE SCHEMA IF NOT EXISTS "master";
CREATE SCHEMA IF NOT EXISTS "user";
------ end create schema ------------