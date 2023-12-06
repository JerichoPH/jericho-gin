package tcpControllers

type AuthController struct{}

// NewAuthController TCP 权鉴
func NewAuthController() *AuthController {
	return &AuthController{}
}

// BindUserUuid 绑定用户uuid到addr
func (receiver AuthController) BindUserUuid(addrToUuid, uuidToAddr map[string]string, uuid, addr string) (map[string]string, map[string]string) {
	addrToUuid[addr] = uuid
	uuidToAddr[uuid] = addr

	return addrToUuid, uuidToAddr
}
