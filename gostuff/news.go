package gostuff

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type NewsProviders struct {
	Status  string
	Sources []NewsSources
}

type NewsSources struct {
	ID              string
	Name            string
	Description     string
	Url             string
	Category        string
	Language        string
	Country         string
	UrlToLogos      []UrlToLogo
	SortByAvailable []string
}

type UrlToLogo struct {
	Small  string
	Medium string
	Large  string
}

type AllNewsProviders struct {
	Providers []NewsProvider
}

type NewsProvider struct {
	Status   string
	Source   string
	Sortby   string
	Articles []NewsArticle
}

type NewsArticle struct {
	Author      string
	Title       string
	Description string
	Url         string
	UrlToImage  string
	PublishedAt string
}

// fetches news from list of news sources and saves them each to their own file
// this will most likely be used as a one time function
func FetchNewsSources() {

	const (
		newsSourceList = "https://newsapi.org/v1/sources?language=en"
	)
	response, err := http.Get(newsSourceList)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var newsProviders NewsProviders

	if err := json.Unmarshal(responseData, &newsProviders); err != nil {
		fmt.Println("Just receieved a message I couldn't decode in news.go FetchNewsSources 1:", string(responseData), err)
	}

	for _, source := range newsProviders.Sources {
		url := "https://newsapi.org/v1/articles?source=" + source.ID + "&apiKey=" + getApiKey()
		//SaveNewsToFile(source.ID, url)
		unmarshalNews(url)
	}
}

// fetches and saves a single source of news
// @param url to fetch JSON of articles from a single news provider
func SaveNewsToFile(filename string, url string) {

	responseData := getHttpResponse(url)
	//fmt.Println(string(responseData))
	newsOutputPath := "data/news/" + filename + ".json"
	err := ioutil.WriteFile(newsOutputPath, responseData, 0666)
	if err != nil {
		fmt.Println(err)
	}
}

// returns the http response of the url in byte
// convert to string for human readable format
func getHttpResponse(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	return responseData
}

// unmarshalls news and returns an array of all articles from a news provider
func unmarshalNews(url string) NewsProvider {

	var newsProvider NewsProvider
	newsData := getHttpResponse(url)
	if err := json.Unmarshal(newsData, &newsProvider); err != nil {
		fmt.Println("Just receieved a message I couldn't decode in news.go FetchNewsSources 1:", string(newsData), err)
	}
	return newsProvider
}

// reads all news from all files that are listed in a textfile
// returns false if there was an error reading news
func ReadAllNews() ([]NewsProvider, bool) {
	// for now we will read one news file, later we will loop through more
	var allArticles []NewsProvider
	article, success := getNewsFromFile("data/news/hacker-news.json")
	if success == false {
		fmt.Println("ReadAllNews 1")
		return allArticles, false
	}
	allArticles = append(allArticles, article)
	//fmt.Println(allArticles)
	return allArticles, true
}

//gets news from file and unmarshalls to be passed to the front end for templating
// returns true if sucessfully reads and unmarshals
func getNewsFromFile(path string) (NewsProvider, bool) {

	newsData, err := ioutil.ReadFile(path)
	var newsProvider NewsProvider

	if err != nil {
		fmt.Println(err)
		return newsProvider, false
	}

	if err := json.Unmarshal(newsData, &newsProvider); err != nil {
		fmt.Println("Just receieved a message I couldn't decode in news.go FetchNewsSources 1:", string(newsData), err)
		return newsProvider, false
	}
	return newsProvider, true
}

// returns news API key
func getApiKey() string {

	// The file path where the news API key is
	const apiKeyPath = "data/newsApiKey.txt"
	config, err := os.Open(apiKeyPath)
	defer config.Close()

	if err != nil {
		fmt.Println("news.go getApiKey 1", err)
	}

	scanner := bufio.NewScanner(config)

	// news API Key
	scanner.Scan()
	return scanner.Text()
}
