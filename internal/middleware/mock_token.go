package middleware

var _ JwtAuther = (*MockAuther)(nil)

type MockAuther struct {
}

func NewMockAuther() *MockAuther {
	return &MockAuther{}
}
func (m MockAuther) RequireLogin(scg StandardClaimsGetter, tokenString string) error {
	return nil
}
