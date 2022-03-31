package model

type Operator[M any] func(M)

type Provider[M any] func() (M, error)

type SliceOperator[M any] func([]M)

type SliceProvider[M any] func() ([]M, error)

func ExecuteForEach[M any](f Operator[M]) SliceOperator[M] {
	return func(models []M) {
		for _, m := range models {
			f(m)
		}
	}
}

type Filter[M any] func(M) bool
