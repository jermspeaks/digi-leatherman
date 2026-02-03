package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

const (
	maxWords       = 1000
	maxSentences   = 100
	maxParagraphs  = 50
	maxCharacters  = 10000
	maxBytes       = 10000
	maxTitleWords  = 20
	maxSlugWords   = 20
	maxCamelWords  = 20
	maxListItems   = 50
	maxHeadings    = 20
	maxHtmlBlocks  = 20
	maxMarkdownBlocks = 20
	maxJsonKeys    = 10
	minCount       = 1
	classicLatin   = "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
)

// LoremOptions holds tool-specific options.
type LoremOptions struct {
	StartWithClassic bool     `json:"startWithClassic"`
	Type             string   `json:"type"`   // words, sentences, paragraphs (generator)
	Vocabulary       string   `json:"vocabulary"` // default, bacon, hipster
	WholeWordsOnly   bool     `json:"wholeWordsOnly"`
	ListStyle        string   `json:"listStyle"`   // bullet, numbered
	HeadingStyle     string   `json:"headingStyle"` // plain, markdown
	Format           string   `json:"format"`   // paragraphs, list, headings (html/markdown)
	Keys             []string `json:"keys"`    // json keys
}

// LoremRequest is the JSON body for the lorem ipsum generator.
type LoremRequest struct {
	Tool    string       `json:"tool"`    // generator, characters, bytes, title, slug, camelCase, list, headings, html, markdown, json
	Type    string       `json:"type"`    // legacy: words, sentences, paragraphs
	Count   int          `json:"count"`
	Options *LoremOptions `json:"options"`
}

// loremWords is a fixed list of lorem-style words for generation.
var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "aliquip", "ex", "ea", "commodo",
	"consequat", "duis", "aute", "irure", "in", "reprehenderit", "voluptate",
	"velit", "esse", "cillum", "fugiat", "nulla", "pariatur", "excepteur", "sint",
	"occaecat", "cupidatat", "non", "proident", "sunt", "culpa", "qui", "officia",
	"deserunt", "mollit", "anim", "id", "est", "laborum",
}

var baconWords = []string{
	"bacon", "ipsum", "dolor", "amet", "short", "ribs", "brisket", "pork",
	"belly", "sausage", "beef", "ribs", "meatball", "turducken", "ham", "hock",
	"pastrami", "fatback", "tongue", "tenderloin", "pancetta", "chicken", "jerky",
	"rump", "shank", "tail", "tri-tip", "kevin", "drumstick", "venison",
}

var hipsterWords = []string{
	"authentic", "artisan", "craft", "vinyl", "retro", "aesthetic", "pinterest",
	"cold-pressed", "succulents", "pour-over", "fixie", "hashtag", "tote", "bag",
	"biodiesel", "meggings", "waistcoat", "chartreuse", "locavore", "sriracha",
	"gentrify", "schlitz", "wolf", "moon", "sustainable", "vegan", "crucifix",
}

func wordsForVocabulary(v string) []string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "bacon":
		return baconWords
	case "hipster":
		return hipsterWords
	default:
		return loremWords
	}
}

// LoremIpsum handles POST /api/lorem-ipsum/generate. Dispatches by tool and returns { result: string }.
func LoremIpsum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req LoremRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	// Backward compatibility: type + count without tool -> generator
	tool := strings.TrimSpace(strings.ToLower(req.Tool))
	if tool == "" {
		req.Type = strings.TrimSpace(strings.ToLower(req.Type))
		if req.Type == "words" || req.Type == "sentences" || req.Type == "paragraphs" {
			tool = "generator"
		} else {
			http.Error(w, "invalid type: must be words, sentences, or paragraphs", http.StatusBadRequest)
			return
		}
	}
	if req.Options == nil {
		req.Options = &LoremOptions{}
	}
	opts := req.Options
	if tool == "generator" && opts.Type == "" {
		opts.Type = req.Type
		if opts.Type == "" {
			opts.Type = "paragraphs"
		}
	}
	opts.Type = strings.TrimSpace(strings.ToLower(opts.Type))
	opts.ListStyle = strings.TrimSpace(strings.ToLower(opts.ListStyle))
	opts.HeadingStyle = strings.TrimSpace(strings.ToLower(opts.HeadingStyle))
	opts.Format = strings.TrimSpace(strings.ToLower(opts.Format))
	opts.Vocabulary = strings.TrimSpace(strings.ToLower(opts.Vocabulary))

	var result string
	var errMsg string
	switch tool {
	case "generator":
		result, errMsg = runGenerator(req.Count, opts)
	case "characters":
		result, errMsg = runCharacters(req.Count, opts)
	case "bytes":
		result, errMsg = runBytes(req.Count, opts)
	case "title":
		result, errMsg = runTitle(req.Count, opts)
	case "slug":
		result, errMsg = runSlug(req.Count, opts)
	case "camelcase":
		result, errMsg = runCamelCase(req.Count, opts)
	case "list":
		result, errMsg = runList(req.Count, opts)
	case "headings":
		result, errMsg = runHeadings(req.Count, opts)
	case "html":
		result, errMsg = runHTML(req.Count, opts)
	case "markdown":
		result, errMsg = runMarkdown(req.Count, opts)
	case "json":
		result, errMsg = runJSON(req.Count, opts)
	default:
		http.Error(w, "invalid tool", http.StatusBadRequest)
		return
	}
	if errMsg != "" {
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	writeJSON(w, StringResponse{Result: result})
}

func runGenerator(count int, opts *LoremOptions) (string, string) {
	typ := opts.Type
	if typ == "" {
		typ = "paragraphs"
	}
	var max int
	switch typ {
	case "words":
		max = maxWords
	case "sentences":
		max = maxSentences
	case "paragraphs":
		max = maxParagraphs
	default:
		return "", "invalid type: must be words, sentences, or paragraphs"
	}
	if count < minCount || count > max {
		return "", "count out of range"
	}
	voc := wordsForVocabulary(opts.Vocabulary)
	out := generateLoremWithVocab(typ, count, voc, opts.StartWithClassic)
	return out, ""
}

func generateLoremWithVocab(typ string, count int, words []string, startWithClassic bool) string {
	nw := len(words)
	var prefix string
	if startWithClassic {
		prefix = classicLatin + " "
	}
	switch typ {
	case "words":
		var b strings.Builder
		b.WriteString(prefix)
		for i := 0; i < count; i++ {
			if i > 0 || prefix != "" {
				if i > 0 {
					b.WriteByte(' ')
				}
			}
			b.WriteString(words[rand.Intn(nw)])
		}
		return strings.TrimSpace(b.String())
	case "sentences":
		sentences := make([]string, count)
		for i := 0; i < count; i++ {
			if startWithClassic && i == 0 {
				sentences[i] = classicLatin
			} else {
				sentences[i] = makeSentenceWithWords(nw, words)
			}
		}
		return strings.Join(sentences, " ")
	case "paragraphs":
		paragraphs := make([]string, count)
		for i := 0; i < count; i++ {
			if startWithClassic && i == 0 {
				paragraphs[i] = classicLatin
			} else {
				paragraphs[i] = makeParagraphWithWords(nw, words)
			}
		}
		return strings.Join(paragraphs, "\n\n")
	default:
		return ""
	}
}

func runCharacters(count int, opts *LoremOptions) (string, string) {
	if count < minCount || count > maxCharacters {
		return "", "count out of range"
	}
	words := wordsForVocabulary(opts.Vocabulary)
	nw := len(words)
	var b strings.Builder
	for b.Len() < count {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(words[rand.Intn(nw)])
	}
	s := b.String()
	if opts.WholeWordsOnly {
		if len(s) > count {
			// Trim to last complete word before count
			idx := strings.LastIndex(s[:count], " ")
			if idx > 0 {
				s = s[:idx]
			} else {
				s = s[:count]
			}
		}
		return s, ""
	}
	if len(s) > count {
		s = s[:count]
	}
	return s, ""
}

func runBytes(count int, opts *LoremOptions) (string, string) {
	if count < minCount || count > maxBytes {
		return "", "count out of range"
	}
	words := wordsForVocabulary(opts.Vocabulary)
	nw := len(words)
	var b strings.Builder
	for len(b.String()) < count {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(words[rand.Intn(nw)])
	}
	s := b.String()
	bs := []byte(s)
	if opts.WholeWordsOnly {
		if len(bs) > count {
			idx := count
			for idx > 0 && bs[idx-1] != ' ' {
				idx--
			}
			if idx > 0 {
				bs = bs[:idx]
			} else {
				bs = bs[:count]
			}
		}
		return string(bs), ""
	}
	if len(bs) > count {
		return string(bs[:count]), ""
	}
	return string(bs), ""
}

func runTitle(count int, opts *LoremOptions) (string, string) {
	if count < minCount || count > maxTitleWords {
		return "", "count out of range"
	}
	words := wordsForVocabulary(opts.Vocabulary)
	nw := len(words)
	var b strings.Builder
	for i := 0; i < count; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		w := words[rand.Intn(nw)]
		b.WriteString(capitalize(w))
	}
	return b.String(), ""
}

func runSlug(count int, opts *LoremOptions) (string, string) {
	if count < minCount || count > maxSlugWords {
		return "", "count out of range"
	}
	words := wordsForVocabulary(opts.Vocabulary)
	nw := len(words)
	parts := make([]string, count)
	for i := 0; i < count; i++ {
		parts[i] = strings.ToLower(words[rand.Intn(nw)])
	}
	return strings.Join(parts, "-"), ""
}

func runCamelCase(count int, opts *LoremOptions) (string, string) {
	if count < minCount || count > maxCamelWords {
		return "", "count out of range"
	}
	words := wordsForVocabulary(opts.Vocabulary)
	nw := len(words)
	var b strings.Builder
	for i := 0; i < count; i++ {
		w := words[rand.Intn(nw)]
		if i == 0 {
			b.WriteString(strings.ToLower(w))
		} else {
			b.WriteString(capitalize(w))
		}
	}
	return b.String(), ""
}

func runList(count int, opts *LoremOptions) (string, string) {
	if count < minCount || count > maxListItems {
		return "", "count out of range"
	}
	words := wordsForVocabulary(opts.Vocabulary)
	nw := len(words)
	style := opts.ListStyle
	if style != "numbered" {
		style = "bullet"
	}
	var b strings.Builder
	for i := 0; i < count; i++ {
		phrase := makeSentenceWithWords(nw, words)
		if style == "numbered" {
			b.WriteString(strconv.Itoa(i + 1))
			b.WriteString(". ")
		} else {
			b.WriteString("- ")
		}
		b.WriteString(phrase)
		b.WriteByte('\n')
	}
	return strings.TrimSuffix(b.String(), "\n"), ""
}

func runHeadings(count int, opts *LoremOptions) (string, string) {
	if count < minCount || count > maxHeadings {
		return "", "count out of range"
	}
	words := wordsForVocabulary(opts.Vocabulary)
	nw := len(words)
	useMarkdown := opts.HeadingStyle == "markdown"
	var b strings.Builder
	for i := 0; i < count; i++ {
		level := (i % 3) + 1
		phrase := makeShortPhrase(nw, words, 3, 7)
		if useMarkdown {
			for j := 0; j < level; j++ {
				b.WriteString("# ")
			}
		}
		b.WriteString(phrase)
		b.WriteByte('\n')
	}
	return strings.TrimSuffix(b.String(), "\n"), ""
}

func runHTML(count int, opts *LoremOptions) (string, string) {
	if count < minCount || count > maxHtmlBlocks {
		return "", "count out of range"
	}
	words := wordsForVocabulary(opts.Vocabulary)
	nw := len(words)
	format := opts.Format
	if format != "list" {
		format = "paragraphs"
	}
	var b strings.Builder
	if format == "paragraphs" {
		for i := 0; i < count; i++ {
			b.WriteString("<p>")
			b.WriteString(makeParagraphWithWords(nw, words))
			b.WriteString("</p>\n")
		}
	} else {
		b.WriteString("<ul>\n")
		for i := 0; i < count; i++ {
			b.WriteString("  <li>")
			b.WriteString(makeSentenceWithWords(nw, words))
			b.WriteString("</li>\n")
		}
		b.WriteString("</ul>")
	}
	return strings.TrimSuffix(b.String(), "\n"), ""
}

func runMarkdown(count int, opts *LoremOptions) (string, string) {
	if count < minCount || count > maxMarkdownBlocks {
		return "", "count out of range"
	}
	words := wordsForVocabulary(opts.Vocabulary)
	nw := len(words)
	format := opts.Format
	if format != "list" && format != "headings" {
		format = "paragraphs"
	}
	var b strings.Builder
	switch format {
	case "paragraphs":
		for i := 0; i < count; i++ {
			b.WriteString(makeParagraphWithWords(nw, words))
			b.WriteString("\n\n")
		}
	case "list":
		for i := 0; i < count; i++ {
			b.WriteString("- ")
			b.WriteString(makeSentenceWithWords(nw, words))
			b.WriteString("\n")
		}
	case "headings":
		for i := 0; i < count; i++ {
			level := (i % 3) + 1
			for j := 0; j < level; j++ {
				b.WriteString("# ")
			}
			b.WriteString(makeShortPhrase(nw, words, 3, 7))
			b.WriteString("\n")
		}
	}
	return strings.TrimSpace(b.String()), ""
}

func runJSON(count int, opts *LoremOptions) (string, string) {
	keys := opts.Keys
	if len(keys) == 0 {
		keys = []string{"title", "body", "summary"}
	}
	if len(keys) > maxJsonKeys {
		return "", "too many keys"
	}
	words := wordsForVocabulary(opts.Vocabulary)
	nw := len(words)
	obj := make(map[string]string)
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		obj[k] = makeSentenceWithWords(nw, words)
	}
	bs, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", "json marshal failed"
	}
	return string(bs), ""
}

func makeShortPhrase(nw int, words []string, minW, maxW int) string {
	numWords := minW + rand.Intn(maxW-minW+1)
	var b strings.Builder
	for i := 0; i < numWords; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		w := words[rand.Intn(nw)]
		if i == 0 {
			w = capitalize(w)
		}
		b.WriteString(w)
	}
	return b.String()
}

func makeSentenceWithWords(nw int, words []string) string {
	numWords := 5 + rand.Intn(11)
	var b strings.Builder
	for i := 0; i < numWords; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		w := words[rand.Intn(nw)]
		if i == 0 {
			w = capitalize(w)
		}
		b.WriteString(w)
	}
	b.WriteByte('.')
	return b.String()
}

func makeParagraphWithWords(nw int, words []string) string {
	numSentences := 3 + rand.Intn(5)
	sentences := make([]string, numSentences)
	for i := 0; i < numSentences; i++ {
		sentences[i] = makeSentenceWithWords(nw, words)
	}
	return strings.Join(sentences, " ")
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := rune(s[0])
	return string(unicode.ToUpper(r)) + s[1:]
}
