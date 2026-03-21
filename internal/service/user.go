package service

import (
	"errors"
	"mini-social/internal/model"
	"mini-social/internal/repository"
	"mini-social/pkg/password"
	"regexp"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUsernameAlreadyExists = errors.New("用户名已存在")
	ErrInvalidUsername       = errors.New("用户名格式不正确：须为3-20位字母、数字或下划线")
)

// 用户名正则表达式
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

type RegisterInput struct {
	Username string
	Password string
}

func (s *UserService) Register(Input RegisterInput) (*model.User, error) {
	//用户名规范性检查
	if !usernameRegex.MatchString(Input.Username) {
		return nil, ErrInvalidUsername
	}
	//查重
	_, err := s.userRepo.GetByUsername(Input.Username)
	if err == nil {
		return nil, ErrUsernameAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err //数据库连不上了或者其他系统错误
	}

	//密码加密
	hashedPassword, err := password.Hash(Input.Password)
	if err != nil {
		return nil, err
	}

	//创建并捕获并发冲突
	user := &model.User{
		Username:     Input.Username,
		PasswordHash: hashedPassword,
	}
	if err := s.userRepo.Create(user); err != nil {
		// 检查是否为MySQL的唯一索引冲突错误(Error 1062)
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrUsernameAlreadyExists
		}
		return nil, err
	}

	return user, nil
}
