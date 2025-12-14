package testutil

import (
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

func (m *MockRadcheckRepository) Create(radcheck *radcheckEntity.Radcheck) error {
	args := m.Called(radcheck)
	return args.Error(0)
}

func (m *MockRadcheckRepository) GetByID(id uint) (*radcheckEntity.Radcheck, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckEntity.Radcheck), args.Error(1)
}

func (m *MockRadcheckRepository) GetByUsernameAndAttribute(username, attribute string) (*radcheckEntity.Radcheck, error) {
	args := m.Called(username, attribute)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckEntity.Radcheck), args.Error(1)
}

func (m *MockRadcheckRepository) GetAll(filter *radcheckDto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
	args := m.Called(filter)
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

func (m *MockRadcheckRepository) Update(radcheck *radcheckEntity.Radcheck) error {
	args := m.Called(radcheck)
	return args.Error(0)
}

func (m *MockRadcheckRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockRadcheckService is a mock implementation of RadcheckService
type MockRadcheckService struct {
	mock.Mock
}

func (m *MockRadcheckService) CreateRadcheck(req *radcheckDto.CreateRadcheckRequest) (*radcheckDto.RadcheckResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckDto.RadcheckResponse), args.Error(1)
}

func (m *MockRadcheckService) GetRadcheckByID(id uint) (*radcheckDto.RadcheckResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckDto.RadcheckResponse), args.Error(1)
}

func (m *MockRadcheckService) GetRadcheckByUsernameAndAttribute(username, attribute string) (*radcheckDto.RadcheckResponse, error) {
	args := m.Called(username, attribute)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckDto.RadcheckResponse), args.Error(1)
}

func (m *MockRadcheckService) ListRadcheck(filter *radcheckDto.RadcheckFilter) (*radcheckDto.ListRadcheckResponse, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckDto.ListRadcheckResponse), args.Error(1)
}

func (m *MockRadcheckService) UpdateRadcheck(id uint, req *radcheckDto.UpdateRadcheckRequest) (*radcheckDto.RadcheckResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*radcheckDto.RadcheckResponse), args.Error(1)
}

func (m *MockRadcheckService) DeleteRadcheck(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockRadreplyRepository is a mock implementation of RadreplyRepository with function fields
type MockRadreplyRepository struct {
	CreateFn                    func(*radreplyEntity.Radreply) error
	GetByIDFn                   func(uint) (*radreplyEntity.Radreply, error)
	GetByUsernameAndAttributeFn func(string, string) (*radreplyEntity.Radreply, error)
	GetAllFn                    func(*radreplyDto.RadreplyFilter) ([]radreplyEntity.Radreply, int64, error)
	UpdateFn                    func(*radreplyEntity.Radreply) error
	DeleteFn                    func(uint) error
}

func NewMockRadreplyRepository() *MockRadreplyRepository {
	return &MockRadreplyRepository{}
}

func (m *MockRadreplyRepository) Create(radreply *radreplyEntity.Radreply) error {
	if m.CreateFn != nil {
		return m.CreateFn(radreply)
	}
	radreply.ID = 1
	return nil
}

func (m *MockRadreplyRepository) GetByID(id uint) (*radreplyEntity.Radreply, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(id)
	}
	return CreateRadreplyFixture(), nil
}

func (m *MockRadreplyRepository) GetByUsernameAndAttribute(username, attribute string) (*radreplyEntity.Radreply, error) {
	if m.GetByUsernameAndAttributeFn != nil {
		return m.GetByUsernameAndAttributeFn(username, attribute)
	}
	return CreateRadreplyFixture(), nil
}

func (m *MockRadreplyRepository) GetAll(filter *radreplyDto.RadreplyFilter) ([]radreplyEntity.Radreply, int64, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(filter)
	}
	return []radreplyEntity.Radreply{*CreateRadreplyFixture()}, 1, nil
}

func (m *MockRadreplyRepository) Update(radreply *radreplyEntity.Radreply) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(radreply)
	}
	return nil
}

func (m *MockRadreplyRepository) Delete(id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

// MockRadreplyService is a mock implementation of RadreplyService with function fields
type MockRadreplyService struct {
	CreateRadreplyFn                    func(*radreplyDto.CreateRadreplyRequest) (*radreplyDto.RadreplyResponse, error)
	GetRadreplyByIDFn                   func(uint) (*radreplyDto.RadreplyResponse, error)
	GetRadreplyByUsernameAndAttributeFn func(string, string) (*radreplyDto.RadreplyResponse, error)
	ListRadreplyFn                      func(*radreplyDto.RadreplyFilter) (*radreplyDto.ListRadreplyResponse, error)
	UpdateRadreplyFn                    func(uint, *radreplyDto.UpdateRadreplyRequest) (*radreplyDto.RadreplyResponse, error)
	DeleteRadreplyFn                    func(uint) error
}

func NewMockRadreplyService() *MockRadreplyService {
	return &MockRadreplyService{}
}

func (m *MockRadreplyService) CreateRadreply(req *radreplyDto.CreateRadreplyRequest) (*radreplyDto.RadreplyResponse, error) {
	if m.CreateRadreplyFn != nil {
		return m.CreateRadreplyFn(req)
	}
	return &radreplyDto.RadreplyResponse{
		ID:        1,
		Username:  req.Username,
		Attribute: req.Attribute,
		Op:        req.Op,
		Value:     req.Value,
	}, nil
}

func (m *MockRadreplyService) GetRadreplyByID(id uint) (*radreplyDto.RadreplyResponse, error) {
	if m.GetRadreplyByIDFn != nil {
		return m.GetRadreplyByIDFn(id)
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

func (m *MockRadreplyService) GetRadreplyByUsernameAndAttribute(username, attribute string) (*radreplyDto.RadreplyResponse, error) {
	if m.GetRadreplyByUsernameAndAttributeFn != nil {
		return m.GetRadreplyByUsernameAndAttributeFn(username, attribute)
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

func (m *MockRadreplyService) ListRadreply(filter *radreplyDto.RadreplyFilter) (*radreplyDto.ListRadreplyResponse, error) {
	if m.ListRadreplyFn != nil {
		return m.ListRadreplyFn(filter)
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

func (m *MockRadreplyService) UpdateRadreply(id uint, req *radreplyDto.UpdateRadreplyRequest) (*radreplyDto.RadreplyResponse, error) {
	if m.UpdateRadreplyFn != nil {
		return m.UpdateRadreplyFn(id, req)
	}
	return &radreplyDto.RadreplyResponse{
		ID:        id,
		Username:  req.Username,
		Attribute: req.Attribute,
		Op:        req.Op,
		Value:     req.Value,
	}, nil
}

func (m *MockRadreplyService) DeleteRadreply(id uint) error {
	if m.DeleteRadreplyFn != nil {
		return m.DeleteRadreplyFn(id)
	}
	return nil
}
