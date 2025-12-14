package testutil

import (
	"context"

	nasDto "github.com/novriyantoAli/freeradius-service/internal/application/nas/dto"
	nasEntity "github.com/novriyantoAli/freeradius-service/internal/application/nas/entity"
	paymentDto "github.com/novriyantoAli/freeradius-service/internal/application/payment/dto"
	paymentEntity "github.com/novriyantoAli/freeradius-service/internal/application/payment/entity"
	radcheckDto "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	radcheckEntity "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	radreplyDto "github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	radreplyEntity "github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	userDto "github.com/novriyantoAli/freeradius-service/internal/application/user/dto"
	userEntity "github.com/novriyantoAli/freeradius-service/internal/application/user/entity"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *userEntity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*userEntity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*userEntity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(filter *userDto.UserFilter) ([]userEntity.User, int64, error) {
	args := m.Called(filter)
	var users []userEntity.User
	if args.Get(0) != nil {
		users = args.Get(0).([]userEntity.User)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return users, count, args.Error(2)
}

func (m *MockUserRepository) Update(user *userEntity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) EmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

// MockPaymentRepository is a mock implementation of PaymentRepository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Create(payment *paymentEntity.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByID(id uint) (*paymentEntity.Payment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paymentEntity.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetAll(filter *paymentDto.PaymentFilter) ([]paymentEntity.Payment, int64, error) {
	args := m.Called(filter)
	var payments []paymentEntity.Payment
	if args.Get(0) != nil {
		payments = args.Get(0).([]paymentEntity.Payment)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return payments, count, args.Error(2)
}

func (m *MockPaymentRepository) Update(payment *paymentEntity.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetByUserID(userID uint) ([]paymentEntity.Payment, error) {
	args := m.Called(userID)
	var payments []paymentEntity.Payment
	if args.Get(0) != nil {
		payments = args.Get(0).([]paymentEntity.Payment)
	}
	return payments, args.Error(1)
}

// MockNASRepository is a mock implementation of NASRepository
type MockNASRepository struct {
	mock.Mock
}

func (m *MockNASRepository) Create(nas *nasEntity.NAS) error {
	args := m.Called(nas)
	return args.Error(0)
}

func (m *MockNASRepository) GetByID(id uint) (*nasEntity.NAS, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nasEntity.NAS), args.Error(1)
}

func (m *MockNASRepository) GetByNASName(nasname string) (*nasEntity.NAS, error) {
	args := m.Called(nasname)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nasEntity.NAS), args.Error(1)
}

func (m *MockNASRepository) GetAll(filter *nasDto.NASFilter) ([]nasEntity.NAS, int64, error) {
	args := m.Called(filter)
	var nasList []nasEntity.NAS
	if args.Get(0) != nil {
		nasList = args.Get(0).([]nasEntity.NAS)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return nasList, count, args.Error(2)
}

func (m *MockNASRepository) Update(nas *nasEntity.NAS) error {
	args := m.Called(nas)
	return args.Error(0)
}

func (m *MockNASRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(req *userDto.CreateUserRequest) (*userDto.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(id uint) (*userDto.UserResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*userDto.UserResponse, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUsers(filter *userDto.UserFilter) (*userDto.UserListResponse, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserListResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(id uint, req *userDto.UpdateUserRequest) (*userDto.UserResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDto.UserResponse), args.Error(1)
}

func (m *MockUserService) UpdateUserPassword(id uint, req *userDto.UpdateUserPasswordRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockNASService is a mock implementation of NASService
type MockNASService struct {
	mock.Mock
}

func (m *MockNASService) CreateNAS(req *nasDto.CreateNASRequest) (*nasDto.NASResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nasDto.NASResponse), args.Error(1)
}

func (m *MockNASService) GetNASByID(id uint) (*nasDto.NASResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nasDto.NASResponse), args.Error(1)
}

func (m *MockNASService) ListNAS(filter *nasDto.NASFilter) (*nasDto.ListNASResponse, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nasDto.ListNASResponse), args.Error(1)
}

func (m *MockNASService) UpdateNAS(id uint, req *nasDto.UpdateNASRequest) (*nasDto.NASResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*nasDto.NASResponse), args.Error(1)
}

func (m *MockNASService) DeleteNAS(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockRadcheckRepository is a mock implementation of RadcheckRepository
type MockRadcheckRepository struct {
	mock.Mock
}

func (m *MockRadcheckRepository) Create(ctx context.Context, radcheck *radcheckEntity.Radcheck) error {
	args := m.Called(ctx, radcheck)
	return args.Error(0)
}

func (m *MockRadcheckRepository) GetByID(ctx context.Context, id uint) (*radcheckEntity.Radcheck, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckEntity.Radcheck), args.Error(1)
}

func (m *MockRadcheckRepository) GetByUsernameAndAttribute(ctx context.Context, username, attribute string) (*radcheckEntity.Radcheck, error) {
	args := m.Called(ctx, username, attribute)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckEntity.Radcheck), args.Error(1)
}

func (m *MockRadcheckRepository) GetAll(ctx context.Context, filter *radcheckDto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
	args := m.Called(ctx, filter)
	var radchecks []radcheckEntity.Radcheck
	if args.Get(0) != nil {
		radchecks = args.Get(0).([]radcheckEntity.Radcheck)
	}

	var count int64
	if args.Get(1) != nil {
		count = args.Get(1).(int64)
	}
	return radchecks, count, args.Error(2)
}

func (m *MockRadcheckRepository) Update(ctx context.Context, radcheck *radcheckEntity.Radcheck) error {
	args := m.Called(ctx, radcheck)
	return args.Error(0)
}

func (m *MockRadcheckRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockRadcheckRepositoryWithFn is a mock implementation of RadcheckRepository with function fields
type MockRadcheckRepositoryWithFn struct {
	CreateFn                    func(context.Context, *radcheckEntity.Radcheck) error
	GetByIDFn                   func(context.Context, uint) (*radcheckEntity.Radcheck, error)
	GetByUsernameAndAttributeFn func(context.Context, string, string) (*radcheckEntity.Radcheck, error)
	GetAllFn                    func(context.Context, *radcheckDto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error)
	UpdateFn                    func(context.Context, *radcheckEntity.Radcheck) error
	DeleteFn                    func(context.Context, uint) error
}

func NewMockRadcheckRepositoryWithFn() *MockRadcheckRepositoryWithFn {
	return &MockRadcheckRepositoryWithFn{}
}

func (m *MockRadcheckRepositoryWithFn) Create(ctx context.Context, radcheck *radcheckEntity.Radcheck) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, radcheck)
	}
	radcheck.ID = 1
	return nil
}

func (m *MockRadcheckRepositoryWithFn) GetByID(ctx context.Context, id uint) (*radcheckEntity.Radcheck, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return CreateRadcheckFixture(), nil
}

func (m *MockRadcheckRepositoryWithFn) GetByUsernameAndAttribute(ctx context.Context, username, attribute string) (*radcheckEntity.Radcheck, error) {
	if m.GetByUsernameAndAttributeFn != nil {
		return m.GetByUsernameAndAttributeFn(ctx, username, attribute)
	}
	return CreateRadcheckFixture(), nil
}

func (m *MockRadcheckRepositoryWithFn) GetAll(ctx context.Context, filter *radcheckDto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx, filter)
	}
	return []radcheckEntity.Radcheck{*CreateRadcheckFixture()}, 1, nil
}

func (m *MockRadcheckRepositoryWithFn) Update(ctx context.Context, radcheck *radcheckEntity.Radcheck) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, radcheck)
	}
	return nil
}

func (m *MockRadcheckRepositoryWithFn) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

// MockRadcheckService is a mock implementation of RadcheckService
type MockRadcheckService struct {
	mock.Mock
}

func (m *MockRadcheckService) CreateRadcheck(ctx context.Context, req *radcheckDto.CreateRadcheckRequest) (*radcheckDto.RadcheckResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckDto.RadcheckResponse), args.Error(1)
}

func (m *MockRadcheckService) GetRadcheckByID(ctx context.Context, id uint) (*radcheckDto.RadcheckResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckDto.RadcheckResponse), args.Error(1)
}

func (m *MockRadcheckService) GetRadcheckByUsernameAndAttribute(ctx context.Context, username, attribute string) (*radcheckDto.RadcheckResponse, error) {
	args := m.Called(ctx, username, attribute)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckDto.RadcheckResponse), args.Error(1)
}

func (m *MockRadcheckService) ListRadcheck(ctx context.Context, filter *radcheckDto.RadcheckFilter) (*radcheckDto.ListRadcheckResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckDto.ListRadcheckResponse), args.Error(1)
}

func (m *MockRadcheckService) UpdateRadcheck(ctx context.Context, id uint, req *radcheckDto.UpdateRadcheckRequest) (*radcheckDto.RadcheckResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckDto.RadcheckResponse), args.Error(1)
}

func (m *MockRadcheckService) DeleteRadcheck(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockRadreplyRepository is a mock implementation of RadreplyRepository with function fields
type MockRadreplyRepository struct {
	CreateFn                    func(context.Context, *radreplyEntity.Radreply) error
	GetByIDFn                   func(context.Context, uint) (*radreplyEntity.Radreply, error)
	GetByUsernameAndAttributeFn func(context.Context, string, string) (*radreplyEntity.Radreply, error)
	GetAllFn                    func(context.Context, *radreplyDto.RadreplyFilter) ([]radreplyEntity.Radreply, int64, error)
	UpdateFn                    func(context.Context, *radreplyEntity.Radreply) error
	DeleteFn                    func(context.Context, uint) error
}

func NewMockRadreplyRepository() *MockRadreplyRepository {
	return &MockRadreplyRepository{}
}

func (m *MockRadreplyRepository) Create(ctx context.Context, radreply *radreplyEntity.Radreply) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, radreply)
	}
	radreply.ID = 1
	return nil
}

func (m *MockRadreplyRepository) GetByID(ctx context.Context, id uint) (*radreplyEntity.Radreply, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return CreateRadreplyFixture(), nil
}

func (m *MockRadreplyRepository) GetByUsernameAndAttribute(ctx context.Context, username, attribute string) (*radreplyEntity.Radreply, error) {
	if m.GetByUsernameAndAttributeFn != nil {
		return m.GetByUsernameAndAttributeFn(ctx, username, attribute)
	}
	return CreateRadreplyFixture(), nil
}

func (m *MockRadreplyRepository) GetAll(ctx context.Context, filter *radreplyDto.RadreplyFilter) ([]radreplyEntity.Radreply, int64, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx, filter)
	}
	return []radreplyEntity.Radreply{*CreateRadreplyFixture()}, 1, nil
}

func (m *MockRadreplyRepository) Update(ctx context.Context, radreply *radreplyEntity.Radreply) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, radreply)
	}
	return nil
}

func (m *MockRadreplyRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

// MockRadreplyService is a mock implementation of RadreplyService with function fields
type MockRadreplyService struct {
	CreateRadreplyFn                    func(context.Context, *radreplyDto.CreateRadreplyRequest) (*radreplyDto.RadreplyResponse, error)
	GetRadreplyByIDFn                   func(context.Context, uint) (*radreplyDto.RadreplyResponse, error)
	GetRadreplyByUsernameAndAttributeFn func(context.Context, string, string) (*radreplyDto.RadreplyResponse, error)
	ListRadreplyFn                      func(context.Context, *radreplyDto.RadreplyFilter) (*radreplyDto.ListRadreplyResponse, error)
	UpdateRadreplyFn                    func(context.Context, uint, *radreplyDto.UpdateRadreplyRequest) (*radreplyDto.RadreplyResponse, error)
	DeleteRadreplyFn                    func(context.Context, uint) error
}

func NewMockRadreplyService() *MockRadreplyService {
	return &MockRadreplyService{}
}

func (m *MockRadreplyService) CreateRadreply(ctx context.Context, req *radreplyDto.CreateRadreplyRequest) (*radreplyDto.RadreplyResponse, error) {
	if m.CreateRadreplyFn != nil {
		return m.CreateRadreplyFn(ctx, req)
	}
	return &radreplyDto.RadreplyResponse{
		ID:        1,
		Username:  req.Username,
		Attribute: req.Attribute,
		Op:        req.Op,
		Value:     req.Value,
	}, nil
}

func (m *MockRadreplyService) GetRadreplyByID(ctx context.Context, id uint) (*radreplyDto.RadreplyResponse, error) {
	if m.GetRadreplyByIDFn != nil {
		return m.GetRadreplyByIDFn(ctx, id)
	}
	fixture := CreateRadreplyFixture()
	return &radreplyDto.RadreplyResponse{
		ID:        fixture.ID,
		Username:  fixture.Username,
		Attribute: fixture.Attribute,
		Op:        fixture.Op,
		Value:     fixture.Value,
	}, nil
}

func (m *MockRadreplyService) GetRadreplyByUsernameAndAttribute(ctx context.Context, username, attribute string) (*radreplyDto.RadreplyResponse, error) {
	if m.GetRadreplyByUsernameAndAttributeFn != nil {
		return m.GetRadreplyByUsernameAndAttributeFn(ctx, username, attribute)
	}
	fixture := CreateRadreplyFixture()
	return &radreplyDto.RadreplyResponse{
		ID:        fixture.ID,
		Username:  fixture.Username,
		Attribute: fixture.Attribute,
		Op:        fixture.Op,
		Value:     fixture.Value,
	}, nil
}

func (m *MockRadreplyService) ListRadreply(ctx context.Context, filter *radreplyDto.RadreplyFilter) (*radreplyDto.ListRadreplyResponse, error) {
	if m.ListRadreplyFn != nil {
		return m.ListRadreplyFn(ctx, filter)
	}
	fixture := CreateRadreplyFixture()
	return &radreplyDto.ListRadreplyResponse{
		Data: []radreplyDto.RadreplyResponse{
			{
				ID:        fixture.ID,
				Username:  fixture.Username,
				Attribute: fixture.Attribute,
				Op:        fixture.Op,
				Value:     fixture.Value,
			},
		},
		Total:     1,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		TotalPage: 1,
	}, nil
}

func (m *MockRadreplyService) UpdateRadreply(ctx context.Context, id uint, req *radreplyDto.UpdateRadreplyRequest) (*radreplyDto.RadreplyResponse, error) {
	if m.UpdateRadreplyFn != nil {
		return m.UpdateRadreplyFn(ctx, id, req)
	}
	return &radreplyDto.RadreplyResponse{
		ID:        id,
		Username:  req.Username,
		Attribute: req.Attribute,
		Op:        req.Op,
		Value:     req.Value,
	}, nil
}

func (m *MockRadreplyService) DeleteRadreply(ctx context.Context, id uint) error {
	if m.DeleteRadreplyFn != nil {
		return m.DeleteRadreplyFn(ctx, id)
	}
	return nil
}
