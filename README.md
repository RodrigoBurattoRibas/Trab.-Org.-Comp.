# Alunos

Rodrigo Buratto Ribas

Erick Marlon Mafra

Daniel Uesler Brito

# Simulador de Memória Cache

Projeto em Go que simula o comportamento de uma memória cache com mapeamento direto e mapeamento por conjunto, conforme a prática M3 de Organização de Computadores.

O programa lê endereços de um arquivo de entrada, aceita valores decimais e hexadecimais com prefixo `0x`, ignora comentários iniciados por `#`, decompõe cada endereço em tag, index e offset, executa a simulação e exibe estatísticas finais de desempenho.

## Passo a Passo no Pop!_OS

O Pop!_OS é baseado em Ubuntu/Debian, então a instalação pode ser feita com `apt`.

1. Abra o Terminal.
2. Atualize a lista de pacotes:

```bash
sudo apt update
```

3. Instale o Go:

```bash
sudo apt install -y golang-go
```

4. Verifique se a instalação funcionou:

```bash
go version
```

5. Entre na pasta do projeto:

```bash
cd "/caminho/para/Trab. Org. Comp."
```

Exemplo, se o projeto estiver em Documentos:

```bash
cd "$HOME/Documentos/Trab. Org. Comp."
```

6. Rode os testes unitários:

```bash
go test ./...
```

7. Execute o simulador sem compilar:

```bash
go run ./cmd/cache-sim
```

8. Ou compile e execute o binário:

```bash
go build -o cache-sim ./cmd/cache-sim
./cache-sim
```

Se o Go reclamar de cache ou de VCS em ambiente restrito, use:

```bash
GOCACHE=/tmp/go-build-cache go test ./...
GOCACHE=/tmp/go-build-cache go build -buildvcs=false -o cache-sim ./cmd/cache-sim
```

## Passo a Passo Universal

1. Instale o Go usando o método adequado para seu sistema operacional.
2. Abra o terminal do sistema: PowerShell/CMD no Windows, Terminal no Linux/macOS.
3. Entre na pasta do projeto.
4. Verifique a instalação com `go version`.
5. Rode os testes com `go test ./...`.
6. Execute com `go run ./cmd/cache-sim`.
7. Se preferir, compile e rode o executável gerado.

Windows:

```powershell
cd "C:\caminho\para\Trab. Org. Comp."
go version
go test ./...
go run ./cmd/cache-sim
```

Compilar e executar no Windows:

```powershell
go build -o cache-sim.exe ./cmd/cache-sim
.\cache-sim.exe
```

Linux/macOS:

```bash
cd "/caminho/para/Trab. Org. Comp."
go version
go test ./...
go run ./cmd/cache-sim
```

Compilar e executar no Linux/macOS:

```bash
go build -o cache-sim ./cmd/cache-sim
./cache-sim
```

## Instalação do Go

Escolha uma das opções abaixo conforme o sistema operacional.

### Windows

Com Winget:

```powershell
winget install GoLang.Go
```

Com Chocolatey:

```powershell
choco install golang -y
```

Com Scoop:

```powershell
scoop install go
```

Também é possível baixar o instalador `.msi` oficial em `https://go.dev/dl/`.

### Ubuntu, Debian, Pop!_OS e Derivados

```bash
sudo apt update
sudo apt install -y golang-go
```

### Fedora

```bash
sudo dnf install -y golang
```

### RHEL, CentOS, Rocky Linux e AlmaLinux

```bash
sudo dnf install -y golang
```

Em versões antigas que usam `yum`:

```bash
sudo yum install -y golang
```

### openSUSE

```bash
sudo zypper install go
```

### Arch Linux e Manjaro

```bash
sudo pacman -S go
```

### macOS

Com Homebrew:

```bash
brew install go
```

Com MacPorts:

```bash
sudo port install go
```

### Instalação Manual Oficial

Também é possível instalar pelo pacote oficial em `https://go.dev/dl/`.

No Windows, baixe o instalador `.msi` e siga o assistente.

No Linux, baixe o arquivo `.tar.gz`, extraia e adicione o diretório `bin` ao `PATH`. Exemplo:

```bash
tar -C "$HOME" -xzf go1.xx.x.linux-amd64.tar.gz
export PATH="$HOME/go/bin:$PATH"
go version
```

## Estrutura

```text
.
├── cmd/cache-sim/main.go          # Ponto de entrada da aplicação
├── internal/cache/menu.go         # Menu interativo do terminal
├── internal/cache/input.go        # Leitura e validação dos endereços
├── internal/cache/validation.go   # Validação da configuração
├── internal/cache/simulator.go    # Lógica da simulação
├── internal/cache/report.go       # Verbose e relatório final
├── internal/cache/types.go        # Structs principais
├── internal/cache/math.go         # Funções auxiliares
├── teste1.txt                     # Arquivo de teste
├── teste2.txt                     # Arquivo de teste
└── go.mod
```

## Como Compilar

Linux/macOS:

```bash
go build -o cache-sim ./cmd/cache-sim
```

Windows:

```powershell
go build -o cache-sim.exe ./cmd/cache-sim
```

Em ambientes restritos, fora de um repositório Git ou com pasta home somente leitura, use:

Linux/macOS:

```bash
GOCACHE=/tmp/go-build-cache go build -buildvcs=false -o cache-sim ./cmd/cache-sim
```

Windows PowerShell:

```powershell
$env:GOCACHE="$env:TEMP\go-build-cache"
go build -buildvcs=false -o cache-sim.exe ./cmd/cache-sim
```

## Como Executar

Executar sem compilar previamente:

```bash
go run ./cmd/cache-sim
```

Em ambientes onde a pasta home esteja somente leitura:

```bash
GOCACHE=/tmp/go-build-cache go run -buildvcs=false ./cmd/cache-sim
```

Ou, após compilar:

Linux/macOS:

```bash
./cache-sim
```

Windows:

```powershell
.\cache-sim.exe
```

O programa não exige parâmetros pela linha de comando. Todas as configurações são solicitadas pelo menu interativo.

## Como Rodar os Testes

Execute:

```bash
go test ./...
```

Em ambiente restrito:

```bash
GOCACHE=/tmp/go-build-cache go test ./...
```

Os testes unitários cobrem:

- cálculo de tag, index e offset para todos os exemplos do PDF;
- leitura de endereços decimais e hexadecimais;
- comentários e linhas vazias no arquivo de entrada;
- rejeição de endereço inválido, arquivo vazio e endereço fora do espaço permitido;
- validações de configuração;
- simulação com mapeamento direto;
- simulação associativa por conjunto com FIFO e LRU.

## Exemplo de Uso

No menu, informe:

```text
1 - Tamanho da cache: 256
2 - Tamanho do bloco: 16
3 - Associatividade: 1
4 - Bits do endereço: 16
5 - Política: LRU ou FIFO
6 - Arquivo: teste1.txt
7 - Executar simulação
```

Para essa configuração, os campos calculados são:

```text
Offset Bits: 4
Index Bits: 4
Tag Bits: 8
```

O modo verbose pode ser ativado ao selecionar o arquivo. Quando ativo, cada acesso mostra endereço original, endereço binário, tag, index, offset, hit ou miss e o estado atual da cache.

## Arquivos de Entrada

Alguns arquivos de teste acompanham o projeto, por exemplo:

- `teste1.txt`: sequência simples para mapeamento direto.
- `teste2.txt`: sequência para cache associativa por conjunto.
- `teste4.txt`: sequência com muitos conflitos para observar substituições na mesma linha/conjunto.
- `acessos_exemplo.txt`: arquivo simples de exemplo para usar no menu.
- `entrada_menu_exemplo.txt`: respostas prontas para testar o menu automaticamente.

Formato aceito:

```text
# Comentários são ignorados
0x0010
0x0024
256
0x00FF # comentário ao lado também é ignorado
```

Cada linha deve conter um endereço decimal ou hexadecimal com prefixo `0x`.



O arquivo `entrada_menu_exemplo.txt` não deve ter comentários, porque cada linha é lida como se o usuário tivesse digitado uma opção no menu.

Para testar rapidamente usando os arquivos de exemplo:

```bash
go run ./cmd/cache-sim < entrada_menu_exemplo.txt
```

Ou, após compilar:

```bash
./cache-sim < entrada_menu_exemplo.txt
``` 

## Políticas de Substituição

FIFO (First In, First Out): quando um conjunto está cheio, remove a linha que entrou primeiro naquele conjunto. A ordem não muda quando ocorre hit.

LRU (Least Recently Used): quando um conjunto está cheio, remove a linha que está há mais tempo sem ser acessada. Em caso de hit, a linha acessada passa a ser considerada a mais recente.

## Validações

O simulador trata erros como cache ou bloco que não sejam potência de 2, associatividade inválida, arquivo inexistente, arquivo vazio, endereço inválido, endereço fora do espaço de endereçamento e configurações inconsistentes.

## Configurações de Validação do PDF

O simulador calcula automaticamente:

```text
offset_bits = log2(block_size)
num_linhas = cache_size / block_size
num_conjuntos = num_linhas / associatividade
index_bits = log2(num_conjuntos)
tag_bits = addr_bits - index_bits - offset_bits
```

Resultados esperados:

| Cache | Bloco | Associatividade | Endereço | Offset | Index | Tag |
| ---:  | ---:  | ---:            | ---:     | ---:   | ---:  | ---:|
| 256   | 16    | 1               | 16 bits  | 4      | 4     | 8   |
| 1024  | 32    | 2-way           | 16 bits  | 5      | 4     | 7   |
| 512   | 8     | 4-way           | 16 bits  | 3      | 4     | 9   |
| 2048  | 64    | 8-way           | 32 bits  | 6      | 2     | 24  |


### Exemplos com `teste4.txt`

O arquivo `teste4.txt` contem uma sequencia com muitos conflitos. Use a mesma configuracao do exemplo do PDF e responda `n` ou `s` ao selecionar o arquivo.

Linux/macOS:

```bash
printf '1\n256\n2\n16\n3\n1\n4\n16\n5\n1\n6\nteste4.txt\nn\n7\n0\n' | GOCACHE=/tmp/go-build-cache go run -buildvcs=false ./cmd/cache-sim
```

Saida resumida:

```text
Particionamento do endereco:
OFFSET(bits): 4
Index(bits): 4
TAG(bits): 8
...
Total de Acessos: 20
Hits: 0
Misses: 20
Taxa de Hit (%): 0.00
Taxa de Miss (%): 100.00
```

Em modo verbose:

```bash
printf '1\n256\n2\n16\n3\n1\n4\n16\n5\n1\n6\nteste4.txt\ns\n7\n0\n' | GOCACHE=/tmp/go-build-cache go run -buildvcs=false ./cmd/cache-sim
```

Trecho:

```text
Particionamento do endereco:
OFFSET(bits): 4
Index(bits): 4
TAG(bits): 8
------------------------
Endereco original: 0x0000
...
```
