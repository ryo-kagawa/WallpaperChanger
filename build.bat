go get github.com/google/go-licenses
go-licenses save . --force --save_path build/licenses
go build -o build/wallpaper.exe
go build -o build/wallpaper-no-console.exe -ldflags -H=windowsgui
COPY config.yaml build\config.yaml
