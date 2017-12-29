package email

import(
    "fmt"
    "net/smtp"
    "log"
)

//email sending functionality
func Sent_mail(email,link string){
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