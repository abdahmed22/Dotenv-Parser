package main

import (
	"fmt"

	dotenv "github.com/codescalersinternships/Dotenv-Abdelrahman-Mahmoud/pkg"
)

func main() {

	parser := dotenv.EnvContent{}

	outputMap, err := parser.LoadFromFile("test2.txt")

	fmt.Println(err)

	fmt.Println(outputMap)

	outputMap, err = parser.LoadFromFiles([]string{"test1.txt", "test2.txt"})

	fmt.Println(err)

	fmt.Println(outputMap)

	err = parser.LoadFromString(`#This is a comment1
#This is a comment2

key1 = value1
key2 = value2
key3 = value3
key4 = value4

key5 = value5
key7 = value7`)

	fmt.Println(parser.GetEnv())

	fmt.Println(parser.Get("key5"))

	fmt.Println(parser.Get("key0"))

	fmt.Println(err)
}
