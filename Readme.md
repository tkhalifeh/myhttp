# Myhttp debugger

Simple tool that fetches html page contents from list of urls and prints the MD5 hash for each url

# Features
* Performs the requests in parallel. 
* The order of the addresses is not maintained when reporting results.
* Ability to limit the number of parallel requests using `-parallel` flag (default 10)
* Example usage `./myhttp -parallel 3 adjust.com google.com facebook.com yahoo.com yandex.com twitter.com
  reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com`

# Getting Started
## Building the project
run `make build`

## Running tests
run `make test`