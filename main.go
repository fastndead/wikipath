package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("You haven't provided TWO args. First should be a start wiki link and the second should be end link.")
		os.Exit(1)
	}

	startLink := normalizeArg(args[0])
	targetLink := normalizeArg(args[1])

	findClosestLink(startLink, targetLink)
}

func findClosestLink(startLink string, targetLink string) []string {
	visitedNodeMap := make(map[string]bool)
	paths := [][]string{{startLink}}
	currentLevel := 1
	fmt.Printf("Current level: %d\n", currentLevel)

	for len(paths) > 0 {
		path := paths[0]
		paths = paths[1:]
		node := path[len(path)-1]

		if len(path) > currentLevel {
			currentLevel += 1
			fmt.Printf("Current level: %d\n", currentLevel)
		}

		if node == targetLink {
			fmt.Println("Found the path!")
			fmt.Println(path)
			return path
		}

		if isWikipediaURL(node) && !visitedNodeMap[node] {
			visitedNodeMap[node] = true
			links := filterLinks(getLinksFromUrl(node))

			for _, link := range links {
				newPath := append(path, normalizeUrl(link))
				paths = append(paths, newPath)
				if link == targetLink {
					fmt.Println("Found the path!")
					fmt.Println(newPath)
					return newPath
				}
			}
		}
	}
	fmt.Println("Didn't find a node!")
	return nil
}

func normalizeArg(urlString string) string {
	index := strings.Index(urlString, "/wiki/")
	if index == -1 {
		fmt.Println("'/wiki/' not found in the string")
		return ""
	}

	urlWithoutDomain := urlString[index:]
	unencodedUrl, _ := url.PathUnescape(urlWithoutDomain)
	return url.PathEscape(unencodedUrl)
}

func normalizeUrl(urlString string) string {
	unencodedUrl, _ := url.PathUnescape(urlString)
	return url.PathEscape(unencodedUrl)
}

func filterLinks(links []string) []string {
	var filtered []string
	for _, link := range links {
		if isWikipediaURL(link) {
			filtered = append(filtered, link)
		}
	}
	return filtered
}

func getPage(url string) (string, error) {
	wikiUrlPrefix := "https://en.wikipedia.org"
	resp, err := http.Get(wikiUrlPrefix + url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	return string(body), nil
}

func getMainElementNode(htmlDoc string) (*html.Node, error) {
	parsedHtmlDoc, err := html.Parse(strings.NewReader(htmlDoc))

	if err != nil {
		return nil, err
	}

	var f func(*html.Node) *html.Node
	f = func(n *html.Node) *html.Node {

		isNodeMain := func(n *html.Node) bool {
			if n.Data == "main" {
				return true
			}

			for _, attr := range n.Attr {
				if attr.Key == "role" && attr.Val == "main" {
					return true
				}
			}

			return false
		}

		if isNodeMain(n) {
			return n
		}

		if n.FirstChild != nil {
			foundNode := f(n.FirstChild)
			if foundNode != nil {
				return foundNode
			}
		}

		if n.NextSibling != nil {
			foundNode := f(n.NextSibling)
			if foundNode != nil {
				return foundNode
			}
		}

		return nil
	}

	return f(parsedHtmlDoc), nil
}

func getAnchorElementsList(sourceNode *html.Node) []*html.Node {
	var anchorsList []*html.Node

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Data == "a" {
			anchorsList = append(anchorsList, n)
		}

		if n.FirstChild != nil {
			f(n.FirstChild)
		}

		if n.NextSibling != nil {
			f(n.NextSibling)
		}
	}

	f(sourceNode)
	return anchorsList
}

func getHrefFromAnchor(anchor *html.Node) string {
	if anchor.Type == html.ElementNode {
		for _, attr := range anchor.Attr {
			if attr.Key == "href" {
				return attr.Val
			}
		}
	}
	return ""
}

func getLinksFromUrl(urlString string) []string {
	//fmt.Printf("Getting from: %s \n", url)
	unescapedUrl, _ := url.PathUnescape(urlString)
	page, err := getPage(unescapedUrl)
	if err != nil {
		fmt.Println("Error happened")
	}

	mainNode, err := getMainElementNode(page)
	if err != nil {
		fmt.Println("Error while parsing")
	}
	if mainNode == nil {
		fmt.Println("Couldn't find main node.")
		fmt.Printf("URL: %s\n", urlString)
	}

	anchorsList := getAnchorElementsList(mainNode)

	var links []string

	for _, anchor := range anchorsList {
		href := getHrefFromAnchor(anchor)

		if href != "" {
			links = append(links, href)
		}
	}

	return links
}

func isWikipediaURL(urlString string) bool {
	urlDecoded, _ := url.PathUnescape(urlString)
	specialLinksPrefixes := []string{"/wiki/Special:", "/wiki/Wikipedia:", "/wiki/Help:", "/wiki/MediaWiki:", "/wiki/User:", "/wiki/Category:", "/wiki/Template:", "/wiki/Template_talk:", "/wiki/File:"}
	isWikiUrl := strings.HasPrefix(urlDecoded, "/wiki/")
	var isSpecialWikiUrl bool
	for _, prefix := range specialLinksPrefixes {
		if strings.HasPrefix(urlDecoded, prefix) {
			isSpecialWikiUrl = true
		}
	}
	return isWikiUrl && !isSpecialWikiUrl
}

func includes(array []string, targetString string) bool {
	for _, el := range array {
		if el == targetString {
			return true
		}
	}
	return false
}
