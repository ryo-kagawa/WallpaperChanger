use windows::{
    Win32::System::Registry::{
        HKEY, HKEY_CURRENT_USER, KEY_CREATE_SUB_KEY, KEY_SET_VALUE, REG_DWORD,
        REG_OPTION_NON_VOLATILE, RegCloseKey, RegCreateKeyExW, RegSetValueExW,
    },
    core::w,
};

mod subcommand;
mod utils;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    if std::env::args().nth(1).as_deref() == Some("version") {
        println!("{}", subcommand::version());
        return Ok(());
    }
    if !utils::is_elevated()? {
        utils::run_as_admin()?;
        return Ok(());
    }
    let result = run()?;
    println!("{}", result);
    Ok(())
}

fn run() -> Result<String, Box<dyn std::error::Error>> {
    unsafe {
        let mut hkey = HKEY::default();
        let result = RegCreateKeyExW(
            HKEY_CURRENT_USER,
            w!("Control Panel\\Desktop"),
            Some(0),
            None,
            REG_OPTION_NON_VOLATILE,
            KEY_CREATE_SUB_KEY | KEY_SET_VALUE,
            None,
            &mut hkey,
            None,
        );
        if result.is_err() {
            return Err(format!("RegCreateKeyExW: {:?}", result).into());
        }
        let result = RegSetValueExW(
            hkey,
            w!("JPEGImportQuality"),
            Some(0),
            REG_DWORD,
            Some(std::slice::from_raw_parts(
                &100u32 as *const u32 as *const u8,
                std::mem::size_of::<u32>(),
            )),
        );
        if result.is_err() {
            return Err(format!("RegSetValueExW: {:?}", result).into());
        }
        let result = RegCloseKey(hkey);
        if result.is_err() {
            return Err(format!("RegCloseKey: {:?}", result).into());
        }
    }
    Ok("".to_string())
}
