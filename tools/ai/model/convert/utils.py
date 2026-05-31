# ---- comment ----
# Shared utilities for embedding evaluation
# ---- comment ----

import yaml
import numpy as np


def load_yaml(path):
    with open(path, "r", encoding="utf-8") as f:
        return yaml.safe_load(f)


def cosine(a: np.ndarray, b: np.ndarray) -> float:
    a = a.astype(np.float32)
    b = b.astype(np.float32)
    na = np.linalg.norm(a)
    nb = np.linalg.norm(b)
    if na == 0 or nb == 0:
        return 0.0
    return float(np.dot(a, b) / (na * nb))


def print_result(label: str, ok: bool, value: float, threshold: float):
    status = "OK" if ok else "FAIL"
    print(f"{label}: {status} ({value:.6f} vs {threshold:.6f})")
