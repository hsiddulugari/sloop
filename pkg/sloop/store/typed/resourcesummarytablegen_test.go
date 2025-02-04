// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

/*
 * Copyright (c) 2019, salesforce.com, inc.
 * All rights reserved.
 * SPDX-License-Identifier: BSD-3-Clause
 * For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause
 */

package typed

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/salesforce/sloop/pkg/sloop/store/untyped"
	"github.com/salesforce/sloop/pkg/sloop/store/untyped/badgerwrap"
	"github.com/stretchr/testify/assert"
)

func helper_ResourceSummary_ShouldSkip() bool {
	// Tests will not work on the fake types in the template, but we want to run tests on real objects
	if "typed.Value"+"Type" == fmt.Sprint(reflect.TypeOf(ResourceSummary{})) {
		fmt.Printf("Skipping unit test")
		return true
	}
	return false
}

func Test_ResourceSummaryTable_SetWorks(t *testing.T) {
	if helper_ResourceSummary_ShouldSkip() {
		return
	}

	untyped.TestHookSetPartitionDuration(time.Hour * 24)
	db, err := (&badgerwrap.MockFactory{}).Open(badger.DefaultOptions(""))
	assert.Nil(t, err)
	err = db.Update(func(txn badgerwrap.Txn) error {
		k := (&ResourceSummaryKey{}).GetTestKey()
		vt := OpenResourceSummaryTable()
		err2 := vt.Set(txn, k, (&ResourceSummaryKey{}).GetTestValue())
		assert.Nil(t, err2)
		return nil
	})
	assert.Nil(t, err)
}

func helper_update_ResourceSummaryTable(t *testing.T, keys []string, val *ResourceSummary) (badgerwrap.DB, *ResourceSummaryTable) {
	b, err := (&badgerwrap.MockFactory{}).Open(badger.DefaultOptions(""))
	assert.Nil(t, err)
	wt := OpenResourceSummaryTable()
	err = b.Update(func(txn badgerwrap.Txn) error {
		var txerr error
		for _, key := range keys {
			txerr = wt.Set(txn, key, val)
			if txerr != nil {
				return txerr
			}
		}
		// Add some keys outside the range
		txerr = txn.Set([]byte("/a/123/"), []byte{})
		if txerr != nil {
			return txerr
		}
		txerr = txn.Set([]byte("/zzz/123/"), []byte{})
		if txerr != nil {
			return txerr
		}
		return nil
	})
	assert.Nil(t, err)
	return b, wt
}

func Test_ResourceSummaryTable_GetUniquePartitionList_Success(t *testing.T) {
	if helper_ResourceSummary_ShouldSkip() {
		return
	}

	db, wt := helper_update_ResourceSummaryTable(t, (&ResourceSummaryKey{}).SetTestKeys(), (&ResourceSummaryKey{}).SetTestValue())
	var partList []string
	var err1 error
	err := db.View(func(txn badgerwrap.Txn) error {
		partList, err1 = wt.GetUniquePartitionList(txn)
		return nil
	})
	assert.Nil(t, err)
	assert.Nil(t, err1)
	assert.Len(t, partList, 3)
	assert.Contains(t, partList, someMinPartition)
	assert.Contains(t, partList, someMiddlePartition)
	assert.Contains(t, partList, someMaxPartition)
}

func Test_ResourceSummaryTable_GetUniquePartitionList_EmptyPartition(t *testing.T) {
	if helper_ResourceSummary_ShouldSkip() {
		return
	}

	db, wt := helper_update_ResourceSummaryTable(t, []string{}, &ResourceSummary{})
	var partList []string
	var err1 error
	err := db.View(func(txn badgerwrap.Txn) error {
		partList, err1 = wt.GetUniquePartitionList(txn)
		return err1
	})
	assert.Nil(t, err)
	assert.Len(t, partList, 0)
}
