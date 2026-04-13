# Requirement: "a multi-threaded video processing framework"

A pipeline that reads frames from a source, runs a worker pool of per-frame transforms, and hands results to a sink in original order.

std
  std.concurrent
    std.concurrent.spawn_worker_pool
      @ (worker_count: i32) -> worker_pool
      + creates a pool of the given size
      # concurrency
    std.concurrent.submit_job
      @ (pool: worker_pool, payload: bytes) -> i64
      + enqueues a job and returns a monotonic job id
      # concurrency
    std.concurrent.collect_result
      @ (pool: worker_pool, job_id: i64) -> result[bytes, string]
      + blocks until the job with that id completes
      - returns error when the pool has been shut down
      # concurrency
    std.concurrent.shutdown_pool
      @ (pool: worker_pool) -> void
      + signals workers to stop and waits for them
      # concurrency

video
  video.open_source
    @ (path: string) -> result[video_source, string]
    + opens a video file and returns a frame source handle
    - returns error when the path cannot be read
    # source
  video.read_frame
    @ (source: video_source) -> optional[video_frame]
    + returns the next frame or none at end of stream
    # source
  video.pipeline_new
    @ (transform: frame_transform, worker_count: i32) -> video_pipeline
    + wires a transform into a worker pool of the given size
    # pipeline
    -> std.concurrent.spawn_worker_pool
  video.pipeline_submit
    @ (pipeline: video_pipeline, frame: video_frame) -> i64
    + submits a frame and returns its sequence id
    # pipeline
    -> std.concurrent.submit_job
  video.pipeline_next_result
    @ (pipeline: video_pipeline) -> optional[video_frame]
    + returns the next processed frame in original submission order
    ? internally reorders out-of-order worker completions
    # pipeline
    -> std.concurrent.collect_result
  video.pipeline_close
    @ (pipeline: video_pipeline) -> void
    + drains pending work and releases worker resources
    # pipeline
    -> std.concurrent.shutdown_pool
