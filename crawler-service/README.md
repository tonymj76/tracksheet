# Crawling Service

- [Overview :notebook_with_decorative_cover:](#overview-notebook_with_decorative_cover)
- [Setting up :gear:](#Setting-up-gear)
- [Accessing end points](#Accessing-end-points)

## Overview :notebook_with_decorative_cover

There is two search Engine implemented.

- The 1st solution takes in the URL or zipped sitemap of the site to crawl. If a URL is provided, it will search the robots.txt file and get the XML sitemap which returns all the product links in the XML files, while if a zipped sitemap is provided it parse it and return the product links.
- The 2nd solution usings recursive page depth to get all the links of the product URL.

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



### Accessing end points

To access the crawl service access the end point GET `/crawler/scrap`, pass the following as queries

* link (link of site to crawl)
* depth (use if using recursive depth search method); an int value. The more the depth, the more itemsretrieved



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



Note: This option should only be used when there is a visible pattern in product url eg site.com/prd/123 and site.com/prd/345

9. To manually provide product url example POST `crawler/scrap`, pass in the following body

```
{ 
    Str1  (product link eg 1)
	Str2   (product link eg 2)
	Str3   (product link eg 3)
	Depth   (depth)
	Link  (url to scrap)
}
```