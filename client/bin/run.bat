@echo off
go build -ldflags "-H windowsgui -X main.Version=%1" -o mstt-client-windows-%1.exe
REM curl -F "file=@mstt-client-windows-%1.exe" http://192.168.20.149/upload.php