package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVParser_ParseCSVFromBytes_Valid(t *testing.T) {
	parser := NewCSVParser(10)
	csvData := []byte(`sku,location,tenant_id,seller_id
sku1,loc1,tenant1,seller1
sku2,loc2,tenant2,seller2`)

	ctx := context.Background()
	result, err := parser.ParseCSVFromBytes(ctx, csvData)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.TotalRows)
	assert.Equal(t, 2, result.ValidRows)
	assert.Equal(t, 0, result.InvalidRows)
	assert.Len(t, result.ValidData, 2)
	assert.Len(t, result.InvalidData, 0)
}

func TestCSVParser_ParseCSVFromBytes_InvalidHeaders(t *testing.T) {
	parser := NewCSVParser(10)
	csvData := []byte(`sku,location,tenant_id
sku1,loc1,tenant1`)

	ctx := context.Background()
	result, err := parser.ParseCSVFromBytes(ctx, csvData)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid CSV headers")
}



func TestCSVParser_ParseCSVFromBytes_InvalidRows(t *testing.T) {
	parser := NewCSVParser(10)
	csvData := []byte(`sku,location,tenant_id,seller_id
,loc1,tenant1,seller1
sku2,loc2,tenant2,`)

	ctx := context.Background()
	result, err := parser.ParseCSVFromBytes(ctx, csvData)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.TotalRows)
	assert.Equal(t, 0, result.ValidRows)
	assert.Equal(t, 2, result.InvalidRows)
	assert.Len(t, result.ValidData, 0)
	assert.Len(t, result.InvalidData, 2)
}
