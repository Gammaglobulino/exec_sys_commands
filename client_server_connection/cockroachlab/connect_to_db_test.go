package cockroachlab

import (
	"../core/handle_connections"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
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

	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	host := localip + ":" + "26257"

	dbName := "bank"
	user := "root"
	sslMode := false

	connectToDB := strings.Builder{}
	connectToDB.WriteString("postgresql://")
	connectToDB.WriteString(user)
	connectToDB.WriteString("@")
	connectToDB.WriteString(host)
	connectToDB.WriteString("/")
	connectToDB.WriteString(dbName)

	if sslMode == false {
		connectToDB.WriteString("?sslmode=disable")
	}
	db, err := sql.Open("postgres",
		connectToDB.String())
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
		log.Println(err)
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
