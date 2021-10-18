# PMLog

**Requirements**
- [Go](https://golang.org/) (version >= 1.15). If you're using Linux, you can use `./scripts/install_go.sh && source /etc/profile`.
- [Make](https://www.gnu.org/software/make/)


### Installation on each node

```shell
git clone https://github.com/nathanieltornow/PMLog.git && cd PMLog && ./scripts/install_go.sh && source /etc/profile && source ~/.profile
```

### Setup benchmark

1. Start the sequencer on one node
```shell
go run cmd/start_sequencer/start_sequencer.go -root -IP :7000
```

2. Startup the shard
```shell
# for ever replica
go run cmd/replica/start_replica.go -order <sequencer-IP>:7000 -IP :4000
```

3. Start a benchmark at the load generator
   1. Modify benchmark/benchmark.config.yaml, that the endpoints are all pointing to the replicas
   2. Modify "clients" benchmark/benchmark.config.yaml, for the number of clients
   3. Modify "appends" and "reads" for the right ratio. (Just the ratio appends/reads is important)
   4. Start the benchmark, scale the throughput with increasing x
   ```shell
   go run benchmark/shared_log/shared_log.go -config benchmark/benchmark.config.yaml -threads x
   ```

