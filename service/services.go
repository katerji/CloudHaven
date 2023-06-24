package service

var (
	userServiceInstance      *userService
	authServiceInstance      *authService
	gcpServiceInstance       *gcpService
	jwtServiceInstance       *jwtService
	fileServiceInstance      *fileService
	fileShareServiceInstance *fileShareService
)

func GetUserService() *userService {
	if userServiceInstance == nil {
		userServiceInstance = initUserService()
	}
	return userServiceInstance
}

func GetAuthService() *authService {
	if authServiceInstance == nil {
		authServiceInstance = &authService{}
	}
	return authServiceInstance
}

func GetGcpService() *gcpService {
	if gcpServiceInstance == nil {
		gcpServiceInstance = &gcpService{}
	}
	return gcpServiceInstance
}

func GetJwtService() *jwtService {
	if jwtServiceInstance == nil {
		jwtServiceInstance = &jwtService{}
	}
	return jwtServiceInstance
}

func GetFileService() *fileService {
	if fileServiceInstance == nil {
		fileServiceInstance = &fileService{}
	}
	return fileServiceInstance
}

func GetFileShareService() *fileShareService {
	if fileShareServiceInstance == nil {
		fileShareServiceInstance = &fileShareService{}
	}
	return fileShareServiceInstance
}
