#!/usr/bin/env python3
# ---- comment ----
# Print all initializer tensor names from an ONNX model
# ---- comment ----

import onnx
import sys

if len(sys.argv) != 2:
    print("Usage: dump_inits.py <model.onnx>")
    sys.exit(1)

path = sys.argv[1]
m = onnx.load(path)

names = sorted(t.name for t in m.graph.initializer)
for n in names:
    print(n)
