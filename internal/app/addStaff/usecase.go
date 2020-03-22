package addStaff

import "context"

type Usecase interface {
	GetQrForStaff(ctx context.Context, idCafe int) (string, error)
}
