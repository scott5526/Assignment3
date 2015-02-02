/*
File: cookieManagement.go
Author: Robinson Thompson

Description:  Manages cookies for timeserver.go

Copyright:  All code was written originally by Robinson Thompson with assistance from various
	    free online resourses.  To view these resources, check out the README
*/
package main

import (
"fmt"
"html/template"
"net/http"
)

func mapSetCookie (newCookie http.Cookie, newUUID string) {
	mutex.Lock()
	cookieMap[newUUID] = newCookie
	mutex.Unlock()
}

func greetingCheck (w http.ResponseWriter, r *http.Request) {
    redirect = true
    for _, currCookie := range r.Cookies() { // check all potential cookies stored by the user for a matching cookie
    	if (currCookie.Name != "") {
	    currCookieVal := currCookie.Value
	    mutex.Lock()
	    mapCookie := cookieMap[currCookieVal]
	    mutex.Unlock()
            if (mapCookie.Value != "") {
    		fmt.Fprintf(w, "Greetings, " + mapCookie.Value)
		redirect = false
	    }
	}
    }
}

func loginCheck (w http.ResponseWriter, r *http.Request) {
   //Ensuring the user does not already have a browser cookie matching a cookie in the local cookie map, if they do
   //redirect the user to the greetings page
    for _, currCookie := range r.Cookies() {  //Run through the range of applicable cookies on the user's browser
    	if (currCookie.Name != "") {
	currCookieVal := currCookie.Value
	mutex.Lock()
	mapCookie := cookieMap[currCookieVal]  //Find the corresponding cookie in the local cookie map
	mutex.Unlock()
        	if (mapCookie.Value != "") {
			path := *templatesPath + "greetingRedirect.html"
    			newTemplate,err := template.New("redirect").ParseFiles(path)  
    			if err != nil {
				fmt.Println("Error running greeting redirect template")
				return
    			}  
    			newTemplate.ExecuteTemplate(w,"greetingRedirectTemplate",portInfoStuff)
		}
    	}
     }
}

func clearMapCookie (r *http.Request) {
   redirect = false // set to true if user cookie is found (they are actually logged in)
   for _, currCookie := range r.Cookies() {  //Run through the range of applicable cookies on the user's browser
    	if (currCookie.Name != "") {
	currCookieVal := currCookie.Value
	mutex.Lock()
	mapCookie := cookieMap[currCookieVal]  //Find the corresponding cookie in the local cookie map
	mutex.Unlock()
        	if (mapCookie.Value != "") {
			redirect = true // user was actually logged in
			mutex.Lock()
    			delete(cookieMap, currCookieVal) //Delete the cleared cookie from the local cookie map
			mutex.Unlock()
			currCookie.MaxAge = -1 //Set the user's cookie's MaxAge to an invalid number to expire it
		}
    	}
    }
}

func getUserName (r *http.Request) string {
    for _, currCookie := range r.Cookies() { //Lookup the user name by cross matching the user cookie's value against the local cookie maps's cookie names
    	if (currCookie.Name != "") {
	currCookieVal := currCookie.Value
	mutex.Lock()
	mapCookie := cookieMap[currCookieVal]
	mutex.Unlock()
        	if (mapCookie.Value != "") {
    			return ", " + mapCookie.Value
		}
    	}
    }
    return ""
}