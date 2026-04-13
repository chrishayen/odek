# Requirement: "a native SSHv2 protocol library"

A full SSHv2 client: key exchange, user authentication, and channel-based session and file transfer operations. Cryptographic primitives live in std.

std
  std.tcp
    std.tcp.dial
      @ (addr: string) -> result[tcp_conn, string]
      + opens a TCP connection to addr
      - returns error when the host is unreachable
      # network
    std.tcp.read
      @ (conn: tcp_conn, n: i32) -> result[bytes, string]
      + reads up to n bytes from a connection
      # network
    std.tcp.write
      @ (conn: tcp_conn, data: bytes) -> result[void, string]
      + writes data to a connection
      # network
  std.crypto
    std.crypto.dh_generate
      @ (group: string) -> tuple[bytes, bytes]
      + returns (private, public) for a Diffie-Hellman group
      # cryptography
    std.crypto.dh_shared
      @ (private: bytes, peer_public: bytes) -> bytes
      + computes a Diffie-Hellman shared secret
      # cryptography
    std.crypto.aes_ctr_encrypt
      @ (key: bytes, iv: bytes, plaintext: bytes) -> bytes
      + encrypts plaintext with AES in CTR mode
      # cryptography
    std.crypto.aes_ctr_decrypt
      @ (key: bytes, iv: bytes, ciphertext: bytes) -> bytes
      + decrypts ciphertext with AES in CTR mode
      # cryptography
    std.crypto.hmac_sha256
      @ (key: bytes, data: bytes) -> bytes
      + computes HMAC-SHA256
      # cryptography
    std.crypto.rsa_sign
      @ (private_key: bytes, data: bytes) -> bytes
      + signs data with an RSA private key
      # cryptography
    std.crypto.rsa_verify
      @ (public_key: bytes, data: bytes, signature: bytes) -> bool
      + returns true when the signature is valid
      # cryptography

ssh2
  ssh2.connect
    @ (host: string, port: i32) -> result[ssh_session, string]
    + opens a TCP connection and exchanges version strings
    - returns error when the peer version is incompatible
    # handshake
    -> std.tcp.dial
    -> std.tcp.read
    -> std.tcp.write
  ssh2.key_exchange
    @ (session: ssh_session) -> result[ssh_session, string]
    + performs Diffie-Hellman group exchange and derives session keys
    - returns error when the peer host key cannot be verified
    # handshake
    -> std.crypto.dh_generate
    -> std.crypto.dh_shared
    -> std.crypto.rsa_verify
  ssh2.auth_password
    @ (session: ssh_session, user: string, password: string) -> result[ssh_session, string]
    + authenticates with a username and password
    - returns error when the server rejects credentials
    # authentication
  ssh2.auth_publickey
    @ (session: ssh_session, user: string, private_key: bytes) -> result[ssh_session, string]
    + authenticates with a public key
    - returns error when the signature is not accepted
    # authentication
    -> std.crypto.rsa_sign
  ssh2.open_channel
    @ (session: ssh_session, kind: string) -> result[ssh_channel, string]
    + opens a named channel over the session
    - returns error when the channel kind is not supported
    # channels
  ssh2.exec
    @ (channel: ssh_channel, command: string) -> result[bytes, string]
    + executes a command on a session channel and returns stdout
    - returns error when the channel is not open
    # exec
    -> std.crypto.aes_ctr_encrypt
    -> std.crypto.aes_ctr_decrypt
    -> std.crypto.hmac_sha256
  ssh2.sftp_open
    @ (session: ssh_session) -> result[sftp_client, string]
    + starts an SFTP subsystem on the session
    # sftp
  ssh2.sftp_get
    @ (client: sftp_client, remote_path: string) -> result[bytes, string]
    + fetches a remote file as bytes
    - returns error when the file does not exist
    # sftp
  ssh2.sftp_put
    @ (client: sftp_client, remote_path: string, data: bytes) -> result[void, string]
    + writes bytes to a remote path
    - returns error when the server denies the write
    # sftp
  ssh2.close
    @ (session: ssh_session) -> void
    + sends disconnect and closes the underlying connection
    # teardown
