use risc0_zkvm::guest::env;
use sha2::{Sha256, Digest};
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize)]
struct MerkleProof {
    leaf_value: String,
    proof_path: Vec<(String, bool)>, // (hash, is_right)
}

#[derive(Serialize, Deserialize)]
struct TreeVisualization {
    level_sizes: Vec<usize>,
    tree_structure: Vec<Vec<String>>,
    data_to_hash_mapping: Vec<(String, String)>,
}

struct MerkleTree {
    leaves: Vec<String>,
    nodes: Vec<Vec<String>>,
    data_values: Vec<String>,  // Store original data values
}

impl MerkleTree {
    fn new() -> Self {
        MerkleTree {
            leaves: Vec::new(),
            nodes: Vec::new(),
            data_values: Vec::new(),
        }
    }

    fn hash(data: &str) -> String {
        let mut hasher = Sha256::new();
        hasher.update(data.as_bytes());
        format!("{:x}", hasher.finalize())
    }

    fn get_visualization(&self) -> TreeVisualization {
        let mut level_sizes = Vec::new();
        let mut tree_structure = Vec::new();
        
        // Add leaves level only once
        level_sizes.push(self.leaves.len());
        tree_structure.push(self.leaves.clone());
        
        // Add internal nodes
        for level in &self.nodes {
            level_sizes.push(level.len());
            tree_structure.push(level.clone());
        }

        // Create data to hash mapping
        let data_to_hash_mapping = self.data_values.iter()
            .zip(self.leaves.iter())
            .map(|(data, hash)| (data.clone(), hash.clone()))
            .collect();

        TreeVisualization {
            level_sizes,
            tree_structure,
            data_to_hash_mapping,
        }
    }

    fn verify_proof(&self, proof: &MerkleProof) -> bool {
        if self.leaves.is_empty() {
            return false;
        }

        let mut current_hash = Self::hash(&proof.leaf_value);

        for (sibling_hash, is_right) in &proof.proof_path {
            let combined = if *is_right {
                format!("{}{}", current_hash, sibling_hash)
            } else {
                format!("{}{}", sibling_hash, current_hash)
            };
            current_hash = Self::hash(&combined);
        }

        Some(current_hash) == self.get_root()
    }

    fn get_root(&self) -> Option<String> {
        self.nodes.last()
            .and_then(|level| level.first())
            .cloned()
            .or_else(|| self.leaves.first().cloned())
    }

    fn verify_external_proof(&self, proof: &MerkleProof, expected_root: &str) -> bool {
        let mut current_hash = Self::hash(&proof.leaf_value);

        for (sibling_hash, is_right) in &proof.proof_path {
            let combined = if *is_right {
                format!("{}{}", current_hash, sibling_hash)
            } else {
                format!("{}{}", sibling_hash, current_hash)
            };
            current_hash = Self::hash(&combined);
        }

        current_hash == expected_root
    }
}

#[derive(Deserialize)]
struct Input {
    operation: String,
    data: Vec<String>,
    proof_request: Option<String>,
    proof: Option<MerkleProof>,
}

#[derive(Serialize)]
struct Output {
    root: Option<String>,
    proof: Option<MerkleProof>,
    verified: Option<bool>,
    visualization: Option<TreeVisualization>,
    receipt: Option<String>,
}

fn main() {
    let input: Input = env::read();
    let tree = MerkleTree::new();
    
    let output = match input.operation.as_str() {
        "verify" => {
            let verified = match (input.proof, input.data.first()) {
                (Some(proof), Some(expected_root)) => {
                    Some(MerkleTree::verify_external_proof(&tree, &proof, expected_root))
                },
                _ => None
            };
            
            Output {
                root: input.data.first().cloned(),
                proof: None,
                verified,
                visualization: None,
                receipt: None,
            }
        },
        _ => Output {
            root: None,
            proof: None,
            verified: None,
            visualization: None,
            receipt: None
        },
    };
    
    env::commit(&output);
}