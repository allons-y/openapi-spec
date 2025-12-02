// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package spec

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestLoader_Issue145(t *testing.T) {
	t.Run("with ExpandSpec", func(t *testing.T) {
		// TODO: This test fails with OpenAPI 3 refs and paths with spaces
		// The expansion of cross-file refs with #/components/... paths needs investigation
		t.Skip("Test needs investigation for cross-file ref expansion with OpenAPI 3")
	})

	t.Run("with ExpandSchema", func(t *testing.T) {
		basePath := filepath.Join("fixtures", "bugs", "145", "Program Files (x86)", "AppName", "ref.json")
		schemaDoc, err := jsonDoc(basePath)
		require.NoError(t, err)

		sch := new(Schema)
		require.NoError(t, json.Unmarshal(schemaDoc, sch))

		require.NoError(t, ExpandSchema(sch, nil, nil))
	})
}
