package smtpd

import (
	"container/list"
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
)

// Take "user+ext@host.com" and return "user", aka the mailbox we'll store it in
func ParseMailboxName(emailAddress string) (result string) {
	result = strings.ToLower(emailAddress)
	if idx := strings.Index(result, "@"); idx > -1 {
		result = result[0:idx]
	}
	if idx := strings.Index(result, "+"); idx > -1 {
		result = result[0:idx]
	}
	return result
}

// Take a mailbox name and hash it into the directory we'll store it in
func HashMailboxName(mailbox string) string {
	h := sha1.New()
	io.WriteString(h, mailbox)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// JoinStringList joins a List containing strings by commas
func JoinStringList(listOfStrings *list.List) string {
	if listOfStrings.Len() == 0 {
		return ""
	}
	s := make([]string, 0, listOfStrings.Len())
	for e := listOfStrings.Front(); e != nil; e = e.Next() {
		s = append(s, e.Value.(string))
	}
	return strings.Join(s, ",")
}

// ValidateDomainPart returns true if the domain part complies to RFC3696, RFC1035
func ValidateDomainPart(domain string) bool {
	if len(domain) == 0 {
		return false
	}
	if len(domain) > 255 {
		return false
	}
	if domain[len(domain)-1] != '.' {
		domain += "."
	}
	prev := '.'
	labelLen := 0
	hasLetters := false

	for _, c := range domain {
		switch {
		case ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '_':
			// Must contain some of these to be a valid label
			hasLetters = true
			labelLen++
		case '0' <= c && c <= '9':
			labelLen++
		case c == '-':
			if prev == '.' {
				// Cannot lead with hyphen
				return false
			}
		case c == '.':
			if prev == '.' || prev == '-' {
				// Cannot end with hyphen or double-dot
				return false
			}
			if labelLen > 63 {
				return false
			}
			if !hasLetters {
				return false
			}
			labelLen = 0
			hasLetters = false
		default:
			// Unknown character
			return false
		}
		prev = c
	}

	return true
}
