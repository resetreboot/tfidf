
package main

import (
   "tfidf"
)

func main() {
    text1 := "This is a text for testing the TFIDF, but we expect the algorithm to crash."
    text2 := "Is that so? We create the world that is put in front of your eyes."

    input := []string{text1, text2}

    tfidf.TfIdf(input)
}
