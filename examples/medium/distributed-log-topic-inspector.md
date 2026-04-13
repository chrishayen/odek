# Requirement: "a library for inspecting topics, partitions, and messages in a distributed log"

Thin wrapper around a message broker client, exposing topic metadata and bounded message reads.

std
  std.broker
    std.broker.connect
      @ (brokers: list[string]) -> result[broker_client, string]
      + opens a connection to one of the seed brokers
      - returns error when none of the brokers accept the connection
      # broker
    std.broker.list_topics
      @ (c: broker_client) -> result[list[string], string]
      + returns every topic name the broker knows about
      # broker
    std.broker.describe_topic
      @ (c: broker_client, topic: string) -> result[list[partition_info], string]
      + returns partition id, leader, and high-water offset for each partition
      - returns error when the topic does not exist
      # broker
    std.broker.fetch
      @ (c: broker_client, topic: string, partition: i32, offset: i64, max_messages: i32) -> result[list[broker_message], string]
      + returns up to max_messages starting at offset
      - returns error when offset is past the high-water mark
      # broker

topic_inspector
  topic_inspector.list_topics
    @ (c: broker_client) -> result[list[string], string]
    + returns topic names in lexicographic order
    # listing
    -> std.broker.list_topics
  topic_inspector.partition_summary
    @ (c: broker_client, topic: string) -> result[list[partition_info], string]
    + returns one row per partition including its message count (high_water - low_water)
    # partitions
    -> std.broker.describe_topic
  topic_inspector.tail
    @ (c: broker_client, topic: string, partition: i32, count: i32) -> result[list[broker_message], string]
    + returns the last count messages on the partition
    - returns error when count is less than 1
    # tailing
    -> std.broker.describe_topic
    -> std.broker.fetch
