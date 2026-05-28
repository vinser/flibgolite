#!/usr/bin/env python3
import torch
from transformers import AutoModel
from pathlib import Path

model_dir = Path("model/bge-small-en-v1.5")
out_path = model_dir / "model.onnx"

print("Loading model from:", model_dir)
model = AutoModel.from_pretrained(model_dir)
model.eval()

dummy = torch.randint(0, 100, (1, 16))

torch.onnx.export(
    model,
    (dummy,),
    out_path.as_posix(),
    input_names=["input_ids"],
    output_names=["last_hidden_state"],
    opset_version=17,
    do_constant_folding=False,
    keep_initializers_as_inputs=True,
    dynamic_axes={"input_ids": {0: "batch", 1: "seq"}},
)

print("Saved clean ONNX →", out_path)
