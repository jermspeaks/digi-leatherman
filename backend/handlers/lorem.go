package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"unicode"
)

const (
	maxWords      = 1000
	maxSentences  = 100
	maxParagraphs = 50
	minCount      = 1
)

// LoremRequest is the JSON body for the lorem ipsum generator.
type LoremRequest struct {
	Type  string `json:"type"` // "words", "sentences", "paragraphs"
	Count int    `json:"count"`
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

// LoremIpsum handles POST /api/lorem-ipsum/generate. Generates placeholder text
// by words, sentences, or paragraphs. Max: words 1000, sentences 100, paragraphs 50.
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
	req.Type = strings.TrimSpace(strings.ToLower(req.Type))
	var max int
	switch req.Type {
	case "words":
		max = maxWords
	case "sentences":
		max = maxSentences
	case "paragraphs":
		max = maxParagraphs
	default:
		http.Error(w, "invalid type: must be words, sentences, or paragraphs", http.StatusBadRequest)
		return
	}
	if req.Count < minCount || req.Count > max {
		http.Error(w, "count out of range", http.StatusBadRequest)
		return
	}
	result := generateLorem(req.Type, req.Count)
	writeJSON(w, StringResponse{Result: result})
}

func generateLorem(typ string, count int) string {
	nw := len(loremWords)
	switch typ {
	case "words":
		var b strings.Builder
		for i := 0; i < count; i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(loremWords[rand.Intn(nw)])
		}
		return b.String()
	case "sentences":
		sentences := make([]string, count)
		for i := 0; i < count; i++ {
			sentences[i] = makeSentence(nw)
		}
		return strings.Join(sentences, " ")
	case "paragraphs":
		paragraphs := make([]string, count)
		for i := 0; i < count; i++ {
			paragraphs[i] = makeParagraph(nw)
		}
		return strings.Join(paragraphs, "\n\n")
	default:
		return ""
	}
}

// makeSentence returns 5–15 random words, capitalized, with a period.
func makeSentence(nw int) string {
	numWords := 5 + rand.Intn(11) // 5 to 15
	var b strings.Builder
	for i := 0; i < numWords; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		w := loremWords[rand.Intn(nw)]
		if i == 0 {
			w = capitalize(w)
		}
		b.WriteString(w)
	}
	b.WriteByte('.')
	return b.String()
}

// makeParagraph returns 3–7 sentences separated by spaces.
func makeParagraph(nw int) string {
	numSentences := 3 + rand.Intn(5) // 3 to 7
	sentences := make([]string, numSentences)
	for i := 0; i < numSentences; i++ {
		sentences[i] = makeSentence(nw)
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
