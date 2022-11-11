# IOiyn
before starting db run this sql statement in your mysql

createSchemaQuery := "create database if not exists game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"


then here in cmd/web/main.go
dsn := flag.String("dsn", "root:password"+
		"@/game?multiStatements=true&parseTime=true", "MySQL data source name")

instead of password write your password of user root in mysql
