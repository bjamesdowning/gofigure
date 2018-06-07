Gofigure: 

Learning Go and how to interact with network infrastructure devices. Gofigure will serve as a front end 'console' to allow
RESTful calls to various boxes. Firstly, using the HTTP (nginx) API provided by Cumulus Linux. Eventually I'd like to
incorporate SSH connections to legacy gear. 

First iteration, this webapp is very basic allowing a user to create an account (stored only in memory) and on the admin
homepage make calls to a static Cumulus switch. Currently everything is built into main.go, as I learn I'll expand to a
structured organization of packages. Goal here is learning how to utilize Go for network configuration & information gathering as I find it easier to complie/distribute/"dockerize" than Python.

Second iteration, same functionality but uses redis databse for user storage. Docker compose used to build out app.
    Work on next 10:
        1. Users to Mongo
        2. Sessions on Redis
        3. Error handling on http response when DB is down
        4. Store results to commands in mongo
        5. Add commands
        6. Integrate http/router
        7. Build tests
        8. Integrate NXAPI
        9. Build SSL front end offload (nginx, HA Proxy?)
        10. Code breakout into packages.
