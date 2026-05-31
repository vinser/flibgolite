
# **Runtime for BGE‑small‑en‑v1.5 (Pure Go Inference)**

This directory contains the **pure Go inference runtime** for the **BGE‑small‑en‑v1.5** embedding model, converted into a compact binary format (`model.bin`) with **INT8 per‑channel quantized linear layers**.

The runtime performs the full forward pass:

- WordPiece tokenization  
- Embedding lookup  
- 12‑layer Transformer encoder  
- Mean pooling  
- L2 normalization  

All operations are implemented in **pure Go**, without CGO or external dependencies.

---

## **1. Model Format (`model.bin`)**

The model is stored in a custom binary format optimized for:

- sequential memory access  
- zero‑copy loading  
- minimal RAM usage  
- fast INT8 matrix multiplications  

### **1.1 Header (32 bytes)**

| Offset | Type      | Meaning |
|--------|-----------|---------|
| 0      | uint32    | Magic = `0x4C414253` (`"LABS"`) |
| 4      | uint32    | Version = 1 |
| 8      | uint32    | Vocabulary size |
| 12     | uint32    | Hidden size (384) |
| 16     | uint32    | Intermediate size (1536) |
| 20     | uint32    | Number of layers (12) |
| 24     | uint32    | Max position embeddings |
| 28     | uint32    | Token type vocab size (2) |

### **1.2 Embeddings Block**

Stored as float32:

- word_embeddings: `[vocab][hidden]`
- position_embeddings: `[max_pos][hidden]`
- token_type_embeddings: `[2][hidden]`
- LayerNorm gamma: `[hidden]`
- LayerNorm beta: `[hidden]`

### **1.3 Encoder Layer Blocks (12 layers)**

Each layer contains **6 INT8 linear projections**:

Order per layer:

1. Q  
2. K  
3. V  
4. O  
5. FFN_up  
6. FFN_down  

Each INT8 linear block:

```
uint32 type = 3
uint32 in_features
uint32 out_features
int8   weights[out][in]
float32 scales[out]
float32 bias[out]
```

### **1.4 LayerNorm Blocks**

Stored as float32:

```
uint32 type = 2
float32 gamma[hidden]
float32 beta[hidden]
```

Each encoder layer has:

- attention.output.LayerNorm
- output.LayerNorm

### **1.5 No Pooler**

BGE‑small‑en‑v1.5 does **not** include:

- pooler  
- final LayerNorm  

The final embedding is produced via **mean pooling**.

---

## **2. Tokenizer**

The runtime uses **BERT WordPiece** tokenization:

- lowercase normalization  
- strip accents  
- WordPiece vocabulary (`vocab.txt`)  
- special tokens: `[CLS]`, `[SEP]`, `[PAD]`, `[UNK]`  
- max sequence length = value from header  

Output:

- token IDs (`[]int32`)  
- attention mask (`[]float32`)

---

## **3. Execution Pipeline**

### **3.1 Embedding Lookup**

For each token:

```
x = word_emb[token] + position_emb[pos] + token_type_emb[0]
x = LayerNorm(x)
```

### **3.2 Encoder Layer (12×)**

Each layer performs:

#### **Self‑Attention**

- Q = Linear(x)
- K = Linear(x)
- V = Linear(x)

Split into 6 heads:

```
head_dim = 64
num_heads = 6
```

Attention:

```
scores = (Q * Kᵀ) / sqrt(64)
scores += mask
weights = softmax(scores)
context = weights * V
```

Output projection:

```
x = Linear(context)
x = LayerNorm(x + residual)
```

#### **Feed‑Forward Network**

```
h = GELU(Linear_up(x))
x = Linear_down(h)
x = LayerNorm(x + residual)
```

### **3.3 Mean Pooling**

```
embedding = mean(hidden_states[token_mask])
```

### **3.4 L2 Normalization**

```
embedding = embedding / ||embedding||
```

Output: `[]float32` of length **384**.

---

## **4. Performance Characteristics**

- pure Go, no CGO  
- no heap allocations inside the hot path  
- sequential memory access  
- int8 GEMM with per‑channel scales  
- deterministic results  

---

## **5. Directory Structure**

```
internal/ai/runtime/
│
├── model.go        # model loading and header parsing
├── layers.go       # Linear(int8), LayerNorm, Attention, FFN
├── tokenizer.go    # WordPiece tokenizer
└── embed.go        # full forward pass (Embed function)
```

---

## **6. Usage Example**

```go
model, _ := LoadModel(modelData)
vec := model.Embed("A book about a young wizard")
```

The result is a normalized 384‑dimensional embedding suitable for:

- semantic search  
- similarity ranking  
- clustering  
- recommendations  

---

## **7. Limitations**

- Only supports BGE‑small‑en‑v1.5 architecture  
- Only INT8 linear layers  
- No batching  
- Max sequence length fixed by model header  

---
