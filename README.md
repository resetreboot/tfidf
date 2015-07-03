# TfIdf Crawler in Go

A simple website crawler and tfidf parsing engine to implement a simple search engine.

## Command Line Parameters

This utility relies in a couple command line arguments.

* -w: Allows to choose the file where the webs to index will be found. If not present, it assumes a *websites.txt* file in the same directory the command is run in.
* -s: Chooses a stopword language. It assumes *english1* as default. See the stopword files for choices.

## Known bugs

* Stopwords should be loaded from the lib directory or implement a path search.
* Probably I messed up some of the tfidf calculations.
