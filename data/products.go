package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/k0kubun/pp"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"gt=0,lte=120"`
	Price       float64 `json:"price" validate:"gt=0"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func GetProducts() Products {
	return products
}

func AddProduct(product *Product) *Products {
	product.Id = getNextId()
	product.CreatedOn = time.Now().UTC().String()
	product.UpdatedOn = time.Now().UTC().String()

	products = append(products, product)

	return &products
}

func UpdateProduct(product *Product, id int) error {
	_, i, err := findProductById(id)

	if err != nil {
		return err
	}

	product.Id = id
	products[i] = product // replace, swap by new product by old i (position)

	pp.Println(products)

	return nil
}

func (p *Products) ToJson(w io.Writer) error {
	return json.NewEncoder(w).Encode(p) // direct return can return nil or error
}

func (p *Product) FromJson(r io.Reader) error {
	return json.NewDecoder(r).Decode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()

	return validate.Struct(p)
}

func getNextId() int {
	nextId := len(products)
	return nextId + 1
}

var ErrProductNotFound = fmt.Errorf("product does not exist")

func findProductById(id int) (*Product, int, error) {
	for i, product := range products {
		if product.Id == id {
			return product, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

var products Products = Products{
	&Product{
		Id:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		Id:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
