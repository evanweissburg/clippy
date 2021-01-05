package mnemonic

import (
	"strings"
	"testing"
	"unicode"
)

func TestPatterns(t *testing.T) {
	for lengthKey, p := range patterns {
		for i, pattern := range p {
			if len(pattern) != lengthKey {
				t.Errorf("len(patterns[%v][%v]) = %v; expected %v", lengthKey, i, len(pattern), lengthKey)
			}

			for _, partOfSpeech := range pattern {
				if _, ok := validPartsOfSpeech[partOfSpeech]; !ok {
					t.Errorf("patterns[%v][%v] contains invalid part of speech %q", lengthKey, i, partOfSpeech)
				}
			}
		}
	}
}

func TestVocabulary(t *testing.T) {
	for partOfSpeech, _ := range validPartsOfSpeech {
		if _, ok := vocabulary[partOfSpeech]; !ok {
			t.Errorf("vocabulary has no entry for part of speech %q", partOfSpeech)
			continue
		}

		for i := 0; i < 26; i++ {
			letter := rune('a' + i)
			if len(vocabulary[partOfSpeech]) <= i {
				t.Errorf("vocabulary[%q] doesn't contain an entry for letter %q", partOfSpeech, letter)
				break
			}
			word := vocabulary[partOfSpeech][i]
			if len(word) == 0 || rune(word[0]) != letter {
				t.Errorf("vocabulary[%q][%v] = %q; expected a word starting with %q", partOfSpeech, i, word, letter)
				break
			}
		}
	}
}

func isValidPhrase(phrase, acronym string) bool {
	if strings.HasSuffix(phrase, " ") {
		return false
	}

	words := strings.Split(phrase, " ")
	for i, letter := range acronym {
		if len(words) <= i || len(words[i]) == 0 || rune(words[i][0]) != unicode.ToLower(letter) {
			return false
		}
	}

	return true
}

func TestCreatePhrase(t *testing.T) {
	validAcronyms := []string{"bc", "dEf", "HiJk", "LMNOP", "vuTsrq"}

	for _, acronym := range validAcronyms {
		phrase, err := CreatePhrase(acronym)
		if err != nil {
			t.Errorf("Error for valid acronym %q: %v", acronym, err)
			continue
		}

		if isValidPhrase(phrase, acronym) {
			t.Logf("%q --> %q", acronym, phrase)
		} else {
			t.Errorf("%q --> %q is invalid phrase", acronym, phrase)
		}
	}

	invalidAcronyms := []string{"", "ab1", "...", "/\\", "aaaaaaaa", "ffĀ", "šaa"}
	for _, acronym := range invalidAcronyms {
		_, err := CreatePhrase(acronym)
		if err == nil {
			t.Errorf("Invalid acronym %q did not produce error", acronym)
		}
	}
}

func TestCreateSentence(t *testing.T) {
	acronyms := []string{"ayjw", "jmwcl", "octguf"}

	for _, acronym := range acronyms {
		sentence, err := CreateSentence(acronym)
		if err != nil {
			t.Errorf("Error for valid acronym %q: %v", acronym, err)
			continue
		}

		if !unicode.IsUpper(rune(sentence[0])) ||
			rune(sentence[len(sentence)-1]) != '.' {
			t.Errorf("%q --> %q does not meet sentence requirements", acronym, sentence)
		} else if phrase := strings.ToLower(sentence[:len(sentence)-1]); !isValidPhrase(phrase, acronym) {
			t.Errorf("%q --> %q is invalid phrase", acronym, phrase)
		} else {
			t.Logf("%q --> %q", acronym, phrase)
		}
	}
}
