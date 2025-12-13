package testutil

import (
	nasDto "github.com/novriyantoAli/freeradius-service/internal/application/nas/dto"
	nasEntity "github.com/novriyantoAli/freeradius-service/internal/application/nas/entity"
	paymentDto "github.com/novriyantoAli/freeradius-service/internal/application/payment/dto"
	paymentEntity "github.com/novriyantoAli/freeradius-service/internal/application/payment/entity"
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
