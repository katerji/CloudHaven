package service

var (
	userServiceInstance      *userService
	authServiceInstance      *authService
	gcsServiceInstance       *gcsService
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

func GetGCSService() *gcsService {
	if gcsServiceInstance == nil {
		gcsServiceInstance = &gcsService{}
	}
	return gcsServiceInstance
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
