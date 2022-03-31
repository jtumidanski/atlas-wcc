package model

type Operator[M any] func(M)

type Provider[M any] func() (M, error)

type SliceOperator[M any] func([]M)

type SliceProvider[M any] func() ([]M, error)

type PreciselyOneFilter[M any] func([]M) (M, error)

func ExecuteForEach[M any](f Operator[M]) SliceOperator[M] {
	return func(models []M) {
		for _, m := range models {
			f(m)
		}
	}
}

type Filter[M any] func(M) bool

func ModelListProviderToModelProviderAdapter[M any](provider SliceProvider[M], preciselyOneFilter PreciselyOneFilter[M]) Provider[M] {
	return func() (M, error) {
		ps, err := provider()
		if err != nil {
			var result M
			return result, err
		}
		return preciselyOneFilter(ps)
	}
}
