# go-mcp-mysql

## Overview

A Model Context Protocol (MCP) server for interacting with MySQL and automation. This server provides tools to list, create, update, and delete MySQL databases and tables.

Please note that this is a work in progress and may not yet be ready for production use.

## Roadmap

- [ ] Implement database listing
- [ ] Implement table listing
- [ ] Implement general CRUD operations
- [ ] Implement read-only mode
- [ ] Implement table creation

## Tools

### Schema Tools

1. `list_database`

    - List all databases in the MySQL server.
    - Parameters: None
    - Returns: A list of matching database names.

2. `list_table`

    - List all tables in the MySQL server.
    - Parameters:
        - `name`: If provided, list tables with the specified name, same as SQL `SHOW TABLES LIKE '%name%'`. Otherwise, list all tables.
    - Returns: A list of matching table names.

3. `create_table`

    - Create a new table in the MySQL server.
    - Parameters:
        - `query`: The SQL query to create the table.
    - Returns: A confirmation message.

4. `desc_table`

    - Describe the structure of a table.
    - Parameters:
        - `name`: The name of the table to describe.
    - Returns: The structure of the table.

### Data Tools

1. `read_query`

    - Execute a read-only SQL query.
    - Parameters:
        - `query`: The SQL query to execute.
    - Returns: The result of the query.

2. `write_query`

    - Execute a write SQL query.
    - Parameters:
        - `query`: The SQL query to execute.
    - Returns: { affected_rows: number, insert_id: number, last_insert_id: number, rows: number, columns: number, err: string }.

3. `update_query`

    - Execute an update SQL query.
    - Parameters:
        - `query`: The SQL query to execute.
    - Returns: { affected_rows: number, insert_id: number, last_insert_id: number, rows: number, columns: number, err: string }.

4. `delete_query`

    - Execute a delete SQL query.
    - Parameters:
        - `query`: The SQL query to execute.
    - Returns: { affected_rows: number, insert_id: number, last_insert_id: number, rows: number, columns: number, err: string }.

## Installation

1. Get the latest [release](https://github.com/Zhwt/go-mcp-mysql/releases) and put it in your `$PATH`.

2. Or if you have Go installed, you can build it from source:

```sh
go install -v github.com/Zhwt/go-mcp-mysql@latest
```

## Usage

### Method A: Using Command Line Arguments

```json
{
  "mcpServers": {
    "mysql": {
      "command": "go-mcp-mysql",
      "args": [
        "-h", "localhost",
        "-u", "root",
        "-p", "password",
        "-P", "3306",
        "-d", "mydb"
      ]
    }
  }
}
```

### Method B: Using DSN With Custom Options

```json
{
  "mcpServers": {
    "mysql": {
      "command": "go-mcp-mysql",
      "args": [
        "--dsn", "username:password@tcp(localhost:3306)/mydb?parseTime=true&loc=Local"
      ]
    }
  }
}
```

Please refer to [MySQL DSN](https://github.com/go-sql-driver/mysql#dsn-data-source-name) for more details.

### Optional Flags

- Add a `--read-only` flag to enable read-only mode. In this mode, only tools beginning with `list`, `read_` and `desc_` are available. Other tool call will result in an immediate error.

## License

MIT
