fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Iterate over all files in proto/ dir
    for entry in std::fs::read_dir("proto")? {
        let entry = entry?;
        let path = entry.path();
        if path.is_file() {
            // Generate Rust code from proto files
            tonic_build::compile_protos(path)?;
        }
    }

    Ok(())
}
