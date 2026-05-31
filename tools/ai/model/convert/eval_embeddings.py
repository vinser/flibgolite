#!/usr/bin/env python3
# ---- comment ----
# Evaluate embedding model quality using phrases.yaml
# ---- comment ----

import sys
from pathlib import Path
import numpy as np
import torch
from transformers import AutoTokenizer, AutoModel

from utils import load_yaml, cosine, print_result


def main():
    if len(sys.argv) != 3:
        print("Usage: eval_embeddings.py <hf_model_name> <output_dir>")
        sys.exit(1)

    hf_name = sys.argv[1]
    out_dir = Path(sys.argv[2])
    out_dir.mkdir(parents=True, exist_ok=True)

    cfg = load_yaml(out_dir / "phrases.yaml")
    thresholds = cfg["thresholds"]

    print("Loading model:", hf_name)
    tokenizer = AutoTokenizer.from_pretrained(hf_name)
    model = AutoModel.from_pretrained(hf_name)
    model.eval()

    # Collect embeddings
    all_embeds = {}
    all_texts = {}

    for locale, data in cfg["locales"].items():
        phrases = data["phrases"]
        all_texts[locale] = phrases
        all_embeds[locale] = []

        for text in phrases:
            x = tokenizer(text, return_tensors="pt")
            with torch.no_grad():
                y = model(**x).last_hidden_state[:, 0, :].cpu().numpy()[0]
            all_embeds[locale].append(y)

    # Save ref_multi.txt
    ref_path = out_dir / "ref_multi.txt"
    with open(ref_path, "w", encoding="utf-8") as f:
        for locale, phrases in all_texts.items():
            for text, emb in zip(phrases, all_embeds[locale]):
                f.write(f"{locale} | {text} | " + " ".join(f"{v:.6f}" for v in emb) + "\n")

    print("Saved:", ref_path)

    # Evaluate
    report_lines = []
    passed = True

    for locale, data in cfg["locales"].items():
        phrases = data["phrases"]
        embeds = all_embeds[locale]

        # similar pairs
        for a, b in data["similar_pairs"]:
            ia = phrases.index(a)
            ib = phrases.index(b)
            c = cosine(embeds[ia], embeds[ib])
            ok = c >= thresholds["similar"]
            print_result(f"[{locale}] similar '{a}' ~ '{b}'", ok, c, thresholds["similar"])
            report_lines.append(f"{locale} similar {a} ~ {b}: {c:.6f}")
            passed &= ok

        # different pairs
        for a, b in data["different_pairs"]:
            ia = phrases.index(a)
            ib = phrases.index(b)
            c = cosine(embeds[ia], embeds[ib])
            ok = c <= thresholds["different"]
            print_result(f"[{locale}] different '{a}' x '{b}'", ok, c, thresholds["different"])
            report_lines.append(f"{locale} different {a} x {b}: {c:.6f}")
            passed &= ok

    # Final result
    print("\nMODEL QUALITY:", "PASS" if passed else "FAIL")

    # Save report
    with open(out_dir / "cosine_report.txt", "w", encoding="utf-8") as f:
        f.write("\n".join(report_lines))
        f.write(f"\n\nRESULT: {'PASS' if passed else 'FAIL'}\n")


if __name__ == "__main__":
    main()
