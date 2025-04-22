use windows::{
    Win32::{
        Foundation::{CloseHandle, HANDLE, HWND},
        Security::{GetTokenInformation, TOKEN_ELEVATION, TOKEN_QUERY, TokenElevation},
        System::Threading::{GetCurrentProcess, OpenProcessToken},
        UI::{Shell::ShellExecuteW, WindowsAndMessaging::SW_SHOWNORMAL},
    },
    core::{PCWSTR, w},
};

pub fn is_elevated() -> Result<bool, Box<dyn std::error::Error>> {
    unsafe {
        let mut token_handle = HANDLE::default();
        OpenProcessToken(GetCurrentProcess(), TOKEN_QUERY, &mut token_handle)
            .map_err(|error| format!("OpenProcessToken {:?}", error))?;
        let mut elevation = TOKEN_ELEVATION::default();
        let mut ret_size = 0u32;
        GetTokenInformation(
            token_handle,
            TokenElevation,
            Some(&mut elevation as *mut _ as *mut _),
            std::mem::size_of::<TOKEN_ELEVATION>() as u32,
            &mut ret_size,
        )
        .map_err(|error| format!("GetTokenInformation: {:?}", error))?;
        CloseHandle(token_handle).map_err(|error| format!("CloseHandle: {:?}", error))?;
        Ok(elevation.TokenIsElevated != 0)
    }
}

pub fn run_as_admin() -> Result<(), Box<dyn std::error::Error>> {
    let verb = w!("runas");
    let exe = std::env::current_exe()?
        .to_string_lossy()
        .encode_utf16()
        .chain(Some(0))
        .collect::<Vec<_>>();
    let args = std::env::args()
        .skip(1)
        .collect::<Vec<_>>()
        .join(" ")
        .encode_utf16()
        .chain(Some(0))
        .collect::<Vec<_>>();
    let cwd = std::env::current_dir()?
        .to_string_lossy()
        .encode_utf16()
        .chain(Some(0))
        .collect::<Vec<_>>();
    let result = unsafe {
        ShellExecuteW(
            Some(HWND::default()),
            verb,
            PCWSTR::from_raw(exe.as_ptr()),
            PCWSTR::from_raw(args.as_ptr()),
            PCWSTR::from_raw(cwd.as_ptr()),
            SW_SHOWNORMAL,
        )
    };
    if result.0 as i32 <= 32 {
        return Err(format!("ShellExecuteW: {:?}", result).into());
    }
    return Ok(());
}
