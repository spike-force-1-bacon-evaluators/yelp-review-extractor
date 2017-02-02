FROM golang:1.7.4-wheezy

COPY . src/github.com/github.com/bacon/yelp-review-extractor

WORKDIR src/github.com/github.com/bacon/yelp-review-extractor

ENTRYPOINT ["go", "test", "./src/...", "-v", "-cover"]
