# L0 AI Go Agent — Operations Guide

## Quick Start
bash
# Run full pipeline manually (export → train → sanity → swap → restart → health)
systemctl start ai-model-retrain.service

# Stack check (timer + health)
 /root/check_retrain_stack.sh

## Timer Management
systemctl list-timers ai-model-retrain.timer
systemctl enable  ai-model-retrain.timer     # enable scheduled retrain
systemctl disable ai-model-retrain.timer     # disable schedule
systemctl start   ai-model-retrain.timer     # trigger timer now
systemctl stop    ai-model-retrain.timer

## Logs & Service Status
# Show last 200 lines of retrain service log
journalctl -u ai-model-retrain.service -n 200 --no-pager

# Systemd unit status
systemctl status ai-model-retrain.service --no-pager

## Health & Inference Test
# Health endpoint
curl -fsS http://localhost:18080/healthz && echo "health OK" || echo "health FAIL"

# Example inference request
curl -s http://localhost:18080/predict -H 'content-type: application/json' \
  -d '{"avg_rtt_ms":55,"min_rtt_ms":22,"max_rtt_ms":130,"avg_loss_pct":0.6,"dup_ratio":0.12,"mempool_size":1520,"sock_queue":3,"overlay_load":0.45}'

## Model Versions & Swap
# Check current model symlink
ls -l /root/latency_model_ext.joblib

# Show last 5 archived models
ls -ltr /root/models | tail -n 5

# Manual atomic swap (symlink rotation)
 /root/ai_model_symlink_swap.sh /tmp/new_model.joblib /root/latency_model_ext.joblib

## Sanity Check
# Threshold configuration
cat /etc/ai-retrain.conf
MAX_TEST_RMSE=0.5
MIN_TEST_R2=0.95
ALLOW_WORSE_PCT=20

# Check a specific train log
 /root/ai_retrain_sanity.sh /var/lib/ai-retrain/train_YYYYmmdd_HHMMSS.log
echo $?   # 0 = OK, !=0 = FAIL

## Node Exporter Metrics
# Check generated metrics file
cat /var/lib/node_exporter/textfile_collector/ai_retrain.prom

# Force re-export of metrics
/root/ai_retrain_exporter.sh
cat /var/lib/node_exporter/textfile_collector/ai_retrain.prom

## Maintenance
# Force log rotation (if configured)
/usr/sbin/logrotate -df /etc/logrotate.d/ai-retrain

# Cleanup old models (>30 days)
find /root/models -type f -name 'latency_model_ext_*.joblib' -mtime +30 -print -delete

## Recommended File Locations
| Path                                         | Purpose                          |
| -------------------------------------------- | -------------------------------- |
| `/root/retrain_and_reload.sh`                | Main retraining pipeline script  |
| `/root/models/`                              | Archived models with timestamps  |
| `/var/lib/ai-retrain/`                       | Train logs and sanity check data |
| `/var/lib/node_exporter/textfile_collector/` | Prometheus textfile metrics      |
| `/etc/ai-retrain.conf`                       | Sanity check thresholds          |




