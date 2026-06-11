package cache

import (
	"strings"
	"testing"
)

func TestReadAddressesSupportsHexDecimalCommentsAndBlankLines(t *testing.T) {
	path := createTempAddressFile(t, `
# comentário
0x0010

256
0x00FF # comentário ao lado
`)

	addresses, err := ReadAddresses(path, 16)
	if err != nil {
		t.Fatalf("ReadAddresses retornou erro: %v", err)
	}
	if len(addresses) != 3 {
		t.Fatalf("quantidade de endereços incorreta: got %d want 3", len(addresses))
	}
	want := []string{"0x0010", "256", "0x00FF"}
	for i := range want {
		if addresses[i].Raw != want[i] {
			t.Fatalf("endereço %d incorreto: got %q want %q", i, addresses[i].Raw, want[i])
		}
	}
}

func TestReadAddressesRejectsInvalidAndOutOfRangeAddresses(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		bits       uint
		wantErrMsg string
	}{
		{name: "texto-invalido", content: "xyz\n", bits: 16, wantErrMsg: "endereço inválido na linha 1"},
		{name: "hexadecimal-sem-prefixo", content: "FF\n", bits: 16, wantErrMsg: "endereço inválido na linha 1"},
		{name: "hexadecimal-com-zero-sem-prefixo", content: "00FF\n", bits: 16, wantErrMsg: "endereço inválido na linha 1"},
		{name: "decimal-comeca-com-letra", content: "abc123\n", bits: 16, wantErrMsg: "endereço inválido na linha 1"},
		{name: "decimal-termina-com-letra", content: "123abc\n", bits: 16, wantErrMsg: "endereço inválido na linha 1"},
		{name: "numero-negativo", content: "-1\n", bits: 16, wantErrMsg: "endereço inválido na linha 1"},
		{name: "hexadecimal-sem-digitos", content: "0x\n", bits: 16, wantErrMsg: "endereço inválido na linha 1"},
		{name: "prefixo-binario-nao-suportado", content: "0b1010\n", bits: 16, wantErrMsg: "endereço inválido na linha 1"},
		{name: "fora-do-espaco", content: "0x10000\n", bits: 16, wantErrMsg: "endereço fora do espaço permitido na linha 1"},
		{name: "vazio", content: "# apenas comentario\n\n", bits: 16, wantErrMsg: "arquivo de entrada vazio"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadAddresses(createTempAddressFile(t, tt.content), tt.bits)
			if err == nil {
				t.Fatalf("ReadAddresses deveria retornar erro")
			}
			if !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Fatalf("erro incorreto: got %q want conter %q", err.Error(), tt.wantErrMsg)
			}
		})
	}
}

func TestReadAddressesReturnsReadError(t *testing.T) {
	longLine := strings.Repeat("1", 70_000)

	_, err := ReadAddresses(createTempAddressFile(t, longLine), 16)
	if err == nil {
		t.Fatalf("ReadAddresses deveria retornar erro de leitura")
	}
	if !strings.Contains(err.Error(), "erro ao ler o arquivo") {
		t.Fatalf("erro incorreto: got %q want conter %q", err.Error(), "erro ao ler o arquivo")
	}
}
