package docs

import (
	"strings"

	"golang.org/x/net/html"
)

// ModuleDoc holds parsed documentation for a module.
type ModuleDoc struct {
	FQCN             string        `json:"fqcn"`
	ShortDescription string        `json:"short_description"`
	VersionAdded     string        `json:"version_added,omitempty"`
	Description      []string      `json:"description,omitempty"`
	Parameters       []Parameter   `json:"parameters,omitempty"`
	Examples         string        `json:"examples,omitempty"`
	ReturnValues     []ReturnValue `json:"return_values,omitempty"`
	Notes            []string      `json:"notes,omitempty"`
	DocsURL          string        `json:"docs_url"`
	Source           string        `json:"source"`
}

// Parameter represents a single module parameter.
type Parameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type,omitempty"`
	Required    bool        `json:"required"`
	Default     string      `json:"default,omitempty"`
	Choices     []string    `json:"choices,omitempty"`
	Description string      `json:"description"`
	Aliases     []string    `json:"aliases,omitempty"`
	Suboptions  []Parameter `json:"suboptions,omitempty"`
}

// ReturnValue represents a value returned by a module.
type ReturnValue struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Returned    string `json:"returned,omitempty"`
	Type        string `json:"type,omitempty"`
	Sample      string `json:"sample,omitempty"`
}

// KeywordEntry is a playbook keyword definition.
type KeywordEntry struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type,omitempty"`
	Default     string `json:"default,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

// SpecialVariable is an Ansible magic variable.
type SpecialVariable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Scope       string `json:"scope,omitempty"`
}

// ParseModuleDocs parses an ansible docs module HTML page into structured data.
func ParseModuleDocs(htmlContent string) (*ModuleDoc, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	result := &ModuleDoc{}
	result.ShortDescription = extractTitle(doc)
	result.Parameters = extractParameters(doc)
	result.Examples = extractExamples(doc)
	result.ReturnValues = extractReturnValues(doc)
	result.Notes = extractNotes(doc)

	return result, nil
}

// ParsePlaybookKeywords parses the playbook keywords reference page.
func ParsePlaybookKeywords(htmlContent string) ([]KeywordEntry, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	return extractKeywords(doc), nil
}

// ParseSpecialVariables parses the special variables reference page.
func ParseSpecialVariables(htmlContent string) ([]SpecialVariable, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	return extractSpecialVars(doc), nil
}

// ParseGenericContent extracts the main text content from an arbitrary docs page.
func ParseGenericContent(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}
	// Find the main content area
	main := findElement(doc, "div", "document")
	if main == nil {
		main = findElement(doc, "article", "")
	}
	if main == nil {
		main = findElement(doc, "main", "")
	}
	if main == nil {
		// Fall back to body
		main = findElement(doc, "body", "")
	}
	if main == nil {
		return "", nil
	}
	return strings.TrimSpace(extractReadableText(main)), nil
}

// extractReadableText returns human-readable text, skipping nav/header/footer.
func extractReadableText(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "nav", "header", "footer", "script", "style", "aside":
				return
			case "p", "li", "h1", "h2", "h3", "h4", "h5", "h6":
				text := strings.TrimSpace(textContent(n))
				if text != "" {
					b.WriteString(text)
					b.WriteString("\n")
				}
				return
			case "pre", "code":
				b.WriteString(textContent(n))
				b.WriteString("\n")
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return b.String()
}

// --- internal helpers ---

func textContent(n *html.Node) string {
	var b strings.Builder
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(n)
	return strings.TrimSpace(b.String())
}

func findElement(n *html.Node, tag, class string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		if class == "" || hasClass(n, class) {
			return n
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findElement(c, tag, class); found != nil {
			return found
		}
	}
	return nil
}

func findAllElements(n *html.Node, tag, class string) []*html.Node {
	var results []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			if class == "" || hasClass(n, class) {
				results = append(results, n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return results
}

func findElementByID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode && getAttr(n, "id") == id {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findElementByID(c, id); found != nil {
			return found
		}
	}
	return nil
}

func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			for _, c := range strings.Fields(attr.Val) {
				if c == class {
					return true
				}
			}
		}
	}
	return false
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func extractTitle(doc *html.Node) string {
	h1 := findElement(doc, "h1", "")
	if h1 != nil {
		t := textContent(h1)
		// Strip "– Ansible Documentation" suffix if present
		if idx := strings.Index(t, " \u2013 "); idx != -1 {
			t = t[:idx]
		}
		if idx := strings.Index(t, " — "); idx != -1 {
			t = t[:idx]
		}
		return strings.TrimSpace(t)
	}
	title := findElement(doc, "title", "")
	if title != nil {
		t := textContent(title)
		if idx := strings.Index(t, " \u2013 "); idx != -1 {
			return strings.TrimSpace(t[:idx])
		}
		if idx := strings.Index(t, " — "); idx != -1 {
			return strings.TrimSpace(t[:idx])
		}
		return t
	}
	return ""
}

func extractParameters(doc *html.Node) []Parameter {
	var params []Parameter

	// Try section#parameters first (modern Ansible docs)
	section := findElementByID(doc, "parameters")
	if section == nil {
		// Try section#id1 or just look for ansible-option list items anywhere
		section = doc
	}

	items := findAllElements(section, "li", "ansible-option")
	for _, item := range items {
		if p := parseParamItem(item); p.Name != "" {
			params = append(params, p)
		}
	}

	if len(params) == 0 {
		// Older docs: dt/dd pairs inside the parameters section
		paramSection := findElementByID(doc, "parameters")
		if paramSection != nil {
			dts := findAllElements(paramSection, "dt", "")
			for _, dt := range dts {
				name := strings.TrimSpace(textContent(dt))
				if name == "" || strings.Contains(name, " ") {
					continue
				}
				p := Parameter{Name: name}
				for sib := dt.NextSibling; sib != nil; sib = sib.NextSibling {
					if sib.Type == html.ElementNode && sib.Data == "dd" {
						p.Description = strings.TrimSpace(textContent(sib))
						break
					}
				}
				params = append(params, p)
			}
		}
	}

	return params
}

func parseParamItem(n *html.Node) Parameter {
	p := Parameter{}

	// Parameter name is usually in a strong or specific span
	nameEl := findElement(n, "strong", "")
	if nameEl != nil {
		p.Name = strings.Trim(textContent(nameEl), " \t\n/")
	}
	if p.Name == "" {
		// Fallback: look for code tag
		code := findElement(n, "code", "")
		if code != nil {
			p.Name = strings.Trim(textContent(code), " \t\n/")
		}
	}

	// Extract type, required, default from text
	allText := textContent(n)
	for _, line := range strings.Split(allText, "\n") {
		line = strings.TrimSpace(line)
		lower := strings.ToLower(line)
		switch {
		case strings.HasPrefix(lower, "type:"):
			p.Type = strings.TrimSpace(line[5:])
		case strings.HasPrefix(lower, "required:"):
			val := strings.TrimSpace(strings.ToLower(line[9:]))
			p.Required = val == "true" || val == "yes"
		case strings.HasPrefix(lower, "default:"):
			p.Default = strings.TrimSpace(line[8:])
		case strings.HasPrefix(lower, "choices:"):
			raw := strings.TrimSpace(line[8:])
			for _, choice := range strings.Split(raw, ",") {
				choice = strings.Trim(strings.TrimSpace(choice), "'\"[]")
				if choice != "" {
					p.Choices = append(p.Choices, choice)
				}
			}
		}
	}

	// Description from paragraph
	pEl := findElement(n, "p", "")
	if pEl != nil {
		p.Description = strings.TrimSpace(textContent(pEl))
	}

	return p
}

func extractExamples(doc *html.Node) string {
	section := findElementByID(doc, "examples")
	if section == nil {
		section = findElementByID(doc, "example")
	}
	if section != nil {
		pre := findElement(section, "pre", "")
		if pre != nil {
			return strings.TrimSpace(textContent(pre))
		}
	}
	return ""
}

func extractReturnValues(doc *html.Node) []ReturnValue {
	var values []ReturnValue

	section := findElementByID(doc, "return-values")
	if section == nil {
		section = findElementByID(doc, "return_values")
	}
	if section == nil {
		return values
	}

	items := findAllElements(section, "li", "ansible-option")
	for _, item := range items {
		rv := ReturnValue{}
		nameEl := findElement(item, "strong", "")
		if nameEl == nil {
			nameEl = findElement(item, "code", "")
		}
		if nameEl != nil {
			rv.Name = strings.TrimSpace(textContent(nameEl))
		}
		if pEl := findElement(item, "p", ""); pEl != nil {
			rv.Description = strings.TrimSpace(textContent(pEl))
		}
		if rv.Name != "" {
			values = append(values, rv)
		}
	}

	return values
}

func extractNotes(doc *html.Node) []string {
	var notes []string
	section := findElementByID(doc, "notes")
	if section == nil {
		return notes
	}
	for _, li := range findAllElements(section, "li", "") {
		if t := strings.TrimSpace(textContent(li)); t != "" {
			notes = append(notes, t)
		}
	}
	return notes
}

func extractKeywords(doc *html.Node) []KeywordEntry {
	var keywords []KeywordEntry
	// Ansible playbook keywords page uses dl/dt/dd structure
	dts := findAllElements(doc, "dt", "")
	for _, dt := range dts {
		name := strings.TrimSpace(textContent(dt))
		if name == "" || strings.ContainsAny(name, "\n\t") {
			continue
		}
		entry := KeywordEntry{Name: name}
		for sib := dt.NextSibling; sib != nil; sib = sib.NextSibling {
			if sib.Type == html.ElementNode && sib.Data == "dd" {
				entry.Description = strings.TrimSpace(textContent(sib))
				break
			}
		}
		keywords = append(keywords, entry)
	}
	return keywords
}

func extractSpecialVars(doc *html.Node) []SpecialVariable {
	var vars []SpecialVariable
	dts := findAllElements(doc, "dt", "")
	for _, dt := range dts {
		name := strings.TrimSpace(textContent(dt))
		if name == "" || strings.ContainsAny(name, "\n\t ") {
			continue
		}
		// Special variables typically have specific naming patterns
		v := SpecialVariable{Name: name}
		for sib := dt.NextSibling; sib != nil; sib = sib.NextSibling {
			if sib.Type == html.ElementNode && sib.Data == "dd" {
				v.Description = strings.TrimSpace(textContent(sib))
				break
			}
		}
		if v.Name != "" {
			vars = append(vars, v)
		}
	}
	return vars
}
