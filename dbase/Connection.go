package dbase

import(
    "log"
    "database/sql"
    _"database/sql/driver/mysql"
)

var Db *sql.DB
func init(){    
    x,err:=sql.Open("mysql","root:sqla@tcp(127.0.0.1:3306)/go")
	if err!=nil{
		log.Fatal(err)
	}
    Db=x
	if err=Db.Ping(); err!=nil{
		log.Fatal(err)
	}
}
