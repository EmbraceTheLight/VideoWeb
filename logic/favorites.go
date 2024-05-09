package logic

func String2Int8(IsPrivate string) int8 {
	var ret int8
	if IsPrivate == "公开" {
		ret = 1
	} else if IsPrivate == "私密" {
		ret = -1
	} else { //处理IsPrivate为空的情况（该情况可能出现在修改收藏夹时没有修改私密性的情况）
		ret = 0
	}
	return ret
}
