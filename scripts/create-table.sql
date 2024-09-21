CREATE TABLE IF NOT EXISTS account (
	account_id INTEGER PRIMARY KEY AUTOINCREMENT,
	account_name TEXT NOT NULL UNIQUE,
	account_password TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trg_account_upd AFTER UPDATE ON account
BEGIN
    UPDATE account
    SET updated_at = DATETIME('now', 'localtime') 
    WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS account_profile (
	account_id INTEGER NOT NULL UNIQUE,
	group_id INTEGER,
	account_role TEXT NOT NULL,
	display_name TEXT,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trg_account_profile_upd AFTER UPDATE ON account_profile
BEGIN
	UPDATE account_profile 
	SET updated_at = DATETIME('now', 'localtime') 
	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS "group" (
	group_id INTEGER PRIMARY KEY AUTOINCREMENT,
	group_name TEXT NOT NULL UNIQUE,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trg_shift_group_upd AFTER UPDATE ON "group"
BEGIN
	UPDATE "group" 
	SET updated_at = DATETIME('now', 'localtime') 
	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS shift (
	group_id INTEGER,
	year INTEGER,
	month INTEGER,
	store_holiday TEXT,
	shift_data TEXT,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	PRIMARY KEY(group_id, shift_year, shift_month)
);

CREATE TRIGGER IF NOT EXISTS trg_shift_upd AFTER UPDATE ON shift
BEGIN
	UPDATE shift
	SET UPDATE_AT = DATETIME('now', 'localtime') 
	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS available_work_days (
	account_id INTEGER,
	year INTEGER,
	month INTEGER,
	available_work_days TEXT,
	memo TEXT,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	PRIMARY KEY(account_id, year, month)
);

CREATE TRIGGER IF NOT EXISTS trg_available_work_days_upd AFTER UPDATE ON available_work_days
BEGIN
	UPDATE available_work_days
		SET updated_at = DATETIME('now', 'localtime') 
		WHERE rowid == NEW.rowid;
END;