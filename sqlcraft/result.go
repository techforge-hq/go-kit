package sqlcraft

// Result represents the result of a SQL query build operation, containing the SQL string and arguments.
type Result struct {
	SQL  string
	Args []any
}
