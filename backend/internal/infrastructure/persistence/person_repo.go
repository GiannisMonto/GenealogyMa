package persistence

import (
	"context"
	"strings"
	"time"

	"github.com/genealogy-ma/platform/internal/domain/person"
	"gorm.io/gorm"
)

// MemberModel GORM模型对应members表
type MemberModel struct {
	MemberID         int64      `gorm:"column:member_id;primaryKey;autoIncrement"`
	LegacyID         string     `gorm:"column:legacy_id;type:char(9);unique;not null"`
	Name             string     `gorm:"column:name;type:varchar(200)"`
	StyleName        string     `gorm:"column:style_name;type:varchar(50)"`
	Gender           string     `gorm:"column:gender;type:gender;not null"`
	Generation       *int16     `gorm:"column:generation;type:smallint"`
	BirthOrder       *string    `gorm:"column:birth_order;type:birth_order"`
	FatherID         *int64     `gorm:"column:father_id"`
	LineagePath      string     `gorm:"column:lineage_path;type:ltree"`
	DetailText       string     `gorm:"column:detail_text;type:text"`
	SonCount         int16      `gorm:"column:son_count;default:0"`
	DaughterCount    int16      `gorm:"column:daughter_count;default:0"`
	AdoptedHeirCount int16      `gorm:"column:adopted_heir_count;default:0"`
	TotalChildren    int16      `gorm:"column:total_children_count;default:0"`
	BirthTimeText    string     `gorm:"column:birth_time_text;type:text"`
	DeathTimeText    string     `gorm:"column:death_time_text;type:text"`
	BirthPlace       string     `gorm:"column:birth_place;type:text"`
	BurialPlace      string     `gorm:"column:burial_place;type:text"`
	BirthGregorian   *time.Time `gorm:"column:birth_gregorian;type:date"`
	BirthYear        *int       `gorm:"column:birth_year"`
	DeathGregorian   *time.Time `gorm:"column:death_gregorian;type:date"`
	DeathYear        *int       `gorm:"column:death_year"`
	PageNumber       *int       `gorm:"column:page_number"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (MemberModel) TableName() string {
	return "members"
}

// SpouseModel GORM模型对应spouses表
type SpouseModel struct {
	SpouseID       int64      `gorm:"column:spouse_id;primaryKey;autoIncrement"`
	MemberID       int64      `gorm:"column:member_id;not null"`
	SpouseType     string     `gorm:"column:spouse_type;type:spouse_type;default:配"`
	Name           string     `gorm:"column:name;type:varchar(200);not null"`
	BirthTimeText  string     `gorm:"column:birth_time_text;type:text"`
	DeathTimeText  string     `gorm:"column:death_time_text;type:text"`
	BirthPlace     string     `gorm:"column:birth_place;type:text"`
	BurialPlace    string     `gorm:"column:burial_place;type:text"`
	BirthGregorian *time.Time `gorm:"column:birth_gregorian;type:date"`
	BirthYear      *int       `gorm:"column:birth_year"`
	DeathGregorian *time.Time `gorm:"column:death_gregorian;type:date"`
	DeathYear      *int       `gorm:"column:death_year"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (SpouseModel) TableName() string {
	return "spouses"
}

// ParentChildRelationModel GORM模型对应parent_child_relations表
type ParentChildRelationModel struct {
	RelationID     int64      `gorm:"column:relation_id;primaryKey;autoIncrement"`
	ParentID       int64      `gorm:"column:parent_id;not null"`
	ChildID        int64      `gorm:"column:child_id;not null"`
	RelationType   string     `gorm:"column:relation_type;type:relation_type;default:biological"`
	BirthOrderNum  *int16     `gorm:"column:birth_order_num"`
	BirthOrderText *string    `gorm:"column:birth_order_text;type:birth_order"`
	IsPrimary      bool       `gorm:"column:is_primary;default:true"`
	Note           string     `gorm:"column:note;type:text"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (ParentChildRelationModel) TableName() string {
	return "parent_child_relations"
}

// PersonRepositoryImpl 人物仓储实现
type PersonRepositoryImpl struct {
	db *gorm.DB
}

// NewPersonRepository 创建人物仓储
func NewPersonRepository(db *gorm.DB) person.Repository {
	return &PersonRepositoryImpl{db: db}
}

// 模型转换函数
func modelToPerson(m *MemberModel) *person.Person {
	var gen int
	if m.Generation != nil {
		gen = int(*m.Generation)
	}

	p := &person.Person{
		ID:               m.MemberID,
		LegacyID:         m.LegacyID,
		Name:             m.Name,
		StyleName:        m.StyleName,
		Gender:           person.Gender(m.Gender),
		Generation:       gen,
		LineagePath:      m.LineagePath,
		DetailText:       m.DetailText,
		SonCount:         int(m.SonCount),
		DaughterCount:    int(m.DaughterCount),
		AdoptedHeirCount: int(m.AdoptedHeirCount),
		TotalChildren:    int(m.TotalChildren),
		BirthTimeText:    m.BirthTimeText,
		DeathTimeText:    m.DeathTimeText,
		BirthPlace:       m.BirthPlace,
		BurialPlace:      m.BurialPlace,
		BirthGregorian:   m.BirthGregorian,
		BirthYear:        m.BirthYear,
		DeathGregorian:   m.DeathGregorian,
		DeathYear:        m.DeathYear,
		PageNumber:       m.PageNumber,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		FatherID:         m.FatherID,
	}

	if m.BirthOrder != nil {
		p.BirthOrder = *m.BirthOrder
	}

	return p
}

func personToModel(p *person.Person) *MemberModel {
	gen := int16(p.Generation)
	return &MemberModel{
		MemberID:         p.ID,
		LegacyID:         p.LegacyID,
		Name:             p.Name,
		StyleName:        p.StyleName,
		Gender:           string(p.Gender),
		Generation:       &gen,
		LineagePath:      p.LineagePath,
		DetailText:       p.DetailText,
		SonCount:         int16(p.SonCount),
		DaughterCount:    int16(p.DaughterCount),
		AdoptedHeirCount: int16(p.AdoptedHeirCount),
		TotalChildren:    int16(p.TotalChildren),
		BirthTimeText:    p.BirthTimeText,
		DeathTimeText:    p.DeathTimeText,
		BirthPlace:       p.BirthPlace,
		BurialPlace:      p.BurialPlace,
		BirthGregorian:   p.BirthGregorian,
		BirthYear:        p.BirthYear,
		DeathGregorian:   p.DeathGregorian,
		DeathYear:        p.DeathYear,
		PageNumber:       p.PageNumber,
	}
}

func modelToSpouse(m *SpouseModel) *person.Spouse {
	return &person.Spouse{
		ID:             m.SpouseID,
		MemberID:       m.MemberID,
		SpouseType:     person.SpouseType(m.SpouseType),
		Name:           m.Name,
		BirthTimeText:  m.BirthTimeText,
		DeathTimeText:  m.DeathTimeText,
		BirthPlace:     m.BirthPlace,
		BurialPlace:    m.BurialPlace,
		BirthGregorian: m.BirthGregorian,
		DeathGregorian: m.DeathGregorian,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// 仓储接口实现

func (r *PersonRepositoryImpl) FindByID(ctx context.Context, id int64) (*person.Person, error) {
	var model MemberModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return modelToPerson(&model), nil
}

func (r *PersonRepositoryImpl) FindByLegacyID(ctx context.Context, legacyID string) (*person.Person, error) {
	var model MemberModel
	if err := r.db.WithContext(ctx).Where("legacy_id = ?", legacyID).First(&model).Error; err != nil {
		return nil, err
	}
	return modelToPerson(&model), nil
}

func (r *PersonRepositoryImpl) Create(ctx context.Context, p *person.Person) error {
	model := personToModel(p)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	p.ID = model.MemberID
	return nil
}

func (r *PersonRepositoryImpl) Update(ctx context.Context, p *person.Person) error {
	model := personToModel(p)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *PersonRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&MemberModel{}, id).Error
}

func (r *PersonRepositoryImpl) Search(ctx context.Context, query *person.SearchQuery) ([]*person.Person, int64, error) {
	db := r.db.WithContext(ctx).Model(&MemberModel{})

	// 关键字搜索
	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR style_name LIKE ? OR detail_text LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 姓名精确搜索
	if query.Name != "" {
		db = db.Where("name = ?", query.Name)
	}

	// 性别过滤
	if query.Gender != nil {
		db = db.Where("gender = ?", string(*query.Gender))
	}

	// 世代过滤
	if query.Generation != nil {
		db = db.Where("generation = ?", *query.Generation)
	}

	// 出生年过滤
	if query.BirthYear != nil {
		db = db.Where("birth_year = ?", *query.BirthYear)
	}

	// 排序
	if query.SortBy != "" {
		order := query.SortBy
		if query.SortDesc {
			order += " DESC"
		}
		db = db.Order(order)
	} else {
		db = db.Order("generation, member_id")
	}

	// 统计总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var models []*MemberModel
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域模型
	persons := make([]*person.Person, len(models))
	for i, m := range models {
		persons[i] = modelToPerson(m)
	}

	return persons, total, nil
}

func (r *PersonRepositoryImpl) FindByGeneration(ctx context.Context, generation int) ([]*person.Person, error) {
	var models []*MemberModel
	if err := r.db.WithContext(ctx).Where("generation = ?", generation).Order("member_id").Find(&models).Error; err != nil {
		return nil, err
	}

	persons := make([]*person.Person, len(models))
	for i, m := range models {
		persons[i] = modelToPerson(m)
	}
	return persons, nil
}

func (r *PersonRepositoryImpl) FindByName(ctx context.Context, name string, fuzzy bool) ([]*person.Person, error) {
	var models []*MemberModel
	db := r.db.WithContext(ctx)

	if fuzzy {
		db = db.Where("name LIKE ?", "%"+name+"%")
	} else {
		db = db.Where("name = ?", name)
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, err
	}

	persons := make([]*person.Person, len(models))
	for i, m := range models {
		persons[i] = modelToPerson(m)
	}
	return persons, nil
}

func (r *PersonRepositoryImpl) FindAncestors(ctx context.Context, personID int64, depth int) ([]*person.Person, error) {
	// 使用PostgreSQL ltree查询祖先
	query := `
		WITH RECURSIVE ancestor_chain AS (
			SELECT *, 0 as level FROM members WHERE member_id = ?
			UNION ALL
			SELECT m.*, ac.level + 1
			FROM members m
			INNER JOIN ancestor_chain ac ON m.member_id = ac.father_id
			WHERE ac.level < ?
		)
		SELECT * FROM ancestor_chain ORDER BY level
	`

	var models []*MemberModel
	if err := r.db.WithContext(ctx).Raw(query, personID, depth).Scan(&models).Error; err != nil {
		return nil, err
	}

	persons := make([]*person.Person, len(models))
	for i, m := range models {
		persons[i] = modelToPerson(m)
	}
	return persons, nil
}

func (r *PersonRepositoryImpl) FindDescendants(ctx context.Context, personID int64, depth int) ([]*person.Person, error) {
	// 使用PostgreSQL ltree查询后代
	query := `
		WITH RECURSIVE descendant_chain AS (
			SELECT *, 0 as level FROM members WHERE member_id = ?
			UNION ALL
			SELECT m.*, dc.level + 1
			FROM members m
			INNER JOIN descendant_chain dc ON m.father_id = dc.member_id
			WHERE dc.level < ?
		)
		SELECT * FROM descendant_chain ORDER BY level, member_id
	`

	var models []*MemberModel
	if err := r.db.WithContext(ctx).Raw(query, personID, depth).Scan(&models).Error; err != nil {
		return nil, err
	}

	persons := make([]*person.Person, len(models))
	for i, m := range models {
		persons[i] = modelToPerson(m)
	}
	return persons, nil
}

func (r *PersonRepositoryImpl) FindSiblings(ctx context.Context, personID int64) ([]*person.Person, error) {
	subQuery := r.db.Model(&MemberModel{}).Select("father_id").Where("member_id = ?", personID)
	var models []*MemberModel
	if err := r.db.WithContext(ctx).Where("father_id = (?)", subQuery).Where("member_id != ?", personID).Find(&models).Error; err != nil {
		return nil, err
	}

	persons := make([]*person.Person, len(models))
	for i, m := range models {
		persons[i] = modelToPerson(m)
	}
	return persons, nil
}

func (r *PersonRepositoryImpl) GetFamilyTree(ctx context.Context, rootID int64, depth int) (*person.TreeResult, error) {
	// 获取根节点
	root, err := r.FindByID(ctx, rootID)
	if err != nil {
		return nil, err
	}

	// 获取所有后代
	descendants, err := r.FindDescendants(ctx, rootID, depth)
	if err != nil {
		return nil, err
	}

	nodes := make([]*person.Person, 0, len(descendants))
	for _, d := range descendants {
		if d.ID != rootID {
			nodes = append(nodes, d)
		}
	}

	return &person.TreeResult{
		Root:       root,
		Nodes:      nodes,
		MaxDepth:   depth,
		TotalNodes: len(nodes) + 1,
	}, nil
}

func (r *PersonRepositoryImpl) FindChildren(ctx context.Context, parentID int64) ([]*person.PersonChild, error) {
	query := `
		SELECT pcr.*, m.*
		FROM parent_child_relations pcr
		JOIN members m ON m.member_id = pcr.child_id
		WHERE pcr.parent_id = ?
		ORDER BY pcr.birth_order_num, pcr.relation_id
	`

	type Result struct {
		ParentChildRelationModel
		MemberModel
	}

	var results []*Result
	if err := r.db.WithContext(ctx).Raw(query, parentID).Scan(&results).Error; err != nil {
		return nil, err
	}

	children := make([]*person.PersonChild, len(results))
	for i, r := range results {
		var birthOrderNum int
		if r.BirthOrderNum != nil {
			birthOrderNum = int(*r.BirthOrderNum)
		}
		children[i] = &person.PersonChild{
			Person:        modelToPerson(&r.MemberModel),
			RelationType:  person.RelationType(r.RelationType),
			BirthOrderNum: birthOrderNum,
			IsPrimary:     r.IsPrimary,
		}
	}
	return children, nil
}

func (r *PersonRepositoryImpl) FindSpouses(ctx context.Context, personID int64) ([]*person.Spouse, error) {
	var models []*SpouseModel
	if err := r.db.WithContext(ctx).Where("member_id = ?", personID).Find(&models).Error; err != nil {
		return nil, err
	}

	spouses := make([]*person.Spouse, len(models))
	for i, m := range models {
		spouses[i] = modelToSpouse(m)
	}
	return spouses, nil
}

func (r *PersonRepositoryImpl) FindParents(ctx context.Context, personID int64) ([]*person.Person, error) {
	// 先查父亲
	var father *person.Person
	p, err := r.FindByID(ctx, personID)
	if err != nil {
		return nil, err
	}

	if p.FatherID != nil {
		f, err := r.FindByID(ctx, *p.FatherID)
		if err == nil {
			father = f
		}
	}

	parents := make([]*person.Person, 0)
	if father != nil {
		parents = append(parents, father)
	}

	return parents, nil
}

func (r *PersonRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&MemberModel{}).Count(&count).Error
	return count, err
}

func (r *PersonRepositoryImpl) CountByGeneration(ctx context.Context) (map[int]int64, error) {
	type Result struct {
		Generation int   `gorm:"column:generation"`
		Count      int64 `gorm:"column:count"`
	}

	var results []*Result
	err := r.db.WithContext(ctx).Model(&MemberModel{}).
		Select("generation, COUNT(*) as count").
		Where("generation IS NOT NULL").
		Group("generation").
		Order("generation").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[int]int64)
	for _, r := range results {
		countMap[r.Generation] = r.Count
	}
	return countMap, nil
}

func (r *PersonRepositoryImpl) BatchCreate(ctx context.Context, persons []*person.Person) error {
	if len(persons) == 0 {
		return nil
	}

	models := make([]*MemberModel, len(persons))
	for i, p := range persons {
		models[i] = personToModel(p)
	}

	return r.db.WithContext(ctx).Create(models).Error
}

func (r *PersonRepositoryImpl) BatchUpdate(ctx context.Context, persons []*person.Person) error {
	if len(persons) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, p := range persons {
			model := personToModel(p)
			if err := tx.Save(model).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
