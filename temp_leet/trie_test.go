package temp_leet

import (
	"math"
	"testing"
)

func TestReSpace(t *testing.T) {

}

func respace(dictionary []string, sentence string) int {
	//这个题目我的思路: 我想到了递归(分而治之)去解题, 也想到了动态规划方法;
	//  想到了dp如何去定义状态,但是状态转移方程没有想出来, 问题就出在我想到了
	//  dp[i]应该由dp[i-1]或dp的前几项去求解, 但是缺乏思考怎么去由dp[i]的前项来推导dp[i]

	//思路: 字典树+动态规划
	//字典树: 这个数据结构相比于map, 在比较word的前缀时, 一旦遇到不存在trie中的树的边时, 会及时终止, 而不用把word遍历完;
	//dp状态转移: min{dp[i-1]+1, dp[j-1](j->i的下标在trie中有匹配的分支)}
	trie := Trie{
		next:  [26]*Trie{},
		isEnd: false,
	}
	//初始化字典树
	for _, value := range dictionary {
		trie.insert(value)
	}
	//定义dp
	dp := make([]int, len(sentence)+1)
	for i := 1; i < len(dp); i++ {
		curNode := &trie
		dp[i] = dp[i-1] + 1
		for j := i; j > 0; j-- {
			pos := int(sentence[j] - 'a')
			if curNode.next[pos] == nil {
				break
			} else if curNode.next[pos].isEnd {
				dp[i] = int(math.Min(float64(dp[i]), float64(dp[j-1])))
			}
			if dp[i] == 0 {
				break
			}
			curNode = curNode.next[pos]
		}
	}
	return dp[len(sentence)]
}

type Trie struct {
	next [26]*Trie //next指针数组的长度, 由charset大小决定
	isEnd bool //判断是否为
}

//str倒序插入trie中
func (this *Trie)insert(str string) {
	curNode := this
	for j := len(str) -1; j >= 0; j-- {
		pos := int(str[j]- 'a')
		if curNode.next[pos] == nil {
			curNode.next[pos] = &Trie{
				next:  [26]*Trie{},
				isEnd: false,
			}
		}
		curNode = curNode.next[pos]
	}
	curNode.isEnd = true
}