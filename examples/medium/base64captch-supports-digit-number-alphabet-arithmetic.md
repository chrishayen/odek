# Requirement: "a CAPTCHA library supporting digit, alphabetic, arithmetic, and audio variants, rendered as base64 images"

Generates a challenge with a known answer, renders it to an image or audio buffer, and exposes a verification step against a stored answer.

std
  std.rand
    std.rand.int_range
      @ (lo: i32, hi: i32) -> i32
      + returns a random integer in [lo, hi]
      # randomness
    std.rand.pick
      @ (choices: list[string]) -> string
      + returns a random element from choices
      # randomness
  std.image
    std.image.new_rgba
      @ (width: i32, height: i32) -> image_buffer
      + returns a blank RGBA image of the given size
      # imaging
    std.image.draw_text
      @ (img: image_buffer, text: string, x: i32, y: i32, size: i32) -> image_buffer
      + draws text onto the image at the specified position
      # imaging
    std.image.encode_png
      @ (img: image_buffer) -> bytes
      + returns the PNG encoding of the image
      # imaging
  std.encoding
    std.encoding.base64_encode
      @ (data: bytes) -> string
      + returns the standard base64 encoding of the input
      # encoding

captcha
  captcha.generate_digits
    @ (length: i32) -> challenge
    + returns a challenge whose answer is a string of length random digits
    # generation
    -> std.rand.int_range
  captcha.generate_alphabet
    @ (length: i32) -> challenge
    + returns a challenge whose answer is an alphabetic string
    # generation
    -> std.rand.pick
  captcha.generate_arithmetic
    @ () -> challenge
    + returns a challenge whose answer is the result of a simple arithmetic expression
    + the prompt shows the expression, the answer is the computed value
    # generation
    -> std.rand.int_range
  captcha.render_image
    @ (ch: challenge) -> string
    + returns a data URI containing a PNG rendering of the challenge prompt
    # rendering
    -> std.image.new_rgba
    -> std.image.draw_text
    -> std.image.encode_png
    -> std.encoding.base64_encode
  captcha.render_audio
    @ (ch: challenge) -> string
    + returns a data URI containing a WAV rendering that speaks each character of the prompt
    # rendering
    -> std.encoding.base64_encode
  captcha.store
    @ (store: captcha_store, id: string, ch: challenge, ttl_ms: i32) -> captcha_store
    + persists the expected answer under id with an expiration
    # storage
  captcha.verify
    @ (store: captcha_store, id: string, answer: string) -> tuple[bool, captcha_store]
    + returns (true, store_without_entry) when answer matches and is not expired
    - returns (false, store_without_entry) when the answer is wrong
    - returns (false, store_unchanged) when the id is unknown
    # verification
