package cache

import (
	"io"
	"path/filepath"
	"testing"
)

func TestCalculateAddressFields(t *testing.T) {
	layout := BitLayout{OffsetBits: 4, IndexBits: 4, TagBits: 8, NumLines: 16, NumSets: 16}
	fields, err := CalculateAddressFields(addressInput{Raw: "0x12A5", Line: 1}, layout, 16)
	if err != nil {
		t.Fatalf("CalculateAddressFields retornou erro: %v", err)
	}
	if fields.Tag != 0x12 || fields.Index != 0xA || fields.Offset != 0x5 {
		t.Fatalf("campos incorretos: tag=%X index=%X offset=%X", fields.Tag, fields.Index, fields.Offset)
	}
	if fields.Binary != "0001001010100101" {
		t.Fatalf("binário incorreto: %s", fields.Binary)
	}
	if fields.TagBits != "00010010" || fields.IndexBits != "1010" || fields.OffsetBits != "0101" {
		t.Fatalf("particionamento incorreto: tag=%s index=%s offset=%s", fields.TagBits, fields.IndexBits, fields.OffsetBits)
	}
}

func TestSimulateDirectMappingCountsHitsAndMisses(t *testing.T) {
	config := Config{
		CacheSize:     256,
		BlockSize:     16,
		Associativity: 1,
		AddressBits:   16,
		Policy:        PolicyLRU,
		InputPath:     "teste",
	}
	layout := BitLayout{OffsetBits: 4, IndexBits: 4, TagBits: 8, NumLines: 16, NumSets: 16}
	addresses := []addressInput{
		{Raw: "0x0010", Line: 1},
		{Raw: "0x0014", Line: 2},
		{Raw: "0x0110", Line: 3},
		{Raw: "0x0010", Line: 4},
	}

	result, err := SimulateDirectMapping(config, layout, addresses, io.Discard)
	if err != nil {
		t.Fatalf("SimulateDirectMapping retornou erro: %v", err)
	}
	if result.TotalAccesses != 4 || result.Hits != 1 || result.Misses != 3 {
		t.Fatalf("resultado incorreto: acessos=%d hits=%d misses=%d", result.TotalAccesses, result.Hits, result.Misses)
	}
}

func TestSetAssociativeLRUAndFIFOReplacement(t *testing.T) {
	layout := BitLayout{OffsetBits: 0, IndexBits: 0, TagBits: 8, NumLines: 2, NumSets: 1}
	addresses := []addressInput{
		{Raw: "0", Line: 1},
		{Raw: "1", Line: 2},
		{Raw: "0", Line: 3},
		{Raw: "2", Line: 4},
		{Raw: "1", Line: 5},
	}
	baseConfig := Config{
		CacheSize:     2,
		BlockSize:     1,
		Associativity: 2,
		AddressBits:   8,
		InputPath:     "teste",
	}

	lruConfig := baseConfig
	lruConfig.Policy = PolicyLRU
	lruResult, err := SimulateSetAssociative(lruConfig, layout, addresses, io.Discard)
	if err != nil {
		t.Fatalf("SimulateSetAssociative LRU retornou erro: %v", err)
	}
	if lruResult.Hits != 1 || lruResult.Misses != 4 {
		t.Fatalf("resultado LRU incorreto: hits=%d misses=%d", lruResult.Hits, lruResult.Misses)
	}

	fifoConfig := baseConfig
	fifoConfig.Policy = PolicyFIFO
	fifoResult, err := SimulateSetAssociative(fifoConfig, layout, addresses, io.Discard)
	if err != nil {
		t.Fatalf("SimulateSetAssociative FIFO retornou erro: %v", err)
	}
	if fifoResult.Hits != 2 || fifoResult.Misses != 3 {
		t.Fatalf("resultado FIFO incorreto: hits=%d misses=%d", fifoResult.Hits, fifoResult.Misses)
	}
}

func TestExampleInputFilesExpectedResults(t *testing.T) {
	tests := []struct {
		name          string
		inputPath     string
		cacheSize     uint64
		blockSize     uint64
		associativity uint64
		wantOffset    uint
		wantIndex     uint
		wantTag       uint
		wantAccesses  int
		wantHits      int
		wantMisses    int
	}{
		{
			name:          "teste1-mapeamento-direto",
			inputPath:     "teste1.txt",
			cacheSize:     256,
			blockSize:     16,
			associativity: 1,
			wantOffset:    4,
			wantIndex:     4,
			wantTag:       8,
			wantAccesses:  8,
			wantHits:      1,
			wantMisses:    7,
		},
		{
			name:          "teste2-2way-lru",
			inputPath:     "teste2.txt",
			cacheSize:     1024,
			blockSize:     32,
			associativity: 2,
			wantOffset:    5,
			wantIndex:     4,
			wantTag:       7,
			wantAccesses:  9,
			wantHits:      3,
			wantMisses:    6,
		},
		{
			name:          "teste3-4way-lru",
			inputPath:     "teste3.txt",
			cacheSize:     512,
			blockSize:     8,
			associativity: 4,
			wantOffset:    3,
			wantIndex:     4,
			wantTag:       9,
			wantAccesses:  4,
			wantHits:      2,
			wantMisses:    2,
		},
		{
			name:          "teste4-conflitos-direto",
			inputPath:     "teste4.txt",
			cacheSize:     256,
			blockSize:     16,
			associativity: 1,
			wantOffset:    4,
			wantIndex:     4,
			wantTag:       8,
			wantAccesses:  20,
			wantHits:      0,
			wantMisses:    20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				CacheSize:     tt.cacheSize,
				BlockSize:     tt.blockSize,
				Associativity: tt.associativity,
				AddressBits:   16,
				Policy:        PolicyLRU,
				InputPath:     filepath.Join("..", "..", tt.inputPath),
			}

			layout, err := ValidateParameters(config)
			if err != nil {
				t.Fatalf("ValidateParameters retornou erro: %v", err)
			}
			if layout.OffsetBits != tt.wantOffset || layout.IndexBits != tt.wantIndex || layout.TagBits != tt.wantTag {
				t.Fatalf("particionamento incorreto: offset=%d index=%d tag=%d", layout.OffsetBits, layout.IndexBits, layout.TagBits)
			}

			addresses, err := ReadAddresses(config.InputPath, config.AddressBits)
			if err != nil {
				t.Fatalf("ReadAddresses retornou erro: %v", err)
			}

			var result SimulationResult
			if config.Associativity == 1 {
				result, err = SimulateDirectMapping(config, layout, addresses, io.Discard)
			} else {
				result, err = SimulateSetAssociative(config, layout, addresses, io.Discard)
			}
			if err != nil {
				t.Fatalf("simulação retornou erro: %v", err)
			}
			if result.TotalAccesses != tt.wantAccesses || result.Hits != tt.wantHits || result.Misses != tt.wantMisses {
				t.Fatalf(
					"resultado incorreto: acessos=%d hits=%d misses=%d; want acessos=%d hits=%d misses=%d",
					result.TotalAccesses,
					result.Hits,
					result.Misses,
					tt.wantAccesses,
					tt.wantHits,
					tt.wantMisses,
				)
			}
		})
	}
}
