#!/usr/bin/env python3
# ---- comment ----
# Convert HF model → model.bin (INT8)
# ---- comment ----

import sys
from pathlib import Path
import struct
import numpy as np
import torch
from transformers import AutoModel


SUPPORTED = {
    "BAAI/bge-small-en-v1.5",
    "intfloat/multilingual-e5-small",
    "microsoft/Multilingual-MiniLM-L12-H384",
}


def qlinear(W):
    W = W.astype(np.float32)
    max_abs = np.max(np.abs(W), axis=1)
    scales = max_abs / 127.0
    scales[scales == 0] = 1e-8
    Wq = np.round(W / scales[:, None]).astype(np.int8)
    return Wq, scales.astype(np.float32)


def write_u32(f, v):
    f.write(struct.pack("<I", v))


def write_f32(f, arr):
    f.write(arr.astype(np.float32).tobytes())


def write_i8(f, arr):
    f.write(arr.astype(np.int8).tobytes())


def save_int8_linear(f, W, b):
    out_f, in_f = W.shape
    Wq, scales = qlinear(W)
    write_u32(f, 3)
    write_u32(f, in_f)
    write_u32(f, out_f)
    write_i8(f, Wq)
    write_f32(f, scales)
    write_f32(f, b)


def save_layernorm(f, gamma, beta):
    write_u32(f, 2)
    write_f32(f, gamma)
    write_f32(f, beta)


def main():
    if len(sys.argv) != 3:
        print("Usage: convert_from_hf.py <hf_model_name> <model.bin>")
        sys.exit(1)

    hf_name = sys.argv[1]
    out_path = Path(sys.argv[2])

    if hf_name not in SUPPORTED:
        raise RuntimeError(f"Unsupported model: {hf_name}")

    print("Loading model:", hf_name)
    model = AutoModel.from_pretrained(hf_name)
    sd = model.state_dict()
    cfg = model.config

    # Embeddings
    word = sd["embeddings.word_embeddings.weight"].numpy()
    pos = sd["embeddings.position_embeddings.weight"].numpy()
    tok = sd["embeddings.token_type_embeddings.weight"].numpy()
    gamma_emb = sd["embeddings.LayerNorm.weight"].numpy()
    beta_emb = sd["embeddings.LayerNorm.bias"].numpy()

    vocab, hidden = word.shape
    max_pos = pos.shape[0]
    type_vocab = tok.shape[0]
    num_layers = cfg.num_hidden_layers
    intermediate = sd["encoder.layer.0.intermediate.dense.weight"].shape[0]

    print("Writing model.bin →", out_path)
    with open(out_path, "wb") as f:
        write_u32(f, 0x4C414253)
        write_u32(f, 1)
        write_u32(f, vocab)
        write_u32(f, hidden)
        write_u32(f, intermediate)
        write_u32(f, num_layers)
        write_u32(f, max_pos)
        write_u32(f, type_vocab)

        write_f32(f, word)
        write_f32(f, pos)
        write_f32(f, tok)
        write_f32(f, gamma_emb)
        write_f32(f, beta_emb)

        for i in range(num_layers):
            p = f"encoder.layer.{i}."

            save_int8_linear(f, sd[p + "attention.self.query.weight"].numpy(),
                                sd[p + "attention.self.query.bias"].numpy())
            save_int8_linear(f, sd[p + "attention.self.key.weight"].numpy(),
                                sd[p + "attention.self.key.bias"].numpy())
            save_int8_linear(f, sd[p + "attention.self.value.weight"].numpy(),
                                sd[p + "attention.self.value.bias"].numpy())
            save_int8_linear(f, sd[p + "attention.output.dense.weight"].numpy(),
                                sd[p + "attention.output.dense.bias"].numpy())

            save_int8_linear(f, sd[p + "intermediate.dense.weight"].numpy(),
                                sd[p + "intermediate.dense.bias"].numpy())
            save_int8_linear(f, sd[p + "output.dense.weight"].numpy(),
                                sd[p + "output.dense.bias"].numpy())

            save_layernorm(f,
                           sd[p + "attention.output.LayerNorm.weight"].numpy(),
                           sd[p + "attention.output.LayerNorm.bias"].numpy())
            save_layernorm(f,
                           sd[p + "output.LayerNorm.weight"].numpy(),
                           sd[p + "output.LayerNorm.bias"].numpy())

    print("Done.")


if __name__ == "__main__":
    main()
