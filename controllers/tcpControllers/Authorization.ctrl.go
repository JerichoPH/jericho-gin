package tcpControllers

type AuthorizationController struct{}

// NewAuthorizationController TCP 权鉴
func NewAuthorizationController() *AuthorizationController {
	return &AuthorizationController{}
}

// BindUserUuid 绑定用户uuid到addr
func (receiver AuthorizationController) BindUserUuid(addrToUuid, uuidToAddr map[string]string, uuid, addr string) (map[string]string, map[string]string) {
	addrToUuid[addr] = uuid
	uuidToAddr[uuid] = addr

	return addrToUuid, uuidToAddr
}
