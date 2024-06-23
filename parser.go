package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func openAndParse(filename string) (*Value, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("there was a problem opening the file: %s", err))
	}

	buf := bytes.NewBuffer(fileBytes)
	tree, err := parse(buf)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("parse error: %s", err))
	}

	return tree, nil
}

// top level parse
func parse(buf *bytes.Buffer) (*Value, error) {
	tree := nullValue
	err := parseExp(&tree, buf)

	if err != nil {
		return nil, err
	}
	return tree, nil
}

func parseExp(head **Value, buf *bytes.Buffer) error {
	for {
		b, err := buf.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		c := rune(b)
		node := nullValue

		if c == '(' {
			// create a new list and recurse
			err = parseExp(&node, buf)
			if err != nil {
				return err
			}
			node = cons(node, nullValue)
		} else if c == ')' {
			// close the current list
			break
		} else if isSpace(c) {
			// this is a space
			continue
		} else if isWhitespace(c) {
			// ignore whitespace
			continue
		} else if c == '"' {
			// reading a string constant
			str, err := buf.ReadBytes(byte('"'))
			if err != nil && err != io.EOF {
				return err
			}
			val := &Value{
				kind: String,
				val:  strings.Trim(string(str), "\""),
			}
			node = cons(val, nullValue)
		} else if c == '\'' {
			// reading a symbol
			var sym strings.Builder
			var prev rune

			for {
				b, err := buf.ReadByte()
				if err != nil && err != io.EOF {
					return err
				}

				t := rune(b)
				// keep reading until whitespace or closing paren
				if isWhitespace(t) || t == ')' {
					if prev != '(' {
						buf.UnreadByte()
					}
					if t == ')' {
						sym.WriteRune(t)
					}
					break
				}
				sym.WriteRune(t)
				prev = t
			}

			val := &Value{
				kind: Symbol,
				val:  sym.String(),
			}
			node = cons(val, nullValue)
		} else if c == ';' {
			// this is a line comment
			// consume the buffer until the newline
			_, err = buf.ReadBytes(byte('\n'))
			if err != nil && err != io.EOF {
				return err
			}
		} else if isChar(c) {
			err = buf.UnreadByte()
			if err != nil {
				return err
			}
			var token strings.Builder

			for {
				b, err := buf.ReadByte()
				if err != nil && err != io.EOF {
					return err
				}

				t := rune(b)
				if !isChar(t) {
					buf.UnreadByte()
					break
				}
				token.WriteRune(t)
			}

			var val *Value
			tok := token.String()

			num, cerr := strconv.ParseInt(tok, 10, 64)
			// try to parse an integer
			if cerr == nil {
				val = &Value{
					kind: Integer,
					val:  num,
				}
			} else {
				// try to parse a float
				fl, cerr := strconv.ParseFloat(tok, 64)
				if cerr == nil {
					val = &Value{
						kind: Float,
						val:  fl,
					}
				} else {
					// try to parse a boolean
					bl, cerr := strconv.ParseBool(tok)
					if cerr == nil {
						val = &Value{
							kind: Boolean,
							val:  bl,
						}
					} else {
						// this is just a name
						val = &Value{
							kind: Name,
							val:  strings.Trim(tok, " "),
						}
					}
				}
			}

			node = cons(val, nullValue)
		} else {
			panic(fmt.Sprintf("unrecognized token: %c", c))
		}

		if !isNull(node) {
			if isNull(*head) {
				// the first item in our parse tree
				*head = node
			} else {
				*head = listAppend(*head, node)
			}
		}
	}

	return nil
}

func isChar(c rune) bool {
	if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '-' || c == '?' || c == '+' || c == '-' || c == '*' || c == '=' || c == '/' || c == '>' || c == '<' || c == '!' || c == '.' {
		return true
	}
	return false
}

func isSpace(c rune) bool {
	if c == ' ' {
		return true
	}
	return false
}

// replace with unicode.IsSpace()
func isWhitespace(c rune) bool {
	switch c {
	case ' ':
		return true
	case '\n':
		return true
	case '\r':
		return true
	case '\t':
		return true
	}
	return false
}
