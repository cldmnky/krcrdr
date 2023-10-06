# TODO Implement an indexer

The indexer will run FrostDB and write entries in the streams

It will run a leader election loop to have a writer instance that persists to object storage in nats

It will subscribe to streams and insert them into FrostDB

It will answer queries from the index to get records from the stream