package runtime

// hack much?
type SystemActions struct {
	actions map[string][]func()
}

func (this *SystemActions) Finish(name string, f func()) *SystemActions {
	arr := this.actions[name]
	arr = append(arr, f)
	return this
}

func (this *SystemActions) Run(name string) {
	if arr, ok := this.actions[name]; ok {
		for _, cb := range arr {
			cb()
		}
	}
}
