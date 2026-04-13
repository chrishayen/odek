# Requirement: "an optical character recognition library supporting many languages"

OCR runs a text detection model to find text regions, then a recognition model per region using a language-specific character set. The library exposes a high-level recognize call.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns the file contents as bytes
      - returns error when the file cannot be read
      # filesystem
  std.image
    std.image.decode
      @ (data: bytes) -> result[image, string]
      + decodes a PNG or JPEG byte buffer into an image
      - returns error when the format is unsupported
      # image
    std.image.to_grayscale
      @ (img: image) -> image
      + returns a single-channel grayscale image
      # image
    std.image.resize
      @ (img: image, width: i32, height: i32) -> image
      + returns a bilinearly resized image
      # image
    std.image.crop
      @ (img: image, x: i32, y: i32, w: i32, h: i32) -> result[image, string]
      + returns a cropped sub-image
      - returns error when the rectangle is out of bounds
      # image
  std.ml
    std.ml.load_model
      @ (path: string) -> result[model_handle, string]
      + loads a neural network model from disk
      - returns error when the file cannot be read or parsed
      # ml
    std.ml.run
      @ (model: model_handle, input: list[f32], input_shape: list[i32]) -> result[list[f32], string]
      + runs the model on an input tensor and returns the output tensor
      - returns error when the input shape does not match the model
      # ml

ocr
  ocr.load
    @ (detector_path: string, recognizer_path: string, language: string) -> result[ocr_engine, string]
    + loads detector and recognizer models and the character set for the language
    - returns error when either model file cannot be loaded
    - returns error when the language is not supported
    # loading
    -> std.ml.load_model
  ocr.preprocess_image
    @ (img: image) -> image
    + converts to grayscale and normalizes for detection input
    # preprocessing
    -> std.image.to_grayscale
    -> std.image.resize
  ocr.detect_regions
    @ (engine: ocr_engine, img: image) -> result[list[text_region], string]
    + returns the bounding boxes of candidate text regions
    - returns error when the detector fails to run
    # detection
    -> std.ml.run
  ocr.recognize_region
    @ (engine: ocr_engine, img: image, region: text_region) -> result[string, string]
    + crops to the region and decodes characters using the recognizer
    - returns error when the crop is out of bounds
    - returns error when the recognizer fails to run
    # recognition
    -> std.image.crop
    -> std.ml.run
  ocr.ctc_decode
    @ (logits: list[f32], alphabet: list[string], seq_len: i32) -> string
    + returns the greedy CTC decode of a logits sequence
    ? collapses repeated characters and removes blank tokens
    # decoding
  ocr.recognize
    @ (engine: ocr_engine, image_path: string) -> result[list[recognized_text], string]
    + returns every recognized text item with its bounding box and confidence
    - returns error when the image cannot be loaded
    - returns error when detection or recognition fails
    # pipeline
    -> std.fs.read_all
    -> std.image.decode
