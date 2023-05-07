package bloom

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
)

// 错误率推导见https://blog.csdn.net/gaoyueace/article/details/90410735
type bloomFilter struct {
	k int64 // hash函数个数,暂时限定8以内
	m int64 // slot个数

	array    []byte // slot数组
	hashFunc hash.Hash
}

func NewBloomFilter(k, m int64) *bloomFilter {
	b := &bloomFilter{k: k, m: m}
	b.array = make([]byte, m/8)
	b.hashFunc = sha256.New()
	return b
}

func (b *bloomFilter) Insert(elem []byte) error {
	hashVals, err := b.hashToUInt64Array(elem)
	if err != nil {
		fmt.Printf("err:%v", err)
		return err
	}
	fmt.Printf("insert, hashvals:%v\n", hashVals)

	for k, v := range hashVals {
		b.setBit(v)
		fmt.Printf("set bit, k:%v, v:%v\n", k, v)
	}

	return nil
}

// 元素是否在集合中
func (b *bloomFilter) Exist(elem []byte) (bool, error) {
	hashVals, err := b.hashToUInt64Array(elem)
	if err != nil {
		fmt.Printf("err:%v", err)
		return false, err
	}

	for k, v := range hashVals {
		if !b.checkBit(v) {
			fmt.Printf("not exist, k:%v, v:%v\n", k, v)
			return false, nil
		}
	}

	fmt.Printf("exist\n")

	return true, nil
}

func (b *bloomFilter) hashToUInt64Array(input []byte) ([]int64, error) {
	// 计算 SHA256 哈希值
	h := sha256.Sum256(input)

	// 按序 2 个 byte 为一组，取前 14 组，并转换为 int64 类型
	var res []int64
	for i := 0; i < int(b.k); i++ {
		// 取出 2 个 byte
		b1 := h[i*2]
		b2 := h[i*2+1]

		// 将 2 个 byte 组成一个 uint16 类型的值
		u := binary.LittleEndian.Uint16([]byte{b1, b2})

		// 将 uint16 类型的值转换为 int64 类型
		res = append(res, int64(u)%b.m)
	}

	return res, nil
}

func (b *bloomFilter) checkBit(pos int64) bool {
	byteIdx := pos / 8
	bitIdx := pos % 8

	val := b.array[byteIdx]

	mask := byte(1 << bitIdx)
	return val&mask != 0
}

func (b *bloomFilter) setBit(pos int64) {
	byteIdx := pos / 8
	bitIdx := pos % 8

	mask := byte(1 << bitIdx)
	b.array[byteIdx] |= mask
}
