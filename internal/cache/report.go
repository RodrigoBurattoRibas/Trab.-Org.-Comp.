package cache

import (
	"fmt"
	"io"
)

func PrintBitLayout(out io.Writer, layout BitLayout) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Particionamento do endereço:")
	fmt.Fprintf(out, "OFFSET(bits): %d\n", layout.OffsetBits)
	fmt.Fprintf(out, "Index(bits): %d\n", layout.IndexBits)
	fmt.Fprintf(out, "TAG(bits): %d\n", layout.TagBits)
	fmt.Fprintln(out)
}

func printVerboseAccess(out io.Writer, fields AddressFields, hit bool, cacheMemory Cache) {
	status := "Miss"
	if hit {
		status = "Hit"
	}

	fmt.Fprintln(out, "------------------------")
	fmt.Fprintf(out, "Endereço original: %s\n", fields.Original)
	fmt.Fprintf(out, "Endereço binário:  %s\n", fields.Binary)
	fmt.Fprintf(out, "Particionamento: TAG=%s | Index=%s | OFFSET=%s\n", printableBits(fields.TagBits), printableBits(fields.IndexBits), printableBits(fields.OffsetBits))
	fmt.Fprintf(out, "Tag: %d | Index: %d | Offset: %d\n", fields.Tag, fields.Index, fields.Offset)
	fmt.Fprintf(out, "Resultado: %s\n", status)
	PrintCacheState(out, cacheMemory)
}

func printableBits(bits string) string {
	if bits == "" {
		return "-"
	}
	return bits
}

// PrintCacheState imprime todos os conjuntos e linhas armazenados no momento da simulação.
func PrintCacheState(out io.Writer, cacheMemory Cache) {
	fmt.Fprintln(out, "Estado atual da cache:")
	for setIndex, set := range cacheMemory.Sets {
		fmt.Fprintf(out, "Conjunto %d:\n", setIndex)
		for lineIndex, line := range set.Lines {
			if line.Valid {
				fmt.Fprintf(out, "  Linha %d: V=1 Tag=%d LastUsed=%d LoadedAt=%d\n", lineIndex, line.Tag, line.LastUsed, line.LoadedAt)
			} else {
				fmt.Fprintf(out, "  Linha %d: V=0 Tag=-\n", lineIndex)
			}
		}
	}
}

// GenerateReport imprime a configuração final e os contadores de desempenho.
func GenerateReport(out io.Writer, result SimulationResult) {
	hitRate := 0.0
	missRate := 0.0
	if result.TotalAccesses > 0 {
		hitRate = float64(result.Hits) * 100 / float64(result.TotalAccesses)
		missRate = float64(result.Misses) * 100 / float64(result.TotalAccesses)
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "========================")
	fmt.Fprintln(out, "CONFIGURAÇÃO")
	fmt.Fprintln(out, "========================")
	fmt.Fprintf(out, "Cache Size: %d bytes\n", result.Config.CacheSize)
	fmt.Fprintf(out, "Block Size: %d bytes\n", result.Config.BlockSize)
	fmt.Fprintf(out, "Associatividade: %d\n", result.Config.Associativity)
	fmt.Fprintf(out, "Bits endereço: %d\n", result.Config.AddressBits)
	fmt.Fprintf(out, "Política: %s\n", result.Config.Policy)
	fmt.Fprintf(out, "Arquivo: %s\n", result.Config.InputPath)
	fmt.Fprintf(out, "Linhas da cache: %d\n", result.Layout.NumLines)
	fmt.Fprintf(out, "Conjuntos: %d\n", result.Layout.NumSets)
	fmt.Fprintf(out, "OFFSET(bits): %d\n", result.Layout.OffsetBits)
	fmt.Fprintf(out, "Index(bits): %d\n", result.Layout.IndexBits)
	fmt.Fprintf(out, "TAG(bits): %d\n", result.Layout.TagBits)

	fmt.Fprintln(out)
	fmt.Fprintln(out, "========================")
	fmt.Fprintln(out, "RESULTADOS")
	fmt.Fprintln(out, "========================")
	fmt.Fprintf(out, "Total de Acessos: %d\n", result.TotalAccesses)
	fmt.Fprintf(out, "Hits: %d\n", result.Hits)
	fmt.Fprintf(out, "Misses: %d\n", result.Misses)
	fmt.Fprintf(out, "Taxa de Hit (%%): %.2f\n", hitRate)
	fmt.Fprintf(out, "Taxa de Miss (%%): %.2f\n", missRate)
}
