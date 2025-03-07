use actix_web::{web, App, HttpServer, HttpResponse};
use methods::{GUEST_ELF, GUEST_ID};
use risc0_zkvm::{default_prover, ExecutorEnv, VerifierContext, ProverOpts};
use serde::{Deserialize, Serialize};
use std::fmt;
use std::io::{self, Read, Write};
use bincode::{serialize_into, deserialize_from};

#[derive(Debug)]
struct ProcessError(String);

impl fmt::Display for ProcessError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}

impl std::error::Error for ProcessError {}

#[derive(Serialize, Deserialize)]
struct NodeHeader {
    length: u32,
    level: u32,
}

#[derive(Serialize, Deserialize)]
struct TreeNode {
    header: NodeHeader,
    data: String,
    hash: String,
}

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
struct TreeVisualization {
    level_sizes: Vec<usize>,
    tree_structure: Vec<Vec<String>>,
    data_to_hash_mapping: Vec<(String, String)>,
}

#[derive(Serialize, Deserialize)]
struct Response {
    root: Option<String>,
    proof: Option<MerkleProof>,
    verified: Option<bool>,
    visualization: Option<TreeVisualization>,
    receipt: Option<String>,
}

impl TreeNode {
    fn new(data: String, hash: String, level: u32) -> Self {
        TreeNode {
            header: NodeHeader {
                length: (data.len() + hash.len()) as u32,
                level,
            },
            data,
            hash,
        }
    }
    fn serialize<W: Write>(&self, writer: &mut W) -> io::Result<()> {
        serialize_into(&mut *writer, &self.header)
            .map_err(|e| io::Error::new(io::ErrorKind::Other, e))?;
        writer.write_all(self.data.as_bytes())?;
        writer.write_all(self.hash.as_bytes())?;
        Ok(())
    }

    fn deserialize<R: Read>(reader: &mut R) -> io::Result<Self> {
        let header: NodeHeader = deserialize_from(&mut *reader)
            .map_err(|e| io::Error::new(io::ErrorKind::Other, e))?;
        
        let mut data_and_hash = vec![0u8; header.length as usize];
        reader.read_exact(&mut data_and_hash)?;
        
        let split_point = data_and_hash.len() / 2;
        let (data_bytes, hash_bytes) = data_and_hash.split_at(split_point);
        
        let data = String::from_utf8(data_bytes.to_vec())
            .map_err(|e| io::Error::new(io::ErrorKind::InvalidData, e))?;
        let hash = String::from_utf8(hash_bytes.to_vec())
            .map_err(|e| io::Error::new(io::ErrorKind::InvalidData, e))?;
        
        Ok(TreeNode { header, data, hash })
    }
}

fn serialize_tree(nodes: &[Vec<TreeNode>]) -> io::Result<Vec<Vec<u8>>> {
    nodes.iter().map(|level| {
        let mut buffer = Vec::new();
        for node in level {
            node.serialize(&mut buffer)?;
        }
        Ok(buffer)
    }).collect()
}

fn deserialize_tree(serialized: &[Vec<u8>]) -> io::Result<Vec<Vec<TreeNode>>> {
    serialized.iter().map(|level_bytes| {
        let mut cursor = io::Cursor::new(level_bytes);
        let mut nodes = Vec::new();
        while cursor.position() < level_bytes.len() as u64 {
            nodes.push(TreeNode::deserialize(&mut cursor)?);
        }
        Ok(nodes)
    }).collect()
}

async fn process(req: web::Json<Request>) -> HttpResponse {
    let request_data = req.0;
    
    let result = web::block(move || {
        println!("Starting request processing...");
        println!("Operation: {}", request_data.operation);
        
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
        
        output.receipt = format!("{}", format!("{:?}", receipt).len()).into();
        
        if let Some(viz) = &output.visualization {
            println!("\nProcessing tree for serialization...");
            let mut tree_nodes = Vec::new();
            
            for (i, level) in viz.tree_structure.iter().enumerate() {
                println!("\nLevel {}", i);
                let mut level_nodes = Vec::new();
                for node_str in level {
                    let hash = viz.data_to_hash_mapping
                        .iter()
                        .find(|(data, _)| data == node_str)
                        .map(|(_, hash)| hash.clone())
                        .unwrap_or_else(|| format!("hash_{}", node_str));
                        
                    let node = TreeNode::new(node_str.clone(), hash.clone(), i as u32);
                    println!("Created node: {} -> {}", node_str, hash);
                    level_nodes.push(node);
                }
                tree_nodes.push(level_nodes);
            }

            println!("\nSerializing tree...");
            match serialize_tree(&tree_nodes) {
                Ok(serialized) => {
                    println!("Serialization successful!");
                    for (i, level) in serialized.iter().enumerate() {
                        println!("Level {}: {} bytes", i, level.len());
                    }
                    
                    println!("\nAttempting deserialization...");
                    match deserialize_tree(&serialized) {
                        Ok(deserialized) => {
                            println!("Deserialization successful!");
                            for (i, level) in deserialized.iter().enumerate() {
                                println!("Level {}: {} nodes", i, level.len());
                                for node in level {
                                    println!("  Node: {} (hash: {})", node.data, node.hash);
                                }
                            }
                        }
                        Err(e) => println!("Deserialization error: {}", e),
                    }
                }
                Err(e) => println!("Serialization error: {}", e),
            }
        }
        
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
    println!("Starting server on port 3000...");
    
    HttpServer::new(|| {
        App::new()
            .route("/process", web::post().to(process))
    })
    .bind("127.0.0.1:3000")?
    .run()
    .await
}