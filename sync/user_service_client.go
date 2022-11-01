package sync

type UserServiceClient struct {
}

type UserInfo struct {
	Username       string
	UserID         string
	OrganizationID string
	IsAdmin        bool
}

type UserKeys struct {
	PublicKey  string
	PrivateKey string
}

const DevPublicKey = "-----BEGIN PUBLIC KEY-----\nMIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgH4PMroKlQ2edizh4JJHxpj7gZJq\nihe9sR9K0q0fYC6rxkwju2fBob0evhLVBy2v8ifedjSf+tOhA4SHhUU2v4sVaClF\n/sfkLUz470bWKkuL58PxQvT+dCVaEsMtONlDmB8Q3/X0MyysUnE0XdIjsth98jKW\nKCgK4P1FenFc0oWbAgMBAAE=\n-----END PUBLIC KEY-----"
const DevPrivateKey = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgH4PMroKlQ2edizh4JJHxpj7gZJqihe9sR9K0q0fYC6rxkwju2fB\nob0evhLVBy2v8ifedjSf+tOhA4SHhUU2v4sVaClF/sfkLUz470bWKkuL58PxQvT+\ndCVaEsMtONlDmB8Q3/X0MyysUnE0XdIjsth98jKWKCgK4P1FenFc0oWbAgMBAAEC\ngYBCTe4pir1hn3KbIufDKTudZdR+VclyuVS7l9h+NN2bTsCLddPxvBg9aDkjoKcY\n8c2WCN31yhvdSniWMc34XNacLQegqjv7I53KTPQH5Paa4a5th7qqoOWr+wysrSPD\nE6co/fLcGLeFnOZy5nAqoeZGZhac251OV+1NRgoeAyJsCQJBAL5cpKx6aXTrYwg+\n6znGwIBZi/GIza3oHBCOHE5qOUnnzIWRNA3zHhy4z9QlEohynV+eNteF1B1ulqkS\nki9Zg78CQQCphoiHgqaMmaivZlieODAt7Fk1j4XKcaLH4/7l/iVLGqSwMFw7u4F0\nJM4Zaqd8DX+Bt54+MOUyOgRgr6QzoEUlAkEAu73s3vp/pUM9UXWUUlAVrMAkB9uv\nVlOz0hQGEMQsqhoFmLmDSDq9OQCAYC8L3yyCzznfxqGDeF+IEUlyiWZUSwJBAIwd\n/krC2hXsC1iuJyDfIDNU3oc+kT66nejJsa03WmuxId3emt1kJaNxqEept7T5EyKM\nOeb9UvMosOWZRwbEuWECQD/3tO9iwL01bTbWQ/IJSKL7+U8vxYk598V7jQTXorCc\ntwaAIAVTj7/dn+Hubt1g0+n8LGNTToF9SMEzEhcN0yk=\n-----END RSA PRIVATE KEY-----"

func NewUserServiceClient() *UserServiceClient {
	return new(UserServiceClient)
}

func (usc *UserServiceClient) GetUserKeys(userID string) *UserKeys {
	return &UserKeys{
		PublicKey:  DevPublicKey,
		PrivateKey: DevPrivateKey,
	}
}

func (usc *UserServiceClient) GetUserInfo(userID string) *UserInfo {
	return &UserInfo{
		Username:       "test",
		UserID:         userID,
		OrganizationID: "Org2MSP",
		IsAdmin:        true,
	}
}
