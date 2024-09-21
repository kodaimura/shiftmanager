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