//Problem link
//
//
//
// https://www.hackerrank.com/challenges/caesar-cipher-1/problem
//
//
//

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"io"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"unicode"
// )

// /*
//  * Complete the 'caesarCipher' function below.
//  *
//  * The function is expected to return a STRING.
//  * The function accepts following parameters:
//  *  1. STRING s
//  *  2. INTEGER k
//  */

// func caesarCipher(s string, k int32) string {
// 	// Write your code here
// 	k = k % 26
// 	encrypted := make([]rune, len(s))

// 	for i, char := range s {
// 		if unicode.IsLetter(char) {
// 			if unicode.IsLower(char) {
// 				encrypted[i] = 'a' + (char-'a'+rune(k))%26
// 			} else if unicode.IsUpper(char) {
// 				encrypted[i] = 'A' + (char-'A'+rune(k))%26
// 			}
// 		} else {
// 			encrypted[i] = char
// 		}
// 	}

// 	return string(encrypted)
// }

// func main() {
// 	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

// 	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
// 	checkError(err)

// 	defer stdout.Close()

// 	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

// 	nTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
// 	checkError(err)
// 	n := int32(nTemp)
// 	fmt.Println(n)

// 	s := readLine(reader)

// 	kTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
// 	checkError(err)
// 	k := int32(kTemp)

// 	result := caesarCipher(s, k)

// 	fmt.Fprintf(writer, "%s\n", result)

// 	writer.Flush()
// }

// func readLine(reader *bufio.Reader) string {
// 	str, _, err := reader.ReadLine()
// 	if err == io.EOF {
// 		return ""
// 	}

// 	return strings.TrimRight(string(str), "\r\n")
// }

// func checkError(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
