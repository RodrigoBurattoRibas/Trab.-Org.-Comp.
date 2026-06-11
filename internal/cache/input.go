package cache

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// carrega endereços decimais ou hexadecimais e ignora comentários/linhas vazias.
func ReadAddresses(path string, addrBits uint) ([]addressInput, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("não foi possível abrir o arquivo: %w", err)
	}
	defer file.Close()

	var addresses []addressInput
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	maxValue := maxAddressValue(addrBits)

	for scanner.Scan() {
		lineNumber++
		text := scanner.Text()
		if commentIndex := strings.Index(text, "#"); commentIndex >= 0 {
			text = text[:commentIndex]
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		value, err := parseAddressValue(text)
		if err != nil {
			return nil, fmt.Errorf("endereço inválido na linha %d: %q", lineNumber, text)
		}
		if value > maxValue {
			return nil, fmt.Errorf("endereço fora do espaço permitido na linha %d: %q excede %d bits", lineNumber, text, addrBits)
		}
		addresses = append(addresses, addressInput{Raw: text, Line: lineNumber})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo: %w", err)
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("arquivo de entrada vazio ou sem endereços válidos")
	}

	return addresses, nil
}

func parseAddressValue(text string) (uint64, error) {
	if strings.HasPrefix(text, "0x") || strings.HasPrefix(text, "0X") {
		if len(text) == 2 {
			return 0, fmt.Errorf("hexadecimal vazio")
		}
		return strconv.ParseUint(text[2:], 16, 64)
	}

	for _, char := range text {
		if char < '0' || char > '9' {
			return 0, fmt.Errorf("decimal inválido")
		}
	}
	return strconv.ParseUint(text, 10, 64)
}
