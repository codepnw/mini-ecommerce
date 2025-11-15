package productusecase

import (
	"context"

	"github.com/codepnw/mini-ecommerce/internal/product"
	productrepository "github.com/codepnw/mini-ecommerce/internal/product/repository"
	"github.com/codepnw/mini-ecommerce/internal/user"
	"github.com/codepnw/mini-ecommerce/internal/utils/consts"
	"github.com/codepnw/mini-ecommerce/internal/utils/errs"
	"github.com/codepnw/mini-ecommerce/pkg/auth"
)

type ProductUsecase interface {
	Create(ctx context.Context, input *product.Product) (*product.Product, error)
	GetByID(ctx context.Context, id int64) (*product.Product, error)
	GetAll(ctx context.Context) ([]*product.Product, error)
	Update(ctx context.Context, input *product.Product) (*product.Product, error)
	Delete(ctx context.Context, productID int64) error
}

type productUsecase struct {
	repo productrepository.ProductRepository
}

func NewProductUsecase(repo productrepository.ProductRepository) ProductUsecase {
	return &productUsecase{
		repo: repo,
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

func (u *productUsecase) Update(ctx context.Context, input *product.Product) (*product.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Check Admin & Product Owner
	if err := u.checkPermissions(ctx, input.ID); err != nil {
		return nil, err
	}

	// Update Product
	resp, err := u.repo.Update(ctx, input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *productUsecase) Delete(ctx context.Context, productID int64) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// TODO: check product in order

	// Check Admin & Product Owner
	if err := u.checkPermissions(ctx, productID); err != nil {
		return err
	}

	// Delete Product
	if err := u.repo.Delete(ctx, productID); err != nil {
		return err
	}
	return nil
}

func (u *productUsecase) checkPermissions(ctx context.Context, productID int64) error {
	currentUser, err := auth.GetCurrentUser(ctx)
	if err != nil {
		return err
	}
	productData, err := u.repo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	// Check Admin & Product Owner
	if currentUser.Role != string(user.RoleAdmin) {
		if currentUser.ID != productData.OwnerID {
			return errs.ErrNoPermissions
		}
	}
	return nil
}
