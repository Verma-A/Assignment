package main
import(
	"log"
	"time"
	"net/http"
    "logout"
	b64 "encoding/base64"  
    "dbase"
	_"database/sql/driver/mysql"
    "account"
    "dbopt"
    "email"
)
func main(){
	http.HandleFunc("/",account.Login)
	http.HandleFunc("/SignUp.html",account.Signup)
	http.HandleFunc("/LogIn.html",account.Login)
	http.HandleFunc("/Account_info.html",account.Acc_info)
	http.HandleFunc("/activate_my_account",account.Activate_account)
    http.HandleFunc("/SignOut",logout.Signout)
    http.HandleFunc("/Update",account.Update)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err:=http.ListenAndServe(":8080",nil)
	if err==nil{
		log.Fatal(err)
	}
	timer:=time.NewTicker(24*time.Hour)		//trigger task in every 24 hours
	go func(){
		for _=range timer.C{
			time_check()
		}
	}()
	time.Sleep(365 * 24 * time.Hour)
	timer.Stop()
}

//for sending email after every 7 days if users account is not activate
func time_check(){
	
	db:=dbase.Db
	t:=time.Now()
	t=t.AddDate(0,0,-7)
	row,err:=dbopt.Get_info(db,"",string(t.Format("2006-2-1")))
	if err!=nil{
		log.Fatal(err)
	}
	for row.Next(){
		var id string
		row.Scan(&id)
        email.Sent_mail(id,"http://localhost:8080/activate_my_account?id="+b64.StdEncoding.EncodeToString([]byte(id)))
	}
}