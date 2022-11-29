package crudsql

// SQLDialect denotes the different dialects which define placeholders differently.
type SQLDialect int

const (
	MySQL SQLDialect = iota
	PostgreSQL
	Oracle
	SQLite
	ODBC
	MariaDB
)
