package domain

import "time"

// Bigtable entity of the bigtable
type Bigtable struct {
	Table string
	Rows  []*Row
}

// Row represent a row of the table
type Row struct {
	Key     string
	Columns []*Column
}

// Column represent a column of the row
type Column struct {
	Family    string
	Qualifier string
	Value     []byte
	Version   time.Time
}
