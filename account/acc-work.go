package account

import(
    "net/http"
    "fmt"
    "log"
    "dbase"
    "net/url"
    b64 "encoding/base64"  
    "strings"
    "time"
    "dbopt"
    "email"
)
//login functioning
func Login(w http.ResponseWriter,r *http.Request){
	if r.Method!="POST"{
		http.ServeFile(w,r,"LogIn.html")
		return
	}
	email:=r.FormValue("id")
	password:=r.FormValue("pwd")
	//db:=db.Db;
    var mail,pwd,status string

	rows,err:=dbopt.Get_info(dbase.Db,email,"login")	
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
	//defer db.Close()
}
//for account information
func Acc_info(w http.ResponseWriter, r *http.Request){
	if r.Method!="POST"{
		db:=dbase.Db
		cookie,_:=r.Cookie("email")
		rows,err:=dbopt.Get_info(db,cookie.Value,"update")
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
func Activate_account(w http.ResponseWriter,r *http.Request){
	db:=dbase.Db
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

//signup functioning
func Signup(w http.ResponseWriter,r *http.Request){
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
		db:=dbase.Db
		db_ins:=dbopt.Insert(db,info)
		
		if (db_ins){
			email.Sent_mail(info[2],"http://localhost:8080/activate_my_account?id="+b64.StdEncoding.EncodeToString([]byte(info[2])))
			fmt.Fprintf(w,"<h1 font-size=200px align='center'>Registered successfully.</h1><br><a href='LogIn.html' align='center'>Click here</a> to log in.")
		}else{
			fmt.Fprintf(w,"<h1 font-size=200px align='center'>Registration unsuccessfull.</h1>")
		}
		return
	}
}

//for updation of account information
func Update(w http.ResponseWriter,r *http.Request){
    if r.Method=="POST"{
        fname:=r.FormValue("fname")
        lname:=r.FormValue("lname")
        dob:=r.FormValue("dob")
        cookie,_:=r.Cookie("email")
        db:=dbase.Db
        if dbopt.Update_db(db,fname,lname,dob,cookie.Value){
            fmt.Fprintf(w,"<h1 align=center>Data updated successfully</h1><br><a href='Account_info.html'>Click here</a> to go back")
        }
    }
}
