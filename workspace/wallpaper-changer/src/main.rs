#![windows_subsystem = "windows"]

use image::{ImageFormat, RgbImage, imageops};
use rand::seq::IndexedRandom;
use windows::{
    Win32::{
        Foundation::{GENERIC_READ, GENERIC_WRITE},
        Storage::FileSystem::{
            CreateFileW, FILE_ATTRIBUTE_NORMAL, FILE_SHARE_READ, FILE_SHARE_WRITE, OPEN_EXISTING,
        },
        System::{
            Console::{ATTACH_PARENT_PROCESS, AllocConsole, AttachConsole},
            Registry::{
                HKEY, HKEY_CURRENT_USER, KEY_QUERY_VALUE, REG_DWORD, REG_VALUE_TYPE, RegOpenKeyExW,
                RegQueryValueExW,
            },
        },
        UI::WindowsAndMessaging::{
            MB_OK, MessageBoxW, SPI_SETDESKWALLPAPER, SPIF_SENDCHANGE, SPIF_UPDATEINIFILE,
            SystemParametersInfoW,
        },
    },
    core::{PCWSTR, w},
};

mod config;
mod subcommand;

fn main() {
    let result = if std::env::args().nth(1).as_deref() == Some("version") {
        Ok(format!("{}", subcommand::version()))
    } else {
        run()
    };
    if match &result {
        Ok(value) if value.is_empty() => false,
        Ok(_) => true,
        Err(_) => true,
    } {
        if unsafe { AttachConsole(ATTACH_PARENT_PROCESS).is_ok() } {
            match result {
                Ok(value) => {
                    println!("{}", value);
                }
                Err(error) => {
                    eprintln!("{:?}", error);
                }
            }
        } else {
            let alloc_console_result =
                unsafe { AllocConsole().map_err(|error| format!("AllocConsole {:?}", error)) };
            if alloc_console_result.is_err() {
                let lptext = format!("AllocConsole {:?}", alloc_console_result)
                    .encode_utf16()
                    .chain(Some(0))
                    .collect::<Vec<_>>();

                unsafe {
                    MessageBoxW(None, PCWSTR::from_raw(lptext.as_ptr()), w!("エラー"), MB_OK);
                }
                return;
            }
            let stdout = unsafe {
                CreateFileW(
                    w!("CONOUT$"),
                    (GENERIC_READ | GENERIC_WRITE).0,
                    FILE_SHARE_READ | FILE_SHARE_WRITE,
                    None,
                    OPEN_EXISTING,
                    FILE_ATTRIBUTE_NORMAL,
                    None,
                )
            };
            if stdout.is_err() {
                let lptext = format!("CreateFileW {:?}", stdout)
                    .encode_utf16()
                    .chain(Some(0))
                    .collect::<Vec<_>>();

                unsafe {
                    MessageBoxW(None, PCWSTR::from_raw(lptext.as_ptr()), w!("エラー"), MB_OK);
                }
                return;
            }
            match result {
                Ok(value) => {
                    println!("{}", value);
                }
                Err(error) => {
                    eprintln!("{:?}", error);
                }
            }
            println!("Enterキーで終了します");
            let _ = std::io::stdin().read_line(&mut String::new());
        }
    }
}

fn run() -> Result<String, Box<dyn std::error::Error>> {
    let mut infos = String::new();

    let config = config::load_config()?;
    let mut file_paths = Vec::new();

    for entry in walkdir::WalkDir::new(&config.image_path) {
        let entry = entry?;
        let path = entry.path().to_path_buf();
        if path.is_dir() {
            continue;
        }

        match ImageFormat::from_path(&path) {
            Ok(_) => file_paths.push(path),
            Err(_) => {
                infos.push_str(&format!("not registers image extension: {}\n", {
                    path.display()
                }));
            }
        }
    }

    let (mut width, mut height) = (0, 0);
    for rectangle in &config.rectangle_list {
        width = std::cmp::max(rectangle.x + rectangle.width, width);
        height = std::cmp::max(rectangle.y + rectangle.height, height);
    }
    let mut output_image = RgbImage::new(width, height);

    let mut rand = rand::rng();
    for rectangle in &config.rectangle_list {
        let file_path = &file_paths
            .as_slice()
            .choose(&mut rand)
            .ok_or("ファイルパス取得エラー")?;
        let img = image::ImageReader::open(file_path)?
            .with_guessed_format()?
            .decode()?
            .into_rgb8();
        let ratio = f64::min(
            rectangle.width as f64 / img.width() as f64,
            rectangle.height as f64 / img.height() as f64,
        );
        let dx = (img.width() as f64 * ratio).ceil() as u32;
        let dy = (img.height() as f64 * ratio).ceil() as u32;
        let offset_x = (rectangle.width - dx) / 2;
        let offset_y = (rectangle.height - dy) / 2;
        let start_x = rectangle.x + offset_x;
        let start_y = rectangle.y + offset_y;
        let resized = imageops::resize(&img, dx, dy, imageops::FilterType::CatmullRom);
        imageops::overlay(&mut output_image, &resized, start_x as i64, start_y as i64);
    }

    let output_file_path = {
        let mut path_buf = std::path::PathBuf::from(std::env::current_exe()?);
        path_buf.pop();
        path_buf.push("wallpaper.bmp");
        path_buf
    };
    output_image
        .save_with_format(&output_file_path, image::ImageFormat::Bmp)
        .map_err(|_| "ファイル出力エラー")?;

    unsafe {
        let mut hkey = HKEY::default();
        let result = RegOpenKeyExW(
            HKEY_CURRENT_USER,
            w!("Control Panel\\Desktop"),
            Some(0),
            KEY_QUERY_VALUE,
            &mut hkey,
        );
        if result.is_err() {
            return Err(format!("RegOpenKeyExA: {}", result.0).into());
        }

        let mut value_type = REG_VALUE_TYPE(0);
        let mut data = [0u8; 4];
        let mut data_size = data.len() as u32;
        let result = RegQueryValueExW(
            hkey,
            w!("JPEGImportQuality"),
            None,
            Some(&mut value_type),
            Some(data.as_mut_ptr()),
            Some(&mut data_size),
        );
        if result.is_err() {
            return Err(format!("RegQueryValueExW: {}", result.0).into());
        }
        if value_type != REG_DWORD || u32::from_le_bytes(data) != 100 {
            infos.push_str(&format!("RegQueryValueExW Control Panel\\Desktop value JPEGImportQuality is Not DWORD value 0x00000064"));
        }

        let wide: Vec<u16> = output_file_path
            .to_string_lossy()
            .encode_utf16()
            .chain(Some(0))
            .collect::<Vec<_>>();

        SystemParametersInfoW(
            SPI_SETDESKWALLPAPER,
            0,
            Some(PCWSTR(wide.as_ptr()).as_ptr() as _),
            SPIF_UPDATEINIFILE | SPIF_SENDCHANGE,
        )
        .map_err(|_| "壁紙の変更に失敗")?;
    }
    Ok(infos)
}
