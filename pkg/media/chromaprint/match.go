package chromaprint

import (
	"fmt"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Match struct {
	Id    string  `json:"id"`
	Score float64 `json:"score"`
	//	Recordings []recording `json:"recording"`
}

/*
type recording struct {
	Id            string         `json:"id"`
	Title         string         `json:"title"`
	Duration      float64        `json:"duration"`
	Artists       []artist       `json:"artists"`
	ReleaseGroups []releasegroup `json:"releasegroups"`
}
*/

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Match) String() string {
	str := "<chromaprint.match"
	if id := this.Id; id != "" {
		str += " id=" + strconv.Quote(id)
	}
	if score := this.Score; score != 0.0 {
		str += " score=" + fmt.Sprint(score)
	}
	return str + ">"
}
