package dependency

// Dependency injection - provides repository instances with their dependencies

import (
	"github.com/alielmi98/go-ecommerce-api/config"
	model "github.com/alielmi98/go-ecommerce-api/domain/models"
	contractRepository "github.com/alielmi98/go-ecommerce-api/domain/repository"
	"github.com/alielmi98/go-ecommerce-api/infra/db"
	infraRepository "github.com/alielmi98/go-ecommerce-api/infra/db/repository"
)

func GetUserRepository(cfg *config.Config) contractRepository.UserRepository {
	return infraRepository.NewUserRepository()
}

func GetCategoryRepository(cfg *config.Config) contractRepository.CategoryRepository {
	var preloads []db.PreloadEntity = []db.PreloadEntity{{Entity: "Products"}}
	return infraRepository.NewBaseRepository[model.Category](preloads)
}

func GetProductRepository(cfg *config.Config) contractRepository.ProductRepository {
	var preloads []db.PreloadEntity = []db.PreloadEntity{{Entity: "Category"}}
	return infraRepository.NewProductRepository(preloads)
}
