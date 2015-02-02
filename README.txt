timeserver.go README

Resources used
---------------------------------------------------------------------------------------------
http://golang.org/pkg/os/
http://golangtutorials.blogspot.com/2011/06/web-programming-with-go-first-web-hello.html
http://stackoverflow.com/questions/10105935/how-to-convert-a-int-value-to-string-in-go
http://golang.org/pkg/net/http/
http://golang.org/pkg/time/
http://grokbase.com/t/gg/golang-nuts/134kenh4xz/go-nuts-time-format-giving-unpredictable-results
http://golang.org/pkg/net/http/#Cookie
http://golang.org/pkg/sync/#RWMutex
http://man7.org/linux/man-pages/man1/uuidgen.1.html
http://stackoverflow.com/questions/12130582/setting-cookies-in-golang-net-http
https://golang.org/ref/spec#Map_types
http://www.w3schools.com/htmL/html_forms.asp
http://stackoverflow.com/questions/12513963/how-to-read-input-from-a-html-form-and-save-it-in-a-file-golang
http://stackoverflow.com/questions/12612409/go-programming-post-formvalue-cant-be-printed
http://msdn.microsoft.com/en-us/library/ie/ms534184(v=vs.85).aspx
http://webmaster.iu.edu/tools-and-guides/maintenance/redirect-meta-refresh.phtml
https://gist.github.com/mschoebel/9398202
http://stackoverflow.com/questions/15130321/is-there-a-method-to-generate-a-uuid-with-go-language
http://www.reddit.com/r/golang/comments/2rkij9/cant_set_a_cookie/
http://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.1.html
https://www.socketloop.com/tutorials/golang-convert-cast-bytes-to-string
http://golang.org/pkg/math/rand/#Int
https://gobyexample.com/mutexes
chrome://settings/cookies
---------------------------------------------------------------------------------------------



Running the timeserver.go file
---------------------------------------------------------------------------------------------
To run timeserver.go, open the Windows command prompt and move to the directory of timeserver.go.  To run the file, use "go run timeserver.go" with any applicable flags.
Also ensure that the login.gtpl and badLogin.gtpl files are in the timeserver.go directory.  login.html doesn't matter

Applicable flags include:

-V ("go run timeserver.go -V)

Runs timeserver.go with the version flag enabled.  Will output the current version of the file and terminate the program with a zero error code.

-port # ("go run timeserver.go -port 9999)

Runs timeserver.go with a specified port (the default port # is 8080).

-p2f ("go run timeserver.go -p2f

Writes accessed URLS to output.txt in addition to the console
---------------------------------------------------------------------------------------------



Accessing the server from a web browser
---------------------------------------------------------------------------------------------
Enter the desired URL.  Some URLS are only accessible when logged in and will redirect otherwise

The supported URLS are (where "(xxx)" is the port #:
http://localhost:(xxx)/
http://localhost:(xxx)/index.html
http://localhost:(xxx)/time
http://localhost:(xxx)/login
http://localhost:(xxx)/logout
---------------------------------------------------------------------------------------------



Caveats
---------------------------------------------------------------------------------------------
When trying to run the server, if the specified port is already in use the program will terminate with a error message on a non-zero error code.

Any URL beyond http://localhost:(port #) that doesn't match the above specified URL will result in a 404 not found web page.

/index.html will redirect to /login if there is no found cookie on the user's webbrowser

/time will add the user's name to the outputted time if a cookie has been found on their webbrowser

output.txt will be wiped whenever restarting the server with the -p2f flag

http://host:port was not used as was asked in the instructions.  Using this URL gave port-in-use errors so it was avoided
---------------------------------------------------------------------------------------------