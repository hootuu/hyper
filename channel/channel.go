package channel

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/htree"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"gorm.io/gorm"
)

func Create(
	ctx context.Context,
	biz collar.Collar,
	name string,
	icon string,
	seq int,
	call func(ctx context.Context, chnM *ChnM) error,
) error {
	if name == "" {
		return fmt.Errorf("require name")
	}
	tx := db(ctx)
	id, err := gChannelIdTree.NextID(gChannelIdTree.Root())
	if err != nil {
		return err
	}
	chnM := &ChnM{
		Biz:    biz.ID(),
		ID:     id,
		Parent: gChannelIdTree.Root(),
		Name:   name,
		Icon:   icon,
		Seq:    seq,
	}
	err = hdb.Create[ChnM](tx, chnM)
	if err != nil {
		return err
	}
	err = call(ctx, chnM)
	if err != nil {
		return err
	}
	return nil
}

func Add(
	ctx context.Context,
	parent ID,
	name string,
	icon string,
	seq int,
	call func(ctx context.Context, chnM *ChnM) error,
) error {
	if parent == 0 {
		return fmt.Errorf("require parent")
	}
	if name == "" {
		return fmt.Errorf("require name")
	}

	tx := db(ctx)
	parentM, err := hdb.Get[ChnM](tx, "id = ?", parent)
	if err != nil {
		return err
	}
	if parentM == nil {
		return fmt.Errorf("parent not found")
	}
	id, err := gChannelIdTree.NextID(parent)
	if err != nil {
		return err
	}
	chnM := &ChnM{
		Biz:    parentM.Biz,
		ID:     id,
		Parent: parentM.ID,
		Name:   name,
		Icon:   icon,
		Seq:    seq,
	}
	err = hdb.Create[ChnM](tx, chnM)
	if err != nil {
		return err
	}
	err = call(ctx, chnM)
	if err != nil {
		return err
	}
	return nil
}

func Get(ctx context.Context, parent ID, deep int) ([]*Channel, error) {
	if deep < 1 || deep > gChannelIdTree.Factory().IdDeep() {
		return nil, fmt.Errorf("invalid deep: %d", deep)
	}
	if parent == Root {
		parent = gChannelIdTree.Root()
	}
	minID, maxID, base, err := gChannelIdTree.Factory().DirectChildren(parent)
	if err != nil {
		return nil, err
	}
	var arr []*Channel
	arr, err = loadChildren(ctx, minID, maxID, base)
	if err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return []*Channel{}, nil
	}
	newDeep := deep - 1
	if newDeep <= 0 {
		return arr, nil
	}
	for _, categ := range arr {
		categ.Children, err = Get(ctx, categ.ID, newDeep)
		if err != nil {
			return nil, err
		}
	}
	return arr, nil
}

func loadChildren(ctx context.Context, minID htree.ID, maxID htree.ID, base htree.ID) ([]*Channel, error) {
	arrM, err := hdb.Find[ChnM](func() *gorm.DB {
		return db(ctx).Where("id % ? = 0 AND id >= ? AND id <= ?", base, minID, maxID)
	})
	if err != nil {
		return []*Channel{}, err
	}
	if len(arrM) == 0 {
		return []*Channel{}, nil
	}
	var arr []*Channel
	for _, item := range arrM {
		arr = append(arr, item.ToChannel())
	}
	return arr, nil
}

func db(ctx context.Context) *gorm.DB {
	tx := hdb.CtxTx(ctx)
	if tx == nil {
		tx = zplt.HelixPgDB().PG()
	}
	return tx
}
