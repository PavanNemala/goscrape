# goscrape

# Start the service as below

go run main.go

# 1. once service is started, you can scrape the amazon website by hitting below url

http://127.0.0.1:8000/amazonscrape?url=https://www.amazon.com/PlayStation-4-Pro-1TB-Console/dp/B01LOP8EZC/

# 2. we can create a document with payload by hitting the below url then it will create a doc file in the same path where this program exists

http://127.0.0.1:8000/createdoc
