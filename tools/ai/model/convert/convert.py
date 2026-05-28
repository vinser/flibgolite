#!/usr/bin/env python3
# ---- comment ----
# BERT-like ONNX (PyTorch 2.3+) → model.bin converter
# All weights come from val_* in strict Q,K,V,O,FFN_up,FFN_down order
# ---- comment ----

import onnx
import numpy as np
import struct
import sys
from pathlib import Path
from onnx import numpy_helper


# ============================================================
#  Helpers
# ============================================================

def qlinear_per_channel(W: np.ndarray):
    W = W.astype(np.float32)
    max_abs = np.max(np.abs(W), axis=1)
    scales = max_abs / 127.0
    scales[scales == 0] = 1e-8
    W_int8 = np.round(W / scales[:, None]).astype(np.int8)
    return W_int8, scales.astype(np.float32)


def write_f32(f, arr):
    f.write(arr.astype(np.float32).tobytes())


def write_i8(f, arr):
    f.write(arr.astype(np.int8).tobytes())


def write_u32(f, v):
    f.write(struct.pack("<I", v))


# ============================================================
#  Load initializers
# ============================================================

def load_initializers(model):
    inits = {}
    for t in model.graph.initializer:
        inits[t.name] = numpy_helper.to_array(t)
    return inits


def get_init(inits, name):
    if name not in inits:
        raise KeyError(name)
    return inits[name]


# ============================================================
#  Save blocks
# ============================================================

def save_int8_linear(f, W, b):
    out_features, in_features = W.shape
    W_int8, scales = qlinear_per_channel(W)

    write_u32(f, 3)
    write_u32(f, in_features)
    write_u32(f, out_features)
    write_i8(f, W_int8)
    write_f32(f, scales)
    write_f32(f, b)


def save_layernorm(f, gamma, beta):
    write_u32(f, 2)
    write_f32(f, gamma)
    write_f32(f, beta)


def save_embeddings(f, inits, vocab_size, hidden, max_pos, type_vocab):
    write_u32(f, 1)

    write_f32(f, get_init(inits, "embeddings.word_embeddings.weight"))
    write_f32(f, get_init(inits, "embeddings.position_embeddings.weight"))
    write_f32(f, get_init(inits, "embeddings.token_type_embeddings.weight"))

    write_f32(f, get_init(inits, "embeddings.LayerNorm.weight"))
    write_f32(f, get_init(inits, "embeddings.LayerNorm.bias"))


# ============================================================
#  Main conversion
# ============================================================

def convert(model_path: Path, out_path: Path):
    print(f"Loading ONNX: {model_path}")
    model = onnx.load(str(model_path))
    inits = load_initializers(model)

    # ---- Collect val_* matrices ----
    val_mats = sorted(
        [(k, v) for k, v in inits.items() if k.startswith("val_") and v.ndim == 2],
        key=lambda x: int(x[0].split("_")[1])
    )
    val_names = [name for name, _ in val_mats]

    # 6 matrices per layer: Q,K,V,O,FFN_up,FFN_down
    assert len(val_names) == 12 * 6, "Unexpected number of val_* matrices"

    # ---- Infer dims ----
    emb_w = get_init(inits, "embeddings.word_embeddings.weight")
    pos_w = get_init(inits, "embeddings.position_embeddings.weight")
    tok_w = get_init(inits, "embeddings.token_type_embeddings.weight")

    vocab_size = emb_w.shape[0]
    hidden = emb_w.shape[1]
    max_pos = pos_w.shape[0]
    type_vocab = tok_w.shape[0]

    # FFN intermediate size = FFN_up rows
    intermediate = val_mats[4][1].shape[0]

    num_layers = 12

    print("Writing model.bin")
    with open(out_path, "wb") as f:

        # ---- Header ----
        write_u32(f, 0x4C414253)
        write_u32(f, 1)
        write_u32(f, vocab_size)
        write_u32(f, hidden)
        write_u32(f, intermediate)
        write_u32(f, num_layers)
        write_u32(f, max_pos)
        write_u32(f, type_vocab)

        # ---- Embeddings ----
        save_embeddings(f, inits, vocab_size, hidden, max_pos, type_vocab)

        # ---- Encoder layers ----
        idx = 0
        for i in range(num_layers):
            prefix = f"encoder.layer.{i}."

            # Q, K, V, O, FFN_up, FFN_down
            Wq = inits[val_names[idx]]; idx += 1
            Wk = inits[val_names[idx]]; idx += 1
            Wv = inits[val_names[idx]]; idx += 1
            Wo = inits[val_names[idx]]; idx += 1
            Wup = inits[val_names[idx]]; idx += 1
            Wdown = inits[val_names[idx]]; idx += 1

            # Biases exist normally
            bq = get_init(inits, prefix + "attention.self.query.bias")
            bk = get_init(inits, prefix + "attention.self.key.bias")
            bv = get_init(inits, prefix + "attention.self.value.bias")
            bo = get_init(inits, prefix + "attention.output.dense.bias")

            bup = get_init(inits, prefix + "intermediate.dense.bias")
            bdown = get_init(inits, prefix + "output.dense.bias")

            # Save attention
            save_int8_linear(f, Wq, bq)
            save_int8_linear(f, Wk, bk)
            save_int8_linear(f, Wv, bv)
            save_int8_linear(f, Wo, bo)

            # Save FFN
            save_int8_linear(f, Wup, bup)
            save_int8_linear(f, Wdown, bdown)

            # LayerNorms
            save_layernorm(
                f,
                get_init(inits, prefix + "attention.output.LayerNorm.weight"),
                get_init(inits, prefix + "attention.output.LayerNorm.bias"),
            )

            save_layernorm(
                f,
                get_init(inits, prefix + "output.LayerNorm.weight"),
                get_init(inits, prefix + "output.LayerNorm.bias"),
            )

    print(f"Done → {out_path}")

# ============================================================
#  Entry point
# ============================================================

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: convert.py model.onnx model.bin")
        sys.exit(1)

    convert(Path(sys.argv[1]), Path(sys.argv[2]))
