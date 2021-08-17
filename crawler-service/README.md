# Crawling Service

- [Crawling Service](#crawling-service)
  - [Setting up :gear](#setting-up-gear)
  - [RUNNING THE CRAWLER](#running-the-crawler)
  - [Accessing end points](#accessing-end-points)
    - [Aamzon setup](#aamzon-setup)



## Setting up :gear

-[] Create an env.json file in the crawler-service dir. The content of the file should be
eg 

```sh
{
"DATABASE_URI":"postgres://postgres:postgres@localhost:5432/tracksheet?sslmode=disable",
"PORT":":3001",
"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
"userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11",
"HTTP_PROXY":"https://proxy"
}

```

* Be sure to have a progresql DB setup
* To start the service run `go run service.go`



<br /> <br />
<hr />

## RUNNING THE CRAWLER
- Ensure you have [GNU Make](https://www.gnu.org/software/make/) installed

```sh
make crawl
```

<hr />
<br /> <br />


## Accessing end points

To access the crawl service access the end point GET `/crawler/scrap`, pass the following as queries

* link (link of site to crawl)
* depth (use if using recursive depth search method); an int value. The more the depth, the more itemsretrieved
* worker (number of workers to configure)


To get the product url details of a url GET `/crawler/product` pass the following queries

* page (table page)
* limit (table limit)
* link (url item which you want to get their product)

To delete a product url GET `crawler/productdelete` pass the following queries

* link (product link you want to delete)


To search for a product GET `/crawler/searchproduct` pass the following queries

* page (table page)
* limit (table limit)
* k (the search string)

To scrap a website sitemap, GET `/crawler/scrapsitemap` pass the following queries
* link (the sitemap to scrap)





### Aamzon setup
* Upload the Go build and `env.json` to the amazon env with all the configuration options
* Run the following command in the env to set up chromium (assuming it is a deb environment)
`wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb`
`sudo apt install ./google-chrome-stable_current_amd64.deb`
For other env visit https://go-rod.github.io/#/compatibility?id=os:

* Start the server with `./goBuild` (avoid using sudo command to start the build)