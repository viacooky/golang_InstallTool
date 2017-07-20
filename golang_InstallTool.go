package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/beevik/etree"
	_ "github.com/mattn/go-adodb"
)

var usageMsg string = `InstallTool
you can use this tool to Test MSSQL Connect/MSSQL DB is exists/Update XML file...

[-type] is must input
	0:Show Messages (default)
	1:TestMsSQLConnect return: 0 connect sueecss; 1 Can not connect 
	2:TestMsSQLInstanceExists return: 0 Instance is not Exists; 1 Instance is Exists
	5:Update XML tag , must input: -file -xpath -tag
	6:Update XML text , must input: -file -xpath -text
	7:Update XML attr value , must input: -file -xpath -key -value

[-h] you can use -h command to see helper
`
var inputToolType = flag.Int("type", 0, "Tool Type")
var inputMSSQLHost = flag.String("mshost", "127.0.0.1", "MS SQL Host")
var inputMSSQLPort = flag.String("mspost", "1433", "MS SQL Port")
var inputMSSQLUsr = flag.String("msusr", "sa", "MS SQL User")
var inputMSSQLpwd = flag.String("mspwd", "", "MS SQL Password")
var inputMSSQLDBName = flag.String("msdbname", "", "MS SQL dbName")
var inputXMLFile = flag.String("file", "", "XML file path")
var inputXMLxpath = flag.String("xpath", "", "xpath")
var inputXMLtag = flag.String("tag", "", "xml tag")
var inputXMLkey = flag.String("key", "", "xml key")
var inputXMLvalue = flag.String("value", "", "xml value")
var inputXMLtext = flag.String("text", "", "xml text")

func init_flag() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageMsg+"\n")
		fmt.Fprintf(os.Stderr, "Usage of :\n\n")
		flag.PrintDefaults()
	}
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

func TestMsSQLConnect() {
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
	var err error
	db.DB, err = sql.Open("adodb", strings.Join(conf, ";"))
	err = db.DB.Ping()
	defer db.Close()
	if err != nil {
		fmt.Println("sql open:", err)
		os.Exit(1)
	}
}

func TestMsSQLInstance() {
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
	var err error
	db.DB, err = sql.Open("adodb", strings.Join(conf, ";"))
	if err != nil {
		fmt.Println("sql open:", err)
		os.Exit(1)
	}
	rows, err := db.DB.Query("select name from master.dbo.sysdatabases where name = '" + db.database + "'")
	defer db.Close()
	if err != nil {
		// fmt.Println("sql query:", err)
		os.Exit(0)
	}
	if rows.Next() {
		os.Exit(1)
	}
}

func UpdateXMLConfigTag() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(*inputXMLFile); err != nil {
		panic(err)
	}
	root := doc.FindElement(*inputXMLxpath)
	root.Tag = *inputXMLtag
	doc.WriteToFile(*inputXMLFile)
}

func UpdateXMLConfigText() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(*inputXMLFile); err != nil {
		panic(err)
	}
	root := doc.FindElement(*inputXMLxpath)
	root.SetText(*inputXMLtext)
	doc.WriteToFile(*inputXMLFile)
}

func UpdateXMLConfigAttr() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(*inputXMLFile); err != nil {
		panic(err)
	}
	root := doc.FindElement(*inputXMLxpath)
	root.SelectAttr(*inputXMLkey).Value = *inputXMLvalue
	doc.WriteToFile(*inputXMLFile)
}

func main() {
	// UpdateXMLConfig()
	init_flag()
	switch *inputToolType {
	case 0:
		fmt.Println(usageMsg)
	case 1:
		TestMsSQLConnect()
	case 2:
		TestMsSQLInstance()
	case 5:
		UpdateXMLConfigTag()
	case 6:
		UpdateXMLConfigText()
	case 7:
		UpdateXMLConfigAttr()
	default:
		fmt.Println(usageMsg)
	}
	os.Exit(0)
}
