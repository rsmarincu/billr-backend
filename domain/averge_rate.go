package domain

type Role string

const (
	RoleFrontendDeveloper  Role = "FRONTEND_DEVELOPER"
	RoleBackendDeveloper   Role = "BACKEND_DEVELOPER"
	RoleFullStackDeveloper Role = "FULL_STACK_DEVELOPER"
	RoleDataEngineer       Role = "DATA_ENGINEER"
	RoleDataScientist      Role = "DATA_SCIENTIST"
	RoleDevopsEngineer     Role = "DEVOPS_ENGINEER"
	RolePlatformEngineer   Role = "PLATFORM_ENGINEER"
	RoleUnknown            Role = "UNKNOWN"
)

type AverageRate struct {
	Id       string
	ClientId string
	Min      float64
	Max      float64
	Median   float64
	Average  float64
	Role     Role
}
