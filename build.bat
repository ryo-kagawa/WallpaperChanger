go get github.com/google/go-licenses
go-licenses save . --force --save_path build/Windows/licenses
go build -o build/WallpaperChanger-Windows/WallpaperChanger.exe -ldflags "-s -w" -trimpath ./cmd/wallpaper-changer
go build -o build/WallpaperChanger-Windows/WallpaperChanger-no-console.exe -ldflags "-s -w -H=windowsgui" -trimpath ./cmd/wallpaper-changer
go build -o build/WallpaperChanger-Windows/ConfigCreater.exe -ldflags "-s -w" -trimpath ./cmd/config-creater
