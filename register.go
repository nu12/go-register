package register

type Register struct {
	Step int
	Err  error
}

func New() *Register {
	return &Register{Step: 0, Err: nil}
}

func (r *Register) Error() error {
	return r.Err
}

func (r *Register) Run(fo ...func()) *Register {
	for _, fi := range fo {
		if r.Err != nil {
			return r
		}
		r.Step++
		fi()
	}
	return r
}

func (r *Register) If(b bool, fo ...func()) *Register {
	if r.Err != nil {
		return r
	}

	if b {
		return r.Run(fo...)
	}

	return r
}

func (r *Register) IfError(f func()) *Register {
	if r.Err != nil {
		f()
	}
	return r
}
