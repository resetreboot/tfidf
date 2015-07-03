package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
    "os"
	"tfidf"
)

func getLinks(text, mainUrl string) []string {
	result := make([]string, 0)
	hrefLocator, _ := regexp.Compile("href=(\"|').*?(\"|')")

	elems := hrefLocator.FindAllString(text, -1)

	for _, link := range elems {
		temp := strings.Replace(link, "href=", "", -1)
		temp = strings.Replace(temp, "\"", "", -1)
		temp = strings.Replace(temp, "'", "", -1)

		switch {
		case strings.Contains(temp, ".png"):
			continue
		case strings.Contains(temp, ".jpg"):
			continue
		case strings.Contains(temp, ".gif"):
			continue
		case strings.Contains(temp, ".xml"):
			continue
		case strings.Contains(temp, ".css"):
			continue
		case strings.Contains(temp, ".js"):
			continue
		case !strings.Contains(temp, mainUrl):
			continue
		default:
			result = append(result, temp)
		}

	}

	return result
}

func fetchAndGetLinks(url string, output chan []string) {
	result := make([]string, 0)

	resp, err := http.Get(url)

	if err != nil {
		output <- result
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		output <- result
		return
	}

	links := getLinks(string(body), url)
	output <- links

}

func fetchText(url string, output chan []string) {
	resp, err := http.Get(url)
	if err != nil {
		output <- nil
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		output <- nil
		return
	}

	if strings.Contains(string(body), "<html") && strings.Contains(string(body), "<body") {
		output <- []string{url, tfidf.ClearHtml(string(body))}
	} else {
		output <- nil
	}
}

func getDocuments(websites []string) ([]string, []string) {
	allDocsUris := make([]string, 0)               // Initial, empty
	allDocsUris = append(allDocsUris, websites...) // The index page IS a document
	results := make([]string, 0)
	sites := make([]string, 0)
	webChannel := make(chan []string, 10)
	textChannel := make(chan []string, 10)

	fmt.Println("Initial link parsing...")

	for _, url := range websites {
		fmt.Printf("Fetching: %s\n", url)
		go fetchAndGetLinks(url, webChannel)
	}

	for count := 0; count < len(websites); count++ {
		links := <-webChannel
		fmt.Printf("Received %d links...\n", len(links))
		allDocsUris = append(allDocsUris, links...)
	}

	fmt.Printf("Now fetching texts from %d urls\n", len(allDocsUris))

	for count, url := range allDocsUris {
		fmt.Printf("%d of %d - Fetching: %s\n", count, len(allDocsUris), url)
		go fetchText(url, textChannel)
	}

	for count := 0; count < len(allDocsUris); count++ {
		fmt.Printf("Received %d web pages.\n", count)
		text := <-textChannel
		if text != nil {
			sites = append(sites, text[0])
			results = append(results, text[1])
		}
	}

	return results, sites
}

func loadWebsites(filename string) []string {
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    elements := strings.Fields(string(content))
    fmt.Printf("Loaded %d websites", len(elements))

    return elements
}

func getCommandLineArgument(option, defaultValue string) string {
    result := defaultValue
    if len(os.Args) > 1 {
        for count := 1; count < len(os.Args); count++ {
            if strings.Contains(os.Args[count], option) {
                result = os.Args[count]
                result = strings.Replace(result, option, "", -1)
                result = strings.TrimSpace(result)

                if result == "" {
                    result = os.Args[count + 1]
                    result = strings.Replace(result, option, "", -1)
                    result = strings.TrimSpace(result)
                }
                break
            }
        }
    }

    return result
}

func getWebsitesArgument() string {
    return getCommandLineArgument("-w", "websites.txt")
}

func getStopWordsArgument() string {
    return getCommandLineArgument("-s", "english1")
}

func main() {
	var searchTerm string
	var matrix tfidf.TfIdfMatrix

	fmt.Println("Init the TfIdf with stopwords")

    stopwordsLanguage := getStopWordsArgument()
	matrix.LoadStopWords(fmt.Sprintf("./src/tfidf/stop-words/stop-words-%s.txt", stopwordsLanguage))

	websites := loadWebsites(getWebsitesArgument())

	input, sites := getDocuments(websites)

	matrix.TfIdf(input)

	running := true

	for running {
		fmt.Println("Introduzca término de búsqueda o q para salir:")
		fmt.Scanf("%s", &searchTerm)

		if searchTerm != "q" {
			results, percentages := matrix.SearchDocumentsByWord(searchTerm)
			fmt.Println("Resultados:")
			for c, value := range results {
				fmt.Printf("%d: %.2f%% - %s\n", c, percentages[c] * 10.0, sites[value])
			}
		} else {
			running = false
		}
	}
}
