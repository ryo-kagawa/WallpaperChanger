use yaml_rust2::YamlLoader;

pub struct Rectangle {
    pub x: u32,
    pub y: u32,
    pub width: u32,
    pub height: u32,
}

pub struct Config {
    pub image_path: String,
    pub rectangle_list: Vec<Rectangle>,
}

pub fn load_config() -> Result<Config, Box<dyn std::error::Error>> {
    let file_path = {
        let mut path_buf = std::path::PathBuf::from(std::env::current_exe()?);
        path_buf.pop();
        path_buf.push("config.yaml");
        path_buf
    };
    let content = std::fs::read_to_string(&file_path)
        .map_err(|_| format!("設定ファイル{}が見つかりません", &file_path.display()))?
        .to_string();
    let docs = YamlLoader::load_from_str(&content)
        .map_err(|_| format!("設定ファイル{}の内容が不正です", &file_path.display()))?;

    let doc = &docs[0];

    let image_path = doc["imagePath"]
        .as_str()
        .ok_or_else(|| format!("設定ファイル{}の内容が不正です", &file_path.display()))?
        .to_string();
    let mut rectangle_list = Vec::new();
    for rectangle in doc["rectangleList"]
        .as_vec()
        .ok_or_else(|| format!("設定ファイル{}の内容が不正です", &file_path.display()))?
    {
        rectangle_list.push(Rectangle {
            x: rectangle["x"]
                .as_i64()
                .ok_or_else(|| format!("設定ファイル{}の内容が不正です", &file_path.display()))?
                as u32,
            y: rectangle["y"]
                .as_i64()
                .ok_or_else(|| format!("設定ファイル{}の内容が不正です", &file_path.display()))?
                as u32,
            width: rectangle["width"]
                .as_i64()
                .ok_or_else(|| format!("設定ファイル{}の内容が不正です", &file_path.display()))?
                as u32,
            height: rectangle["height"]
                .as_i64()
                .ok_or_else(|| format!("設定ファイル{}の内容が不正です", &file_path.display()))?
                as u32,
        })
    }
    Ok(Config {
        image_path,
        rectangle_list,
    })
}
