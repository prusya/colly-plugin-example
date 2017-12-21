# colly-plugin-example

## Purpose

Showcase of a web scraper used as a go plugin

## Build status [![Build Status](https://travis-ci.org/prusya/colly-plugin-example.svg?branch=master)](https://travis-ci.org/prusya/colly-plugin-example)

## Description

Use https://github.com/gocolly/colly in a plugin

## Requirements

go 1.8+

## Installation

```bash
go get -t -v github.com/prusya/colly-plugin-example
```

## Compilation

Inside your GOPATH directory

```bash
cd src/github.com/prusya/colly-plugin-example
go build -buildmode=plugin -o plugins/bitcq/bitcq.so plugins/bitcq/main.go
go build .
```

## Testing

```bash
go test -v ./...
```
