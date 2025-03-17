package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	Host string
	User string
	Pass string
	Port int
	Db   string

	DSN string

	ReadOnly bool
)

func main() {
	flag.StringVar(&Host, "host", "localhost", "MySQL hostname")
	flag.StringVar(&User, "user", "root", "MySQL username")
	flag.StringVar(&Pass, "pass", "", "MySQL password")
	flag.IntVar(&Port, "port", 3306, "MySQL port")
	flag.StringVar(&Db, "db", "", "MySQL database")

	flag.StringVar(&DSN, "dsn", "", "MySQL DSN")

	flag.BoolVar(&ReadOnly, "read-only", true, "Enable read-only mode")
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
		mcp.WithString("name",
			mcp.Description("If provided, list tables with the specified name"),
		),
	)

	createTableTool := mcp.NewTool(
		"create_table",
		mcp.WithDescription("Create a new table in the MySQL server"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to create the table"),
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
		mcp.WithDescription("Execute a read-only SQL query"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to execute"),
		),
	)

	writeQueryTool := mcp.NewTool(
		"write_query",
		mcp.WithDescription("Execute a write SQL query"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to execute"),
		),
	)

	updateQueryTool := mcp.NewTool(
		"update_query",
		mcp.WithDescription("Execute an update SQL query"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to execute"),
		),
	)

	deleteQueryTool := mcp.NewTool(
		"delete_query",
		mcp.WithDescription("Execute a delete SQL query"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The SQL query to execute"),
		),
	)

	// Add tools with placeholder handlers
	s.AddTool(listDatabaseTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultError("Not implemented"), nil
	})

	s.AddTool(listTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultError("Not implemented"), nil
	})

	s.AddTool(createTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if ReadOnly {
			return mcp.NewToolResultError("Table creation is disabled in read-only mode"), nil
		}

		return mcp.NewToolResultError("Not implemented"), nil
	})

	s.AddTool(descTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultError("Not implemented"), nil
	})

	s.AddTool(readQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultError("Not implemented"), nil
	})

	s.AddTool(writeQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if ReadOnly {
			return mcp.NewToolResultError("Write queries are disabled in read-only mode"), nil
		}

		return mcp.NewToolResultError("Not implemented"), nil
	})

	s.AddTool(updateQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if ReadOnly {
			return mcp.NewToolResultError("Update queries are disabled in read-only mode"), nil
		}

		return mcp.NewToolResultError("Not implemented"), nil
	})

	s.AddTool(deleteQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if ReadOnly {
			return mcp.NewToolResultError("Delete queries are disabled in read-only mode"), nil
		}

		return mcp.NewToolResultError("Not implemented"), nil
	})

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
