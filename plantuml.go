package main

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	formatTXT      = "txt"
	formatPNG      = "png"
	formatSVG      = "svg"
	outputTypeHASH = "hash"
	outputTypeLINK = "link"
	outputTypeSAVE = "save"
	mapper         = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
)

type Options struct {
	Server     string
	Format     string
	OutputType string
	FileNames  []string
}

func main() {
	opts := parseArgs()

	if len(opts.FileNames) == 0 {
		outFileName := "uml_out"
		err := process(&opts, os.Stdin, outFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error processing stdin %s: %v\n", outFileName, err)
		}
	}
	for _, filename := range opts.FileNames {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening file %s: %v\n", filename, err)
		}
		err = process(&opts, f, strings.TrimSuffix(filename, filepath.Ext(filename)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error processing file %s: %v\n", filename, err)
		}
	}
}

func process(options *Options, file *os.File, basename string) error {
	textFormatB, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	textFormat := encodeAsTextFormat(textFormatB)

	if options.OutputType == outputTypeHASH {
		fmt.Printf("%s: %s\n", basename, textFormat)
	} else if options.OutputType == outputTypeLINK {
		fmt.Printf("%s: %s/%s/~1%s\n", basename, options.Server, options.Format, textFormat)
	} else if options.OutputType == outputTypeSAVE {
		link := fmt.Sprintf("%s/%s/~1%s", options.Server, options.Format, textFormat)

		fileName := fmt.Sprintf("%s.%s", basename, options.Format)
		output, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		response, err := http.Get(link)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if response.StatusCode != 200 {
			return fmt.Errorf("Error in Fetching `%s`: %s", link, response.Status)
		}

		io.Copy(output, response.Body)
		response.Body.Close()
		output.Close()
	}
	return nil
}

func encodeAsTextFormat(raw []byte) string {
	compressed := deflate(raw)
	return base64_encode(compressed)
}

func deflate(input []byte) []byte {
	var b bytes.Buffer
	w, _ := zlib.NewWriterLevel(&b, zlib.BestCompression)
	w.Write(input)
	w.Close()
	return b.Bytes()
}

func base64_encode(input []byte) string {
	var buffer bytes.Buffer
	inputLength := len(input)
	for i := 0; i < 3-inputLength%3; i++ {
		input = append(input, byte(0))
	}

	for i := 0; i < inputLength; i += 3 {
		b1, b2, b3, b4 := input[i], input[i+1], input[i+2], byte(0)

		b4 = b3 & 0x3f
		b3 = ((b2 & 0xf) << 2) | (b3 >> 6)
		b2 = ((b1 & 0x3) << 4) | (b2 >> 4)
		b1 = b1 >> 2

		for _, b := range []byte{b1, b2, b3, b4} {
			buffer.WriteByte(byte(mapper[b]))
		}
	}
	return string(buffer.Bytes())
}

func parseArgs() Options {
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)
	server := flag.String("server", "http://plantuml.com/plantuml", "Plantuml `server` address. Used when generating link or extracting output")
	format := flag.String("format", formatPNG, fmt.Sprintf("Output `format` type. (Options: %s,%s,%s)", formatPNG, formatSVG, formatTXT))
	outputType := flag.String("type", outputTypeSAVE, fmt.Sprintf("Indicates if output type. (Options: %s,%s,%s)", outputTypeSAVE, outputTypeLINK, outputTypeHASH))
	help := flag.Bool("help", false, "Show help (this) text")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *format != formatPNG && *format != formatSVG && *format != formatTXT {
		fmt.Println("Invalid format")
		os.Exit(1)
	}

	if *outputType != outputTypeHASH && *outputType != outputTypeLINK && *outputType != outputTypeSAVE {
		fmt.Println("Invalid output type")
		os.Exit(1)
	}

	opts := Options{
		Server:     *server,
		Format:     *format,
		OutputType: *outputType,
		FileNames:  flag.Args()}
	return opts
}
