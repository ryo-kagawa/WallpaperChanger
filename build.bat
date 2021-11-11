go get github.com/google/go-licenses
go-licenses save . --force --save_path build/Windows/licenses
go build -o build/Windows/WallpaperChanger.exe -trimpath
go build -o build/Windows/WallpaperChanger-no-console.exe -ldflags -H=windowsgui -trimpath
COPY config.yaml build\config.yaml
