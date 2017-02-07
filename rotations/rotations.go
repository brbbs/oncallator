package rotations

import(
	"fmt"
	"time"
)

type Rotation struct {
	Start time.Time
	Primary string
	Secondary string
}

func (r Rotation) String() string {
	return fmt.Sprintf("%s %s %s", r.Start.Format(time.RFC3339), r.Primary, r.Secondary)
}

type Rotations []Rotation

type TerraformLayer struct {
	Name string `json:"name"`
	Start string `json:"start"`
	Users []string `json:"users"`
}

type TerraformLayers struct {
	Primary []TerraformLayer
	Secondary []TerraformLayer
}

func (r Rotations) TerraformLayers() TerraformLayers {
	l := TerraformLayers{}
	for _, rot := range r {
		start := rot.Start.Format(time.RFC3339)
		// TODO(brb): do we _need_ to supply |end|?
		primary := TerraformLayer{
			Name: "Primary",
			Start: start,
			Users: []string{rot.Primary},
		}
		l.Primary = append(l.Primary, primary)
		secondary := TerraformLayer{
			Name: "Secondary",
			Start: start,
			Users: []string{rot.Secondary},
		}
		l.Secondary = append(l.Secondary, secondary)
	}
	return l
}

func New(users []string, start time.Time, duration time.Duration, n int) Rotations {
	var rots Rotations

	if len(users) == 0 {
		return rots
	}

	for i := 0; i < n; i++ {
		primary := users[i % len(users)]
		secondary := users[(i+1) % len(users)]
		rot := Rotation{
			Start: start,
			Primary: primary,
			Secondary: secondary,
		}
		rots = append(rots, rot)
		start = start.Add(duration)
	}
	return rots
}

func Num(start, end time.Time, duration time.Duration) int {
	length := end.Sub(start)
	if length <= 0 {
		return 0
	}

	rotations := int(length / duration)
	if length % duration > 0 {
		rotations += 1
	}
	return rotations
}
