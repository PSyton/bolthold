// Copyright 2016 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package bolthold_test

import (
	"testing"

	bh "github.com/psyton/bolthold"
	bolt "go.etcd.io/bbolt"
)

func TestIndexSlice(t *testing.T) {
	testWrap(t, func(store *bh.Store, t *testing.T) {
		var testData = []ItemTest{
			{
				Key:  0,
				Name: "John",
				Tags: []string{"red", "green", "blue"},
			},
			{
				Key:  1,
				Name: "Bill",
				Tags: []string{"red", "purple"},
			},
			{
				Key:  2,
				Name: "Jane",
				Tags: []string{"red", "orange"},
			},
			{
				Key:  3,
				Name: "Brian",
				Tags: []string{"red", "purple"},
			},
		}

		for _, data := range testData {
			ok(t, store.Insert(data.Key, data))
		}

		b := store.Bolt()

		ok(t, b.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("_index_v2:ItemTest:Tags"))
			assert(t, bucket != nil, "No index bucket found for Tags index")

			indexCount := 0
			_ = bucket.ForEach(func(k, v []byte) error {
				indexCount++
				return nil
			})

			// 4 - red, 2 - purple, 1 — orange, 1 - blue, 1 - green
			equals(t, indexCount, 9)
			return nil
		}))

	})
}

func Test85SliceIndex(t *testing.T) {
	type Event struct {
		Id         uint64
		Type       string   `boltholdIndex:"Type"`
		Categories []string `boltholdSliceIndex:"Categories"`
	}

	testWrap(t, func(store *bh.Store, t *testing.T) {
		e1 := &Event{Id: 1, Type: "Type1", Categories: []string{"Cat 1", "Cat 2"}}
		e2 := &Event{Id: 2, Type: "Type1", Categories: []string{"Cat 3"}}

		ok(t, store.Insert(e1.Id, e1))
		ok(t, store.Insert(e2.Id, e2))

		var es []*Event
		ok(t, store.Find(&es, bh.Where("Categories").Contains("Cat 1").Index("Categories")))
		equals(t, len(es), 1)
	})
}

func Test87SliceIndex(t *testing.T) {
	type Event struct {
		Id         uint64
		Type       string   `boltholdIndex:"Type"`
		Categories []string `boltholdSliceIndex:"Categories"`
	}

	testWrap(t, func(store *bh.Store, t *testing.T) {
		e1 := &Event{Id: 1, Type: "Type1", Categories: []string{"Cat 1", "Cat 2"}}
		e2 := &Event{Id: 2, Type: "Type1", Categories: []string{"Cat 3"}}

		ok(t, store.Insert(e1.Id, e1))
		ok(t, store.Insert(e2.Id, e2))
		var es []*Event
		ok(t, store.Find(&es, bh.Where("Categories").ContainsAny("Cat 1").Index("Categories")))
		equals(t, len(es), 1)
	})
}

func TestSliceIndexWithPointers(t *testing.T) {
	type Event struct {
		Id         uint64
		Type       string    `boltholdIndex:"Type"`
		Categories []*string `boltholdSliceIndex:"Categories"`
	}

	testWrap(t, func(store *bh.Store, t *testing.T) {
		cat1 := "Cat 1"
		cat2 := "Cat 2"
		cat3 := "Cat 3"

		e1 := &Event{Id: 1, Type: "Type1", Categories: []*string{&cat1, &cat2}}
		e2 := &Event{Id: 2, Type: "Type1", Categories: []*string{&cat3}}

		ok(t, store.Insert(e1.Id, e1))
		ok(t, store.Insert(e2.Id, e2))

		var es []*Event
		ok(t, store.Find(&es, bh.Where("Categories").ContainsAll("Cat 1").Index("Categories")))
		equals(t, len(es), 1)
	})
}
