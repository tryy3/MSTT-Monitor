@echo off
echo "Compiling windows 64bit"
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-H windowsgui -X main.Version=%1" -o %1/mstt-client-%GOOS%-%GOARCH%-%1.exe

echo "Compiling Linux 64bit"
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-H windowsgui -X main.Version=%1" -o %1/mstt-client-%GOOS%-%GOARCH%-%1.exe

echo "Compiling OSX 64bit"
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "-H windowsgui -X main.Version=%1" -o %1/mstt-client-%GOOS%-%GOARCH%-%1.exe
REM curl -F "file=@mstt-client-windows-%1.exe" http://192.168.20.149/upload.php