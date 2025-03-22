package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	StatementTypeNoExplainCheck = ""
	StatementTypeSelect         = "SELECT"
	StatementTypeInsert         = "INSERT"
	StatementTypeUpdate         = "UPDATE"
	StatementTypeDelete         = "DELETE"
)

var (
	Host string
	User string
	Pass string
	Port int
	Db   string

	DSN string

	ReadOnly         bool
	WithExplainCheck bool

	DB *sqlx.DB
)

type ExplainResult struct {
	Id           *string `db:"id"`
	SelectType   *string `db:"select_type"`
	Table        *string `db:"table"`
	Partitions   *string `db:"partitions"`
	Type         *string `db:"type"`
	PossibleKeys *string `db:"possible_keys"`
	Key          *string `db:"key"`
	KeyLen       *string `db:"key_len"`
	Ref          *string `db:"ref"`
	Rows         *string `db:"rows"`
	Filtered     *string `db:"filtered"`
	Extra        *string `db:"Extra"`
}

type ShowCreateTableResult struct {
	Table       string `db:"Table"`
	CreateTable string `db:"Create Table"`
}

func main() {
	flag.StringVar(&Host, "host", "localhost", "MySQL hostname")
	flag.StringVar(&User, "user", "root", "MySQL username")
	flag.StringVar(&Pass, "pass", "", "MySQL password")
	flag.IntVar(&Port, "port", 3306, "MySQL port")
	flag.StringVar(&Db, "db", "", "MySQL database")

	flag.StringVar(&DSN, "dsn", "", "MySQL DSN")

	flag.BoolVar(&ReadOnly, "read-only", false, "Enable read-only mode")
	flag.BoolVar(&WithExplainCheck, "with-explain-check", false, "Check query plan with `EXPLAIN` before executing")
	flag.Parse()

	if len(DSN) == 0 {
		DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", User, Pass, Host, Port, Db)
	}

	s := server.NewMCPServer(
		"go-mcp-mysql",
		"0.1.0",
	)

	// Schema Tools
	listDatabaseTool := mcp.NewTool(
		"list_database",
		mcp.WithDescription("List all databases in the MySQL server"),
	)

	listTableTool := mcp.NewTool(
		"list_table",
		mcp.WithDescription("List all tables in the MySQL server"),
	)

	createTableTool := mcp.NewTool(
		"create_table",
		mcp.WithDescription("Create a new table in the MySQL server. Make sure you have added proper comments for each column and the table itself"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to create the table"),
		),
	)

	alterTableTool := mcp.NewTool(
		"alter_table",
		mcp.WithDescription("Alter an existing table in the MySQL server. Make sure you have updated comments for each modified column. DO NOT drop table or existing columns!"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to alter the table"),
		),
	)

	descTableTool := mcp.NewTool(
		"desc_table",
		mcp.WithDescription("Describe the structure of a table"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the table to describe"),
		),
	)

	// Data Tools
	readQueryTool := mcp.NewTool(
		"read_query",
		mcp.WithDescription("Execute a read-only SQL query. Make sure you have knowledge of the table structure before writing WHERE conditions. Call `desc_table` first if necessary"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to execute"),
		),
	)

	writeQueryTool := mcp.NewTool(
		"write_query",
		mcp.WithDescription("Execute a write SQL query. Make sure you have knowledge of the table structure before executing the query. Make sure the data types match the columns' definitions"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to execute"),
		),
	)

	updateQueryTool := mcp.NewTool(
		"update_query",
		mcp.WithDescription("Execute an update SQL query. Make sure you have knowledge of the table structure before executing the query. Make sure there is always a WHERE condition. Call `desc_table` first if necessary"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to execute"),
		),
	)

	deleteQueryTool := mcp.NewTool(
		"delete_query",
		mcp.WithDescription("Execute a delete SQL query. Make sure you have knowledge of the table structure before executing the query. Make sure there is always a WHERE condition. Call `desc_table` first if necessary"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to execute"),
		),
	)

	s.AddTool(listDatabaseTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := HandleQuery("SHOW DATABASES", StatementTypeNoExplainCheck)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(listTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := HandleQuery("SHOW TABLES", StatementTypeNoExplainCheck)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(result), nil
	})

	if !ReadOnly {
		s.AddTool(createTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := HandleExec(request.Params.Arguments["query"].(string), StatementTypeNoExplainCheck)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(result), nil
		})
	}

	if !ReadOnly {
		s.AddTool(alterTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := HandleExec(request.Params.Arguments["query"].(string), StatementTypeNoExplainCheck)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(result), nil
		})
	}

	s.AddTool(descTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := HandleDescTable(request.Params.Arguments["name"].(string))
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(readQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := HandleQuery(request.Params.Arguments["query"].(string), StatementTypeSelect)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(result), nil
	})

	if !ReadOnly {
		s.AddTool(writeQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := HandleExec(request.Params.Arguments["query"].(string), StatementTypeInsert)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(result), nil
		})
	}

	if !ReadOnly {
		s.AddTool(updateQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := HandleExec(request.Params.Arguments["query"].(string), StatementTypeUpdate)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(result), nil
		})
	}

	if !ReadOnly {
		s.AddTool(deleteQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := HandleExec(request.Params.Arguments["query"].(string), StatementTypeDelete)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(result), nil
		})
	}

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func GetDB() (*sqlx.DB, error) {
	if DB != nil {
		return DB, nil
	}

	db, err := sqlx.Connect("mysql", DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to establish database connection: %v", err)
	}

	DB = db

	return DB, nil
}

func HandleQuery(query, expect string) (string, error) {
	result, headers, err := DoQuery(query, expect)
	if err != nil {
		return "", err
	}

	s, err := MapToCSV(result, headers)
	if err != nil {
		return "", err
	}

	return s, nil
}

func DoQuery(query, expect string) ([]map[string]interface{}, []string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, nil, err
	}

	if len(expect) > 0 {
		if err := HandleExplain(query, expect); err != nil {
			return nil, nil, err
		}
	}

	rows, err := db.Queryx(query)
	if err != nil {
		return nil, nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	result := []map[string]interface{}{}
	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			return nil, nil, err
		}

		resultRow := map[string]interface{}{}
		for i, col := range cols {
			switch v := row[i].(type) {
			case []byte:
				resultRow[col] = string(v)
			default:
				resultRow[col] = v
			}
		}
		result = append(result, resultRow)
	}

	return result, cols, nil
}

func HandleExec(query, expect string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	if len(expect) > 0 {
		if err := HandleExplain(query, expect); err != nil {
			return "", err
		}
	}

	result, err := db.Exec(query)
	if err != nil {
		return "", err
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	switch expect {
	case StatementTypeInsert:
		li, err := result.LastInsertId()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%d rows affected, last insert id: %d", ra, li), nil
	default:
		return fmt.Sprintf("%d rows affected", ra), nil
	}
}

func HandleExplain(query, expect string) error {
	if !WithExplainCheck {
		return nil
	}

	db, err := GetDB()
	if err != nil {
		return err
	}

	rows, err := db.Queryx(fmt.Sprintf("EXPLAIN %s", query))
	if err != nil {
		return err
	}

	result := []ExplainResult{}
	for rows.Next() {
		var row ExplainResult
		if err := rows.StructScan(&row); err != nil {
			return err
		}
		result = append(result, row)
	}

	if len(result) != 1 {
		return fmt.Errorf("unable to check query plan, denied")
	}

	match := false
	switch expect {
	case StatementTypeInsert:
		fallthrough
	case StatementTypeUpdate:
		fallthrough
	case StatementTypeDelete:
		if *result[0].SelectType == expect {
			match = true
		}
	default:
		// for SELECT type query, the select_type will be multiple values
		// here we check if it's not INSERT, UPDATE or DELETE
		match = true
		for _, typ := range []string{StatementTypeInsert, StatementTypeUpdate, StatementTypeDelete} {
			if *result[0].SelectType == typ {
				match = false
				break
			}
		}
	}

	if !match {
		return fmt.Errorf("query plan does not match expected pattern, denied")
	}

	return nil
}

func HandleDescTable(name string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Queryx(fmt.Sprintf("SHOW CREATE TABLE %s", name))
	if err != nil {
		return "", err
	}

	result := []ShowCreateTableResult{}
	for rows.Next() {
		var row ShowCreateTableResult
		if err := rows.StructScan(&row); err != nil {
			return "", err
		}
		result = append(result, row)
	}

	if len(result) == 0 {
		return "", fmt.Errorf("table %s does not exist", name)
	}

	return result[0].CreateTable, nil
}

func MapToCSV(m []map[string]interface{}, headers []string) (string, error) {
	var csvBuf strings.Builder
	writer := csv.NewWriter(&csvBuf)

	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write headers: %v", err)
	}

	for _, item := range m {
		row := make([]string, len(headers))
		for i, header := range headers {
			value, exists := item[header]
			if !exists {
				return "", fmt.Errorf("key '%s' not found in map", header)
			}
			row[i] = fmt.Sprintf("%v", value)
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write row: %v", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("error flushing CSV writer: %v", err)
	}

	return csvBuf.String(), nil
}
