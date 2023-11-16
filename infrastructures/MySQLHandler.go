package infrastructures

import (
	"database/sql"
	"fmt"
)

type MySQLHandler struct {
	Conn *sql.DB
}

func (handler *MySQLHandler) SetConn(Connection *sql.DB) {
	handler.Conn = Connection
}

func (handler *MySQLHandler) Execute(statement string, args ...any) (sql.Result, error) {

	return handler.Conn.Exec(statement, args...)
}

func (handler *MySQLHandler) Query(statement string, args ...any) (*sql.Rows, error) {
	rows, err := handler.Conn.Query(statement, args...)

	if err != nil {
		fmt.Println(err)
		return rows, err
	}

	return rows, nil
}

func (handler *MySQLHandler) QueryRow(statement string, args ...any) *sql.Row {
	row := handler.Conn.QueryRow(statement, args...)

	return row
}
