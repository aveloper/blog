package config

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

func readString(reader *bufio.Reader, prompt, _default string, required bool) string {
	fmt.Printf("%s(%s): ", prompt, _default)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Panicf("Failed to read %s: %v", prompt, err)
		return ""
	}

	input = strings.TrimSpace(input)

	if input == "" && required {
		return readString(reader, prompt, _default, required)
	}

	if input == "" {
		return _default
	}

	return input
}

func readInt(reader *bufio.Reader, prompt string, _default int, required bool) int {
	s := readString(reader, prompt, fmt.Sprintf("%d", _default), required)

	input, err := strconv.ParseInt(s, 10, 32)
	if err != nil && errors.Is(err, strconv.ErrSyntax) {
		fmt.Println("Invalid integer for port")
		return readInt(reader, prompt, _default, required)
	}

	if err != nil {
		log.Panicf("Failed to read %s: %v", prompt, err)
	}

	return int(input)
}

func readPassword(prompt string) string {
	fmt.Printf("%s: ", prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Panicf("Failed to read %s: %v", prompt, err)
	}

	fmt.Println()

	return string(password)
}
