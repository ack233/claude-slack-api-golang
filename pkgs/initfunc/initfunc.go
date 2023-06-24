package initfunc

type initFunc func()

var mfunc = []initFunc{}

func InitFun() {
	for _, f := range mfunc {
		f()
	}
}

func RegisterInitFunc(f ...initFunc) {
	mfunc = append(mfunc, f...)
}
