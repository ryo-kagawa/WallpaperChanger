fn main() {
    let rustc_version = String::from_utf8_lossy(
        &std::process::Command::new("rustc")
            .arg("--version")
            .output()
            .expect("エラー: \"rustc --version\"")
            .stdout,
    )
    .split_whitespace()
    .nth(1)
    .unwrap()
    .to_string();

    let metadata_dependencies = serde_json::from_slice::<serde_json::Value>(
        &std::process::Command::new("cargo")
            .args(["metadata", "--format-version=1"])
            .output()
            .unwrap()
            .stdout,
    )
    .unwrap()
    .get("packages")
    .unwrap()
    .as_array()
    .unwrap()
    .iter()
    .map(|package| {
        format!(
            "{} {}",
            package["name"].as_str().unwrap(),
            package["version"].as_str().unwrap()
        )
    })
    .collect::<Vec<_>>()
    .join("\\n");

    println!(
        "cargo:rustc-env=VERSION_INFO={}",
        [
            format!("WallpaperRegistry.exe {}", env!("CARGO_PKG_VERSION")),
            format!("rustc {}", rustc_version),
            format!("{}", metadata_dependencies)
        ]
        .join("\\n")
    );

    if std::env::var("PROFILE").unwrap() == "release" {
        println!("cargo:rustc-cfg=release_build");
    }
}
