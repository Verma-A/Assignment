package logout
import "net/http"
import "time"

//for deleting the cookies and logout from the account
func Signout(w http.ResponseWriter,r *http.Request){
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
