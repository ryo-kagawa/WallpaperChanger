「ConfigCreater.exe」を実行してconfig.yamlを作成してください
「WallpaperChanger.exe」を実行するとconfig.yamlに記載された内容で壁紙を設定します
「WallpaperChanger-no-console.exe」を実行するとコンソールログなしで実行します
タスクスケジューラーに登録する場合は「WallpaperChanger-no-console.exe」を登録してください

config.yamlは自動生成したものを利用することもできますが、自分で編集することも可能です
config.yamlの説明は次の通りです
・YAML形式
・imagePathには壁紙として設定したいフォルダーを指定してください
・rectangleListには矩形を配列で指定してください、ここで指定された開始座標から指定された大きさで表示します
・rectangleList[].x: 開始座標のx軸
・rectangleList[].y: 開始座標のy軸
・rectangleList[].width: 画像の横幅
・rectangleList[].height: 画像の縦幅
