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
)

// User fixtures
func CreateUserFixture() *userEntity.User {
	return &userEntity.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "$2a$10$example.hashed.password",
	}
}

func CreateUserRequestFixture() *userDto.CreateUserRequest {
	return &userDto.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
}

func CreateUpdateUserRequestFixture() *userDto.UpdateUserRequest {
	return &userDto.UpdateUserRequest{
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}
}

// Payment fixtures
func CreatePaymentFixture() *paymentEntity.Payment {
	return &paymentEntity.Payment{
		ID:          1,
		Amount:      100.50,
		Currency:    "USD",
		Status:      paymentEntity.PaymentStatusPending,
		Description: "Test payment",
		UserID:      1,
	}
}

func CreatePaymentRequestFixture() *paymentDto.CreatePaymentRequest {
	return &paymentDto.CreatePaymentRequest{
		Amount:      100.50,
		Currency:    "USD",
		Description: "Test payment",
		UserID:      1,
	}
}

func CreateUpdatePaymentRequestFixture() *paymentDto.UpdatePaymentRequest {
	return &paymentDto.UpdatePaymentRequest{
		Status:      paymentEntity.PaymentStatusCompleted.String(),
		Description: "Payment completed",
	}
}

func CreatePaymentFilterFixture() *paymentDto.PaymentFilter {
	return &paymentDto.PaymentFilter{
		Status:   "pending",
		Currency: "USD",
		UserID:   1,
		Page:     1,
		PageSize: 10,
	}
}

// NAS fixtures
func CreateNASFixture() *nasEntity.NAS {
	ports := 1812
	return &nasEntity.NAS{
		ID:              1,
		NASName:         "test-nas-01",
		ShortName:       "test-nas",
		Type:            "other",
		Ports:           ports,
		Secret:          "testing123",
		Server:          "192.168.1.1",
		Community:       "public",
		Description:     "Test NAS",
		RequireMa:       "auto",
		LimitProxyState: "auto",
	}
}

func CreateNASRequestFixture() *nasDto.CreateNASRequest {
	ports := 1812
	return &nasDto.CreateNASRequest{
		NASName:         "test-nas-01",
		ShortName:       "test-nas",
		Type:            "other",
		Ports:           &ports,
		Secret:          "testing123",
		Server:          "192.168.1.1",
		Community:       "public",
		Description:     "Test NAS",
		RequireMa:       "auto",
		LimitProxyState: "auto",
	}
}

func CreateUpdateNASRequestFixture() *nasDto.UpdateNASRequest {
	ports := 1813
	return &nasDto.UpdateNASRequest{
		NASName:     "updated-nas-01",
		Description: "Updated NAS",
		Ports:       &ports,
	}
}

func CreateNASFilterFixture() *nasDto.NASFilter {
	return &nasDto.NASFilter{
		NASName:     "",
		ShortName:   "",
		Type:        "",
		Description: "",
		Page:        1,
		PageSize:    10,
	}
}

// Radcheck fixtures
func CreateRadcheckFixture() *radcheckEntity.Radcheck {
	return &radcheckEntity.Radcheck{
		ID:        1,
		Username:  "testuser",
		Attribute: "User-Password",
		Op:        ":=",
		Value:     "testing123",
	}
}

func CreateRadcheckRequestFixture() *radcheckDto.CreateRadcheckRequest {
	return &radcheckDto.CreateRadcheckRequest{
		Username:  "testuser",
		Attribute: "User-Password",
		Op:        ":=",
		Value:     "testing123",
	}
}

func CreateUpdateRadcheckRequestFixture() *radcheckDto.UpdateRadcheckRequest {
	return &radcheckDto.UpdateRadcheckRequest{
		Username:  "testuser",
		Attribute: "User-Password",
		Op:        ":=",
		Value:     "newpassword123",
	}
}

func CreateRadcheckFilterFixture() *radcheckDto.RadcheckFilter {
	return &radcheckDto.RadcheckFilter{
		Username:  "",
		Attribute: "",
		Page:      1,
		PageSize:  10,
	}
}

// Radreply fixtures
func CreateRadreplyFixture() *radreplyEntity.Radreply {
	return &radreplyEntity.Radreply{
		ID:        1,
		Username:  "john",
		Attribute: "Reply-Message",
		Op:        "=",
		Value:     "Welcome",
	}
}

func CreateRadreplyRequestFixture() *radreplyDto.CreateRadreplyRequest {
	return &radreplyDto.CreateRadreplyRequest{
		Username:  "john",
		Attribute: "Reply-Message",
		Op:        "=",
		Value:     "Welcome",
	}
}

func CreateUpdateRadreplyRequestFixture() *radreplyDto.UpdateRadreplyRequest {
	return &radreplyDto.UpdateRadreplyRequest{
		Username:  "john",
		Attribute: "Reply-Message",
		Op:        "=",
		Value:     "Updated",
	}
}

func CreateRadreplyFilterFixture() *radreplyDto.RadreplyFilter {
	return &radreplyDto.RadreplyFilter{
		Username:  "",
		Attribute: "",
		Page:      1,
		PageSize:  10,
	}
}
