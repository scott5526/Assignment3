/*
File: timeserver.go
Author: Robinson Thompson

Description: Runs a simple timeserver to pull up a URL page displaying the current time.  Support was verified for Windows 7 OS.  Support has not been tested for other OS

Copyright:  All code was written originally by Robinson Thompson with assistance from various
	    free online resourses.  To view these resources, check out the README
*/
package main

import (
"flag"
"fmt"
"html/template"
"math/rand"
"net/http"
"sync"
"os"
//"os/exec"
"strconv"
"time"
)

var currUser string
var templatesPath *string
var redirect bool
var portNO *int
var printToFile int
var writeFile *os.File
var cookieMap = make(map[string]http.Cookie)
var mutex = &sync.Mutex{}
var portInfoStuff PortInfo

type PortInfo struct {
	PortNum string
}

type TimeInfo struct {
	Name string
	LocalTime string
	UTCTime string
	PortNum string
}

/*
Greeting Redirect 1

Redirects to greetingHandler with a saved URL "/"
*/

func greetingRedirect1(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
	badHandler(w,r) // check if the URL is valid
	return
    }

    fmt.Println("localhost:" + strconv.Itoa(*portNO) + "/")

    if  printToFile == 1 { // make sure p2f is enabled
        currentWrite := []byte("localhost:" + strconv.Itoa(*portNO) + "/" + "\r\n")
	writeFile.Write(currentWrite)
    }
    greetingHandler(w,r)
}

/*
Greeting Redirect 2

Redirects to greetingHandler with a saved URL "/index.html"
*/

func greetingRedirect2(w http.ResponseWriter, r *http.Request) {
    fmt.Println("localhost:" + strconv.Itoa(*portNO) + "/index.html")

    if  printToFile == 1 { // make sure p2f is enabled
        currentWrite := []byte("localhost:" + strconv.Itoa(*portNO) + "/index.html" + "\r\n")
	writeFile.Write(currentWrite)
    }
    greetingHandler(w,r)
}

/*
Greeting message

Presents the user with a login message if a cookie is found for them, otherwise redirects to the login page
*/
func greetingHandler(w http.ResponseWriter, r *http.Request) {
    greetingCheck(w, r)

    if redirect == true { //If no matching cookie was found in the cookie map, redirect
	path := (*templatesPath + "loginRedirect.html")
    	newTemplate,err := template.New("redirect").ParseFiles(path) 
    	if err != nil {
		fmt.Println("Error running login redirect template")
		return
    	}   
    	newTemplate.ExecuteTemplate(w,"loginRedirectTemplate",portInfoStuff)
    }
}

/*
Login handler.  
Displays a html generated login form for the user to provide a name.  
Creates a cookie for the user name and redirects them to the home page if a valid user name was provided.  
If no valid user name was provided, outputs an error message
*/
func loginHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("localhost:" + strconv.Itoa(*portNO) + "/login")

    if  printToFile == 1 {
        currentWrite := []byte("localhost:" + strconv.Itoa(*portNO) + "/login" + "\r\n")	
	writeFile.Write(currentWrite)
	
    }
  
    loginCheck(w,r)


    // Unique ID generation below

    //tempUUID,_ := exec.Command("uuidgen").Output()
    // uncomment me (^^^^^^^^^) when testing on linux!!!

    newUUID := strconv.Itoa(rand.Int())
    // comment me (^^^^^^^^^) when testing on linux!!!
    //newUUID := string(tempUUID[:])
    // uncomment me (^^^^^^^^^) when testing on linux!!!

    expDate := time.Now()
    expDate.AddDate(1,0,0)

    //Generate & set browser cookie
    cookie := http.Cookie{Name: "localhost", Value: newUUID, Expires: expDate, HttpOnly: true, MaxAge: 100000, Path: "/"}
    http.SetCookie(w,&cookie)

    path := *templatesPath + "login.html"
    newTemplate,err := template.ParseFiles(path)   
    if err != nil {
	fmt.Println("Error running login template")
	return;
    } 
    newTemplate.Execute(w,"loginTemplate")

    r.ParseForm()
    name := r.PostFormValue("name")
    submit := r.PostFormValue("submit") 

    if submit == "Submit" { // check if the user hit the "submit" button
    	if name == "" {
		path = *templatesPath + "/badLogin.html"
    		newTemplate,_ := template.New("outputUpdate").ParseFiles(path)   
    		newTemplate.ExecuteTemplate(w,"badLoginTemplate",nil)
    	} else {
		//generate cookie map's cookie
		mapCookie := http.Cookie{
		Name: newUUID, 
		Value: name, 
		Path: "/", 
		Domain: "localhost", 
		Expires: expDate,
 		HttpOnly: true, 
		MaxAge: 100000,
		}
		//lock the cookie map while it's being written to
		mapSetCookie(mapCookie, newUUID)

		fmt.Println("localhost:" + strconv.Itoa(*portNO) + "/login?name=" + name)

    		if  printToFile == 1 { // check if the p2f flag was set
        		currentWrite := []byte("localhost:" + strconv.Itoa(*portNO) + "/login?name=" + name + "\r\n")
			writeFile.Write(currentWrite)
    		}

		//Redirect to greetings (home) page
		path = *templatesPath + "greetingRedirect.html"
    		newTemp,err := template.New("redirect").ParseFiles(path)   
    		if err != nil {
			fmt.Println("Error running greeting redirect template")
			return;
    		} 
    		newTemp.ExecuteTemplate(w,"greetingRedirectTemplate",portInfoStuff)
    	}
    }
}

/*
Logout handler.  

Clears user cookie, displays goodbye message for 10 seconds, then redirects user to login form
*/
func logoutHandler(w http.ResponseWriter, r *http.Request) {
   fmt.Println("localhost:" + strconv.Itoa(*portNO) + "/logout")

   if  printToFile == 1 { //Check if p2f flag is set
        currentWrite := []byte("localhost:" + strconv.Itoa(*portNO) + "/logout" + "\r\n")	
	writeFile.Write(currentWrite)
   }

   clearMapCookie(r)

    // User wasn't actually logged in, redirect them to login page
    if !redirect {
	path := (*templatesPath + "loginRedirect.html")
    	newTemplate,err := template.New("redirect").ParseFiles(path) 
    	if err != nil {
		fmt.Println("Error running login redirect template")
		return
    	}   
    	newTemplate.ExecuteTemplate(w,"loginRedirectTemplate",portInfoStuff)
    }

    //Redirect to the login page
    path := *templatesPath + "logoutToLoginRedirect.html"
    newTemplate,err := template.New("redirect").ParseFiles(path)  
    if err != nil {
	fmt.Println("Error running login redirect template")
	return;
    }  
    newTemplate.ExecuteTemplate(w,"loginRedirectTemplate",portInfoStuff)
}


/*
Handler for time requests.  

Outputs the current time in the format:
Hour:Minute:Second PM/AM
*/
func timeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("localhost:" + strconv.Itoa(*portNO) + "/time")

    if  printToFile == 1 { //Check if the p2f flag is set
        currentWrite := []byte("localhost:" + strconv.Itoa(*portNO) + "/time" + "\r\n")
	writeFile.Write(currentWrite)
    }

    user := getUserName(r)

    currTime := time.Now().Format("03:04:05 PM")
    utcTime := time.Now().UTC()
    utcTime = time.Date(
        time.Now().UTC().Year(),
        time.Now().UTC().Month(),
        time.Now().UTC().Day(),
        time.Now().UTC().Hour(),
        time.Now().UTC().Minute(),
        time.Now().UTC().Second(),
        time.Now().UTC().Nanosecond(),
        time.UTC,
    )

    utcTime.UTC()
    //utcTime.Format("03:04:05 07")

    currTimeInfo := TimeInfo {
    	Name: user,
	LocalTime: currTime,
	UTCTime: utcTime.Format("03:04:05"),
	PortNum: strconv.Itoa(*portNO),
    }

    path := *templatesPath + "time.html"
    newTemplate,err := template.New("timeoutput").ParseFiles(path)  
    if err != nil {
	fmt.Println("Error running time template")
	return;
    } 
    newTemplate.ExecuteTemplate(w,"timeTemplate",currTimeInfo)
}

/*
Menu handler.  

Displays menu consisting of Home, Time, Logout, and About us
*/
func menuHandler(w http.ResponseWriter, r *http.Request) {
   fmt.Println("localhost:" + strconv.Itoa(*portNO) + "/menu")

   if  printToFile == 1 { //Check if p2f flag is set
        currentWrite := []byte("localhost:" + strconv.Itoa(*portNO) + "/menu" + "\r\n")	
	writeFile.Write(currentWrite)
   }

    //Redirect to the menu page
    path := *templatesPath + "menu.html"
    newTemplate,err := template.New("redirect").ParseFiles(path)  
    if err != nil {
	fmt.Println("Error running menu redirect template")
	return;
    }  
    newTemplate.ExecuteTemplate(w,"menuTemplate",portInfoStuff)
}

/*
Handler for invalid requests.  Outputs a 404 error message and a cheeky message
*/
func badHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/index.html" {
	return
    } else if r.URL.Path == "/login" {
	return
    } else if r.URL.Path == "/logout" {
	return
    } else if r.URL.Path == "/time" {
	return
    } else if r.URL.Path == "/menu" {
	return
    }

    http.NotFound(w, r)
    w.Write([]byte("These are not the URLs you're looking for."))
    return
}

/*
Main
*/
func main() {
    fmt.Println("Starting new server")
    //Version output & port selection
    version := flag.Bool("V", false, "Version 3.4.1") //Create a bool flag for version  
    						    //and default to no false

    portNO = flag.Int("port", 8080, "")	    //Create a int flag for port selection
					            //and default to port 8080

    p2f := flag.Bool("p2f", false, "") //flag to output to file

    templatesPath = flag.String("templates", "Templates/", "")


    printToFile = 0 // set to false

    flag.Parse()

    if *version == true {		//If version outputting selected, output version and 
        fmt.Println("Version 3.4.1")	//terminate program with 0 error code
        os.Exit(0)
    }

    if *p2f == true {
	writeFile,_ = os.Create("output.txt")
	printToFile = 1 // set to true
    }
	
    portInfoStuff = PortInfo{
	PortNum: strconv.Itoa(*portNO),
    }

    // URL handling
    http.HandleFunc("/", greetingRedirect1)
    http.HandleFunc("/index.html", greetingRedirect2)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/logout", logoutHandler)
    http.HandleFunc("/time", timeHandler)
    http.HandleFunc("/menu", menuHandler)
    
    //Check host:(specified port #) for incomming connections
    error := http.ListenAndServe("localhost:" + strconv.Itoa(*portNO), nil)

    if error != nil {				// If the specified port is already in use, 
	fmt.Println("Port already in use")	// output a error message and exit with a 
	os.Exit(1)				// non-zero error code
    }
}
