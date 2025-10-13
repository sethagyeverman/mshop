package biz

import (
	"context"
	"encoding/json"
	"log"
	"mshop/pkg/errx"
	pb "mshop/service/goods/api/goods/v1"
)

// GetAllCategorysList 获取所有一级及子分类
func (uc *GoodsUsecase) GetAllCategorysList(ctx context.Context, req *pb.Empty) (resp *pb.CategoryListResponse, err error) {
	var category []*Category

	// 查询所有一级分类及其二级、三级子分类
	result := uc.db.Preload("SubCategories.SubCategories").Where("level = ?", 1).Find(&category)
	if result.Error != nil {
		log.Printf("[GetAllCategorysList] database error: %v", result.Error)
		return nil, errx.ErrorDatabaseError("db error: %v", result.Error)
	}
	v, marshalErr := json.Marshal(category)
	if marshalErr != nil {
		log.Printf("[GetAllCategorysList] marshal error: %v", marshalErr)
		return nil, errx.ErrorDatabaseError("marshal error: %v", marshalErr)
	}
	return &pb.CategoryListResponse{
		JsonData: string(v),
	}, nil
}

// GetSubCategory 查询某分类的直接子分类，需健壮地处理异常
func (uc *GoodsUsecase) GetSubCategory(ctx context.Context, req *pb.CategoryListRequest) (resp *pb.SubCategoryListResponse, err error) {
	var category Category

	// 查询带子分类
	result := uc.db.Preload("SubCategories").First(&category, req.Id)
	if result.Error != nil {
		log.Printf("[GetSubCategory] database error: %v", result.Error)
		return nil, errx.ErrorDatabaseError("db error: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("[GetSubCategory] category not found, id=%v", req.Id)
		return nil, errx.ErrorCategoryNotFound("category not found")
	}
	if category.SubCategories == nil {
		category.SubCategories = []*Category{}
	}
	resp = &pb.SubCategoryListResponse{
		Total: int32(len(category.SubCategories)),
	}
	resp.Info = &pb.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		ParentCategory: category.ParentCategoryID,
		IsTab:          category.IsTab,
	}
	resp.SubCategorys = make([]*pb.CategoryInfoResponse, 0, len(category.SubCategories))
	for _, v := range category.SubCategories {
		if v == nil {
			continue
		}
		resp.SubCategorys = append(resp.SubCategorys, &pb.CategoryInfoResponse{
			Id:             v.ID,
			Name:           v.Name,
			Level:          v.Level,
			ParentCategory: v.ParentCategoryID,
			IsTab:          v.IsTab,
		})
	}
	return resp, nil
}

// CreateCategory 创建分类，有详细的输入校验和日志
func (uc *GoodsUsecase) CreateCategory(ctx context.Context, req *pb.CategoryInfoRequest) (resp *pb.CategoryInfoResponse, err error) {
	// 校验必要参数
	if req.Name == "" {
		log.Printf("[CreateCategory] category name invalid, input empty")
		return nil, errx.ErrorCategoryNameEmpty("category name invalid")
	}
	// 分类名唯一性校验
	if result := uc.db.Where("name = ?", req.Name).First(&Category{}); result.Error != nil {
		log.Printf("[CreateCategory] db error in name check: %v", result.Error)
		return nil, errx.ErrorDatabaseError("db error: %v", result.Error)
	} else if result.RowsAffected != 0 {
		log.Printf("[CreateCategory] category name exists: %s", req.Name)
		return nil, errx.ErrorCategoryNameExists("category already exists")
	}
	// 父类校验，仅在有父类时检查
	if req.ParentCategory != 0 {
		if result := uc.db.First(&Category{}, req.ParentCategory); result.Error != nil {
			log.Printf("[CreateCategory] db error in parent check: %v", result.Error)
			return nil, errx.ErrorDatabaseError("db error: %v", result.Error)
		} else if result.RowsAffected == 0 {
			log.Printf("[CreateCategory] parent category not found, id=%v", req.ParentCategory)
			return nil, errx.ErrorCategoryParentInvalid("category parent category not exists")
		}
	}
	if req.Level <= 0 || req.Level > 3 {
		log.Printf("[CreateCategory] category level invalid: %d", req.Level)
		return nil, errx.ErrorCategoryLevelInvalid("category level invalid")
	}

	// 创建分类
	category := &Category{
		Name:             req.Name,
		ParentCategoryID: req.ParentCategory,
		Level:            req.Level,
		IsTab:            req.IsTab,
	}
	createResult := uc.db.Create(category)
	if createResult.Error != nil {
		log.Printf("[CreateCategory] create db error: %v", createResult.Error)
		return nil, errx.ErrorDatabaseError("db error: %v", createResult.Error)
	}
	resp = &pb.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		ParentCategory: category.ParentCategoryID,
		Level:          category.Level,
		IsTab:          category.IsTab,
	}
	return resp, nil
}

// DeleteCategory 递归删除分类及所有下级，加强健壮性
func (uc *GoodsUsecase) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (_ *pb.Empty, err error) {
	var deleteF func(*Category)

	deleteF = func(category *Category) {
		if category == nil {
			return
		}
		// 重新获取最新的子分类
		uc.db.Preload("SubCategories").Find(category)
		if len(category.SubCategories) != 0 {
			for _, v := range category.SubCategories {
				deleteF(v)
			}
		}
		// 日志记录删除动作
		log.Printf("[DeleteCategory] deleting category, id=%v", category.ID)
		uc.db.Delete(category)
	}

	var category Category
	result := uc.db.Preload("SubCategories").First(&category, req.Id)
	if result.Error != nil {
		log.Printf("[DeleteCategory] database error: %v", result.Error)
		return nil, errx.ErrorDatabaseError("db error: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("[DeleteCategory] category not found, id=%v", req.Id)
		return nil, errx.ErrorCategoryNotFound("category not found")
	}

	deleteF(&category)
	return &pb.Empty{}, nil
}

// UpdateCategory 更新分类
func (uc *GoodsUsecase) UpdateCategory(ctx context.Context, req *pb.CategoryInfoRequest) (_ *pb.Empty, err error) {
	// 1. 检查请求参数
	if req.Name == "" {
		log.Printf("[UpdateCategory] category name invalid: empty")
		return nil, errx.ErrorCategoryNameEmpty("category name invalid")
	}

	// 2. 检查分类是否存在
	var category Category
	if result := uc.db.First(&category, req.Id); result.Error != nil {
		log.Printf("[UpdateCategory] database error on category find: %v", result.Error)
		return nil, errx.ErrorDatabaseError("db error: %v", result.Error)
	} else if result.RowsAffected == 0 {
		log.Printf("[UpdateCategory] category not found, id=%v", req.Id)
		return nil, errx.ErrorCategoryNotFound("category not found")
	}

	// 3. 检查新分类名冲突
	if result := uc.db.Where("name = ? AND id != ?", req.Name, req.Id).First(&Category{}); result.Error != nil {
		log.Printf("[UpdateCategory] db error on name check: %v", result.Error)
		return nil, errx.ErrorDatabaseError("db error: %v", result.Error)
	} else if result.RowsAffected != 0 {
		log.Printf("[UpdateCategory] category name exists: %s", req.Name)
		return nil, errx.ErrorCategoryNameExists("category already exists")
	}

	// 4. 检查父类合法性（非顶级分类时）
	if req.ParentCategory != 0 && req.ParentCategory != category.ParentCategoryID {
		var parentCategory Category
		if result := uc.db.First(&parentCategory, req.ParentCategory); result.Error != nil {
			log.Printf("[UpdateCategory] db error on parent check: %v", result.Error)
			return nil, errx.ErrorDatabaseError("db error: %v", result.Error)
		} else if result.RowsAffected == 0 {
			log.Printf("[UpdateCategory] parent category not found, id=%v", req.ParentCategory)
			return nil, errx.ErrorCategoryParentInvalid("category parent category not exists")
		}
		category.ParentCategoryID = req.ParentCategory
		category.Level = parentCategory.Level + 1
	} else if req.ParentCategory == 0 {
		category.ParentCategoryID = 0
		category.Level = 1
	}

	// 5. 检查层级范围
	if req.Level <= 0 || req.Level > 3 {
		log.Printf("[UpdateCategory] category level invalid: %d", req.Level)
		return nil, errx.ErrorCategoryLevelInvalid("category level invalid")
	}

	category.Name = req.Name
	category.IsTab = req.IsTab

	// 6. 更新数据
	if result := uc.db.Updates(&category); result.Error != nil {
		log.Printf("[UpdateCategory] update db error: %v", result.Error)
		return nil, errx.ErrorDatabaseError("db error: %v", result.Error)
	}

	return &pb.Empty{}, nil
}
