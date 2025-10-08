# L0 AI Go Agent — Operations Guide

---

## 📚 Table of Contents
- [🚀 Quick Start](#-quick-start)  
- [⏱ Timer Management](#-timer-management)  
- [📝 Logs & Service Status](#-logs--service-status)  
- [🧪 Health & Inference Test](#-health--inference-test)  
- [📌 Model Versions & Swap](#-model-versions--swap)  
- [🧭 Sanity Check](#-sanity-check)  
- [📊 Node Exporter Metrics](#-node-exporter-metrics)  
- [🧹 Maintenance](#-maintenance)  
- [🔒 Immutable Protection](#-immutable-protection)

---

## 🚀 Quick Start

```bash
# Run full pipeline manually (export → train → sanity → swap → restart → health)
systemctl start ai-model-retrain.service

# Stack check (timer + health)
 /root/check_retrain_stack.sh


