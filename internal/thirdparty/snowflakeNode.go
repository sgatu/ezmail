package thirdparty

import "github.com/bwmarrin/snowflake"

type BaseSnowflakeNode interface {
	Generate() snowflake.ID
}

type SnowflakeNode struct {
	*snowflake.Node
}
