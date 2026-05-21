package db

type RowTransformable[T any, R any] interface {
	*T
	FromRow(row R)
}

func TransformRow[T any, R any, PtrT RowTransformable[T, R]](row R) T {
	var next T
	PtrT(&next).FromRow(row)
	return next
}

type BatchRowTransformable[T any, R any] interface {
	*T
	FromBatchRow(row R)
}

func TransformBatchRow[T any, R any, PtrT BatchRowTransformable[T, R]](row R) T {
	var next T
	PtrT(&next).FromBatchRow(row)
	return next
}
