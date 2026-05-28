# **Binary Model Format (`model.bin`)**

This directory contains the binary representation of the **BGE‑small‑en‑v1.5** embedding model, converted from ONNX into a compact format optimized for:

- embedding directly into the Go binary (`go:embed`)
- zero‑dependency loading
- fast sequential parsing
- pure Go inference (no CGO, no external libraries)
- INT8 per‑channel quantized linear layers

The format is intentionally minimal, deterministic, and stable.

---

## **1. File Layout**

```
model.bin
│
├── Header (32 bytes)
└── Weights (mixed float32 + int8 blocks)
```

The header contains model metadata.  
The weights section contains all tensors in a strict, fixed order.

---

## **2. Header Format (32 bytes)**

All fields are **little‑endian**.

| Offset | Size | Type   | Name                    | Description |
|--------|------|--------|-------------------------|-------------|
| 0      | 4    | uint32 | Magic                   | `"LABS"` = 0x4C414253 |
| 4      | 4    | uint32 | Version                 | Format version (1) |
| 8      | 4    | uint32 | VocabSize               | Token vocabulary size |
| 12     | 4    | uint32 | HiddenSize              | Hidden dimension (384) |
| 16     | 4    | uint32 | IntermediateSize        | FFN dimension (1536) |
| 20     | 4    | uint32 | NumLayers               | Number of encoder layers (12) |
| 24     | 4    | uint32 | MaxPositionEmbeddings   | Maximum positional index |
| 28     | 4    | uint32 | TokenTypeVocabSize      | Usually 2 |

The weights section begins immediately after the 32‑byte header.

---

## **3. Weights Section**

The weights section contains:

- **float32** tensors (embeddings, LayerNorm)
- **INT8 linear blocks** (Q/K/V/O, FFN_up, FFN_down)

The order is **strict** and must match exactly.

---

## **3.1 Embedding Matrices (float32)**

1. `embeddings.word_embeddings.weight`  
   Shape: `[VocabSize][Hidden]`

2. `embeddings.position_embeddings.weight`  
   Shape: `[MaxPositionEmbeddings][Hidden]`

3. `embeddings.token_type_embeddings.weight`  
   Shape: `[TokenTypeVocabSize][Hidden]`

4. `embeddings.LayerNorm.weight`  
   Shape: `[Hidden]`

5. `embeddings.LayerNorm.bias`  
   Shape: `[Hidden]`

---

## **3.2 Encoder Layers (12 layers)**

Each encoder layer contains **6 INT8 linear projections** followed by **2 LayerNorm blocks**.

### **Order per layer (strict):**

1. **Q** projection (INT8 Linear)  
2. **K** projection (INT8 Linear)  
3. **V** projection (INT8 Linear)  
4. **O** projection (INT8 Linear)  
5. **FFN_up** (INT8 Linear)  
6. **FFN_down** (INT8 Linear)  
7. **attention.output.LayerNorm** (float32)  
8. **output.LayerNorm** (float32)

---

## **3.3 INT8 Linear Block Format**

Each INT8 linear block is stored as:

```
uint32 type = 3
uint32 in_features
uint32 out_features
int8    weights[out][in]
float32 scales[out]
float32 bias[out]
```

- `weights` are per‑channel quantized  
- `scales` are float32 per output channel  
- `bias` is float32  

This matches the converter and the Go runtime.

---

## **3.4 LayerNorm Block Format**

Each LayerNorm block:

```
uint32 type = 2
float32 gamma[Hidden]
float32 beta[Hidden]
```

---

## **3.5 No Final LayerNorm / No Pooler**

BGE‑small‑en‑v1.5 does **not** include:

- pooler  
- final LayerNorm  

The final embedding is produced via **mean pooling** in the runtime.

---

## **4. Data Types**

- Embeddings: **float32**
- LayerNorm: **float32**
- Linear weights: **int8**
- Linear scales: **float32**
- Linear bias: **float32**

All values are little‑endian.

---

## **5. Format Versioning**

The `Magic` field identifies the format:

- `"LABS"` — current version (INT8 linear, MiniLM architecture)

Future versions may introduce:

- `"LAB2"` — float16 support  
- `"LAB3"` — full int8 model  
- `"LAB4"` — RoPE positional embeddings  

---

## **6. Converting ONNX → model.bin**

The converter is located at:

```
tools/ai/model/convert/convert.py
```

Usage:

```
python3 convert.py model.onnx model.bin
```

The converter:

- loads the ONNX model  
- extracts all tensors  
- maps `val_*` matrices into Q/K/V/O/FFN_up/FFN_down  
- writes the 32‑byte header  
- writes all weights in the strict order described above  

---

## **7. Loading the Model in Go**

The model is embedded using:

```go
//go:embed ../model/model.bin
var modelData []byte
```

Loading:

```go
m, err := LoadModel(modelData)
```

The loader:

- validates the `"LABS"` header  
- reads metadata  
- creates typed slices referencing the weight buffer  
- performs zero allocations  

---

## **8. Compatibility Requirements**

- The model must follow the BGE‑small‑en‑v1.5 architecture  
- All tensor shapes must match the header  
- The weight order must be exact  
- The runtime assumes:
  - 12 layers  
  - 6 attention heads  
  - hidden = 384  
  - intermediate = 1536  
  - INT8 linear blocks  

---

## **9. Notes**

- ONNX files are not stored in the repository  
- `model.bin` is the only artifact required at runtime  
- The format is optimized for:
  - sequential parsing  
  - minimal memory overhead  
  - deterministic inference  

---
