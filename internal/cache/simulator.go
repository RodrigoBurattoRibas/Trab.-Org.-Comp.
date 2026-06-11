package cache

import (
	"fmt"
	"io"
)

func newCache(config Config, layout BitLayout) Cache {
	sets := make([]CacheSet, layout.NumSets)
	for i := range sets {
		sets[i].Lines = make([]CacheLine, config.Associativity)
	}
	return Cache{Sets: sets, Config: config, Layout: layout}
}

// CalculateAddressFields decompõe um endereço físico em tag, index e offset.
func CalculateAddressFields(input addressInput, layout BitLayout, addrBits uint) (AddressFields, error) {
	value, err := parseAddressValue(input.Raw)
	if err != nil {
		return AddressFields{}, fmt.Errorf("endereço inválido na linha %d: %q", input.Line, input.Raw)
	}

	offsetMask := uint64(0)
	if layout.OffsetBits > 0 {
		offsetMask = (uint64(1) << layout.OffsetBits) - 1
	}
	indexMask := uint64(0)
	if layout.IndexBits > 0 {
		indexMask = (uint64(1) << layout.IndexBits) - 1
	}

	binary := addressBinary(value, addrBits)
	tagBits, indexBits, offsetBits := partitionAddressBinary(binary, layout)

	fields := AddressFields{
		Original:   input.Raw,
		Value:      value,
		Binary:     binary,
		TagBits:    tagBits,
		IndexBits:  indexBits,
		OffsetBits: offsetBits,
		Offset:     value & offsetMask,
		Index:      (value >> layout.OffsetBits) & indexMask,
		Tag:        value >> (layout.OffsetBits + layout.IndexBits),
	}
	return fields, nil
}

func partitionAddressBinary(binary string, layout BitLayout) (string, string, string) {
	tagEnd := int(layout.TagBits)
	indexEnd := tagEnd + int(layout.IndexBits)

	return binary[:tagEnd], binary[tagEnd:indexEnd], binary[indexEnd:]
}

// SimulateDirectMapping executa a simulação de cache com mapeamento direto.
func SimulateDirectMapping(config Config, layout BitLayout, addresses []addressInput, out io.Writer) (SimulationResult, error) {
	return simulate(config, layout, addresses, out)
}

// SimulateSetAssociative executa a simulação de cache associativa por conjunto N-way.
func SimulateSetAssociative(config Config, layout BitLayout, addresses []addressInput, out io.Writer) (SimulationResult, error) {
	return simulate(config, layout, addresses, out)
}

func simulate(config Config, layout BitLayout, addresses []addressInput, out io.Writer) (SimulationResult, error) {
	cacheMemory := newCache(config, layout)
	result := SimulationResult{Config: config, Layout: layout}

	for _, input := range addresses {
		fields, err := CalculateAddressFields(input, layout, config.AddressBits)
		if err != nil {
			return result, err
		}
		hit := accessCache(&cacheMemory, fields)
		result.TotalAccesses++
		if hit {
			result.Hits++
		} else {
			result.Misses++
		}

		if config.Verbose {
			printVerboseAccess(out, fields, hit, cacheMemory)
		}
	}

	return result, nil
}

func accessCache(cacheMemory *Cache, fields AddressFields) bool {
	cacheMemory.Clock++
	set := &cacheMemory.Sets[fields.Index]

	for i := range set.Lines {
		line := &set.Lines[i]
		if line.Valid && line.Tag == fields.Tag {
			if cacheMemory.Config.Policy == PolicyLRU {
				line.LastUsed = cacheMemory.Clock
			}
			return true
		}
	}

	replacementIndex := firstInvalidLine(set)
	if replacementIndex == -1 {
		if cacheMemory.Config.Associativity == 1 {
			replacementIndex = 0
		} else if cacheMemory.Config.Policy == PolicyFIFO {
			replacementIndex = applyFIFO(set)
		} else {
			replacementIndex = applyLRU(set)
		}
	}

	set.Lines[replacementIndex] = CacheLine{
		Valid:    true,
		Tag:      fields.Tag,
		LastUsed: cacheMemory.Clock,
		LoadedAt: cacheMemory.Clock,
	}
	if cacheMemory.Config.Policy == PolicyFIFO && replacementIndex == set.NextFIFO {
		set.NextFIFO = (set.NextFIFO + 1) % len(set.Lines)
	}

	return false
}

func firstInvalidLine(set *CacheSet) int {
	for i, line := range set.Lines {
		if !line.Valid {
			return i
		}
	}
	return -1
}

// applyLRU seleciona a linha usada menos recentemente dentro de um conjunto cheio.
func applyLRU(set *CacheSet) int {
	selected := 0
	for i := 1; i < len(set.Lines); i++ {
		if set.Lines[i].LastUsed < set.Lines[selected].LastUsed {
			selected = i
		}
	}
	return selected
}

// applyFIFO seleciona a linha que entrou primeiro dentro de um conjunto cheio.
func applyFIFO(set *CacheSet) int {
	return set.NextFIFO
}
