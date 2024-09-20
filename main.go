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
		"Usage of %s:\n"+
			"\tReads in json from stdin or from a file,\n"+
			"\tprints out DynamoDB formatted JSON\n\n", exeName)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr,
		"\nExamples:\n"+
			"\t%s -f <json-file>\n"+
			"\t%s < file.json\n"+
			"\techo '{ \"a\": \"b\" }' | %s \n", exeName, exeName, exeName)
}

func main() {
	var readIn []byte
	var err error
	fileIn := flag.String("file", "", "Read in from file")
	flag.Usage = Usage
	flag.Parse()
	if *fileIn != "" {
		readIn, err = os.ReadFile(*fileIn)
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
