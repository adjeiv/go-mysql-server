// Copyright 2023 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memo

import "github.com/dolthub/go-mysql-server/sql"

type scalarProps struct {
	Cols          sql.ColSet
	Tables        sql.FastIntSet
	nullRejecting bool
}

func newScalarProps(e ScalarExpr) *scalarProps {
	var tables sql.FastIntSet
	var cols sql.ColSet
	nullRejecting := true
	switch e := e.(type) {
	case *ColRef:
		tables.Add(int(e.Table) - 1)
		cols.Add(e.Col)
	case *NullSafeEq:
		nullRejecting = false
	case *Hidden:
		cols = cols.Union(e.Cols)
		tables = tables.Union(e.Tables)
	default:
	}

	for _, c := range e.Children() {
		tables = tables.Union(c.ScalarProps().Tables)
		cols = cols.Union(c.ScalarProps().Cols)
	}
	return &scalarProps{
		nullRejecting: nullRejecting,
		Tables:        tables,
		Cols:          cols,
	}
}
