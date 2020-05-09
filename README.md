# PlantUML Go Client

This project provides a handy CLI for PlantUML users.

## Motivation

* Self-contained tool
* Non-Java
* Able to work with hosted PlantUML server
* Produces "Text Format"
* Produces Link
* Produced Images (Wow!)
* Useful app I wanted to improve forked from [this repo](https://github.com/yogendra/plantuml-go)

## Usage

Get go package first.

```shell
go get github.com/acarlson99/goplantuml
goplantuml my-uml.puml
```

```shell
$ echo "@startuml
a -> b : hello world
@enduml" | goplantuml                  # reads from stdin, outputs to `uml_out.png`
$ goplantuml test.puml                 # reads from file, outputs to `test.png`
$ goplantuml -format txt test.puml     # reads from file, outputs to `test.txt`
$ cat test.txt
     ┌───┐          ┌─────┐
     │Bob│          │Alice│
     └─┬─┘          └──┬──┘
       │    hello      │   
       │──────────────>│   
     ┌─┴─┐          ┌──┴──┐
     │Bob│          │Alice│
     └───┘          └─────┘
$ goplantuml -help
  -format format
    	Output format type. (Options: png,svg,txt) (default "png")
  -help
    	Show help (this) text
  -server server
    	Plantuml server address. Used when generating link or extracting output (default "http://plantuml.com/plantuml")
  -type string
    	Indicates if output type. (Options: save,link,hash) (default "save")
```
