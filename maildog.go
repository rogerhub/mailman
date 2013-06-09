/**
 *  Helper functions for mailman.
 */

package mailman

import (
	"strings"
)

/**
 *  Because I didn't read rfc5322
 */
func ValidateEmail (e string) bool {
	return true
}

/**
 *  Internal list data type
 */
func ListContains (item string, list string) bool {
	elements := strings.Split(list, " ")
	if strings.Contains(item, "@") {
		itemParts := strings.Split(item, "@")
		for i := 0; i < len(elements); i ++ {
			if strings.HasPrefix(elements[i], "@") {
				if elements[i] == "@" + itemParts[len(itemParts) - 1] {
					return true
				}
			} else if strings.Contains(elements[i], "@") {
				if elements[i] == item {
					return true
				}
			} else {
				if elements[i] == itemParts[0] {
					return true
				}
			}
		}
	} else {
		for i := 0; i < len(elements); i ++ {
			if elements[i] == item {
				return true
			}
		}
	}
	return false
}
