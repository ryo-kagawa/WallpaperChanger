use std::{ffi::OsString, mem, os::windows::ffi::OsStringExt};

use windows::{
    Win32::{
        Foundation::{FALSE, HWND, LPARAM, MAX_PATH, RECT, TRUE},
        Graphics::Gdi::{EnumDisplayMonitors, GetMonitorInfoW, HDC, HMONITOR, MONITORINFOEXW},
        System::Com::{COINIT_APARTMENTTHREADED, CoInitializeEx, CoUninitialize},
        UI::Shell::{BIF_NEWDIALOGSTYLE, BROWSEINFOW, SHBrowseForFolderW, SHGetPathFromIDListW},
    },
    core::{BOOL, w},
};

mod config;
mod subcommand;

pub struct ComInitializer;
impl ComInitializer {
    pub fn new() -> Result<Self, Box<dyn std::error::Error>> {
        let result = unsafe { CoInitializeEx(None, COINIT_APARTMENTTHREADED) };
        if result.is_err() {
            return Err(format!("CoInitializeEx: {:?}", result).into());
        }
        Ok(Self)
    }
}
impl Drop for ComInitializer {
    fn drop(&mut self) {
        unsafe {
            CoUninitialize();
        }
    }
}

pub struct Rectangle {
    pub x: u64,
    pub y: u64,
    pub width: u64,
    pub height: u64,
}

pub struct MonitorData {
    rectangles: Vec<Rectangle>,
}

unsafe extern "system" fn monitor_enum_proc(
    hmonitor: HMONITOR,
    _hdc: HDC,
    _lprect: *mut RECT,
    lparam: LPARAM,
) -> BOOL {
    let mut info: MONITORINFOEXW = unsafe { mem::zeroed() };
    info.monitorInfo.cbSize = mem::size_of::<MONITORINFOEXW>() as u32;
    let data = unsafe { &mut *(lparam.0 as *mut MonitorData) };
    if !unsafe { GetMonitorInfoW(hmonitor, &mut info.monitorInfo).as_bool() } {
        return FALSE;
    }
    let rect = &info.monitorInfo.rcMonitor;
    data.rectangles.push(Rectangle {
        x: rect.left as u64,
        y: rect.top as u64,
        width: (rect.right - rect.left) as u64,
        height: (rect.bottom - rect.top) as u64,
    });
    TRUE
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    if std::env::args().nth(1).as_deref() == Some("version") {
        println!("{}", subcommand::version());
        return Ok(());
    }
    let result = run()?;
    println!("{}", result);
    Ok(())
}

fn run() -> Result<String, Box<dyn std::error::Error>> {
    let _com = ComInitializer::new()?;
    let browse_info = BROWSEINFOW {
        hwndOwner: HWND::default(),
        lpszTitle: w!("壁紙フォルダーを選択してください"),
        ulFlags: BIF_NEWDIALOGSTYLE,
        ..Default::default()
    };
    let pidl = unsafe { SHBrowseForFolderW(&browse_info) };
    if pidl.is_null() {
        return Ok("フォルダーを選択しませんでした".to_string());
    }
    let mut path_buffer = [0u16; MAX_PATH as usize];
    let success = unsafe { SHGetPathFromIDListW(pidl, &mut path_buffer) }.as_bool();
    if !success {
        return Err("SHGetPathFromIDListWでエラー".into());
    }
    let len = path_buffer
        .iter()
        .position(|&value| value == 0)
        .unwrap_or(MAX_PATH as usize);
    let image_path = OsString::from_wide(&path_buffer[..len])
        .to_string_lossy()
        .to_string();

    let mut data = MonitorData {
        rectangles: Vec::new(),
    };
    let result = unsafe {
        EnumDisplayMonitors(
            Some(HDC::default()),
            None,
            Some(monitor_enum_proc),
            LPARAM(&mut data as *mut _ as isize),
        )
    };
    if !result.as_bool() {
        return Err("モニター情報の取得に失敗しました".into());
    }

    config::Config::new(
        image_path,
        data.rectangles
            .iter()
            .map(|rectangle| {
                config::Rectangle::new(
                    rectangle.x as u32,
                    rectangle.y as u32,
                    rectangle.width as u32,
                    rectangle.height as u32,
                )
            })
            .collect(),
    )
    .save()?;

    Ok("".to_string())
}
