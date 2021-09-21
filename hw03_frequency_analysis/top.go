package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type Pair struct {
	key string
	val int
}

type SortPair struct {
	key []string
	val int
}

var zp = regexp.MustCompile(`[\s!.,?]+`)

func Top10(text string) []string {
	dictionary := make(map[string]int)
	for _, elem := range zp.Split(text, -1) {
		elem = strings.ToLower(elem)
		_, found := dictionary[elem]
		if found {
			dictionary[elem]++
		} else if elem != "-" {
			dictionary[elem] = 1
		}
	}

	dicArr := make([]Pair, 0, len(dictionary)*2)

	for key, val := range dictionary {
		dicArr = append(dicArr, Pair{key, val})
	}

	sort.Slice(dicArr, func(i, j int) bool {
		return dicArr[i].val > dicArr[j].val
	})

	sortt := make(map[int][]string)

	for _, val := range dicArr {
		sortt[val.val] = append(sortt[val.val], val.key)
	}

	for key := range sortt {
		sort.Slice(sortt[key], func(i, j int) bool {
			return sortt[key][i] < sortt[key][j]
		})
	}

	sortedArr := make([]SortPair, 0, len(sortt)*2)

	for key, val := range sortt {
		sortedArr = append(sortedArr, SortPair{val, key})
	}

	sort.Slice(sortedArr, func(i, j int) bool {
		return sortedArr[i].val > sortedArr[j].val
	})

	result := make([]string, 0, 10)

	for _, elem := range sortedArr {
		result = append(result, elem.key...)
	}

	if len(result) > 10 {
		return result[:10]
	}

	return result[0:0:0]
}
