package sf

import (
	"github.com/bwmarrin/snowflake"
)

func NewSnowflake(
	opt *Options,
) (*snowflake.Node, error) {
	sf, err := snowflake.NewNode(opt.NodeNumber)
	if err != nil {
		return nil, err
	}

	return sf, nil
}
