package cache

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestPrintMenuShowsDerivedCacheFields(t *testing.T) {
	var out strings.Builder
	config := Config{
		CacheSize:     256,
		BlockSize:     16,
		Associativity: 1,
		AddressBits:   16,
		Policy:        PolicyLRU,
		InputPath:     "teste1.txt",
	}

	printMenu(&out, config)
	text := out.String()
	wantFragments := []string{
		"Linhas: 16 | Conjuntos: 16",
		"Campos: TAG=8 bits | Index=4 bits | OFFSET=4 bits",
	}
	for _, fragment := range wantFragments {
		if !strings.Contains(text, fragment) {
			t.Fatalf("menu não contém %q\nmenu:\n%s", fragment, text)
		}
	}
}

func TestPrintMenuShowsVerboseEnabled(t *testing.T) {
	var out strings.Builder
	config := Config{
		CacheSize:     256,
		BlockSize:     16,
		Associativity: 1,
		AddressBits:   16,
		Policy:        PolicyLRU,
		InputPath:     "teste1.txt",
		Verbose:       true,
	}

	printMenu(&out, config)
	if !strings.Contains(out.String(), "Verbose: Sim") {
		t.Fatalf("menu não indica verbose ativado:\n%s", out.String())
	}
}

func TestReadUintRejectsInvalidMenuNumbers(t *testing.T) {
	tests := []string{
		"0\n",
		"-1\n",
		"abc\n",
		"12abc\n",
		"0x10\n",
		"\n",
	}

	for _, input := range tests {
		t.Run(strings.TrimSpace(input), func(t *testing.T) {
			_, err := readUint(bufio.NewReader(strings.NewReader(input)), io.Discard, "valor: ")
			if err == nil {
				t.Fatalf("readUint deveria rejeitar %q", input)
			}
			if !strings.Contains(err.Error(), "número inteiro positivo") {
				t.Fatalf("erro incorreto: got %q", err.Error())
			}
		})
	}
}

func TestReadPolicyRejectsInvalidInput(t *testing.T) {
	tests := []string{
		"3\n",
		"random\n",
		"\n",
	}

	for _, input := range tests {
		t.Run(strings.TrimSpace(input), func(t *testing.T) {
			_, err := readPolicy(bufio.NewReader(strings.NewReader(input)), io.Discard)
			if err == nil {
				t.Fatalf("readPolicy deveria rejeitar %q", input)
			}
			if !strings.Contains(err.Error(), "política inválida") {
				t.Fatalf("erro incorreto: got %q", err.Error())
			}
		})
	}
}

func TestReadYesNoRejectsInvalidInput(t *testing.T) {
	tests := []string{
		"talvez\n",
		"1\n",
		"\n",
	}

	for _, input := range tests {
		t.Run(strings.TrimSpace(input), func(t *testing.T) {
			_, err := readYesNo(bufio.NewReader(strings.NewReader(input)), io.Discard, "verbose: ")
			if err == nil {
				t.Fatalf("readYesNo deveria rejeitar %q", input)
			}
			if !strings.Contains(err.Error(), "Sim ou Não") {
				t.Fatalf("erro incorreto: got %q", err.Error())
			}
		})
	}
}

func TestReadYesNoAcceptsValidInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantVerb bool
	}{
		{name: "sim-curto", input: "s\n", wantVerb: true},
		{name: "sim-completo", input: "sim\n", wantVerb: true},
		{name: "yes", input: "yes\n", wantVerb: true},
		{name: "nao", input: "nao\n", wantVerb: false},
		{name: "nao-com-acento", input: "não\n", wantVerb: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readYesNo(bufio.NewReader(strings.NewReader(tt.input)), io.Discard, "verbose: ")
			if err != nil {
				t.Fatalf("readYesNo retornou erro: %v", err)
			}
			if got != tt.wantVerb {
				t.Fatalf("readYesNo retornou %v, want %v", got, tt.wantVerb)
			}
		})
	}
}
