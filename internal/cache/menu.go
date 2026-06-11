package cache

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// RunInteractiveMenu apresenta o menu interativo exigido e executa as simulações.
func RunInteractiveMenu(in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	config := Config{Policy: PolicyLRU}

	for {
		printMenu(out, config)
		choice, err := readLine(reader, out, "Escolha uma opção: ")
		if err != nil {
			return err
		}

		switch choice {
		case "1":
			value, err := readUint(reader, out, "Tamanho da cache (bytes): ")
			if err != nil {
				fmt.Fprintf(out, "Erro: %v\n", err)
				continue
			}
			config.CacheSize = value
		case "2":
			value, err := readUint(reader, out, "Tamanho do bloco (bytes): ")
			if err != nil {
				fmt.Fprintf(out, "Erro: %v\n", err)
				continue
			}
			config.BlockSize = value
		case "3":
			value, err := readUint(reader, out, "Grau de associatividade: ")
			if err != nil {
				fmt.Fprintf(out, "Erro: %v\n", err)
				continue
			}
			config.Associativity = value
		case "4":
			value, err := readUint(reader, out, "Quantidade de bits do endereço: ")
			if err != nil {
				fmt.Fprintf(out, "Erro: %v\n", err)
				continue
			}
			config.AddressBits = uint(value)
		case "5":
			policy, err := readPolicy(reader, out)
			if err != nil {
				fmt.Fprintf(out, "Erro: %v\n", err)
				continue
			}
			config.Policy = policy
		case "6":
			path, err := readLine(reader, out, "Caminho do arquivo de entrada: ")
			if err != nil {
				return err
			}
			config.InputPath = path
			verbose, err := readYesNo(reader, out, "Ativar modo verbose? (s/n): ")
			if err != nil {
				fmt.Fprintf(out, "Erro: %v\n", err)
				continue
			}
			config.Verbose = verbose
		case "7":
			if err := executeSimulation(config, out); err != nil {
				fmt.Fprintf(out, "Erro: %v\n", err)
			}
		case "0":
			fmt.Fprintln(out, "Encerrando.")
			return nil
		default:
			fmt.Fprintln(out, "Opção inválida.")
		}
	}
}

func printMenu(out io.Writer, config Config) {
	verbose := "Não"
	if config.Verbose {
		verbose = "Sim"
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "=================================")
	fmt.Fprintln(out, " SIMULADOR DE MEMÓRIA CACHE")
	fmt.Fprintln(out, "=================================")
	fmt.Fprintf(out, "Cache: %d bytes | Bloco: %d bytes | Associatividade: %d\n", config.CacheSize, config.BlockSize, config.Associativity)
	fmt.Fprintf(out, "Bits endereço: %d | Política: %s | Verbose: %s\n", config.AddressBits, config.Policy, verbose)
	fmt.Fprintf(out, "Arquivo: %s\n", valueOrDash(config.InputPath))
	if layout, err := DeriveBitLayout(config); err == nil {
		fmt.Fprintf(out, "Linhas: %d | Conjuntos: %d\n", layout.NumLines, layout.NumSets)
		fmt.Fprintf(out, "Campos: TAG=%d bits | Index=%d bits | OFFSET=%d bits\n", layout.TagBits, layout.IndexBits, layout.OffsetBits)
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "1 - Informar tamanho da cache")
	fmt.Fprintln(out, "2 - Informar tamanho do bloco")
	fmt.Fprintln(out, "3 - Informar associatividade")
	fmt.Fprintln(out, "4 - Informar bits do endereço")
	fmt.Fprintln(out, "5 - Selecionar política")
	fmt.Fprintln(out, "6 - Selecionar arquivo")
	fmt.Fprintln(out, "7 - Executar simulação")
	fmt.Fprintln(out, "0 - Sair")
	fmt.Fprintln(out)
}

func executeSimulation(config Config, out io.Writer) error {
	layout, err := ValidateParameters(config)
	if err != nil {
		return err
	}
	addresses, err := ReadAddresses(config.InputPath, config.AddressBits)
	if err != nil {
		return err
	}

	PrintBitLayout(out, layout)

	var result SimulationResult
	if config.Associativity == 1 {
		result, err = SimulateDirectMapping(config, layout, addresses, out)
	} else {
		result, err = SimulateSetAssociative(config, layout, addresses, out)
	}
	if err != nil {
		return err
	}

	GenerateReport(out, result)
	return nil
}

func readPolicy(reader *bufio.Reader, out io.Writer) (ReplacementPolicy, error) {
	fmt.Fprintln(out, "1 - LRU")
	fmt.Fprintln(out, "2 - FIFO")
	value, err := readLine(reader, out, "Selecione a política: ")
	if err != nil {
		return "", err
	}
	switch strings.ToLower(value) {
	case "1", "lru":
		return PolicyLRU, nil
	case "2", "fifo":
		return PolicyFIFO, nil
	default:
		return "", fmt.Errorf("política inválida")
	}
}

func readUint(reader *bufio.Reader, out io.Writer, prompt string) (uint64, error) {
	text, err := readLine(reader, out, prompt)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseUint(text, 10, 64)
	if err != nil || value == 0 {
		return 0, fmt.Errorf("informe um número inteiro positivo")
	}
	return value, nil
}

func readYesNo(reader *bufio.Reader, out io.Writer, prompt string) (bool, error) {
	text, err := readLine(reader, out, prompt)
	if err != nil {
		return false, err
	}
	switch strings.ToLower(text) {
	case "s", "sim", "y", "yes":
		return true, nil
	case "n", "nao", "não", "no":
		return false, nil
	default:
		return false, fmt.Errorf("responda com Sim ou Não")
	}
}

func readLine(reader *bufio.Reader, out io.Writer, prompt string) (string, error) {
	fmt.Fprint(out, prompt)
	text, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	if err == io.EOF && text == "" {
		return "", io.EOF
	}
	return strings.TrimSpace(text), nil
}

func valueOrDash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
