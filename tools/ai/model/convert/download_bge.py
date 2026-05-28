#!/usr/bin/env python3
# ---- comment ----
# Download BGE-small-en-v1.5 locally without optimizations
# ---- comment ----

from transformers import AutoModel, AutoTokenizer

name = "BAAI/bge-small-en-v1.5"
out = "model/bge-small-en-v1.5"

print("Downloading:", name)
model = AutoModel.from_pretrained(name)
tokenizer = AutoTokenizer.from_pretrained(name)

model.save_pretrained(out)
tokenizer.save_pretrained(out)

print("Done →", out)
