package person

import (
	"context"
	"fmt"
)

// Service 人物领域服务
type Service struct {
	repo Repository
}

// NewService 创建人物领域服务
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetPerson 获取人物详情（含关联信息）
func (s *Service) GetPerson(ctx context.Context, id int64, withRelations bool) (*Person, error) {
	person, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if withRelations {
		// 加载配偶
		spouses, err := s.repo.FindSpouses(ctx, id)
		if err == nil {
			person.Spouses = spouses
		}

		// 加载子女
		children, err := s.repo.FindChildren(ctx, id)
		if err == nil {
			person.Children = children
		}

		// 加载父亲
		if person.FatherID != nil {
			father, err := s.repo.FindByID(ctx, *person.FatherID)
			if err == nil {
				person.Father = father
			}
		}
	}

	return person, nil
}

// CreatePerson 创建新人物
func (s *Service) CreatePerson(ctx context.Context, person *Person) error {
	// 验证领域规则
	if err := person.Validate(); err != nil {
		return err
	}

	return s.repo.Create(ctx, person)
}

// UpdatePerson 更新人物信息
func (s *Service) UpdatePerson(ctx context.Context, person *Person) error {
	if err := person.Validate(); err != nil {
		return err
	}

	return s.repo.Update(ctx, person)
}

// DeletePerson 删除人物
func (s *Service) DeletePerson(ctx context.Context, id int64) error {
	// 检查是否有子女
	children, err := s.repo.FindChildren(ctx, id)
	if err == nil && len(children) > 0 {
		return fmt.Errorf("cannot delete person with %d children", len(children))
	}

	return s.repo.Delete(ctx, id)
}

// SearchPersons 搜索人物
func (s *Service) SearchPersons(ctx context.Context, query *SearchQuery) ([]*Person, int64, error) {
	// 设置默认分页参数
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	return s.repo.Search(ctx, query)
}

// GetFamilyTree 获取上下五代族谱树
func (s *Service) GetFamilyTree(ctx context.Context, personID int64, upDepth, downDepth int) (*FamilyTreeResult, error) {
	// 获取目标人物
	person, err := s.repo.FindByID(ctx, personID)
	if err != nil {
		return nil, err
	}

	// 向上获取祖先
	ancestors, err := s.repo.FindAncestors(ctx, personID, upDepth)
	if err != nil {
		return nil, err
	}

	// 向下获取后代
	descendants, err := s.repo.FindDescendants(ctx, personID, downDepth)
	if err != nil {
		return nil, err
	}

	// 构建结果
	result := &FamilyTreeResult{
		CenterPerson: person,
		Ancestors:    ancestors,
		Descendants:  descendants,
		Generations:  upDepth + downDepth + 1,
		TotalNodes:   len(ancestors) + len(descendants) + 1,
	}

	return result, nil
}

// AddParentChildRelation 添加父子关系
func (s *Service) AddParentChildRelation(ctx context.Context, parentID, childID int64, relationType RelationType, birthOrderNum int, isPrimary bool) error {
	// 检查是否会形成循环
	if err := s.checkCycle(ctx, parentID, childID); err != nil {
		return err
	}

	// TODO: 实现关系创建逻辑
	return nil
}

// checkCycle 检查是否会形成循环引用
func (s *Service) checkCycle(ctx context.Context, parentID, childID int64) error {
	if parentID == childID {
		return fmt.Errorf("self reference not allowed")
	}

	// 检查child的后代中是否有parent
	descendants, err := s.repo.FindDescendants(ctx, childID, 30)
	if err != nil {
		return err
	}

	for _, d := range descendants {
		if d.ID == parentID {
			return fmt.Errorf("cycle detected: person %d is already a descendant of %d", parentID, childID)
		}
	}

	return nil
}

// GetStatistics 获取统计数据
func (s *Service) GetStatistics(ctx context.Context) (*Statistics, error) {
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	byGeneration, err := s.repo.CountByGeneration(ctx)
	if err != nil {
		return nil, err
	}

	return &Statistics{
		TotalPersons:  total,
		ByGeneration:  byGeneration,
	}, nil
}

// FamilyTreeResult 族谱树结果
type FamilyTreeResult struct {
	CenterPerson *Person   `json:"center_person"`
	Ancestors    []*Person `json:"ancestors"`
	Descendants  []*Person `json:"descendants"`
	Generations  int       `json:"generations"`
	TotalNodes   int       `json:"total_nodes"`
}

// Statistics 统计数据
type Statistics struct {
	TotalPersons int64         `json:"total_persons"`
	ByGeneration map[int]int64 `json:"by_generation"`
}
