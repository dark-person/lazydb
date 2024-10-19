# LazyDB
Go SQLite collection for lazy people.

Include these feature:
- Create SQLite database file
- Migration with `embed.fs`
- Backup when migration

Note: You must has `CGO` enabled to compile this project.

## Get Started

An example as:

```go
//go:embed all:schema
var schema embed.FS

func main() {
    // Init db
    db := lazydb.New(
        lazydb.DbPath("path/data.db"),
	    lazydb.Migrate(schema, "schema"),
	    lazydb.BackupDir("./backup"),
    )

    // Connect to db, which will create file if necessary
	err = db.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Migration performed
	backup, err := db.Migrate()
	if err != nil {
		panic(err)
	}

	if backup != "" {
		fmt.Println("Auto backup as migration performed: " + backup)
	}

    // Usage here...
}
```