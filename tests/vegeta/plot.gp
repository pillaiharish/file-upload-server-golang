set terminal dumb
set title "Requests per Second"
set xlabel "Time (s)"
set ylabel "Requests"
set autoscale

plot '< tail -f /tmp/results.bin | vegeta report -type=json | jq -r ".buckets[] | [(.start / 1000000000), .latencies.mean / 1000000] | @csv"' using 1:2 with lines title "Latency (ms)"
