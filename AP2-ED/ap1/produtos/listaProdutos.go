package produtos

import (
	m "mcronalds/metricas"
	"strings"
)

type NodoProduto struct {
	Produto Produto
	Proximo *NodoProduto
}

type ListaProdutos struct {
	Cabeca *NodoProduto
}

const maxProdutos = 50

var totalProdutos = 0

var ListaProdutosEncadeada ListaProdutos // Variável global para a lista encadeada

func tentarCriar(nome, descricao string, preco float64, id int) *NodoProduto {
	if id != -1 {
		_, idProcurado := BuscarId(id)
		if idProcurado != -1 {
			return nil
		}
	}

	return &NodoProduto{Produto: criar(nome, descricao, preco, id), Proximo: nil}
}

/*
Adiciona um produto com nome, descrição e preço à lista de produtos.
Adiciona o produto ao final da lista.
Caso já exista um produto com o mesmo id, não adiciona e retorna -3.
Caso já exista um produto com o mesmo nome, não adiciona e retorna erro -2.
Retorna -1 caso a lista esteja cheia, ou o número de produtos cadastrados em caso de sucesso.
*/
func AdicionarUnico(nome, descricao string, preco float64, id int) int {
	if totalProdutos == maxProdutos {
		return -1 // Overflow
	}

	novoNodo := tentarCriar(nome, descricao, preco, id)
	if novoNodo == nil {
		return -3
	}

	// Adiciona o novo nó ao final da lista
	if ListaProdutosEncadeada.Cabeca == nil {
		ListaProdutosEncadeada.Cabeca = novoNodo
	} else {
		ultimoNodo := obterUltimoNodo(ListaProdutosEncadeada.Cabeca)
		ultimoNodo.Proximo = novoNodo
	}

	totalProdutos++
	m.M.SomaProdutosCadastrados(1)
	return totalProdutos
}

// Função auxiliar para obter o último nodo na lista encadeada
func obterUltimoNodo(cabeca *NodoProduto) *NodoProduto {
	atual := cabeca
	for atual.Proximo != nil {
		atual = atual.Proximo
	}
	return atual
}

/*
Localiza um produto a partir do seu id.
Retorna o produto encontrado e a sua posição na lista, em caso de sucesso.
Retorna um produto vazio e -1 em caso de erro.
*/
func BuscarId(id int) (Produto, int) {
	atual := ListaProdutosEncadeada.Cabeca
	indice := 0

	for atual != nil {
		if atual.Produto.Id == id {
			return atual.Produto, indice
		}

		atual = atual.Proximo
		indice++
	}

	return Produto{}, -1
}

/*
Localiza produtos que iniciem com a string passada.
Retorna um slice com todos os produtos encontrados, e o tamanho do slice.
*/
func BuscarNome(comecaCom string) ([]Produto, int) {
	var produtosEncontrados []Produto

	atual := ListaProdutosEncadeada.Cabeca

	for atual != nil {
		if strings.HasPrefix(atual.Produto.Nome, comecaCom) {
			produtosEncontrados = append(produtosEncontrados, atual.Produto)
		}
		atual = atual.Proximo
	}

	return produtosEncontrados, len(produtosEncontrados)
}

/*
Exibe todos os produtos cadastrados.
*/
func Exibir() {
	atual := ListaProdutosEncadeada.Cabeca

	for atual != nil {
		atual.Produto.Exibir()
		atual = atual.Proximo
	}
}

/*
Remove um produto da lista a partir do seu id.
Retorna -2 caso não haja produtos na lista.
Retorna -1 caso não haja um produto com o id passado, ou 0 em caso de sucesso.
*/
func Excluir(id int) int {
	if totalProdutos == 0 {
		return -2
	}

	if ListaProdutosEncadeada.Cabeca == nil {
		return -1
	}

	if ListaProdutosEncadeada.Cabeca.Produto.Id == id {
		// Remove o primeiro nó se for o procurado
		ListaProdutosEncadeada.Cabeca = ListaProdutosEncadeada.Cabeca.Proximo
		totalProdutos--
		m.M.SomaProdutosCadastrados(-1)
		return 0
	}

	anterior := ListaProdutosEncadeada.Cabeca
	atual := anterior.Proximo

	for atual != nil {
		if atual.Produto.Id == id {
			// Remove o nó encontrado
			anterior.Proximo = atual.Proximo
			totalProdutos--
			m.M.SomaProdutosCadastrados(-1)
			return 0
		}
		anterior = atual
		atual = atual.Proximo
	}

	return -1
}

func AtualizarPrecoProduto(id int, novoPreco float64) bool {
	atual := ListaProdutosEncadeada.Cabeca

	for atual != nil {
		if atual.Produto.Id == id {
			atual.Produto.Preco = novoPreco
			return true
		}
		atual = atual.Proximo
	}

	return false
}

func ExibirOrdenadoPorNome() {
	// Converte a lista encadeada para um slice para facilitar a ordenação
	produtos := listaParaSlice(ListaProdutosEncadeada.Cabeca)

	// Ordena os produtos por nome
	ordenarProdutosPorNome(produtos)

	for _, produto := range produtos {
		produto.Exibir()
	}
}

// Função auxiliar para converter uma lista encadeada em um slice
func listaParaSlice(cabeca *NodoProduto) []Produto {
	var produtos []Produto

	atual := cabeca
	for atual != nil {
		produtos = append(produtos, atual.Produto)
		atual = atual.Proximo
	}

	return produtos
}

// Função auxiliar para ordenar os produtos por nome
func ordenarProdutosPorNome(produtos []Produto) {
	for i := 0; i < len(produtos)-1; i++ {
		for j := i + 1; j < len(produtos); j++ {
			if produtos[i].Nome > produtos[j].Nome {
				// Troca os produtos de posição se estiverem fora de ordem
				produtos[i], produtos[j] = produtos[j], produtos[i]
			}
		}
	}
}
