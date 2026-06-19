fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Le indicamos dónde está el archivo .proto que creamos en el paso 1
    tonic_build::compile_protos("../proto/quiniela.proto")?;
    Ok(())
}