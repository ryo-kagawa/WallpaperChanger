pub fn version() -> String {
    env!("VERSION_INFO").replace("\\n", "\n")
}
