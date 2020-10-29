package psql

import (
	"fmt"
	"testing"
)

func TestReplaceDBName(t *testing.T) {

	inputs := []string{
		"postgresql://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca",
		"postgresql://jack:secret@pg.example.com:5432/mydb",
		"postgresql://jack:secret@pg.example.com:5432/",
		"postgresql://jack:secret@pg.example.com:5432/?a=b",
		"postgresql://jack:secret@pg.example.com:5432",
	}

	outputs := []string{
		"postgresql://jack:secret@pg.example.com:5432/newDB?sslmode=verify-ca",
		"postgresql://jack:secret@pg.example.com:5432/newDB",
		"postgresql://jack:secret@pg.example.com:5432/newDB",
		"postgresql://jack:secret@pg.example.com:5432/newDB?a=b",
		"postgresql://jack:secret@pg.example.com:5432/newDB",
	}

	for i := range inputs {
		newDsn, err := ReplaceDbName(inputs[i], "newDB")
		if (err != nil) {
			fmt.Println(err)
			t.Fail()
		}

		if (newDsn != outputs[i]) {
			fmt.Printf("FAIL interation %v\n", i)
			fmt.Println(inputs[i])
			fmt.Println(newDsn)
			t.Fail()
		}
	}
}
