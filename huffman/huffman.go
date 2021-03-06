/*
 * @Author: zhanghao
 * @Date: 2018-11-20 15:28:57
 * @Last Modified by: zhanghao
 * @Last Modified time: 2018-11-27 02:09:05
 */

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	huffman := NewHuffman("./huffman.txt")
	huffman.printWeights()
	fmt.Println()
	huffman.printEncodeMap()
	fmt.Println()
	fmt.Println("用0-1字符串表示文件编码结果(用0-1表示bit位):")
	fmt.Println(huffman.encodedBinStr())
	fmt.Println()
	fmt.Printf("压缩比例: %f", huffman.compressRate)
	fmt.Println()
}

type CHuffman struct {
	filename     string
	weights      map[byte]uint
	tree         *binaryTree
	encodeMap    map[byte]string
	compressRate float64
}

func NewHuffman(filename string) *CHuffman {
	wghts := getASICCWeights(filename)
	t := generateHuffmantree(wghts)
	enmap := make(map[byte]string)
	getEncodeMap(t, "", enmap)
	return &CHuffman{
		filename:  filename,
		weights:   wghts,
		tree:      t,
		encodeMap: enmap,
	}
}

func (hm *CHuffman) printWeights() {
	fmt.Println("asicc码及出现次数: ")
	for k, _ := range hm.weights {
		fmt.Print("["+string(k)+":", hm.weights[k], "]")
	}
	fmt.Println()
}

func (hm *CHuffman) printEncodeMap() {
	fmt.Println("Huffman编码对照表(用0-1表示bit位): ")
	for k, _ := range hm.encodeMap {
		fmt.Print("["+string(k)+":", hm.encodeMap[k], "]  ")
	}
	fmt.Println()
}

func (hm *CHuffman) encodedBinStr() string {
	fileBytes, err := ioutil.ReadFile(hm.filename)
	if err != nil {
		panic(err)
	}
	var encodedBinStr string
	for i, _ := range fileBytes {
		encodedBinStr += hm.encodeMap[fileBytes[i]]
	}
	var finfo os.FileInfo
	finfo, err = os.Stat(hm.filename)
	hm.compressRate = float64(len(encodedBinStr)) / float64((finfo.Size() * 8))
	return encodedBinStr
}

func getEncodeMap(ht *binaryTree, code string, m map[byte]string) {
	if ht.left == nil {
		m[ht.c] = code
		return
	}
	getEncodeMap(ht.left, code+"0", m)
	getEncodeMap(ht.right, code+"1", m)
}

func getASICCWeights(filename string) (weights map[byte]uint) {
	weights = make(map[byte]uint)
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	for i, _ := range fileBytes {
		weights[fileBytes[i]]++
	}
	return weights
}

type binaryTree struct {
	c      byte
	weight uint
	left   *binaryTree
	right  *binaryTree
}

// 创建huffman tree.
func generateHuffmantree(weights map[byte]uint) *binaryTree {
	// 创建森林 []*binaryTree
	forest := make([]*binaryTree, len(weights))
	i := 0
	for k, _ := range weights {
		forest[i] = &binaryTree{
			c:      k,
			weight: weights[k],
		}
		i++
	}

	minIndex := [2]byte{}
	minBT := [2]*binaryTree{}
	finish := false
	for !finish {
		j := 0
		for ; j < 2; j++ {
			for i, _ := range forest {
				if forest[i] != nil {
					if minBT[j] == nil {
						minBT[j] = forest[i]
						minIndex[j] = byte(i)
					} else if forest[i].weight < minBT[j].weight {
						minBT[j] = forest[i]
						minIndex[j] = byte(i)
					}
				}
			}
			if j == 0 {
				forest[minIndex[0]] = nil
			} else if minBT[1] == nil {
				finish = true
			} else {
				forest[minIndex[1]] = &binaryTree{
					weight: minBT[0].weight + minBT[1].weight,
					left:   minBT[0],
					right:  minBT[1],
				}
				minBT[0], minBT[1] = nil, nil
			}
		}
	}
	return minBT[0]
}
