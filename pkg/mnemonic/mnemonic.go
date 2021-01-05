package mnemonic

import (
	"errors"
	"math/rand"
	"strings"
	"unicode"
)

var validPartsOfSpeech = map[rune]bool{'N': true, 'V': true, 'A': true}

var patterns = map[int][]string{
	2: {"AN"},
	3: {"AAN", "VAN", "NVN"},
	4: {"NVAN", "ANVN"},
	5: {"ANVAN"},
	6: {"AANVAN", "ANVAAN"},
}

var vocabulary = map[rune][]string{
	'N': {"airplanes", "books", "cheeses", "dequeues", "eigenvalues", "frogs", "golfers", "hashmaps", "igloos", "jellyfish", "keyboards", "laptops", "maps", "noodles", "owls", "poodles", "queues", "rodents", "stamps", "towels", "unicorns", "vowels", "waffles", "xylocarps", "yaks", "zealots"},
	'A': {"adjoint", "bored", "consistent", "diagonalizable", "energetic", "fancy", "gorgeous", "hermitian", "inverted", "jumpy", "keynesian", "lame", "marxist", "nonsingular", "orthogonal", "perpendicular", "quirky", "romantic", "slow", "triangular", "uniform", "vainglorious", "well-ordered", "xylotomous", "young-at-heart", "zany"},
	'V': {"ate", "borrow", "create", "destroy", "eat", "forgot", "grew", "help", "ignore", "juggle", "knit", "like", "mention", "need", "outlast", "peruse", "quash", "reuse", "surprise", "torment", "utilize", "vaccinate", "want", "xerox", "yeet", "zap"},
}

func CreatePhrase(acronym string) (string, error) {
	possiblePatterns, ok := patterns[len(acronym)]
	if !ok {
		return "", errors.New("No patterns of correct length")
	}

	pattern := possiblePatterns[rand.Intn(len(possiblePatterns))]

	var sb strings.Builder
	for i, letter := range strings.ToLower(acronym) {
		letterIndex := letter - 'a'
		if letterIndex < 0 || letterIndex >= 26 {
			return "", errors.New("Input contains nonletter character(s)")
		}

		partOfSpeech := rune(pattern[i])
		sb.WriteString(vocabulary[partOfSpeech][letterIndex])

		if i != len(acronym)-1 {
			sb.WriteString(" ")
		}
	}

	return sb.String(), nil
}

func CreateSentence(acronym string) (string, error) {
	phrase, err := CreatePhrase(acronym)
	if err != nil {
		return "", err
	}

	return string(unicode.ToUpper(rune(phrase[0]))) + phrase[1:] + ".", nil
}
