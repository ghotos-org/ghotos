package migrations

import (
	"database/sql"
	"os"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(Up_20211031022801, Down_20211031022801)
}

func Up_20211031022801(txn *sql.Tx) error {
	path := "./data"
	if os.Getenv("MOUNT_DEFAULT") != "" {
		path = os.Getenv("MOUNT_DEFAULT")
	}
	_, err := txn.Exec("INSERT INTO mounts (active,path,type) VALUE(1,'" + path + "', 1);")
	return err

}

func Down_20211031022801(txn *sql.Tx) error {
	path := "./data"
	if os.Getenv("MOUNT_DEFAULT") != "" {
		path = os.Getenv("MOUNT_DEFAULT")
	}
	_, err := txn.Exec("DELETE FROM mounts WHERE active = 1 AND path = '" + path + "';")
	return err
}
