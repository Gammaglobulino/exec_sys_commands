package cockroachlab

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestConnectToDb(t *testing.T) {

	db, err := sql.Open("postgres",
		"postgresql://root@192.168.177.1:26257/defaultdb")
	if err != nil {
		t.Fatal("error connecting to the database: ", err)
	}
	assert.NotNil(t, db)
	defer db.Close()
}
func TestExecuteSqlCommand(t *testing.T) {
	db, err := sql.Open("postgres",
		"postgresql://root@192.168.177.1:26257/bank?sslmode=disable")
	if err != nil {
		t.Fatal("error connecting to the database: ", err)
	}
	assert.NotNil(t, db)
	defer db.Close()

	// Create the "accounts" table.
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS accounts (id INT PRIMARY KEY, balance INT)"); err != nil {
		log.Fatal(err)
	}

	// Insert two rows into the "accounts" table.
	if _, err := db.Exec(
		"INSERT INTO accounts (id, balance) VALUES (1, 1000), (2, 250)"); err != nil {
		log.Fatal(err)
	}

	// Print out the balances.
	rows, err := db.Query("SELECT id, balance FROM accounts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("Initial balances:")
	for rows.Next() {
		var id, balance int
		if err := rows.Scan(&id, &balance); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d %d\n", id, balance)
	}

}
