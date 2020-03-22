package addStaff

import "context"

type Repository interface {
	Add(ctx context.Context, uuid string, id int) error
}
