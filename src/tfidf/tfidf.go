package tfidf

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"sort"
	"strings"
)

type TfIdfMatrix struct {
	Stopwords []string
	TotalFreq map[string]float64
	DocFreqs  []map[string]float64
}

// totalFreq := make(map[string]float64)
// docFreqs := make([]map[string]float64, 0)

func (matrix *TfIdfMatrix) LoadStopWords(filename string) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	elements := strings.Fields(string(contents))

	fmt.Printf("Loaded %d stopwords.\n", len(elements))
	matrix.Stopwords = append(matrix.Stopwords, elements...)
}

func (matrix *TfIdfMatrix) checkStopWords(word string) bool {
	result := false

	for _, stopword := range matrix.Stopwords {
		if word == stopword {
			result = true
			break
		}
	}

	return result
}

func ClearHtml(text string) string {
	htmlCleaner, _ := regexp.Compile("<.*?>")
	result := string(htmlCleaner.ReplaceAllString(text, ""))

	return result
}

func (matrix *TfIdfMatrix) TfCalc(text string) map[string]float64 {
	tfMatrix := make(map[string]float64)

	punctuation := regexp.MustCompile("(,|;|:|!||¡|¿)")
	noStrangeChars := regexp.MustCompile("^[A-Za-z]*?$")

	temp := strings.Replace(text, ".", " ", -1)
	temp = strings.Replace(temp, "?", " ", -1)

	temp = string(punctuation.ReplaceAllString(temp, ""))

	temp = strings.ToLower(temp)

	for _, value := range strings.Fields(temp) {
		if noStrangeChars.MatchString(value) && !matrix.checkStopWords(value) {
			if amount, ok := tfMatrix[value]; ok {
				tfMatrix[value] = amount + 1.0
			} else {
				tfMatrix[value] = 1.0
			}
		}
	}

	return tfMatrix
}

func (matrix *TfIdfMatrix) DfUpdate(docFreqs map[string]float64) {
	if matrix.TotalFreq == nil {
		matrix.TotalFreq = make(map[string]float64)
	}

	for word, quantity := range docFreqs {
		if amount, ok := matrix.TotalFreq[word]; ok {
			matrix.TotalFreq[word] = amount + quantity
		} else {
			matrix.TotalFreq[word] = quantity
		}
	}
}

func (matrix *TfIdfMatrix) CalculateIdf(totalDocNo int) {
	for word, amount := range matrix.TotalFreq {
		matrix.TotalFreq[word] = math.Log(float64(totalDocNo) / amount)
	}
}

func (matrix *TfIdfMatrix) PrintTfIdf(documentFreqs map[string]float64) {
	maxTfIdf := 0.0

	res := make(map[string]float64, 0)

	for word, amount := range documentFreqs {
		if value, ok := matrix.TotalFreq[word]; ok {
			temp := amount * value
			res[word] = temp
			if temp > maxTfIdf {
				maxTfIdf = temp
			}
		} else {
			temp := amount
			res[word] = temp
			if temp > maxTfIdf {
				maxTfIdf = temp
			}
		}
	}

	for word, amount := range res {
		normalized := (amount / maxTfIdf) * 50.0 // Normalize and extend to 50 range
		if !(normalized < 1.0) {
			fmt.Printf("%s:\n", word)
			fmt.Println(strings.Repeat("#", int(normalized)))
		}
	}
}

func (matrix *TfIdfMatrix) SearchDocumentsByWord(text string) ([]int, []float64) {
	match := make([]float64, 0)
	docs := make([]int, 0)
	result := make([]int, 0)

	for count, document := range matrix.DocFreqs {
		if value, ok := document[text]; ok {
			temp := matrix.TotalFreq[text]
			fr := temp * value
			if fr > 0 {
				match = append(match, temp*value)
				docs = append(docs, count)
			}
		}
	}

	temp := make([]float64, len(match))
	_ = copy(temp, match)

	sort.Sort(sort.Reverse(sort.Float64Slice(temp)))

	for cnt, freq := range temp {
		if cnt > 5 {
			break
		}
		count := 0
		value := match[count]
		for value != freq {
			count++
			value = match[count]
		}
		match[count] = -1.0 // Mark this result as used
		result = append(result, count)
	}

	if len(result) > 5 {
        return result[:5], temp[:5]
	} else {
		return result, temp
	}
}

func (matrix *TfIdfMatrix) TfIdf(texts []string) {
	totalDocuments := len(texts)

	for _, text := range texts {
		tf := matrix.TfCalc(text)
		matrix.DocFreqs = append(matrix.DocFreqs, tf)
		matrix.DfUpdate(tf)
	}

	matrix.CalculateIdf(totalDocuments)
}
