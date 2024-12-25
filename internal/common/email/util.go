package email

import (
    "bytes"
    "context"
    "fmt"
    "html/template"
    "path/filepath"
    "regexp"
    "strings"
    "golang.org/x/net/html" 
    "github.com/Boostport/mjml-go"
)

func RenderTemplate(templatePath string, data map[string]interface{}) (string, string, error) {
    tmpl, err := template.ParseFiles(templatePath)
    if err != nil {
        return "", "", fmt.Errorf("failed to parse template: %w", err)
    }

    var mjmlBuf bytes.Buffer
    if err := tmpl.Execute(&mjmlBuf, data); err != nil {
        return "", "", fmt.Errorf("failed to execute template: %w", err)
    }

    // Create a context for MJML conversion
    ctx := context.Background()
    
    html, err := mjml.ToHTML(ctx, mjmlBuf.String(), mjml.WithMinify(true))
    if err != nil {
        return "", "", fmt.Errorf("failed to convert MJML to HTML: %w", err)
    }

    plainText := StripHTML(html)
    return html, plainText, nil
}


func StripHTML(htmlContent string) string {
    // Pre-process: replace common HTML entities
    content := strings.ReplaceAll(htmlContent, "&nbsp;", " ")
    content = strings.ReplaceAll(content, "&amp;", "&")
    content = strings.ReplaceAll(content, "&lt;", "<")
    content = strings.ReplaceAll(content, "&gt;", ">")
    content = strings.ReplaceAll(content, "&quot;", "\"")
    content = strings.ReplaceAll(content, "&#39;", "'")

    doc, err := html.Parse(strings.NewReader(content))
    if err != nil {
        // Fallback to simple regex-based stripping if parsing fails
        return stripHTMLRegex(content)
    }

    var buf bytes.Buffer
    extractText(doc, &buf)

    // Post-process the text
    text := buf.String()
    
    // Remove extra whitespace
    text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
    
    // Handle line breaks properly
    text = strings.ReplaceAll(text, "<br>", "\n")
    text = strings.ReplaceAll(text, "<br/>", "\n")
    text = strings.ReplaceAll(text, "<br />", "\n")
    
    // Remove multiple newlines
    text = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(text, "\n\n")
    
    // Trim spaces at start and end
    text = strings.TrimSpace(text)

    return text
}

func extractText(n *html.Node, buf *bytes.Buffer) {
    if n.Type == html.TextNode {
        // Write the text content
        buf.WriteString(strings.TrimSpace(n.Data))
        if len(n.Data) > 0 {
            buf.WriteString(" ")
        }
        return
    }
    
    // Handle special tags
    if n.Type == html.ElementNode {
        switch n.Data {
        case "br", "p", "div", "tr":
            buf.WriteString("\n")
        case "h1", "h2", "h3", "h4", "h5", "h6":
            buf.WriteString("\n\n")
        case "li":
            buf.WriteString("\n- ")
        case "a":
            // For links, we might want to include the href
            var href string
            for _, attr := range n.Attr {
                if attr.Key == "href" {
                    href = attr.Val
                    break
                }
            }
            if href != "" {
                linkText := extractLinkText(n)
                if linkText != href {
                    buf.WriteString(linkText + " (" + href + ") ")
                    return
                }
            }
        }
    }

    // Recursively process child nodes
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        extractText(c, buf)
    }

    // Add appropriate spacing after block elements
    if n.Type == html.ElementNode {
        switch n.Data {
        case "p", "div", "h1", "h2", "h3", "h4", "h5", "h6":
            buf.WriteString("\n")
        }
    }
}

func extractLinkText(n *html.Node) string {
    var buf bytes.Buffer
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if c.Type == html.TextNode {
            buf.WriteString(c.Data)
        }
    }
    return strings.TrimSpace(buf.String())
}

// Fallback function using regex for simple HTML stripping
func stripHTMLRegex(html string) string {
    // First, replace line break tags with newlines
    html = regexp.MustCompile(`<br[^>]*>`).ReplaceAllString(html, "\n")
    
    // Replace other block-level elements with newlines
    html = regexp.MustCompile(`</?(div|p|h\d|table|tr|ul|ol|li)[^>]*>`).ReplaceAllString(html, "\n")
    
    // Remove all other HTML tags
    html = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(html, "")
    
    // Replace HTML entities
    html = strings.ReplaceAll(html, "&nbsp;", " ")
    html = strings.ReplaceAll(html, "&amp;", "&")
    html = strings.ReplaceAll(html, "&lt;", "<")
    html = strings.ReplaceAll(html, "&gt;", ">")
    html = strings.ReplaceAll(html, "&quot;", "\"")
    html = strings.ReplaceAll(html, "&#39;", "'")
    
    // Clean up whitespace
    html = regexp.MustCompile(`\s+`).ReplaceAllString(html, " ")
    html = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(html, "\n\n")
    
    return strings.TrimSpace(html)
}

func GetMimeType(filename string) string {
    ext := filepath.Ext(filename)
    switch ext {
    case ".pdf":
        return "application/pdf"
    case ".jpg", ".jpeg":
        return "image/jpeg"
    case ".png":
        return "image/png"
    case ".gif":
        return "image/gif"
    default:
        return "application/octet-stream"
    }
}