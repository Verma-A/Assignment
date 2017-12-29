package dbopt

import(
    "log"
    "strings"
    "time"
    "database/sql"
)
//for updating data into the database
func Update_db(db *sql.DB,fname,lname,dob,email string)(bool){
    _,err:=db.Exec("update users set fname='"+fname+"', lname='"+lname+"', dob='"+dob+"' where email='"+email+"'")
    if err!=nil{
        log.Fatal(err)
        return false
    }else{
        return true
    }
}

//getting data from database
func Get_info(db *sql.DB,id string,work string)(*sql.Rows,error){
	var rows *sql.Rows
	var err error
	if (0==strings.Compare(work,"login")){
	rows,err=db.Query("Select email,password,status from users where email='"+id+"'")
	}else if(0==strings.Compare("update",work)){
		rows,err=db.Query("select fname,lname,dob from users where email='"+id+"'")
	}else{
		rows,err=db.Query("select email from users where status<>'active' and status='"+work+"'")
	}
	return rows,err
}

//insertion of data into database
func Insert(db *sql.DB,info [5]string)bool{
	date:=time.Now().Local().Format("2006-2-1")
	_, err:=db.Exec("insert into users values('"+info[0]+"','"+info[1]+"','"+info[2]+"','"+info[3]+"','"+info[4]+"','"+date+"')")
	if err!=nil{
		//log.Fatal(err)
		return false
	}else{
		return true
	}
}
