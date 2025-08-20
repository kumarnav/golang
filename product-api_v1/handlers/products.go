package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"product-api/data"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// ServeHTTP is the main entry point for the handler and staisfies the http.Handler
// interface
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle the request for a list of products
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		p.getProducts(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		//expect the id in the URI
		p.l.Println("PUT", r.URL.Path)
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URL, more than 1 id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URL, more than 1 capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI unable to convert to number", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.l.Println("got id", id)
		p.updateProduct(id, rw, r)
		p.getProducts(rw, r)
	}

	// catch all
	// if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getProducts returns the products from the data store
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	// fetch the products from the datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// addProduct return the new product added
func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Prodcuts")
	//var NewProduct data.Product
	NewProduct := &data.Product{}
	err := NewProduct.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}
	data.AddProduct(NewProduct)
	rw.WriteHeader(http.StatusCreated)
	//rw.Write([]byte(NewProduct))
	//fmt.Fprintf(rw, "User added: %+v\n", http.StatusAccepted)
}

// addProduct return the new product added
func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Prodcuts")
	//var NewProduct data.Product
	NewProduct := &data.Product{}
	err := NewProduct.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = data.UpdateProduct(id, NewProduct)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
