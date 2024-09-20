package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func Usage() {
	exeName := path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Convert JSON to DynamoDB format\n\n"+
			"Reads JSON input from a file, standard input, or directly from the command line,\n"+
			"and converts it to DynamoDB-compatible JSON format.\n\n"+
			"USAGE:\n"+
			"  %s [OPTIONS] [FILE | JSON_STRING]\n\n"+
			"OPTIONS:\n", exeName)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr,
		"\nARGUMENTS:\n"+
			"  FILE                    Path to the input JSON file\n"+
			"  JSON_STRING             JSON data provided directly on the command line\n\n"+
			"NOTES:\n"+
			"  * Only one input method is allowed at a time.\n"+
			"  * Precedence: `--file` > command-line JSON string > standard input\n\n"+
			"EXAMPLES:\n"+
			"  # From a file:\n"+
			"  %s -f input.json\n\n"+
			"  # From command-line arguments (JSON string):\n"+
			"  %s '{\"a\": \"b\"}'\n\n"+
			"  # From standard input (piped):\n"+
			"  cat input.json | %s \n", exeName, exeName, exeName)
}

func main() {
	var readIn []byte
	var err error
	fileIn := flag.String("file", "", "Read in from file")
	flag.Usage = Usage
	flag.Parse()
	if *fileIn != "" {
		readIn, err = os.ReadFile(*fileIn)
	} else if len(os.Args) > 1 {
		readIn = []byte(os.Args[1])
	} else {
		readIn, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		panic(err)
	}
	var out []byte
	if readIn[0] == '[' {
		log.Fatal("Top level arrays are not supported in DynamoDB JSON format")
	} else {
		var jsonIn any

		err = json.Unmarshal(readIn, &jsonIn)
		if err != nil {
			panic(err)
		}

		attrMap, err := dynamodbattribute.Marshal(jsonIn)
		if err != nil {
			panic(err)
		}
		anyMap := AttribValueMapToAnyMap(attrMap.M)

		out, err = json.Marshal(anyMap)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(string(out))
}

// AttribValueToAnyMap takes in a dynamodb.AttributeValue
// and returns a map[string]any of it's fields, recursively.
// It removes any fields with values set to nil
func AttribValueToAnyMap(v *dynamodb.AttributeValue) map[string]any {
	mp := make(map[string]any)
	if v.B != nil {
		mp["B"] = v.B
	}
	if v.BOOL != nil {
		mp["BOOL"] = *v.BOOL
	}
	if v.BS != nil {
		mp["BS"] = v.BS
	}
	if v.L != nil {
		mp["L"] = AttribValueArrToAnyMap(v.L)
	}
	if v.M != nil {
		mp["M"] = AttribValueMapToAnyMap(v.M)
	}
	if v.N != nil {
		mp["N"] = *v.N
	}
	if v.NS != nil {
		mp["NS"] = v.NS
	}
	if v.NULL != nil {
		mp["NULL"] = *v.NULL
	}
	if v.S != nil {
		mp["S"] = *v.S
	}
	if v.SS != nil {
		mp["SS"] = v.SS
	}

	return mp
}

func AttribValueArrToAnyMap(va []*dynamodb.AttributeValue) []map[string]any {
	arr := []map[string]any{}
	for _, v := range va {
		if v != nil {
			arr = append(arr, AttribValueToAnyMap(v))
		}
	}
	return arr
}

func AttribValueMapToAnyMap(va map[string]*dynamodb.AttributeValue) map[string]any {
	mp := make(map[string]any)
	for k, v := range va {
		if v != nil {
			mp[k] = AttribValueToAnyMap(v)
		}
	}
	return mp
}
