package productusecase

import (
	"context"

	"github.com/codepnw/mini-ecommerce/internal/product"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	userusecase "github.com/codepnw/mini-ecommerce/internal/user/usecase"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
)

type ProductUsecase interface {
	Create(ctx context.Context, input *product.Product) (*product.Product, error)
	GetByID(ctx context.Context, id int64) (*product.Product, error)
	GetAll(ctx context.Context) ([]*product.Product, error)
	Update(ctx context.Context, userID int64, input *product.Product) (*product.Product, error)
	Delete(ctx context.Context, userID, productID int64) error
}

type productUsecase struct {
	repo   productrepository.ProductRepository
	userUc userusecase.UserUsecase
}

func NewProductUsecase(repo productrepository.ProductRepository, userUc userusecase.UserUsecase) ProductUsecase {
	return &productUsecase{
		repo:   repo,
		userUc: userUc,
	}
}

func (u *productUsecase) Create(ctx context.Context, input *product.Product) (*product.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	if input.Stock < 0 {
		return nil, errs.ErrProductStockInvalid
	}
	if input.Price < 0 {
		return nil, errs.ErrProductPriceInvalid
	}

	// Check SKU
	exists, err := u.repo.SKUExists(ctx, input.SKU)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errs.ErrProductSKUExists
	}

	resp, err := u.repo.Insert(ctx, input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *productUsecase) GetByID(ctx context.Context, id int64) (*product.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	resp, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *productUsecase) GetAll(ctx context.Context) ([]*product.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// TODO: filter products

	resp, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *productUsecase) Update(ctx context.Context, userID int64, input *product.Product) (*product.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Check Product Owner
	userResp, err := u.userUc.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	prodResp, err := u.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if userResp.ID != prodResp.OwnerID {
		return nil, errs.ErrNoPermissions
	}

	// Update Product
	resp, err := u.repo.Update(ctx, input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *productUsecase) Delete(ctx context.Context, userID, productID int64) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// TODO: check product in order

	// Check Product Owner
	userResp, err := u.userUc.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	prodResp, err := u.repo.FindByID(ctx, productID)
	if err != nil {
		return err
	}
	if userResp.ID != prodResp.OwnerID {
		return errs.ErrNoPermissions
	}

	// Delete Product
	if err := u.repo.Delete(ctx, prodResp.ID); err != nil {
		return err
	}
	return nil
}
