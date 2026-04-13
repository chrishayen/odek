# Requirement: "a library for creating, previewing, and sending custom email templates"

Templates are stored as named sources with a subject line and a body; rendering substitutes variables; sending delegates to a pluggable transport.

std
  std.template
    std.template.render
      @ (source: string, vars: map[string, string]) -> result[string, string]
      + substitutes {{name}} placeholders with values from vars
      - returns error on unclosed placeholders
      - returns error when a referenced variable is missing
      # templating
  std.smtp
    std.smtp.send
      @ (host: string, port: i32, from: string, to: list[string], subject: string, body: string) -> result[void, string]
      + delivers a single message over SMTP
      - returns error on connection or protocol failure
      # mail_transport

email_templates
  email_templates.new_store
    @ () -> template_store
    + creates an empty template store
    # construction
  email_templates.register
    @ (store: template_store, name: string, subject: string, body: string) -> result[void, string]
    + stores a template keyed by name
    - returns error when name is empty
    - returns error when a template with that name already exists
    # registration
  email_templates.render
    @ (store: template_store, name: string, vars: map[string, string]) -> result[rendered_email, string]
    + returns the rendered subject and body
    - returns error when the template name is unknown
    # rendering
    -> std.template.render
  email_templates.preview
    @ (store: template_store, name: string, vars: map[string, string]) -> result[string, string]
    + returns a human-readable preview combining subject and body
    - returns error when rendering fails
    # preview
  email_templates.send
    @ (store: template_store, name: string, to: list[string], vars: map[string, string], transport: smtp_config) -> result[void, string]
    + renders the template and dispatches it through the transport
    - returns error when rendering fails
    - returns error when transport delivery fails
    # delivery
    -> std.smtp.send
