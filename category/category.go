package category

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/components/htree"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Category struct {
	Code string
	tree *htree.Tree
}

func newCategory(code string, flag uint, cfg []uint) (*Category, error) {
	ctg := &Category{
		Code: code,
		tree: nil,
	}
	var err error
	ctg.tree, err = htree.NewTree(fmt.Sprintf("category_%s", code), flag, cfg)
	if err != nil {
		hlog.Err("hyper.category.newCategory: NewTree", zap.Error(err))
		return nil, err
	}
	return ctg, nil
}

func (c *Category) Add(parent htree.ID, name string, icon string) (htree.ID, error) {
	if name == "" {
		return 0, errors.New("require name")
	}
	if parent == Root {
		parent = c.tree.Root()
	}
	b, err := hpg.Exist[CtgM](c.db().PG().Table(c.tableName()), "parent = ? AND name = ?", parent, name)
	if err != nil {
		return -1, err
	}
	if b {
		return -1, fmt.Errorf("exists: parent=%d,name=%s", parent, name)
	}
	var newID htree.ID
	err = c.tree.Next(parent, func(id htree.ID) error {
		newID = id
		return nil
	})
	if err != nil {
		return -1, err
	}
	ctgM := &CtgM{
		ID:     newID,
		Parent: parent,
		Name:   name,
		Icon:   icon,
	}
	err = hpg.Create[CtgM](c.db().PG().Table(c.tableName()), ctgM)
	if err != nil {
		fmt.Println(hjson.MustToString(ctgM))
		return -1, err
	}
	return ctgM.ID, nil
}

func (c *Category) Mut(id htree.ID, name string, icon string) error {
	if name == "" {
		return errors.New("require name")
	}
	dbM, err := hpg.MustGet[CtgM](c.db().PG().Table(c.tableName()), "id = ?", id)
	if err != nil {
		return err
	}
	mut := make(map[string]any)
	if dbM.Name != name {
		mut["name"] = name
	}
	if dbM.Icon != icon {
		mut["icon"] = icon
	}
	if len(mut) == 0 {
		return nil
	}
	b, err := hpg.Exist[CtgM](c.db().PG().Table(c.tableName()), "parent = ? AND name = ? AND id <> ?", dbM.Parent, name, id)
	if err != nil {
		return err
	}
	if b {
		return fmt.Errorf("exists: parent=%d,name=%s", dbM.Parent, name)
	}

	err = hpg.Update[CtgM](c.db().PG().Table(c.tableName()), mut, "id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Category) Get(parent htree.ID, deep int) ([]*Categ, error) {
	if deep < 1 || deep > c.tree.Factory().IdDeep() {
		return nil, fmt.Errorf("invalid deep: %d", deep)
	}
	if parent == Root {
		parent = c.tree.Root()
	}
	minID, maxID, base, err := c.tree.Factory().DirectChildren(parent)
	if err != nil {
		return nil, err
	}
	var arr []*Categ
	arr, err = c.loadChildren(minID, maxID, base)
	if err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return []*Categ{}, nil
	}
	newDeep := deep - 1
	if newDeep <= 0 {
		return arr, nil
	}
	for _, categ := range arr {
		categ.Children, err = c.Get(categ.ID, newDeep)
		if err != nil {
			return nil, err
		}
	}
	return arr, nil
}

func (c *Category) loadChildren(minID htree.ID, maxID htree.ID, base htree.ID) ([]*Categ, error) {
	arrM, err := hpg.Find[CtgM](func() *gorm.DB {
		return c.db().PG().Table(c.tableName()).
			Where("id % ? = 0 AND id >= ? AND id <= ?", base, minID, maxID)
	})
	if err != nil {
		return []*Categ{}, err
	}
	if len(arrM) == 0 {
		return []*Categ{}, nil
	}
	var arr []*Categ
	for _, item := range arrM {
		arr = append(arr, item.ToCateg())
	}
	return arr, nil
}

func (c *Category) Helix() helix.Helix {
	return helix.BuildHelix(c.Code, func() (context.Context, error) {
		err := hpg.AutoMigrateWithTable(c.db().PG(), hpg.NewTable(c.tableName(), &CtgM{}))
		if err != nil {
			hlog.Err("hyper.category.Helix: AutoMigrateWithTable", zap.Error(err))
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	})
}

func (c *Category) tableName() string {
	return c.Code
}

func (c *Category) db() *hpg.Database {
	return zplt.HelixPgDB()
}
