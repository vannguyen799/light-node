use actix_web::{web, App, HttpServer, HttpResponse};
use methods::{GUEST_ELF, GUEST_ID};
use risc0_zkvm::{default_prover, ExecutorEnv, VerifierContext, ProverOpts};
use serde::{Deserialize, Serialize};
use std::fmt;

#[derive(Debug)]
struct ProcessError(String);

impl fmt::Display for ProcessError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}

impl std::error::Error for ProcessError {}

#[derive(Serialize, Deserialize)]
struct Request {
    operation: String,
    data: Vec<String>,
    proof_request: Option<String>,
    proof: Option<MerkleProof>,
}

#[derive(Serialize, Deserialize)]
struct MerkleProof {
    leaf_value: String,
    proof_path: Vec<(String, bool)>,
}

#[derive(Serialize, Deserialize)]
struct Response {
    root: Option<String>,
    proof: Option<MerkleProof>,
    verified: Option<bool>,
    receipt: Option<String>,
}

async fn process(req: web::Json<Request>) -> HttpResponse {
    let request_data = req.0;
    
    let result = web::block(move || {
        println!("Starting request processing...");
        println!("Operation: {}", request_data.operation);
        println!("Input data: {:?}", request_data.data);

        let env = ExecutorEnv::builder()
            .write(&request_data)
            .unwrap()
            .build()
            .map_err(|e| {
                println!("Environment build error: {}", e);
                ProcessError(e.to_string())
            })?;

        println!("Generating proof...");
        
        let receipt = default_prover()
        .prove_with_ctx(
            env,
            &VerifierContext::default(),
            GUEST_ELF,
            &ProverOpts::groth16(),
        )
        .map_err(|e| ProcessError(e.to_string()))?
        .receipt;

        println!("Verifying receipt...");
        receipt.verify(GUEST_ID)
            .map_err(|e| {
                println!("Receipt verification error: {}", e);
                ProcessError(e.to_string())
            })?;

        println!("Decoding journal...");
        let mut output: Response = receipt.journal.decode()
            .map_err(|e| {
                println!("Journal decode error: {}", e);
                ProcessError(e.to_string())
            })?;
        
            let journal_hex = hex::encode(&receipt.journal.bytes);
            output.receipt = Some(journal_hex);
        
        Ok::<Response, ProcessError>(output)
    })
    .await;

    match result {
        Ok(Ok(output)) => HttpResponse::Ok().json(output),
        Ok(Err(e)) => {
            println!("Processing error: {}", e);
            HttpResponse::InternalServerError().json(format!("Processing error: {}", e))
        },
        Err(e) => {
            println!("Blocking error: {}", e);
            HttpResponse::InternalServerError().json(format!("Server error: {}", e))
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv::dotenv().ok();
    println!("Starting server on port 3001...");
    
    HttpServer::new(|| {
        App::new()
            .route("/process", web::post().to(process))
    })
    .bind("127.0.0.1:3001")?
    .run()
    .await
}