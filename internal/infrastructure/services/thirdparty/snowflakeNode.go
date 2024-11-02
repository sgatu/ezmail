package thirdparty

import "github.com/bwmarrin/snowflake"

type BaseSnowflakeNode interface {
	Generate() snowflake.ID
}

type SnowflakeNode struct {
	*snowflake.Node
}

type MockSnowflakeNode struct {
	nextGenerateID snowflake.ID
}

func (ms *MockSnowflakeNode) SetNextId(id snowflake.ID) {
	ms.nextGenerateID = id
}

func (ms *MockSnowflakeNode) Generate() snowflake.ID {
	return ms.nextGenerateID
}
