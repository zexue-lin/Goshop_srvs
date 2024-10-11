package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"goshop_srvs/goods_srv/global"
	"goshop_srvs/goods_srv/model"
	"goshop_srvs/goods_srv/proto"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// 只需要实现user.proto里的方法即可
// 再生成的pb.go文件中找 interface 关键字，找到那几个API

type UserServer struct{}

// 方法1 查询用户-用户列表
// 指向 UserServer 结构体的指针;第二个参数PageInfo 消息类型的指针;返回值是一个指向 UserListResponse 消息类型的指针和一个错误
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	// 获取用户列表（所有）
	var users []model.User
	result := global.DB.Find(&users) // gorm查询所有用户
	if result.Error != nil {
		return nil, result.Error
	}

	// 开始构建返回，创建响应消息
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected) // 因为 user.proto 文件中定义的类型是int32

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users) // 分页查询

	for _, user := range users {
		userInfoRsp := ModelToResponse(user) // 将这里users里面取出的model对象转换成文件user.pb.go里的UserListResponse对象
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

// 分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	// gorm分页
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// 将查出来的对象转成proto里的对象
func ModelToResponse(user model.User) proto.UserInfoResponse {
	// 在grpc的message字段中，如果有默认值不能随便赋值nil，容易出错
	// 要搞清哪些有默认值，birthday注册的时候可能是nil，
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.Birthday = uint64(user.Birthday.Unix()) // 得到Unix时间戳转换为uint64赋值给birthday
	}
	return userInfoRsp
}

/*
定义了 UserServer 结构体，并实现了 GetUserList 方法，该方法用于从数据库中获取用户列表并返回分页的用户数据。
辅助函数 Paginate 和 ModelToResponse 来帮助处理分页和模型到响应消息的转换

req *proto.PageInfo：一个指向 PageInfo 消息类型的指针，包含了分页信息（如页码 Pn 和每页大小 PSize）
*/

// 方法2 根据手机号查询用户
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	// gorm的基础理解，这句sql相当于：条件 where 模型定义的Mobile = req.Mobile
	// 最终 SQL 类似于：SELECT * FROM users WHERE mobile = 'xxxx' LIMIT 1;
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// 通过id查询用户
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User

	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// 新建用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// 严谨的做法: 新建之前先查询一下用户是否已经存在
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName

	// 密码加密
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.Password, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// 个人中心更新用户
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	// 如果将一个int类型转换成一个Time类型，重点是时间转换
	birthday := time.Unix(int64(req.Birthday), 0)
	user.NickName = req.NickName
	user.Birthday = &birthday // 是一个指针类型
	user.Gender = req.Gender

	result = global.DB.Save(user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}

// 检查用户密码
func (s *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")

	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options) // 第0个是空格，第2个是salt
	return &proto.CheckResponse{Success: check}, nil
}
