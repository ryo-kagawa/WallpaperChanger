「ConfigCreater.exe」を実行してconfig.yamlを作成してください
「WallpaperChanger-no-console.exe」を実行するとconfig.yamlに記載された内容で壁紙を設定します

config.yamlの記述方法
・YAML形式
・imagePathには壁紙として設定したいフォルダーを指定してください
・rectangleListには矩形を配列で指定してください、ここで指定された開始座標から指定された大きさで表示します
・rectangleList[].x: 開始座標のx軸
・rectangleList[].y: 開始座標のy軸
・rectangleList[].width: 画像の横幅
・rectangleList[].height: 画像の縦幅
