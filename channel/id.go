package channel

import (
	"github.com/hootuu/helix/components/htree"
)

var gChannelIdTree *htree.Tree

func initChannelIdTree() error {
	var err error
	gChannelIdTree, err = htree.NewTree("hyper_channel_tree", 1,
		[]uint{3, 3, 3, 3, 3})
	if err != nil {
		return err
	}
	return nil
}
