package mas
import "time"

type Task struct {
	Difficulty  int
}


func (t Task) processTime(agents ...*Agent) int {
	abilitySum := 0
	for _, a := range agents {
		abilitySum += a.Ability
	}
	return 100 * t.Difficulty * int(time.Millisecond) / abilitySum / abilitySum
}