# Wikipedia Shortest Path CLI

Have you ever wondered what's the shortest number of clicks that would get you from one page on Wikipedia to another? Me neither. But I made this CLI tool to do exactly that anyway. This is a command-line interface (CLI) program written in Go that finds the shortest path between two Wikipedia pages. Given two Wikipedia URLs, the program will crawl through the links on each page to find the shortest path between them.

## Usage

To use this program, run the following command:

```
go run main.go <start_url> <end_url>
```

Replace <start_url> and <end_url> with the URLs of the Wikipedia pages you want to find the shortest path between.

For example, if you want to find the shortest path between the Wikipedia pages for "Albert Einstein" and "Isaac Newton", you would run:

```
go run main.go https://en.wikipedia.org/wiki/Albert_Einstein https://en.wikipedia.org/wiki/Isaac_Newton
```
