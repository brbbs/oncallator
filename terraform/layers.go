package terraform

import (
	"time"

	"github.com/websdev/oncallator/schedule"
)

type Layer struct {
	Name string `json:"name"`
	Start string `json:"start"`
	Users []string `json:"users"`
}

type Layers struct {
	Primary []Layer
	Secondary []Layer
}

func NewLayers(rs []schedule.Rotation) Layers {
	l := Layers{}
	for _, r := range rs {
		start := r.Start.Format(time.RFC3339)
		// TODO(brb): do we _need_ to supply |end|?
		primary := Layer{
			Name: "Primary",
			Start: start,
			Users: []string{r.Primary},
		}
		l.Primary = append(l.Primary, primary)
		secondary := Layer{
			Name: "Secondary",
			Start: start,
			Users: []string{r.Secondary},
		}
		l.Secondary = append(l.Secondary, secondary)
	}
	return l
}
