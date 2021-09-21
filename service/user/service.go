package user

import (
	"context"
	"encoding/json"
	"github.com/stasBigunenko/monorepa/pkg/storage/newStorage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"unicode"

	"github.com/google/uuid"
	"github.com/stasBigunenko/monorepa/model"
)

type UsrService struct {
	Us      model.UserHTTP
	storage newStorage.NewStore
}

func NewUsrService(s newStorage.NewStore) *UsrService {
	return &UsrService{
		storage: s,
	}
}

func (u *UsrService) Get(_ context.Context, id uuid.UUID) (model.UserHTTP, error) {

	res, err := u.storage.Get(context.Background(), id)
	if res == nil && err == nil || err != nil {
		return model.UserHTTP{}, status.Error(codes.Internal, "internal problem")
	}

	err = json.Unmarshal(res, &u.Us)
	if err != nil {
		return model.UserHTTP{}, status.Error(codes.Aborted, "wrong with unmarshal")
	}

	return u.Us, nil
}

func (u *UsrService) GetAll(_ context.Context) ([]model.UserHTTP, error) {

	res, err := u.storage.GetAll(context.Background())
	if err != nil {
		return nil, status.Error(codes.Internal, "storage problem")
	}

	ac := []model.UserHTTP{}

	for _, val := range res {
		err := json.Unmarshal(val, &u.Us)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal problem")
		}

		ac = append(ac, u.Us)
	}
	return ac, nil
}

func (u *UsrService) Create(_ context.Context, b string) (model.UserHTTP, error) {

	for _, r := range b {
		if !unicode.IsLetter(r) {
			return model.UserHTTP{}, status.Error(codes.InvalidArgument, "invalid username")
		}
	}

	if len(b) <= 2 {
		return model.UserHTTP{}, status.Error(codes.InvalidArgument, "invalid name")
	}

	bt, err := json.Marshal(b)
	if err != nil {
		return model.UserHTTP{}, status.Error(codes.Internal, "internal problem")
	}

	res, id, err := u.storage.Create(context.Background(), bt)
	if err != nil {
		return model.UserHTTP{}, status.Error(codes.Internal, "internal problem")
	}

	err = json.Unmarshal(res, &u.Us)
	if err != nil {
		return model.UserHTTP{}, status.Error(codes.Internal, "internal problem")
	}

	u.Us.ID = id

	return u.Us, nil
}

func (u *UsrService) Update(_ context.Context, user model.UserHTTP) (model.UserHTTP, error) {
	id := user.ID

	bt, err := json.Marshal(user)
	if err != nil {
		return model.UserHTTP{}, status.Error(codes.Internal, "internal problem")
	}

	res, err := u.storage.Update(context.Background(), id, bt)
	if err != nil || res == nil && err == nil {
		return model.UserHTTP{}, status.Error(codes.Internal, "internal problem")
	}

	err = json.Unmarshal(res, &u.Us)
	if err != nil {
		return model.UserHTTP{}, status.Error(codes.Internal, "internal problem")
	}

	return u.Us, nil
}

func (u *UsrService) Delete(_ context.Context, id uuid.UUID) error {
	b, err := u.storage.Delete(context.Background(), id)
	if err != nil || !b {
		return status.Error(codes.Internal, "not found")
	}

	return nil
}
