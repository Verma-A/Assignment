package main
import(
	"log"
	//"html/template"
	"strings"
	"time"
	"fmt"
	"net/smtp"
	"net/http"
	"net/url"
	b64 "encoding/base64"  
	"database/sql"
	_"database/sql/driver/mysql"
)
func main(){
	http.HandleFunc("/",login)
	http.HandleFunc("/SignUp.html",signup)
	http.HandleFunc("/LogIn.html",login)
	http.HandleFunc("/Account_info.html",acc_info)
	http.HandleFunc("/activate_my_account",activate_account)
    http.HandleFunc("/SignOut",signout)
    http.HandleFunc("/Update",update)
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

//for updation of account information
func update(w http.ResponseWriter,r *http.Request){
    if r.Method=="POST"{
        fname:=r.FormValue("fname")
        lname:=r.FormValue("lname")
        dob:=r.FormValue("dob")
        cookie,_:=r.Cookie("email")
        db:=getConnection()
        if update_db(db,fname,lname,dob,cookie.Value){
            fmt.Fprintf(w,"<h1 align=center>Data updated successfully</h1><br><a href='Account_info.html'>Click here</a> to go back")
        }
    }
}

//for updating data into the database
func update_db(db *sql.DB,fname,lname,dob,email string)(bool){
    _,err:=db.Exec("update users set fname='"+fname+"', lname='"+lname+"', dob='"+dob+"' where email='"+email+"'")
    if err!=nil{
        log.Fatal(err)
        return false
    }else{
        return true
    }
}

//for deleting the cookies and logout from the account
func signout(w http.ResponseWriter,r *http.Request){
    c:=&http.Cookie{
        Name:   "fname",
        Value:  "",
        Path:   "",
        Expires:time.Unix(0,0),
        HttpOnly:true,
    }
    http.SetCookie(w, c)
    c=&http.Cookie{
        Name:   "lname",
        Value:  "",
        Path:   "",
        Expires:time.Unix(0,0),
        HttpOnly:true,
    }
    http.SetCookie(w, c)
    c=&http.Cookie{
        Name:   "dob",
        Value:  "",
        Path:   "",
        Expires:time.Unix(0,0),
        HttpOnly:true,
    }
    http.SetCookie(w, c)
    c=&http.Cookie{
        Name:   "status",
        Value:  "",
        Path:   "",
        Expires:time.Unix(0,0),
        HttpOnly:true,
    }
    http.SetCookie(w,c)
    http.ServeFile(w,r,"LogIn.html")
}

//for sending email after every 7 days if users account is not activate
func time_check(){
	
	db:=getConnection()
	t:=time.Now()
	t=t.AddDate(0,0,-7)
	row,err:=get_info(db,"",string(t.Format("2006-2-1")))
	if err!=nil{
		log.Fatal(err)
	}
	for row.Next(){
		var email string
		row.Scan(&email)
		sent_mail(email,"http://localhost:8080/activate_my_account?id="+b64.StdEncoding.EncodeToString([]byte(email)))
	}
}
/*type info struct{
	fname string
}*/
func acc_info(w http.ResponseWriter, r *http.Request){
	if r.Method!="POST"{
		db:=getConnection();
		cookie,_:=r.Cookie("email")
		rows,err:=get_info(db,cookie.Value,"update")
		if err!=nil{
			log.Fatal(err)
		}else{
			var fname,lname,dob string
			rows.Next()
			if err:=rows.Scan(&fname,&lname,&dob); err!=nil{
				log.Fatal(err)
			}else{
				exp:=time.Now().Add(365 * 24 * time.Hour)
				c:=http.Cookie{Name:"fname",Value:fname,Expires:exp}
				http.SetCookie(w,&c)
				c=http.Cookie{Name:"lname",Value:lname,Expires:exp}
				http.SetCookie(w,&c)
				c=http.Cookie{Name:"dob",Value:dob,Expires:exp}
				http.SetCookie(w,&c)
                http.ServeFile(w,r,"Account_info.html")
			}
		}
	}
}

//for activating user account
func activate_account(w http.ResponseWriter,r *http.Request){
	db:=getConnection()
	 url1:=r.URL.RequestURI()
	 u,_:=url.Parse(url1)
	m,_:=url.ParseQuery(u.RawQuery)
	enc:=m["id"][0]
	id,_:=b64.URLEncoding.DecodeString(enc)
	_, err:=db.Exec("update users set status='active' where email='"+string(id)+"'")
	if err!=nil{
		log.Fatal(err)
	}else{
		fmt.Fprintf(w,"<h1 font-size=200px align='center' color='green'>Your Account has been activated successfully</h1><br><a href='LogIn.html'>Click here</a> to login.")
	}
}

//login functioning
func login(w http.ResponseWriter,r *http.Request){
	if r.Method!="POST"{
		http.ServeFile(w,r,"LogIn.html")
		return
	}
	email:=r.FormValue("id")
	password:=r.FormValue("pwd")
	db:=getConnection();
    var mail,pwd,status string
	rows,err:=get_info(db,email,"login")	
	if err!=nil{
		fmt.Fprintf(w,"Server is under maintainance")
	}else{
		if rows.Next(){
			if err:=rows.Scan(&mail,&pwd,&status);err!=nil{
				log.Fatal(err)
			}else{
				u_pwd:=b64.StdEncoding.EncodeToString([]byte(password))
				if(0==strings.Compare(u_pwd,pwd)){
					exp:=time.Now().Add(365*24*time.Hour)
					cookie:=http.Cookie{Name:"email",Value:mail,Expires:exp}		//cookie created
					http.SetCookie(w,&cookie)                                       //set cookie
					cookie=http.Cookie{Name:"status",Value:status,Expires:exp}		//set cookie for reminder pop up
                    http.SetCookie(w,&cookie)
                    //http.ServeFile(w,r,"Account_info.html")
                    http.Redirect(w,r,"Account_info.html",301)
				}else{
					http.ServeFile(w,r,"LogIn.html")
				}
			}
		}else{
			fmt.Fprintf(w,"Incorrect Email.")
		}
	}
	defer rows.Close()
	defer db.Close()
}

//getting data from database
func get_info(db *sql.DB,id string,work string)(*sql.Rows,error){
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

//signup functioning
func signup(w http.ResponseWriter,r *http.Request){
	if r.Method!="POST"{
		http.ServeFile(w,r,"SignUp.html")
		return
	}else{
		var info [5]string
		info[0]=r.FormValue("fname")
		info[1]=r.FormValue("lname")
		info[2]=r.FormValue("id")
		info[3]=r.FormValue("dob")
		info[4]=b64.StdEncoding.EncodeToString([]byte(r.FormValue("pwd")))
		db:=getConnection()
		db_ins:=insert(db,info)
		
		if (db_ins){
			sent_mail(info[2],"http://localhost:8080/activate_my_account?id="+b64.StdEncoding.EncodeToString([]byte(info[2])))
			fmt.Fprintf(w,"<h1 font-size=200px align='center'>Registered successfully.</h1><br><a href='LogIn.html' align='center'>Click here</a> to log in.")
		}else{
			fmt.Fprintf(w,"<h1 font-size=200px align='center'>Registration unsuccessfull.</h1>")
		}
		return
	}
}

//getting connection of database
func getConnection()(*sql.DB){
	db,err:=sql.Open("mysql","root:sqla@tcp(127.0.0.1:3306)/go")
	if err!=nil{
		log.Fatal(err)
	}
	if err=db.Ping(); err!=nil{
		log.Fatal(err)
	}
	return db;
}

//insertion of data into database
func insert(db *sql.DB,info [5]string)bool{
	date:=time.Now().Local().Format("2006-2-1")
	_, err:=db.Exec("insert into users values('"+info[0]+"','"+info[1]+"','"+info[2]+"','"+info[3]+"','"+info[4]+"','"+date+"')")
	if err!=nil{
		//log.Fatal(err)
		return false
	}else{
		return true
	}
}

//email sending functionality
func sent_mail(email,link string){
	auth:=smtp.PlainAuth("","gofirsttime@gmail.com","qwerty@12345","smtp.gmail.com")
	to:=[]string{email}
	msg:=[]byte("To: "+email+"\r\n"+
				"Subject: Activation Link\r\n"+
				"\r\n"+"Hi,\nWelcome to our world. To enjoy joyfull future please activate your account by clicking the following link.\n\n"+
				""+link+"\r\n")
	err:=smtp.SendMail("smtp.gmail.com:587",auth,"gofirsttime@gmail.com",to,msg)
	if err!=nil{
		log.Fatal(err)
	}else{
		fmt.Println("Mail sent successfully")
	}
}