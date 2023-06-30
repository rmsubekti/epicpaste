-- this query will be executed after gorm table migration

----- start creating field trigger -----
-- CREATE OR REPLACE FUNCTION FieldID() RETURNS TRIGGER AS $$
-- 	declare
-- 		rmnum text := right(concat('0000',(floor(random() * 9999)::int)::text),4);
-- 	BEGIN
-- 		new.id := UPPER(concat('F-', (new.grid), '-', rmnum::varchar, '-',rmluhn(concat(new.grid, rmnum::varchar)))::varchar);
-- 		RETURN NEW;
-- 	END;
-- $$ LANGUAGE plpgsql;

-- DROP TRIGGER IF EXISTS generateFieldID ON master.field;
-- CREATE TRIGGER generateFieldID BEFORE insert ON master.field FOR EACH ROW EXECUTE PROCEDURE FieldID();
-- ----- end field trigger------
