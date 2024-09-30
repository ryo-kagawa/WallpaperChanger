CHCP 65001

SET /P version=バージョンを入力: 
go install github.com/google/go-licenses@latest
go-licenses save ./... --force --save_path build/WallpaperChanger/licenses
go build -o build/WallpaperChanger/ConfigCreater.exe -ldflags "-s -w -X=github.com/ryo-kagawa/WallpaperChanger/subcommand.name=ConfigCreater.exe -X=github.com/ryo-kagawa/WallpaperChanger/subcommand.version=%version%" -trimpath ./cmd/config-creater
go build -o build/WallpaperChanger/WallpaperChanger.exe -ldflags "-s -w -H=windowsgui -X=github.com/ryo-kagawa/WallpaperChanger/subcommand.name=WallpaperChanger.exe -X=github.com/ryo-kagawa/WallpaperChanger/subcommand.version=%version%" -trimpath ./cmd/wallpaper-changer
go build -o build/WallpaperChanger/WallpaperRegistry.exe -ldflags "-s -w -X=github.com/ryo-kagawa/WallpaperChanger/subcommand.name=WallpaperRegistry.exe -X=github.com/ryo-kagawa/WallpaperChanger/subcommand.version=%version%" -trimpath ./cmd/wallpaper-registry
COPY README-windows.txt build\WallpaperChanger\README.txt
