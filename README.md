# Searxng-go

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/zvirgilx/searxng-go/deploy-push.yml)](https://github.com/zvirgilx/searxng-go/actions)
![GitHub Release](https://img.shields.io/github/v/release/zvirgilx/searxng-go?sort=semver)
[![Go Report Card](https://goreportcard.com/badge/github.com/zvirgilx/searxng-go/kernel)](https://goreportcard.com/report/github.com/zvirgilx/searxng-go/kernel)
![Go Version](https://img.shields.io/badge/go_version-%3E%3D1.21-blue)


---

Searxng-go is a metasearch engine written in Go, inspired by searxng, 
is designed to retrieve search results from multiple search engines simultaneously
and aggregate them into a result.

Being written in Go, Searxng-go benefits from the language's performance, concurrency, and simplicity.
Go's strong type system and functional standard library provide us with the tools to build a reliable 
and efficient metasearch engine.

Searxng-go can be customized to fit individual preferences and needs.
Searxng-go allows choosing the search engines, configuring the ranking algorithm,
and applying various filters to the search results, tailoring search experience according to specific requirements.
In addition, within this framework, you can easily and quickly develop customized search engines and display them in search results.

--- 

## Features

Searxng-go implements most of the impressive features of searxng, including but not limited to the following:

* protect user privacy from online tracking.
* enable multiple search engines, and aggregate their search result.
* score, sort and display results flexibly.
* support the extension of custom search engine.


## To start using Searxng-go

The architecture of searxng-go follows a front-end separation approach, 
where the front end is implemented using Vue in [Web](web), 
while the back end is implemented using Gin in [Kernel](kernel).

### API server 
To start a api server in develop environment:
```shell
cd kernel
go run main.go -l debug api -m debug -a :9999
```

Search keyword in command line:
```shell
go run main.go -l debug search superman
```

### Web server

The web server can be started by referring to [WebServer Starter](web/README.md).

### All in one

Use `docker-compose` to start both the `API server` and `Web server`.
```shell
docker compose up --build
```
Searxng-go will be deployed locally, you can access it by going to http://127.0.0.1:3000.

## Development

No matter what your skills and interests are, we welcome your participation and contribution.
Join the Searxng-go to build a better search engine framework!!!

### Setting up dev environment

Searxng-go is based on golang, so you need to install golang (version >= 1.21.5) from [Go](https://go.dev/dl/) on your computer.

After doing that, you can clone Searxng-go on your computer:
```shell
git clone https://github.com/zvirgilx/searxng-go.git
```

Installation related dependencies:
```shell
cd searxng-go/kernel
go mod tidy
```

Check if searxng-go can be started:
```shell
go run main.go -l debug search superman
```

Congratulations! You have completed the deployment of the Searxng-go, you can open the IDE for development and debugging as you want.

### Extending custom search engine

Searxng-go provides a convenient framework that allows developers to extend and customize their own search engines.

We use an interface to define the behavior of a search engine.
```go
type Engine interface {
	// Request reports how the engine initiates a request.
	Request(context.Context, *Options) error
	// Response reports how the engine parse the response.
	Response(context.Context, *Options, []byte) (*result.Result, error)
	// GetName returns engine name.
	GetName() string
	// ApplyConfig apply custom configuration to engine.
	ApplyConfig(config Config) error
}
```

When Searxng-go receives a search request, it first calls the **Request** method of each engine to generate a corresponding request.
Searxng-go initiates these requests asynchronously and returns the response to the engine's **Response** method for processing.

Custom search engine should implement how to convert a search request to a http url.
For example, the search engine extracts the search keywords and generates a http request.
```go
func (e *example) Request(ctx context.Context, opt *Options) error {
	// ... 
	
	// https://www.example.com/search?q=keyword 
	base, _ := url.parse("https://www.example.com")
	r := g.client.Get().Base(base).Path("search").Param("q", opts.Query)
	opts.Request = r
	
	// ...
}


```
Custom engine needs to parse the results of the response and generate search results.

```go

func (e *example) Response(ctx context.Context, opts *engine.Options, body []byte) (*result.Result, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, errors.New("error parsing document")
	}

	res := result.CreateResult(EngineNameExample, opts.PageNo)
	doc.Find("div.g").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3").First().Text()
		link, _ := s.Find("a").First().Attr("href")
		content := s.Find(".VwiC3b").First().Text()

		// ignore empty title or content result
		if title == "" || content == ""{
			return
		}
		
		// ...

		res.AppendData(&result.Data{
			Engine:  EngineNameExample,
			Title:   title,
			Url:     link,
			Content: content,
			Query:   opts.Query,
		})
	})
	
	return res, nil
}
```

## Customizing your searxng-go

The configuration file for Searxng-go is located in [configuration](kernel/config/default.yaml).
You can load your own configuration file by importing it at startup, it will overwrite the default configuration.

```shell
go run main.go api -c <your_config.yaml> 
```

### Custom search engine

You can customize the configuration for each engine.
Here is an example of configuring ElasticSearch for host, index and network proxies.

```yaml
elastic_search:
  enable: true  # enable this engine
  extra: # additional parameters required by the engine
    base_url: http://127.0.0.1:9200 # elasticsearch service
    index: test-1
    query_type: multi_match
    query_fields: ["title","description"]
  network:
    timeout: 3s
    proxy_url: https://www.proxy.com/your_own_proxy
```


### Custom scoring rule

Searxng-go provides a flexible scoring rule system that allows for scoring and sorting results from various search engines.
With searxng-go, you can customize scoring rules and assign scores to search results based on different criteria. 
These criteria can include keyword relevance, engine weight and more. 
By setting up appropriate scoring rules, you can personalize the sorting of search results based on user preferences and requirements.

Here is an example of How to score a result by the weight of search engines.

```yaml
rules:
  - name: "engine_weight" # rule name.
    score: 100 # score if match this rule.
    enable: true # enable this rule
    conditions: # A rule is matched only when all conditions are matched.
      - field: "engine"
        operator: "in"
        expects : ["imdb","elastic_search"]
```
imdb and elasticsearch engines will get 100 points, which make their search results will rank ahead of others.

You can set multiple conditions for a rule.

```yaml
rules:
  - name: "match_query_content"
    score: 5
    enable: true
    conditions:
      - field: "title"
        operator: "containAny"
        expects: ["$QUERY"]
      - field: "content"
        operator: "containAny"
        expects: ["$QUERY"]
```
The results containing query keywords in the title and content will receive an additional 5 points.


## License

The project is licensed under the [MIT License](LICENSE).