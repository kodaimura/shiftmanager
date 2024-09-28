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
	account_id INTEGER PRIMARY KEY,
	account_role TEXT NOT NULL,
	display_name TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trg_account_profile_upd AFTER UPDATE ON account_profile
BEGIN
	UPDATE account_profile 
	SET updated_at = DATETIME('now', 'localtime') 
	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS shift (
	year INTEGER,
	month INTEGER,
	store_holiday TEXT,
	shift_data TEXT,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	PRIMARY KEY(year, month)
);

CREATE TRIGGER IF NOT EXISTS trg_shift_upd AFTER UPDATE ON shift
BEGIN
	UPDATE shift
	SET UPDATE_AT = DATETIME('now', 'localtime') 
	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS shift_preferred (
	account_id INTEGER,
	year INTEGER,
	month INTEGER,
	dates TEXT,
	notes TEXT,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	PRIMARY KEY(account_id, year, month)
);

CREATE TRIGGER IF NOT EXISTS trg_shift_preferred_upd AFTER UPDATE ON shift_preferred
BEGIN
	UPDATE shift_preferred
	SET updated_at = DATETIME('now', 'localtime') 
	WHERE rowid == NEW.rowid;
END;