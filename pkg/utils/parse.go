package utils

import (
    "strings"
    "strconv"
)

func ParseRange(str string) [][]int {
    var ranges[][]int
    str = strings.Replace(str, " ", "", -1)
    str = str[0:len(str)-1]
    res := strings.Split(str, ",")
    for _, r := range res {
        if strings.Contains(r, "-") {
            s := strings.Split(r, "-")
            start, _ := strconv.Atoi(s[0])
            end, _ := strconv.Atoi(s[1])
            ranges = append(ranges, []int{start, end})
        } else {
            start, _ := strconv.Atoi(r)
            ranges = append(ranges, []int{start, start})
        }
    }
    return ranges
}

func InsideRange(ranges [][]int, val int) bool {
    for _, r := range ranges {
        if val >= r[0] && val <= r[1] {
            return true
        }
    }
    return false
}
