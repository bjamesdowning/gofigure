Gofigure: 

Learning Go and how to interact with network infrastructure devices. Gofigure will serve as a front end 'console' to allow
RESTful calls to various boxes. Firstly, using the HTTP (nginx) API provided by Cumulus Linux. Eventually I'd like to
incorporate SSH connections to legacy gear. 

First iteration, this webapp is very basic allowing a user to create an account (stored only in memory) and on the admin
homepage make calls to a static Cumulus switch. Currently everything is built into main.go, as I learn I'll expand to a
structured organization of packages. Goal here is learning how to utilize Go for network configuration & information gathering
as I find it easier to complie/distribute/"dockerize" than Python.
