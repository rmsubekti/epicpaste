package model

import (
	"epicpaste/system/config"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = config.GetDB()

	// init extensions, schemas, and functions
	db.Exec(`
	-------- enable uuid extension -----------
	--CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	-------- end uuid extension -----------

	------ create schemas ---------------
	CREATE SCHEMA IF NOT EXISTS "master";
	CREATE SCHEMA IF NOT EXISTS "user";
	------ end create schema ------------
	`)

	db.AutoMigrate(
		&Account{},
		&User{},
		&Tag{},
		&Category{},
		&Application{},
		&OperatingSystem{},
		&Language{},
		&Paste{},
		&CodeDetail{},
	)

	// init dummy data
	createDummy(db)

	// init table function and trigger relations
	// 	db.Exec(`

	// ----- start creating field trigger -----
	// CREATE OR REPLACE FUNCTION FieldID() RETURNS TRIGGER AS $$
	// 	declare
	// 		rmnum text := right(concat('0000',(floor(random() * 9999)::int)::text),4);
	// 	BEGIN
	// 		new.id := UPPER(concat('F-', (new.grid), '-', rmnum::varchar, '-',rmluhn(concat(new.grid, rmnum::varchar)))::varchar);
	// 		RETURN NEW;
	// 	END;
	// $$ LANGUAGE plpgsql;

	// DROP TRIGGER IF EXISTS generateFieldID ON master.field;
	// CREATE TRIGGER generateFieldID BEFORE insert ON master.field FOR EACH ROW EXECUTE PROCEDURE FieldID();
	// ----- end field trigger------

	// 	`)
}
