package test

import (
	"VideoWeb/Utilities"
	"VideoWeb/logic"
	"fmt"
	"path"
	"testing"
)

func TestVideoTime(t *testing.T) {

	time, _ := logic.GetVideoDuration(
		"/home/zey/ZeyGO/project/VideoWeb/resources/Videos/2024-04-07T011915.avi")
	//if err != nil {
	//	t.Error(err)
	//}
	fmt.Println(Utilities.SecondToTime(time))
}

func TestPath(t *testing.T) {
	p := "/home/zey/ZeyGO/project/VideoWeb/resources/Videos/2024-04-07T011915.avi"
	fmt.Println(path.Dir(p))
}
