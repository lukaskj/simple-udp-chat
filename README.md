# Simple UDP server/client chat made in Go

First project in Go done for study purpose


## Usage

### Dependencies
```
go get -u github.com/satori/go.uuid
```

### Usage
#### Server
* Build
```
go build -o ./build/server ./server/main.go ./server/server.go
```

* Run server
```
./build/server
```
> Currently port 8080 is fixed in code (plans to add port via arguments)

#### Commands
* **q** or **quit**: shutdown server and send it to all clients
* **users**: list of connected users
* **debug**: toggle debug mode (prints all receiving messages to console)

#### Client
* Build
```
go build -o ./build/client ./client/main.go
```

* Run
```
./build/client
```
> Currently localhost:8080 is fixed in code (plans to add ip:port via arguments)

#### Commands
* **:q** or **:quit**: disconnect from server
