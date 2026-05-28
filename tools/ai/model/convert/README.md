
# **BGE Model Conversion Pipeline (`tools/ai/model/convert`)**

This directory contains the full toolchain required to download, export, inspect, and convert the **BGE‑small‑en‑v1.5** model into the compact binary format (`model.bin`) used by the pure Go runtime.

The pipeline consists of:

0. **Prepare Python virtual environment**  
1. **Download** the HuggingFace model  
2. **Export** it into a clean ONNX representation  
3. **Inspect** ONNX initializers (optional)  
4. **Convert** ONNX → `model.bin` (INT8 + float32 mixed format)

The final output (`model.bin`) is embedded into the Go binary via `go:embed`.

---

# **0. Prepare Python Virtual Environment (required)**

Create a clean venv:

```bash
python3 -m venv venv
source venv/bin/activate
```

Upgrade pip:

```bash
pip install --upgrade pip
```

Install dependencies:

```bash
pip install -r requirements.txt
```

### **Contents of `requirements.txt`:**

```
torch
transformers
onnx
onnxscript
```

These versions are compatible with:

- PyTorch 2.3+ ONNX exporter  
- HuggingFace BGE models  
- ONNX opset used by BGE‑small‑en‑v1.5  

---

# **1. Download HuggingFace Model**

```
python3 download_bge.py
```

This script downloads into:

```
model/bge-small-en-v1.5/
    config.json
    model.safetensors
    tokenizer.json
    tokenizer_config.json
```

These files are required for ONNX export.

---

# **2. Export Clean ONNX**

```
python3 export_bge_onnx.py
```

This script:

- loads the safetensors model  
- exports a **clean ONNX** using PyTorch 2.3+  
- writes:

```
model/bge-small-en-v1.5/model.onnx
model/bge-small-en-v1.5/model.onnx.data
```

The ONNX file contains:

- named tensors (embeddings, LayerNorm, biases)  
- unnamed tensors (`val_*`) for all linear weights  

This is expected and required.

---

# **3. Inspect ONNX Initializers (optional)**

```
python3 dump_inits.py
```

This prints:

- all initializer names  
- shapes  
- mapping of `val_*` matrices  
- MatMul graph structure  

Useful for verifying:

- Q  
- K  
- V  
- O  
- FFN_up  
- FFN_down  

---

# **4. Convert ONNX → model.bin**

```
python3 convert.py model/bge-small-en-v1.5/model.onnx model.bin
```

This script:

- loads ONNX  
- extracts all initializers  
- maps `val_*` matrices into the correct order  
- quantizes linear layers to **INT8 per‑channel**  
- writes the **32‑byte LABS header**  
- writes all weights in strict order  
- produces `model.bin`

The resulting file is consumed by the Go runtime.

---

# **5. Output Format (`model.bin`)**

The binary format is documented in:

```
tools/ai/model/README.md
```

Summary:

- **Header (32 bytes)**  
- **float32 embeddings**  
- **float32 LayerNorm**  
- **INT8 linear blocks**  
- **no final LayerNorm**  
- **no pooler**  

Order per encoder layer:

1. Q  
2. K  
3. V  
4. O  
5. FFN_up  
6. FFN_down  
7. attention LayerNorm  
8. output LayerNorm  

---

# **6. Notes**

- The converter does **not** use safetensors directly — only ONNX.  
- The ONNX file is **not** stored in the repository.  
- The converter must be re‑run whenever the model or ONNX exporter changes.  
- The Go runtime assumes the exact binary layout produced by this converter.  
- The mapping of `val_*` matrices is deterministic and based on ONNX MatMul order.

---

# **7. Directory Structure**

```
tools/ai/model/convert/
│
├── download_bge.py
├── export_bge_onnx.py
├── dump_inits.py
├── convert.py
│
├── model/
│   └── bge-small-en-v1.5/
│       ├── config.json
│       ├── model.safetensors
│       ├── tokenizer.json
│       ├── tokenizer_config.json
│       ├── model.onnx
│       └── model.onnx.data
│
├── model.bin
└── requirements.txt
```

---
