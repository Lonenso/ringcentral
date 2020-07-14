## Ringcentral Interview

## How to run it?
### Release
- [linux/amd64]()
- [win/am64]()
### Build from source
```
git clone https://github.com/Lonenso/ringcentral.git
go mod tidy && go mod vendor 
go build main.go 
```
windows
```
./main.exe 
```
linux
```
./main
```
it should fail if CGO_ENABLED=0, make sure you have gcc installed.

### Run
```
./main.exe
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.newFile (4 handlers)
[GIN-debug] POST   /                         --> main.newFile (4 handlers)
[GIN-debug] GET    /text/:id                 --> main.downloadFile (4 handlers)
[GIN-debug] GET    /text                     --> main.viewFile (4 handlers)
[GIN-debug] GET    /draft/:id                --> main.editFile (4 handlers)
[GIN-debug] POST   /draft/:id                --> main.editFile (4 handlers)
[GIN-debug] GET    /debug                    --> main.showKV (4 handlers)
[GIN-debug] Listening and serving HTTP on :8080
```