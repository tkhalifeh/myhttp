# Myhttp debugger

Simple tool that fetches html page contents from list of urls and prints the MD5 hash for each url

# Features
* Performs the requests in parallel. 
* The order or the addresses is not maintained when reporting results.
* Ability to limit the number of parallel requests. use `-parallel` flag (default 10)
* example usage `./myhttp -parallel 3 adjust.com google.com facebook.com yahoo.com yandex.com twitter.com
  reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com`

# Getting Started
## Building project
run `make build`

## Runnings tests
run `make test`