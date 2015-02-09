package dbr

import (
	// "database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionReal(t *testing.T) {
	s := createRealSessionWithFixtures()

	tx, err := s.Begin()
	assert.NoError(t, err)

	b := tx.InsertInto("dbr_people").Columns("name", "email").Values("Barack", "obama@whitehouse.gov")
	id, err := execAndGetID(b)

	assert.NoError(t, err)

	// TODO: RowsAffected isn't available for Postgres, because we Query rather than Exec
	// rowsAff, err := res.RowsAffected()
	// assert.NoError(t, err)

	assert.True(t, id > 0)
	// assert.Equal(t, rowsAff, 1)

	var person dbrPerson
	err = tx.Select("*").From("dbr_people").Where("id = ?", id).LoadStruct(&person)
	assert.NoError(t, err)

	assert.Equal(t, person.Id, id)
	assert.Equal(t, person.Name, "Barack")
	assert.Equal(t, person.Email.Valid, true)
	assert.Equal(t, person.Email.String, "obama@whitehouse.gov")

	err = tx.Commit()
	assert.NoError(t, err)
}

func TestTransactionRollbackReal(t *testing.T) {
	// Insert by specifying values
	s := createRealSessionWithFixtures()

	tx, err := s.Begin()
	assert.NoError(t, err)

	var person dbrPerson
	err = tx.Select("*").From("dbr_people").Where("email = ?", "jonathan@uservoice.com").LoadStruct(&person)
	assert.NoError(t, err)
	assert.Equal(t, person.Name, "Jonathan")

	err = tx.Rollback()
	assert.NoError(t, err)
}
