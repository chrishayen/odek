# Requirement: "a sitemap generator"

Builds a sitemap document from URL entries and serializes it to the standard xml format.

std
  std.xml
    std.xml.escape_text
      @ (s: string) -> string
      + replaces `&`, `<`, `>`, `"`, `'` with their xml entities
      # serialization

sitemap
  sitemap.new
    @ () -> sitemap_doc
    + creates an empty sitemap document
    # construction
  sitemap.add_url
    @ (doc: sitemap_doc, loc: string, lastmod: optional[string], changefreq: optional[string], priority: optional[f32]) -> result[sitemap_doc, string]
    + appends the url entry
    - returns error when loc is empty or priority is outside 0.0..1.0
    # entries
  sitemap.to_xml
    @ (doc: sitemap_doc) -> string
    + serializes the document as a urlset xml string
    + omits optional fields that were not provided
    # serialization
    -> std.xml.escape_text
