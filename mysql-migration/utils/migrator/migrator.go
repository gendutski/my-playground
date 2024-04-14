package migrator

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/gendutski/my-playground/mysql-migration/utils/readdir"
	"github.com/gendutski/my-playground/mysql-migration/utils/terminal"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/term"
)

const (
	defaultMysqlHost      string = "localhost"
	defaultMysqlPort      string = "3306"
	defaultMysqlDBName    string = "mbpro"
	defaultMysqlCharset   string = "utf8mb4"
	defaultMysqlLoc       string = "Local"
	defaultMysqlParseTime bool   = true
	defaultMysqlLogMode   int    = 0

	mysqlHostString       string = "MYSQL_HOST"
	mysqlPortString       string = "MYSQL_PORT"
	mysqlUserNameString   string = "MYSQL_USERNAME"
	mysqlDBNameString     string = "MYSQL_DB_NAME"
	mysqlPasswordString   string = "MYSQL_PASSWORD"
	mysqlDBNotExistsErrNo uint16 = 1049
	migrationsFolder      string = "migrations"
)

func Run(r *bufio.Reader) (map[string]string, error) {
	var mysqlHost, mysqlPort, mysqlUserName, dbName string
	var pass []byte
	var err error

	mysqlHost, err = terminal.ReadString(fmt.Sprintf("Enter mysql host (default %s): ", defaultMysqlHost), defaultMysqlHost, r)
	if err != nil {
		return nil, err
	}

	mysqlPort, err = terminal.ReadString(fmt.Sprintf("Enter mysql port (default %s): ", defaultMysqlPort), defaultMysqlPort, r)
	if err != nil {
		return nil, err
	}

	mysqlUserName, err = terminal.ReadString("Enter mysql user name: ", "", r)
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter mysql user password: ")
	pass, err = term.ReadPassword(syscall.Stdin)
	if err != nil {
		return nil, err
	}
	password := strings.TrimSpace(string(pass))

	dbName, err = terminal.ReadString(fmt.Sprintf("\nEnter database name (default %s): ", defaultMysqlDBName), defaultMysqlDBName, r)
	if err != nil {
		return nil, err
	}

	// validate mysql connection
	err = testMysqlConnection(mysqlUserName, password, mysqlHost, mysqlPort, dbName, r)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		mysqlHostString:     mysqlHost,
		mysqlPortString:     mysqlPort,
		mysqlUserNameString: mysqlUserName,
		mysqlDBNameString:   dbName,
		mysqlPasswordString: password,
	}, nil
}

func testMysqlConnection(userName, password, host, port, dbName string, r *bufio.Reader) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", userName, password, host, port))
	if err != nil {
		return err
	}
	defer db.Close()

	err = useDB(db, dbName)
	if err != nil {
		if driveErr, ok := err.(*mysql.MySQLError); ok && driveErr.Number == mysqlDBNotExistsErrNo {
			if terminal.Confirm("Database not exists, create one (y/n)? ", r) {
				return migrate(true, dbName, db)
			}
		}
		return err
	}

	return migrate(false, dbName, db)
}

func useDB(db *sql.DB, dbName string) error {
	log.Println("execute: USE " + dbName)
	_, err := db.Exec("USE " + dbName)
	return err
}

func createDB(db *sql.DB, dbName string) error {
	log.Println("execute: CREATE DATABASE " + dbName)
	_, err := db.Exec("CREATE DATABASE " + dbName)
	return err
}

func migrate(autoCreateDB bool, dbName string, db *sql.DB) error {
	var err error
	if autoCreateDB {
		err = createDB(db, dbName)
		if err != nil {
			return err
		}
		err = useDB(db, dbName)
		if err != nil {
			return err
		}
	}

	// get all migration files
	sqls, err := readdir.Run(migrationsFolder, true)
	if err != nil {
		return err
	}

	for _, sql := range sqls {
		log.Println("Read:", sql.FileName)
		_, err = db.Exec(sql.Content)
		if err != nil {
			return err
		}
	}
	return nil
}
