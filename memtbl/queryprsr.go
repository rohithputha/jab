package memtbl

import (
	"errors"
	"strings"
)

// right now supports only simple commands like Get, Put, Delete etc.
// should support complex sql parsing in the future
/*
	examples: (simple KV ops)
1. Put key1 value1 In Table1
2. Get Key1 From Table1
3. Get keyName > 50 From Table1 : phase 2
4. Schema Table1 (key1 key type, value1 type, value2 type.....) : phase 2

	phase 3 and beyond (support for SQL queries)

1. Select * From Table1
	Where table1.key1 > 100
... more examples to come.
*/

// parses the get command
// takes the token string array, returns the remaining tokens, keyName, TableName, expected command length
func parseGetCommand(token []string) ([]string, string, string, int, error) {
	const commandLength = 4
	if strings.ToUpper(token[0]) != "GET" {
		return nil, "", "", commandLength, errors.New("invalid GET command")
	}
	getCommandTokens := token[1:commandLength]
	if len(getCommandTokens) < commandLength-1 && strings.ToUpper(getCommandTokens[1]) != "FROM" {
		return nil, "", "", commandLength, errors.New("invalid GET command")
	}
	return token[commandLength:], getCommandTokens[0], getCommandTokens[2], commandLength, nil
}

// parses the put command
// takes the token string array, returns remainign tokens, keyname, value, table name, expected command len, error
func parsePutCommand(token []string) ([]string, string, string, string, int, error) {
	const commandLength = 5
	if strings.ToUpper(token[0]) != "PUT" {
		return nil, "", "", "", commandLength, errors.New("invalid PUT command")
	}
	putCommandTokens := token[1:commandLength]
	if len(putCommandTokens) < commandLength-1 && strings.ToUpper(putCommandTokens[1]) != "IN" {
		return nil, "", "", "", commandLength, errors.New("invalid PUT command")
	}
	return token[commandLength:], putCommandTokens[0], putCommandTokens[1], putCommandTokens[3], commandLength, nil
}

type parsedContent struct {
	tokens  []string
	command string
	tables  []string
	reqVars []string
	cond    []string // not used right now
	input   []string
}

// / returns command,
func parseTokens(tokens []string) ([]string, []*parsedContent, error) {
	prsConts := make([]*parsedContent, 0)
	for len(tokens) != 0 {
		switch strings.ToUpper(tokens[0]) {
		case "GET":
			newToks, key, table, comlen, err := parseGetCommand(tokens)
			if err != nil {
				return tokens, nil, err
			}
			prsConts = append(prsConts, &parsedContent{
				tokens:  tokens[:comlen],
				command: "GET",
				tables:  []string{table},
				reqVars: []string{key},
				cond:    nil,
				input:   nil,
			})
			tokens = newToks
			break
		case "PUT":
			newToks, key, value, table, comlen, err := parsePutCommand(tokens)
			if err != nil {
				return tokens, nil, err
			}

			prsConts = append(prsConts, &parsedContent{
				tokens:  tokens[:comlen],
				command: "PUT",
				tables:  []string{table},
				reqVars: []string{key},
				cond:    nil,
				input:   []string{value},
			})
			tokens = newToks
			break

		}
	}
	return tokens, prsConts, nil
}
func getParsedContent(command string) ([]*parsedContent, error) {
	tokens := strings.Fields(command)
	_, prsConts, err := parseTokens(tokens)
	if err != nil {
		return nil, err
	}
	return prsConts, nil
}
