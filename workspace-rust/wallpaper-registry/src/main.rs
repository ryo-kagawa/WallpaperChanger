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
    windows_registry::CURRENT_USER
        .create("Control Panel\\Desktop")
        .map_err(|error| format!("RegCreateKeyExW: {}", error))?
        .set_u32("JPEGImportQuality", 100)
        .map_err(|error| format!("RegSetValueExW: {}", error))?;
    Ok("".to_string())
}
