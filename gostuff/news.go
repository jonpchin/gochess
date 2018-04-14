package gostuff

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
	PageTitle string
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

	logFile, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer logFile.Close()
	log := log.New(logFile, "", log.LstdFlags|log.Lshortfile)

	client := TimeOutHttp(5)
	response, err := client.Get(newsSourceList)
	if response == nil {
		log.Println("FetchNewsSources URL time out for ", newsSourceList)
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var newsProviders NewsProviders

	if err := json.Unmarshal(responseData, &newsProviders); err != nil {
		log.Println("Just receieved a message I couldn't decode in news.go FetchNewsSources 1:", string(responseData), err)
		return
	}

	apiKey := getApiKey()
	for _, source := range newsProviders.Sources {
		url := "https://newsapi.org/v1/articles?source=" + source.ID + "&apiKey=" + apiKey
		saveNewsToFile(source.ID, url)
	}
}

// fetches and saves a single source of news
// @param url to fetch JSON of articles from a single news provider
func saveNewsToFile(filename string, url string) {

	responseData := getHttpResponse(url)
	if responseData == nil {
		return
	}
	newsOutputPath := "privatedata/news/" + filename + ".json"
	err := ioutil.WriteFile(newsOutputPath, responseData, 0666)
	if err != nil {
		fmt.Println(err)
	}
}

// returns the http response of the url in byte
// convert to string for human readable format
func getHttpResponse(url string) []byte {

	client := TimeOutHttp(5)
	response, err := client.Get(url)
	if err != nil {
		fmt.Println("getHttpResponse 0", err)
		return nil
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("getHttpResponse 1", err)
		return nil
	}
	return responseData
}

// reads all news from all files that are listed in a textfile
func (allArticles *AllNewsProviders) ReadAllNews() {

	const newsConfigPath = "privatedata/newsConfig.txt"
	config, err := os.Open(newsConfigPath)
	defer config.Close()

	if err != nil {
		fmt.Println("news.go ReadAllNews 0", err)
		return
	}

	allArticles.PageTitle = "News"

	scanner := bufio.NewScanner(config)

	for scanner.Scan() {
		fileName := scanner.Text()
		var newsProvider NewsProvider
		success := newsProvider.getNewsFromFile("privatedata/news/" + fileName + ".json")
		if success == false {
			fmt.Println("error reading news source for ", fileName)
		}
		newsProvider.convertToHttps()
		allArticles.Providers = append(allArticles.Providers, newsProvider)
	}
}

//gets news from file and unmarshalls to be passed to the front end for templating
// returns true if successfully reads and unmarshals
func (newsProvider *NewsProvider) getNewsFromFile(path string) bool {

	newsData, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println("getNewsFromFile", err)
		return false
	}

	if err := json.Unmarshal(newsData, &newsProvider); err != nil {
		fmt.Println("Just receieved a message I couldn't decode in news.go getNewsFromFile 1:",
			string(newsData), err)
		return false
	}
	return true
}

// updates all news files that are listed in the newsConfig.txt
func UpdateNewsFromConfig() {

	const newsConfigPath = "privatedata/newsConfig.txt"
	config, err := os.Open(newsConfigPath)
	defer config.Close()

	if err != nil {
		fmt.Println("news.go UpdateNewsFromConfig 1", err)
	}
	scanner := bufio.NewScanner(config)
	apiKey := getApiKey()

	for scanner.Scan() {
		fileName := scanner.Text()
		url := "https://newsapi.org/v1/articles?source=" + fileName + "&apiKey=" + apiKey
		saveNewsToFile(fileName, url)
	}
	// creates a cached news file, moved to template.go
	parseNewsCache()
}

// makes image url that are http into https if its valid
// in a NewsProvider, othwerise convert back to http if https times out
func (newsProvider *NewsProvider) convertToHttps() {

	//log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	for index, article := range newsProvider.Articles {

		// Check to make sure link is not empty
		if newsProvider.Articles[index].UrlToImage != "" {
			newsProvider.Articles[index].UrlToImage = strings.Replace(article.UrlToImage,
				"http://", "https://", 1)

			client := TimeOutHttp(5)
			client.Get(newsProvider.Articles[index].UrlToImage)
			// Below is commented out to reduce logging spam
			//if response == nil {
			//log.Println("convertToHttps URL time out for ",
			//newsProvider.Articles[index].UrlToImage, err)
			//}
		}
	}
}

// returns news API key
func getApiKey() string {

	// The file path where the news API key is
	const apiKeyPath = "privatedata/newsApiKey.txt"
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
