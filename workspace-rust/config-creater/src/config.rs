use std::fs::File;
use std::io::Write;

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

impl Rectangle {
    pub fn new(x: u32, y: u32, width: u32, height: u32) -> Self {
        Self {
            x,
            y,
            width,
            height,
        }
    }
}

impl Config {
    pub fn new(image_path: String, rectangle_list: Vec<Rectangle>) -> Self {
        Self {
            image_path,
            rectangle_list,
        }
    }
    pub fn save(&self) -> Result<(), Box<dyn std::error::Error>> {
        let mut output = String::new();
        output.push_str(&format!("{}\n", self.image_path));
        for rectangle in &self.rectangle_list {
            output.push_str(&format!("  - x: {}\n", rectangle.x));
            output.push_str(&format!("    y: {}\n", rectangle.y));
            output.push_str(&format!("    width: {}\n", rectangle.width));
            output.push_str(&format!("    height: {}\n", rectangle.height));
        }
        let file_path = {
            let mut path_buf = std::path::PathBuf::from(std::env::current_exe()?);
            path_buf.pop();
            path_buf.push("config.yaml");
            path_buf
        };
        let mut output_file = File::create(file_path)?;
        write!(output_file, "{}", output)?;
        Ok(())
    }
}
