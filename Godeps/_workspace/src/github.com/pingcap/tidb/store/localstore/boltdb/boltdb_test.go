// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package boltdb

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
	. "github.com/pingcap/check"
	"github.com/pingcap/tidb/store/localstore/engine"
)

func TestT(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testSuite{})

type testSuite struct {
	db engine.DB
}

const testPath = "/tmp/test-tidb-boltdb"

func (s *testSuite) SetUpSuite(c *C) {
	var (
		d   Driver
		err error
	)
	s.db, err = d.Open(testPath)
	c.Assert(err, IsNil)
}

func (s *testSuite) TearDownSuite(c *C) {
	s.db.Close()
	os.Remove(testPath)
}

func (s *testSuite) TestPutNilAndDelete(c *C) {
	d := s.db
	rawDB := d.(*db).DB
	// nil as value
	b := d.NewBatch()
	b.Put([]byte("aa"), nil)
	err := d.Commit(b)
	c.Assert(err, IsNil)

	v, err := d.Get([]byte("aa"))
	c.Assert(err, IsNil)
	c.Assert(len(v), Equals, 0)

	found := false
	rawDB.View(func(tx *bolt.Tx) error {
		b1 := tx.Bucket(bucketName)
		c := b1.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if string(k) == "aa" {
				found = true
			}
		}
		return nil
	})
	c.Assert(found, Equals, true)

	// real delete
	b = d.NewBatch()
	b.Delete([]byte("aa"))
	err = d.Commit(b)
	c.Assert(err, IsNil)

	found = false
	rawDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if string(k) == "aa" {
				found = true
			}

		}
		return nil
	})
	c.Assert(found, Equals, false)
}

func (s *testSuite) TestDB(c *C) {
	db := s.db

	b := db.NewBatch()
	b.Put([]byte("a"), []byte("1"))
	b.Put([]byte("b"), []byte("2"))
	b.Delete([]byte("c"))

	err := db.Commit(b)
	c.Assert(err, IsNil)

	v, err := db.Get([]byte("c"))
	c.Assert(err, IsNil)
	c.Assert(v, IsNil)

	v, err = db.Get([]byte("a"))
	c.Assert(err, IsNil)
	c.Assert(v, DeepEquals, []byte("1"))

	iter, err := db.Seek(nil)
	c.Assert(err, IsNil)
	c.Assert(iter.Next(), Equals, true)
	c.Assert(iter.Key(), DeepEquals, []byte("a"))
	c.Assert(iter.Next(), Equals, true)
	c.Assert(iter.Key(), DeepEquals, []byte("b"))
	c.Assert(iter.Next(), Equals, false)
	iter.Release()

	iter, err = db.Seek([]byte("b"))
	c.Assert(err, IsNil)
	c.Assert(iter.Next(), Equals, true)
	c.Assert(iter.Key(), DeepEquals, []byte("b"))
	c.Assert(iter.Value(), DeepEquals, []byte("2"))
	c.Assert(iter.Next(), Equals, false)
	iter.Release()
}
