package dbr

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

//
// Test helpers
//

// Returns a session that's not backed by a database
func createFakeSession() *Session {
	cxn := NewConnection(nil, nil)
	return cxn.NewSession(nil)
}

func createRealSession() *Session {
	cxn := NewConnection(realDb(), nil)
	return cxn.NewSession(nil)
}

func createRealSessionWithFixtures() *Session {
	sess := createRealSession()
	installFixtures(sess.cxn.Db)
	return sess
}

func quoteSQL(sqlFmt string, cols ...string) string {
	args := make([]interface{}, len(cols))

	for i := range cols {
		args[i] = quoteColumn(cols[i])
	}

	return fmt.Sprintf(sqlFmt, args...)
}

func realDb() *sql.DB {
	driver := os.Getenv("DBR_TEST_DRIVER")
	if driver == "" {
		driver = "mysql"
	}

	dsn := os.Getenv("DBR_TEST_DSN")
	if dsn == "" {
		dsn = "root:unprotected@unix(/tmp/mysql.sock)/uservoice_development?charset=utf8&parseTime=true"
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalln("Mysql error ", err)
	}

	return db
}

type dbrPerson struct {
	Id    int64
	Name  string
	Email NullString
	Key   NullString
}

func installFixtures(db *sql.DB) {
	createTablePeople := `
		CREATE TABLE dbr_people (
			id SERIAL PRIMARY KEY,
			name varchar(255) NOT NULL,
			email varchar(255),
			key varchar(255)
		)
	`

	sqlToRun := []string{
		"DROP TABLE IF EXISTS dbr_people",
		createTablePeople,
		"INSERT INTO dbr_people (name,email) VALUES ('Jonathan', 'jonathan@uservoice.com')",
		"INSERT INTO dbr_people (name,email) VALUES ('Dmitri', 'zavorotni@jadius.com')",
	}

	for _, v := range sqlToRun {
		_, err := db.Exec(v)
		if err != nil {
			log.Fatalln("Failed to execute statement: ", v, " Got error: ", err)
		}
	}
}
