# crawling-service

-[] Create a .env in the crawler-service dir. The content of the file should be 

```
DATABASE_URI = your_postgres_url
```

### Accessing end points

1. To access the crawl service access the end point GET `/crawler/scrap`, pass the following as queries
* accept (ignore to use default)
* user_aget (ignore to use daefult)
* link (link of site to crawl)
* depth (use if using recursive depth search method); an int value. The more the depth, the more itemsretrieved


2. To get all urls in DB GET `/crawler/url` pass the following queries
* page (table page)
*limit (table limit)
This would return a list of all url items in the DB


3. To get the product url details of a url GET `/crawler/product` pass the following queries
* page (table page)
* limit (table limit)
* link (url item which you want to get their product)


4. To delete a product url GET `crawler/productdelete` pass the following queries
* link (product link you want to delete)


5. To delete a url item, and all related  GET `/crawler/urldelete` pass the following queries
* link (url link to delete)

5. To search for a url GET `/crawler/searchurl` pass the following queries
* page (table page)
* limit (table limit)
* k (the search string)

5. To search for a product GET `/crawler/searchproduct` pass the following queries
* page (table page)
* limit (table limit)
* k (the search string)

6. To get path maps GET `/crawler/getpaths` pass the following queries
* page (table page)
* limit (table limit)
