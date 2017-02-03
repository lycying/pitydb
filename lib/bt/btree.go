package btree

// Copyright (c) 2011 Alexander Sychev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package btree implements B-trees with fixed size keys are saved in a specified storage, http://en.wikipedia.org/wiki/Btree.

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

var (
	bo = binary.LittleEndian
)

// Errors in addition to IO errors.
var (
	NoReader        = errors.New("reader is not specified")
	NoWriter        = errors.New("write operations are assumed - io.ReadWriteSeeker has to be specified insead of io.ReadSeeker")
	OddCapacity     = errors.New("capacity must be even")
	MagicMismatch   = errors.New("magic mismatch")
	KeySizeMismatch = errors.New("key size mismatch")
)

// An interface a key must support to be stored in the tree.
// A key must be exact the len of b
// All fields of the key must be exportable.
// Compare has to return a value less that 0, equal of 0 or more that 0 if k is less, equal or more that an underlying key in b. An error must be returned if something is wrong.
// Size returns a size of key in bytes, the size has to be immutable
// Read reads a key from b
// Write writes a key to b
type Key interface {
	Compare(b []byte) (int, error)
	Size() uint
	Read(b []byte) error
	Write(b []byte) error
}

// A basic interface for tree operations.
type Tree interface {
	Find(key Key) (Key, error)
	Insert(key Key) (Key, error)
	Update(key Key) (Key, error)
	Delete(key Key) (Key, error)
	Enum(key Key) func() (Key, error)
	ReverseEnum(key Key) func() (Key, error)
}

// header of file with the tree
type fileHeader struct {
	Magic      [16]byte // file magic
	KeySize    uint32   // size of key in bytes
	Capacity   uint32   // size of btree node in datas
	Root       int64    // offset of the root node
	EmptyNodes int64    // offset of empty nodes offsets
}

// node struct.
type node struct {
	size   int   // size of one data element
	offset int64 // offset of node
	count  int
	raw    []byte // raw node in bytes
	datas  []byte // datas
}

// an internal representation of B-tree.
type BTree struct {
	header  fileHeader    // header of index file
	storage io.ReadSeeker // storage interface
	empty   []int64       // empty nodes offsets
	key     Key           // key interface
}

// NewBTree creates a new B-tree with magic like a file magic, key like a key type and capacity like a number of elements per data.
// The capacity must be even.
// It returns a pointer to the new tree and an error, if any
func NewBTree(storage io.ReadWriteSeeker, magic [16]byte, key Key, capacity uint) (Tree, error) {
	if storage == nil {
		return nil, NoWriter
	}
	if capacity%2 == 1 {
		return nil, OddCapacity
	}
	this := new(BTree)
	this.storage = storage
	this.header.Magic = magic
	this.header.KeySize = uint32(key.Size())
	this.header.Capacity = uint32(capacity)
	this.header.Root = -1
	this.header.EmptyNodes = -1
	if _, err := storage.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	if err := this.header.write(storage, nil); err != nil {
		return nil, err
	}
	this.key = key
	return this, nil
}

// OpenBTree opens an existing B-tree.
// The file magic and magic must be the same, the type and the size of key and of the key in the tree must be the same too.
// If changing of the tree is planning, storage has to be io.ReadWriteSeeker.
// It returns a pointer to the new tree and an error, if any.
func OpenBTree(storage io.ReadSeeker, magic [16]byte, key Key) (Tree, error) {
	if storage == nil {
		return nil, NoReader
	}
	this := new(BTree)
	this.storage = storage
	if en, err := this.header.read(this.storage); err != nil {
		return nil, err
	} else {
		this.empty = en
	}
	if !bytes.Equal(this.header.Magic[:], magic[:]) {
		return nil, MagicMismatch
	}
	if uint(this.header.KeySize) != key.Size() {
		return nil, KeySizeMismatch
	}
	this.key = key
	return this, nil
}

// Find searches a key in the tree.
// It returns the key if it is found or nil if the key is not found and an error, if any.
func (this *BTree) Find(key Key) (Key, error) {
	if uint(this.header.KeySize) != key.Size() {
		return nil, KeySizeMismatch
	}
	_, _, b, err := this.find(key)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, nil
	}
	if err := this.key.Read(b); err != nil {
		return nil, err
	}
	return this.key, nil
}

// Find searches a key in the tree.
// It returns the key if it is already exists or nil if the key is inserted and an error, if any.
func (this *BTree) Insert(key Key) (Key, error) {
	if uint(this.header.KeySize) != key.Size() {
		return nil, KeySizeMismatch
	}
	b, err := this.insert(key)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, nil
	}
	if err := this.key.Read(b); err != nil {
		return nil, err
	}
	return this.key, nil
}

// Update updates a key in the tree.
// This is useful if the key is complex type with additional information.
// It returns an old value of the key if the key is updated or nil if the key is not found and an error, if any.
func (this *BTree) Update(key Key) (Key, error) {
	if this.storage.(io.WriteSeeker) == nil {
		return nil, NoWriter
	}
	if uint(this.header.KeySize) != key.Size() {
		return nil, KeySizeMismatch
	}
	n, i, b, err := this.find(key)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, nil
	}
	if err := this.key.Read(b); err != nil {
		return nil, err
	}
	n.setKey(i, key)
	return this.key, this.writeNode(n)
}

// BUG(santucco): The storage with B-tree can't be reduced even after erasing all of its keys. All empty nodes are stored and reused when necessary

// Delete deletes a key from the tree.
// It returns the key if it is deleted or nil if the key is not found and an error, if any.
func (this *BTree) Delete(key Key) (Key, error) {
	if this.storage.(io.WriteSeeker) == nil {
		return nil, NoWriter
	}
	if uint(this.header.KeySize) != key.Size() {
		return nil, KeySizeMismatch
	}
	b, err := this.delete(key)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, nil
	}
	if err := this.key.Read(b); err != nil {
		return nil, err
	}
	return this.key, nil
}

// Enum returns a function-iterator to process enumeration entire the tree from lower to bigger keys.
// Enumerating starts with key, if it is specified, or with lowest key otherwise.
// The iterator returns the key or nil if the end of the tree is reached and an error, if any.
func (this *BTree) Enum(key Key) func() (Key, error) {
	nodes := make([]*node, 0, 10)
	offset := this.header.Root
	return func() (Key, error) {
		for true {
			if offset == -1 {
				if len(nodes) == 0 {
					return nil, nil
				}
				l := len(nodes) - 1
				p := nodes[l]
				if p.count > 0 {
					offset = p.getOffset(0)
					if err := this.key.Read(p.getKey(0)); err != nil {
						return nil, err
					}
					p.datas = p.datas[p.size:]
					p.count--
					return this.key, nil
				}
				nodes = nodes[:l]
			} else {
				var p node
				p.init(this)
				if err := p.read(this.storage, offset); err != nil {
					return nil, err
				}
				if key != nil {
					if idx, _, less, err := p.find(key); err != nil {
						return nil, err
					} else if idx != -1 {
						offset = -1
						p.datas = p.datas[p.size*idx:]
						p.count -= idx
						key = nil
					} else if less == -1 {
						offset = p.getLeast()
					} else {
						offset = p.getOffset(less)
						p.datas = p.datas[p.size*(less+1):]
						p.count -= less + 1
					}
				} else {
					offset = p.getLeast()
				}
				nodes = append(nodes, &p)
			}
		}
		return nil, nil
	}
}

// ReverseEnum returns a function-iterator to process enumeration entire the tree from bigger to lower keys.
// Enumerating starts with key, if it is specified, or with biggest key otherwise.
// The iterator returns the key or nil if the end of the tree is reached and an error, if any.
func (this *BTree) ReverseEnum(key Key) func() (Key, error) {
	nodes := make([]*node, 0, 10)
	offset := this.header.Root
	return func() (Key, error) {
		for true {
			if offset == -1 {
				if len(nodes) == 0 {
					return nil, nil
				}
				l := len(nodes) - 1
				p := nodes[l]
				if p.count > 0 {
					if err := this.key.Read(p.getKey(int(p.count - 1))); err != nil {
						return nil, err
					}
					p.count--
					p.datas = p.datas[:int(p.count)*p.size]
					if p.count > 0 {
						offset = p.getOffset(int(p.count - 1))
					} else {
						offset = p.getLeast()
					}
					return this.key, nil
				}
				nodes = nodes[:l]
			} else {
				var p node
				p.init(this)
				if err := p.read(this.storage, offset); err != nil {
					return nil, err
				}
				if key != nil {
					if idx, _, less, err := p.find(key); err != nil {
						return nil, err
					} else if idx != -1 {
						offset = -1
						p.setOffset(idx, -1)
						idx++
						p.datas = p.datas[:p.size*idx]
						p.count = idx
						key = nil
					} else if less != -1 {
						offset = p.getOffset(less)
						less++
						p.datas = p.datas[:p.size*less]
						p.count = less
					} else {
						offset = p.getLeast()
						p.count = 0
					}
				} else {
					offset = p.getOffset(int(p.count - 1))
				}
				nodes = append(nodes, &p)
			}
		}
		return nil, nil
	}
}

// KeySize returns a size in bytes of the underlying key.
func (this BTree) KeySize() uint {
	return uint(this.header.KeySize)
}

// Capacity returns a number of keys per node.
func (this BTree) Capacity() uint {
	return uint(this.header.Capacity)
}

// Magic returns a magic of the B-tree storage.
func (this BTree) Magic() [16]byte {
	return this.header.Magic
}

// writeNode writes node to a buffer and writes the buffer to the storage.
// It returns nil on success or an error.
func (this *BTree) writeNode(p *node) error {
	writer := this.storage.(io.WriteSeeker)
	if p.offset != -1 {
		if _, err := writer.Seek(p.offset, os.SEEK_SET); err != nil {
			return err
		}
	} else {
		if len(this.empty) != 0 {
			p.offset = this.empty[len(this.empty)-1]
			this.empty[len(this.empty)-1] = 0
			this.empty = this.empty[:len(this.empty)-1]
			if err := this.header.write(writer, this.empty); err != nil {
				return err
			}
			if _, err := writer.Seek(p.offset, os.SEEK_SET); err != nil {
				return err
			}
		} else if this.header.EmptyNodes != -1 {
			p.offset = this.header.EmptyNodes
			if off, err := writer.Seek(0, os.SEEK_END); err != nil {
				return err
			} else if off-this.header.EmptyNodes > int64(cap(p.raw)) {
				this.header.EmptyNodes = off
			} else {
				this.header.EmptyNodes = -1
			}
			if err := this.header.write(writer, this.empty); err != nil {
				return err
			}
			if _, err := writer.Seek(p.offset, os.SEEK_SET); err != nil {
				return err
			}
		} else if off, err := writer.Seek(0, os.SEEK_END); err != nil {
			return err
		} else {
			p.offset = off
		}
	}

	return p.write(writer)
}

// find finds a key.
// It returns pointer to a node with the key and an index of the key in the node and an error, if any.
func (this *BTree) find(key Key) (*node, int, []byte, error) {
	curoff := this.header.Root
	for curoff > 0 {
		var p node
		p.init(this)
		if err := p.read(this.storage, curoff); err != nil {
			return nil, -1, nil, err
		}
		if idx, b, less, err := p.find(key); err != nil {
			return nil, -1, nil, err
		} else if idx != -1 {
			return &p, idx, b, nil
		} else if less == -1 {
			curoff = p.getLeast()
		} else {
			curoff = p.getOffset(less)
		}
	}
	return nil, -1, nil, nil
}

// insert inserts a key in the tree.
// It returns the key if it is already exists or nil if the key is inserted and an error, if any.
func (this *BTree) insert(key Key) ([]byte, error) {
	offset := this.header.Root
	offsets := make([]int64, 0, 10)
	c := 0
	var p node
	p.init(this)
	var less int = -1
	for true {
		var b []byte
		var err error
		if offset != -1 {
			if err := p.read(this.storage, offset); err != nil {
				return nil, err
			}
			_, b, less, err = p.find(key)
		}
		if err != nil {
			return nil, err
		}
		if uint32(p.count) == this.header.Capacity {
			c++
		} else {
			c = 0
		}
		if b != nil {
			return b, nil
		}
		least := p.getLeast()
		if less == -1 && least != -1 {
			offset = least
			offsets = append(offsets, p.offset)
			continue
		} else if less != -1 {
			if off := p.getOffset(less); off != -1 {
				offset = off
				offsets = append(offsets, p.offset)
				continue
			}
		}
		break
	}
	// if current and previous nodes are full, trying to reserve a space in the storage
	// and put nodes in the list of empty nodes
	if c > len(this.empty) {
		c -= len(this.empty)
		var e node
		e.init(this)
		for i := 0; i < c; i++ {
			e.offset = -1
			if err := this.writeNode(&e); err != nil {
				return nil, err
			}
			this.empty = append(this.empty, e.offset)
		}
	}
	p.insert(less, key, -1)
	full := uint32(p.count) > this.header.Capacity
	if !full {
		if err := this.writeNode(&p); err != nil {
			return nil, err
		}
	}
	if this.header.Root == -1 {
		this.header.Root = p.offset
		return nil, this.header.write(this.storage.(io.WriteSeeker), this.empty)
	}
	if !full {
		return nil, nil
	}
	return nil, this.split(&p, offsets)
}

// split splits a node pointed by p on two nodes, inserts a middle data in a parent node and saves all changed nodes.
// It returns nil on success or an error.
func (this *BTree) split(p *node, offsets []int64) error {
	var root int64 = -1
	for true {
		var np node
		np.init(this)
		np.count = p.count / 2
		np.datas = np.datas[:int(np.count)*p.size]
		copy(np.datas, p.datas[int(np.count+1)*p.size:])
		off := p.getOffset(int(np.count))
		np.setLeast(off)
		if err := this.writeNode(&np); err != nil {
			return err
		}
		offset := np.offset
		if err := this.key.Read(p.getKey(int(np.count))); err != nil {
			return err
		}
		p.count /= 2
		if err := this.writeNode(p); err != nil {
			return err
		}
		if len(offsets) == 0 {
			var rp node
			rp.init(this)
			rp.insert(-1, this.key, offset)
			rp.setLeast(p.offset)
			if err := this.writeNode(&rp); err != nil {
				return err
			}
			root = rp.offset
		} else {
			off := offsets[len(offsets)-1]
			offsets = offsets[:len(offsets)-1]
			p.read(this.storage, off)
			_, _, less, err := p.find(this.key)
			if err != nil {
				return err
			}
			_, err = p.insert(less, this.key, offset)
			if err != nil {
				return err
			}
		}
		if uint32(p.count) <= this.header.Capacity {
			if err := this.writeNode(p); err != nil {
				return err
			}
			break
		}
	}
	if root == -1 {
		return nil
	}
	this.header.Root = root
	return this.header.write(this.storage.(io.WriteSeeker), this.empty)
}

// delete deletes a key from the tree.
// It returns nil on success or an error.
func (this *BTree) delete(key Key) ([]byte, error) {
	offset := this.header.Root
	index := -1
	var pnode *node
	var poff []byte

	var p *node
	// loking for the key
	var buf []byte
	for true {
		tmp := new(node)
		tmp.init(this)
		less := -1
		var b []byte
		if offset != -1 {
			var err error
			if err = tmp.read(this.storage, offset); err != nil {
				return nil, err
			}
			if index, b, less, err = tmp.find(key); err != nil {
				return nil, err
			}
		}
		if b != nil {
			buf = append(buf, b...)
			p = tmp
			break
		}
		pnode = tmp
		least := tmp.getLeast()
		if less == -1 && least != -1 {
			offset = least
			poff = tmp.raw[0:4]
		} else if less != -1 {
			if off := tmp.getOffset(less); off != -1 {
				offset = off
				poff = tmp.datas[less*pnode.size : less*pnode.size+4]
				continue
			}
		} else {
			return nil, nil
		}
	}
	// removing data
	for true {
		offset = p.getOffset(index)
		if index < int(p.count) {
			copy(p.datas[index*p.size:], p.datas[(index+1)*p.size:])
		}
		p.count--
		index--
		if offset == -1 {
			if uint32(p.count) >= this.header.Capacity/2 {
				return buf, this.writeNode(p)
			}
			index = -1
			offset = p.getLeast()
			if offset != -1 {
				pnode = p
				poff = p.raw[0:4]
			}
			for i := 0; offset == -1 && i < int(p.count); i++ {
				off := p.getOffset(i)
				if off != -1 {
					offset = off
					poff = p.datas[i*p.size : i*p.size+4]
					index = i
				}
			}
		}
		if offset == -1 {
			// it is a leaf node
			if err := this.writeNode(p); err != nil {
				return buf, err
			}
			if p.count > 0 {
				return buf, nil
			}
			// it is an empty leaf node
			if poff != nil {
				copy(poff, []byte{0xff, 0xff, 0xff, 0xff})
				if err := this.writeNode(pnode); err != nil {
					return buf, err
				}
			}
			this.empty = append(this.empty, p.offset)
			if p.offset == this.header.Root {
				this.header.Root = -1
			}
			return buf, this.header.write(this.storage.(io.WriteSeeker), this.empty)
		}
		if poff != nil {
			copy(poff, []byte{0xff, 0xff, 0xff, 0xff})
		}
		var np *node
		for off := offset; ; {
			tmp := new(node)
			tmp.init(this)
			if err := tmp.read(this.storage, off); err != nil {
				return buf, err
			}
			least := tmp.getLeast()
			if least == -1 {
				np = tmp
				break
			}
			off = least
			pnode = tmp
			poff = tmp.raw[0:4]
		}
		if err := this.key.Read(np.getKey(0)); err != nil {
			return nil, err
		}
		po, err := p.insert(index, this.key, offset)
		if err != nil {
			return nil, err
		}
		if offset == np.offset {
			poff = po
			pnode = p
		}
		index = 0
		if err := this.writeNode(p); err != nil {
			return buf, err
		}
		p = np

	}
	return buf, nil
}

// bufferSize returns a size of a buffer in bytes will be read from the reader at time.
func (this fileHeader) bufferSize() uint32 {
	return /*size of least*/ 4 +
		/*size of count*/ 4 +
		this.Capacity*(this.KeySize+ /*size of offset*/ 4)
}

// read reads the file header from reader.
// It retuns a slice of offsets of empty nodes and nil or nil and an error.
func (this *fileHeader) read(reader io.ReadSeeker) ([]int64, error) {
	if _, err := reader.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, bo, this); err != nil {
		return nil, err
	}
	if this.EmptyNodes == -1 {
		return nil, nil
	}
	if _, err := reader.Seek(this.EmptyNodes, os.SEEK_SET); err != nil {
		return nil, err
	}
	var count int32
	if err := binary.Read(reader, bo, &count); err != nil {
		return nil, err
	}
	e := make([]byte, count*8)
	if _, err := io.ReadFull(reader, e); err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(e)
	empty := make([]int64, count)
	for i := 0; i < int(count); i++ {
		if err := binary.Read(b, bo, &empty[i]); err != nil {
			return nil, err
		}
	}
	return empty, nil
}

// read writes the file header to writer
// It retuns nil on success or an error
func (this *fileHeader) write(writer io.WriteSeeker, empty []int64) error {
	if len(empty) != 0 {
		if this.EmptyNodes == -1 {
			if off, err := writer.Seek(0, os.SEEK_END); err != nil {
				return err
			} else {
				this.EmptyNodes = off
			}
		} else {
			if _, err := writer.Seek(this.EmptyNodes, os.SEEK_SET); err != nil {
				return err
			}
		}
		var count int32 = int32(len(empty))
		if err := binary.Write(writer, bo, count); err != nil {
			return err
		}
		e := empty[:cap(empty)]
		b := bytes.NewBuffer(nil)
		for i := 0; i < len(e); i++ {
			if err := binary.Write(b, bo, e[i]); err != nil {
				return err
			}
		}
		if _, err := b.WriteTo(writer); err != nil {
			return err
		}
	}

	if _, err := writer.Seek(0, os.SEEK_SET); err != nil {
		return err
	}
	return binary.Write(writer, bo, this)
}

// init inits the node with default values.
func (this *node) init(tree *BTree) {
	this.size = int(tree.header.KeySize) + /*size of offset*/ 4
	this.offset = -1
	this.count = 0
	this.raw = make([]byte, tree.header.bufferSize(), tree.header.bufferSize()+uint32(this.size))
	this.datas = this.raw[8:]
	copy(this.raw[0:4], []byte{0xff, 0xff, 0xff, 0xff})
}

// find finds a key in the node.
// It returns an index of the key in the node and -1 or -1 and an index of a nearest key is less of the key.
func (this node) find(key Key) (int, []byte, int, error) {
	min := 1
	max := int(this.count)
	for max >= min {
		i := (max + min) / 2
		b := this.getKey(i - 1)
		r, err := key.Compare(b)
		if err != nil {
			return -1, nil, -1, err
		}
		switch {
		case r < 0:
			max = i - 1
		case r == 0:
			return i - 1, b, -1, nil
		case r > 0:
			min = i + 1
		}
	}
	return -1, nil, max - 1, nil
}

// getOffset returns an offset is got from the datas by index i
func (this *node) getOffset(i int) int64 {
	var offset int32 = int32(bo.Uint32(this.datas[i*this.size : i*this.size+4]))
	return int64(offset)
}

// setOffset writes an offset to the datas by index i
func (this *node) setOffset(i int, off int64) {
	var offset int32 = int32(off)
	bo.PutUint32(this.datas[i*this.size:i*this.size+4], uint32(offset))
}

// getKey returns a key data is got from the datas by index i. It is using k like a key sample
func (this *node) getKey(i int) []byte {
	return this.datas[i*this.size+4 : (i+1)*this.size]
}

// setKey writes a key to the datas by index i
func (this *node) setKey(i int, k Key) error {
	return k.Write(this.datas[i*this.size+4 : (i+1)*this.size])

}

// getLeast returns an offset of the least node
func (this *node) getLeast() int64 {
	var least int32 = int32(bo.Uint32(this.raw[0:4]))
	return int64(least)
}

// setLeast sets an offset of the least node
func (this *node) setLeast(l int64) {
	var least int32 = int32(l)
	bo.PutUint32(this.raw[0:4], uint32(least))
}

// insert inserts a key in the node after index position with offset.
// It returns a pointer to the filled offset.
func (this *node) insert(index int, key Key, offset int64) ([]byte, error) {
	this.count++
	this.datas = this.raw[8 : 8+int(this.count)*this.size]
	if index+2 < int(this.count) {
		copy(this.datas[(index+2)*this.size:], this.datas[(index+1)*this.size:])
	}
	this.setOffset(index+1, offset)
	if err := this.setKey(index+1, key); err != nil {
		return nil, err
	}
	return this.datas[(index+1)*this.size : (index+1)*this.size+4], nil
}

// read reads the node from reader with keys of type t and sets a node offset to offset.
// It retuns nil on success or an error.
func (this *node) read(reader io.ReadSeeker, offset int64) error {
	if _, err := reader.Seek(offset, os.SEEK_SET); err != nil {
		return err
	}
	if _, err := io.ReadFull(reader, this.raw); err != nil {
		return err
	}
	this.count = int(bo.Uint32(this.raw[4:]))
	this.offset = offset
	this.datas = this.raw[8 : 8+int(this.count)*this.size]
	return nil
}

// write writes the node to writer with keys of type t.
// It retuns nil on success or an error.
func (this *node) write(writer io.Writer) error {
	bo.PutUint32(this.raw[4:8], uint32(this.count))
	_, err := writer.Write(this.raw)
	return err
}
