package produtos

import (
	m "mcronalds/metricas"
	"strings"
)

const maxProdutos = 50

var Produtos [maxProdutos]Produto
var totalProdutos = 0

func tentarCriar(nome, descricao string, preco float64, id int) Produto {
	if id != -1 {
		_, idProcurado := BuscarId(id)
		if idProcurado != -1 {
			return Produto{}
		}
	}

	return criar(nome, descricao, preco, id)
}

/*
Adiciona um produto com nome, descrição e preço à lista de produtos.
Adiciona o produto primeiro espaço vazio da lista.
Caso já exista um produto com o mesmo id, não adiciona e retorna -3.
Caso já exista um produto com o mesmo nome, não adiciona e retorna erro -2.
Retorna -1 caso a lista esteja cheia, ou o número de produtos cadastrados em caso de sucesso.
*/
func AdicionarUnico(nome, descricao string, preco float64, id int) int {
	if totalProdutos == maxProdutos {
		return -1 // Overflow
	}

	for i, produto := range Produtos {
		if (produto == Produto{}) {
			break
		}

		if produto.Id == id {
			// Atualiza o preço do produto existente
			AtualizarPrecoProduto(id, preco)
			return i + 1 // Retorna o número de produtos cadastrados
		}
	}

	// Se o produto não existe, adiciona um novo
	produtoCriado := tentarCriar(nome, descricao, preco, id)
	if (produtoCriado == Produto{}) {
		return -3
	}

	Produtos[totalProdutos] = produtoCriado
	totalProdutos++
	m.M.SomaProdutosCadastrados(1)
	return totalProdutos
}

/*
Localiza um produto a partir do seu id.
Retorna o produto encontrado e a sua posição na lista, em caso de sucesso.
Retorna um produto vazio e -1 em caso de erro.
*/
func BuscarId(id int) (Produto, int) {
	for ind, produto := range Produtos {
		if (produto == Produto{}) {
			break
		}
		if produto.Id == id {
			return produto, ind
		}
	}

	return Produto{}, -1
}

/*
Localiza produtos que iniciem com a string passada.
Retorna um slice com todos os produtos encontrados, e o tamanho do slice.
*/
func BuscarNome(comecaCom string) ([]Produto, int) {
	var produtosEncontrados []Produto

	for _, produto := range Produtos {
		if (produto == Produto{}) {
			break
		}

		if strings.HasPrefix(produto.Nome, comecaCom) {
			produtosEncontrados = append(produtosEncontrados, produto)
		}
	}
	return produtosEncontrados, len(produtosEncontrados)
}

/*
Exibe todos os produtos cadastrados.
*/
func Exibir() {
	for _, produto := range Produtos {
		if (produto == Produto{}) {
			break
		}
		produto.Exibir()
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

	_, ind := BuscarId(id)
	if ind == -1 {
		return -1
	}

	for i := ind; i < totalProdutos-1; i++ {
		Produtos[i] = Produtos[i+1]
	}
	totalProdutos--
	Produtos[totalProdutos] = Produto{}
	m.M.SomaProdutosCadastrados(-1)
	return 0
}

func AtualizarPrecoProduto(id int, novoPreco float64) bool {
	for i, produto := range Produtos {
		if (produto == Produto{}) {
			break
		}

		if produto.Id == id {
			Produtos[i].Preco = novoPreco
			return true
		}
	}

	return false
}