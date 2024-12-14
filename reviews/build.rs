fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Create dir src/gen/proto if it doesn't exist
    std::fs::create_dir_all("src/gen/proto").expect("Failed to create dir src/gen/proto");

    println!("Starting to build protos");

    // Iterate over all files in proto/ dir
    for entry in std::fs::read_dir("proto").expect("Failed to read proto dir") {
        let entry = entry?;
        let path = entry.path();
        if path.is_file() {
            println!("Building {:?}.", path);
            // Generate Rust code from proto files
            tonic_build::configure()
                .build_client(false)
                .out_dir("src/gen/proto")
                .compile_protos(&[path], &["proto"])
                .expect("Failed to build proto file!");
        }
    }

    println!("Protofiles have been built!");
    Ok(())
}
