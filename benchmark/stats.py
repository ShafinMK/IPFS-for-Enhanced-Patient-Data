#!/usr/bin/env python3
"""Compute mean ± 95% CI, CV, and Mann–Whitney U for benchmark results (Table 6)."""
import argparse
import pandas as pd
import numpy as np
from scipy import stats


def ci95(x, t=2.045):
    n = len(x)
    s = np.std(x, ddof=1)
    return t * s / np.sqrt(n)


def analyse(df):
    rows = []
    for (op, vol), g in df.groupby(["op", "volume"]):
        tps = g["tps"].values
        lat = g["mean_latency_ms"].values
        rows.append({
            "operation": op, "volume": vol, "n": len(tps),
            "tps_mean": np.mean(tps), "tps_ci": ci95(tps), "tps_cv": np.std(tps, ddof=1)/np.mean(tps)*100,
            "lat_mean": np.mean(lat), "lat_ci": ci95(lat), "lat_cv": np.std(lat, ddof=1)/np.mean(lat)*100,
        })
    return pd.DataFrame(rows)


def mann_whitney(df, op_a, op_b, vol):
    a = df[(df.op == op_a) & (df.volume == vol)]["tps"].values
    b = df[(df.op == op_b) & (df.volume == vol)]["tps"].values
    u, p = stats.mannwhitneyu(a, b, alternative="two-sided")
    return u, p


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--input", required=True)
    ap.add_argument("--output", required=True)
    a = ap.parse_args()
    df = pd.read_csv(a.input)
    out = analyse(df)
    out.to_csv(a.output, index=False)
    print(out.to_string(index=False))
    print("\nMann–Whitney U (query vs invoke, vol=15000):", mann_whitney(df, "query", "invoke", 15000))
    print("Mann–Whitney U (query vs invoke, vol=15):   ", mann_whitney(df, "query", "invoke", 15))


if __name__ == "__main__":
    main()