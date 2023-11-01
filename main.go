package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)


func init(){
	if err:=godotenv.Load();err != nil {
		log.Fatalf("error loading .env")
	}
}


var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	// TODO:randomize
	randomState="random"
)

func main() {
	fmt.Println(os.Getenv("GOOGLE_CLIENT_ID"))
	http.HandleFunc("/",handleHome)
	http.HandleFunc("/login",handleLogin)
	http.HandleFunc("/callback",handleCallback)

	http.ListenAndServe(":8080", nil)
}
func handleHome(w http.ResponseWriter, r * http.Request){
	html:=`<html><body><a href="/login">Google Log In</a></body></html>`
	fmt.Fprint(w,html)
}


func handleLogin(w  http.ResponseWriter, r *http.Request){
	url:=googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w,r,url,http.StatusTemporaryRedirect)
}


func handleCallback(w http.ResponseWriter, r *http.Request){
	state:=r.FormValue("state")
	if state != randomState{
		fmt.Println("state is not valid")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}

	code:=r.FormValue("code")
	token,err:=googleOauthConfig.Exchange(oauth2.NoContext,code)
	if err != nil {
		fmt.Printf("could not get token: %s\n",err.Error())
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}


	resp,err:=http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token="+token.AccessToken)

if err != nil {
	fmt.Println("couldn't create get request")
	http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
	return
}

defer resp.Body.Close()

content,err:=ioutil.ReadAll(resp.Body)
if err != nil {
	fmt.Println("could not parse response")
	http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
	return
}

fmt.Fprintf(w,"response:%s",content)
}