package cache

import (
	"fmt"
	"os"
)

// DeriveBitLayout valida os parâmetros numéricos e calcula os campos da cache.
func DeriveBitLayout(config Config) (BitLayout, error) {
	if config.CacheSize == 0 {
		return BitLayout{}, fmt.Errorf("o tamanho da cache deve ser maior que zero")
	}
	if !isPowerOfTwo(config.CacheSize) {
		return BitLayout{}, fmt.Errorf("o tamanho da cache deve ser potência de 2")
	}
	if config.BlockSize == 0 {
		return BitLayout{}, fmt.Errorf("o tamanho do bloco deve ser maior que zero")
	}
	if !isPowerOfTwo(config.BlockSize) {
		return BitLayout{}, fmt.Errorf("o tamanho do bloco deve ser potência de 2")
	}
	if config.BlockSize > config.CacheSize {
		return BitLayout{}, fmt.Errorf("o bloco não pode ser maior que a cache")
	}
	if config.Associativity == 0 {
		return BitLayout{}, fmt.Errorf("a associatividade deve ser maior que zero")
	}
	if !isPowerOfTwo(config.Associativity) {
		return BitLayout{}, fmt.Errorf("a associatividade deve ser potência de 2")
	}
	if config.AddressBits == 0 || config.AddressBits > 64 {
		return BitLayout{}, fmt.Errorf("a quantidade de bits do endereço deve estar entre 1 e 64")
	}
	if config.Policy != PolicyLRU && config.Policy != PolicyFIFO {
		return BitLayout{}, fmt.Errorf("política de substituição inválida: use LRU ou FIFO")
	}

	numLines := config.CacheSize / config.BlockSize
	if config.Associativity > numLines {
		return BitLayout{}, fmt.Errorf("a associatividade não pode ser maior que o número de linhas da cache (%d)", numLines)
	}
	if numLines%config.Associativity != 0 {
		return BitLayout{}, fmt.Errorf("configuração inconsistente: linhas (%d) não divisíveis pela associatividade (%d)", numLines, config.Associativity)
	}

	numSets := numLines / config.Associativity
	if !isPowerOfTwo(numSets) {
		return BitLayout{}, fmt.Errorf("número de conjuntos deve ser potência de 2")
	}

	layout := BitLayout{
		OffsetBits: log2PowerOfTwo(config.BlockSize),
		IndexBits:  log2PowerOfTwo(numSets),
		NumLines:   numLines,
		NumSets:    numSets,
	}
	if config.AddressBits < layout.OffsetBits+layout.IndexBits {
		return BitLayout{}, fmt.Errorf("bits de endereço insuficientes para index e offset")
	}
	layout.TagBits = config.AddressBits - layout.IndexBits - layout.OffsetBits

	return layout, nil
}

// verifica a consistência da configuração antes da simulação.
func ValidateParameters(config Config) (BitLayout, error) {
	layout, err := DeriveBitLayout(config)
	if err != nil {
		return BitLayout{}, err
	}
	if config.InputPath == "" {
		return BitLayout{}, fmt.Errorf("informe o caminho do arquivo de entrada")
	}
	info, err := os.Stat(config.InputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return BitLayout{}, fmt.Errorf("arquivo de entrada inexistente: %s", config.InputPath)
		}
		return BitLayout{}, fmt.Errorf("não foi possível acessar o arquivo: %w", err)
	}
	if info.IsDir() {
		return BitLayout{}, fmt.Errorf("o caminho informado é um diretório, não um arquivo")
	}

	return layout, nil
}
