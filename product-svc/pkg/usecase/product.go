package usecase_prodcut_svc

import (
	"errors"

	config_product_svc "github.com/vajid-hussain/mobile-mart-microservice-product/pkg/config"
	requestmodel_product_svc "github.com/vajid-hussain/mobile-mart-microservice-product/pkg/model/requestmodel"
	responsemodel_product_svc "github.com/vajid-hussain/mobile-mart-microservice-product/pkg/model/responsemodel"
	interfaceRepositoryProductService "github.com/vajid-hussain/mobile-mart-microservice-product/pkg/repository/interface"
	interfaceUseCase_product_svc "github.com/vajid-hussain/mobile-mart-microservice-product/pkg/usecase/interface"
	helper_product_svc "github.com/vajid-hussain/mobile-mart-microservice-product/pkg/utils"
)

type categoryUseCase struct {
	repo   interfaceRepositoryProductService.ICategoryRepository
	config *config_product_svc.Config
}

func NewCategoryUseCase(repository interfaceRepositoryProductService.ICategoryRepository, config *config_product_svc.Config) interfaceUseCase_product_svc.ICategoryUseCase {
	return &categoryUseCase{repo: repository,
		config: config}
}

func (r *categoryUseCase) NewCategory(categoryDetails *requestmodel_product_svc.Category) (*responsemodel_product_svc.CategoryDetails, error) {

	category, err := r.repo.InsertCategory(categoryDetails)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *categoryUseCase) GetAllCategory(page string, limit string) (*[]responsemodel_product_svc.CategoryDetails, error) {

	offSet, limits, err := helper_product_svc.Pagination(page, limit)
	if err != nil {
		return nil, err
	}

	categoryDetails, err := r.repo.GetAllCategory(offSet, limits)
	if err != nil {
		return nil, err
	}

	return categoryDetails, nil
}

//brand

func (r *categoryUseCase) CreateBrand(brandDetails *requestmodel_product_svc.Brand) (*responsemodel_product_svc.BrandDetails, error) {

	brand, err := r.repo.InsertBrand(brandDetails)
	if err != nil {
		return nil, err
	}

	return brand, nil
}

func (r *categoryUseCase) GetAllBrand(page string, limit string) (*[]responsemodel_product_svc.BrandDetails, error) {

	offSet, limits, err := helper_product_svc.Pagination(page, limit)
	if err != nil {
		return nil, err
	}

	brandDetails, err := r.repo.GetAllBrand(offSet, limits)
	if err != nil {
		return nil, err
	}

	return brandDetails, nil
}

func (d *categoryUseCase) AddInventory(inventory *requestmodel_product_svc.InventoryReq) (*responsemodel_product_svc.InventoryRes, error) {

	sess := helper_product_svc.CreateSession(d.config)

	ImageURL, err := helper_product_svc.UploadImageToS3(inventory.Image, sess)
	if err != nil {
		return nil, err
	}

	inventory.ImageURL = ImageURL
	discountedPrice := helper_product_svc.FindDiscount(float64(inventory.Mrp), float64(inventory.Discount))
	inventory.Saleprice = discountedPrice

	product, err := d.repo.CreateProduct(inventory)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *categoryUseCase) GetAllInventory(page string, limit string) (*[]responsemodel_product_svc.InventoryShowcase, error) {

	offSet, limits, err := helper_product_svc.Pagination(page, limit)
	if err != nil {
		return nil, err
	}

	inventories, err := r.repo.GetInventory(offSet, limits)
	if err != nil {
		return nil, err
	}

	return inventories, nil
}

func (r *categoryUseCase) CreateCart(cart *requestmodel_product_svc.Cart) (*responsemodel_product_svc.Cart, error) {

	count, err := r.repo.IsInventoryExistInCart(cart.InventoryID, cart.UserID)
	if err != nil {
		return nil, err
	}

	if count >= 1 {
		return nil, errors.New("inverntory alrady exist in cart now you can purchase")
	}

	cart.Quantity = 1

	inserCart, err := r.repo.InsertToCart(cart)
	if err != nil {
		return nil, err
	}
	return inserCart, nil
}

// func (r *cartUseCase) ShowCart(userID string) (*responsemodel.UserCart, error) {

// 	cart := &responsemodel.UserCart{}

// 	quantity, err := r.repo.GetCartCriteria(userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cart.InventoryCount = quantity
// 	cart.UserID = userID

// 	cartInventories, err := r.repo.GetCart(userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for i, inventory := range *cartInventories {

// 		// price, err := r.repo.GetNetAmoutOfCart(userID, inventory.InventoryID)
// 		// if err != nil {
// 		// 	return nil, err
// 		// }
// 		(*cartInventories)[i].FinalPrice = helper.FindDiscount(float64(inventory.Price), float64(inventory.Discount+inventory.CategoryDiscount)) * inventory.Quantity
// 		fmt.Println("**", (*cartInventories)[i].FinalPrice)
// 		cart.TotalPrice += (*cartInventories)[i].FinalPrice
// 	}

// 	cart.Cart = *cartInventories
// 	return cart, nil
// }
