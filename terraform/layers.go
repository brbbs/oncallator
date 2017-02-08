package terraform

import (
	"time"

	"github.com/websdev/oncallator/schedule"
)

type Layer struct {
	Start string `json:"start"`
	Users []string `json:"users"`
	RotationVirtualStart string `json:"rotation_virtual_start"`
	RotationTurnLengthSeconds int `json:"rotation_turn_length_seconds"`
}

type Layers struct {
	Primary []Layer
	Secondary []Layer
}

func NewLayers(s *schedule.Schedule) Layers {
	l := Layers{}
	for _, r := range s.Rotations {
		start := r.Start.Format(time.RFC3339)
		primary := Layer{
			Start: start,
			Users: []string{r.Primary},
			RotationVirtualStart: start,
			RotationTurnLengthSeconds: int(s.RotationDuration.Seconds()),
		}
		l.Primary = append(l.Primary, primary)
		secondary := Layer{
			Start: start,
			Users: []string{r.Secondary},
			RotationVirtualStart: start,
			RotationTurnLengthSeconds: int(s.RotationDuration.Seconds()),
		}
		l.Secondary = append(l.Secondary, secondary)
	}
	return l
}
