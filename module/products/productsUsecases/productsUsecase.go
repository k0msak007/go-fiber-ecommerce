package productsUsecases

import (
	"math"

	"github.com/k0msak007/go-fiber-ecommerce/module/entities"
	"github.com/k0msak007/go-fiber-ecommerce/module/products"
	"github.com/k0msak007/go-fiber-ecommerce/module/products/productsRepositories"
)

type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProduct(req *products.ProductFilter) *entities.PaginateRes
	AddProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
	UpdateProduct(req *products.Product) (*products.Product, error)
}

type productsUsecase struct {
	productsRepository productsRepositories.IProductsRepository
}

func ProductsUsecase(productsRepository productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productsRepository: productsRepository,
	}
}

func (u *productsUsecase) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.productsRepository.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productsUsecase) FindProduct(req *products.ProductFilter) *entities.PaginateRes {
	products, count := u.productsRepository.FindProduct(req)

	return &entities.PaginateRes{
		Data:      products,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *productsUsecase) AddProduct(req *products.Product) (*products.Product, error) {
	return u.productsRepository.InsertProduct(req)
}

func (u *productsUsecase) DeleteProduct(productId string) error {
	return u.productsRepository.DeleteProduct(productId)
}

func (u *productsUsecase) UpdateProduct(req *products.Product) (*products.Product, error) {
	return u.productsRepository.UpdateProduct(req)
}
