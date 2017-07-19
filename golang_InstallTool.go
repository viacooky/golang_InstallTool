package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-adodb"
)

var toolTypeDoc string = `
0:do noting
1:TestMsSQLConnect return: 1 can not connect 
2:TestMsSQLInstanceExists return: 2 Exists
`
var inputToolType = flag.Int("type", 0, toolTypeDoc)
var inputMSSQLHost = flag.String("mshost", "127.0.0.1", "MS SQL Host")
var inputMSSQLPort = flag.String("mspost", "1433", "MS SQL Port")
var inputMSSQLUsr = flag.String("msusr", "sa", "MS SQL User")
var inputMSSQLpwd = flag.String("mspwd", "", "MS SQL Password")
var inputMSSQLDBName = flag.String("msdbname", "", "MS SQL dbName")

func init_flag() {
	flag.Parse()
}

type MSSQL struct {
	*sql.DB
	dataSource string
	database   string
	windows    bool
	sa         SA
}

type SA struct {
	user   string
	passwd string
}

func SpliceConnectString(db MSSQL) []string {
	var conf []string
	conf = append(conf, "Provider=SQLOLEDB")
	conf = append(conf, "Data Source="+db.dataSource)
	conf = append(conf, "Initial Catalog="+db.database)
	// Integrated Security=SSPI 这个表示以当前WINDOWS系统用户身去登录SQL SERVER服务器
	// (需要在安装sqlserver时候设置)，
	// 如果SQL SERVER服务器不支持这种方式登录时，就会出错。
	if db.windows {
		conf = append(conf, "integrated security=SSPI")
	} else {
		conf = append(conf, "user id="+db.sa.user)
		conf = append(conf, "password="+db.sa.passwd)
	}
	return conf
}

func TestMsSQLConnect() (err error) {
	db := MSSQL{
		dataSource: *inputMSSQLHost,
		database:   "",
		windows:    false,
		sa: SA{
			user:   *inputMSSQLUsr,
			passwd: *inputMSSQLpwd,
		},
	}
	conf := SpliceConnectString(db)
	db.DB, err = sql.Open("adodb", strings.Join(conf, ";"))
	err = db.DB.Ping()
	if err != nil {
		fmt.Println("sql open:", err)
		os.Exit(1)
		return err
	}
	defer db.Close()
	return nil
}

func TestMsSQLInstance() (err error) {
	db := MSSQL{
		dataSource: *inputMSSQLHost,
		database:   *inputMSSQLDBName,
		windows:    false,
		sa: SA{
			user:   *inputMSSQLUsr,
			passwd: *inputMSSQLpwd,
		},
	}
	conf := SpliceConnectString(db)
	db.DB, err = sql.Open("adodb", strings.Join(conf, ";"))
	rows, err := db.DB.Query("")
	if rows.Next() {

	}
	if err != nil {
		fmt.Println("sql open:", err)
		os.Exit(1)
		return err
	}
	defer db.Close()
	return nil
}

func main() {

	init_flag()
	switch *inputToolType {
	case 0:
		fmt.Println("help message")
	case 1:
		TestMsSQLConnect()
	case 2:
		TestMsSQLInstance()
	default:
		fmt.Println("222222222")
	}
	os.Exit(0)

}
