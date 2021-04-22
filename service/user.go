package service

import (
	"crypto/sha256"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/middleware"
	"github.com/biningo/boil-gin/model"
	"github.com/biningo/boil-gin/utils"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"strings"
	"time"
)

/**
*@Author lyer
*@Date 4/16/21 23:10
*@Describe
**/

func GetUserById(uid int) (model.User, error) {
	db := global.G_DB
	result, err := db.Query("SELECT username,bio,avatar_id FROM boil_user WHERE id=?", uid)
	if err != nil {
		return model.User{}, err
	}
	defer result.Close()
	user := model.User{}
	result.Next()
	result.Scan(&user.UserName, &user.Bio, &user.AvatarID)
	result.Close()
	return user, nil
}

func GetUserByName(username string) (user model.User, err error) {
	db := global.G_DB
	exec, _ := db.Prepare("SELECT id,username,password,bio,avatar_id,salt FROM boil_user WHERE username=?")
	result, err := exec.Query(username)
	if err != nil {
		return
	}
	result.Next()
	result.Scan(&user.ID, &user.UserName, &user.PassWord, &user.Bio, &user.AvatarID, &user.Salt)
	result.Close()
	return
}

func UpdateUserBio(uid int, bio string) error {
	db := global.G_DB
	_, err := db.Exec("UPDATE boil_user SET bio=? WHERE id=?", bio, uid)
	return err
}

func InsertUser(user model.User) error {
	db := global.G_DB
	_, err := db.Exec(
		"INSERT INTO boil_user(username,password,avatar_id,salt,bio) VALUE(?,?,?,?,?)",
		user.UserName, user.PassWord, user.AvatarID, user.Salt, user.Bio)
	return err
}

func CountUserByName(username string) (count int, err error) {
	db := global.G_DB
	result, err := db.Query("SELECT COUNT(*) FROM boil_user WHERE username=?", username)
	if err != nil {
		return 0, nil
	}
	result.Next()
	result.Scan(&count)
	result.Close()
	return
}

func UserRegistryVoToUser(userRegistryVo model.UserRegistryVo) (user model.User) {
	salt := utils.Rand5Str()
	pwd256Token := sha256.Sum256([]byte(userRegistryVo.PassWord + salt))
	if userRegistryVo.AvatarID <= 0 || userRegistryVo.AvatarID > 5 {
		userRegistryVo.AvatarID = utils.RandInt(1, 5)
	}
	if userRegistryVo.Bio == "" {
		userRegistryVo.Bio = "这个家伙很懒,什么也没留下"
	}
	user.UserName = userRegistryVo.UserName
	user.PassWord = string(pwd256Token[:])
	user.Salt = salt
	user.Bio = userRegistryVo.Bio
	user.AvatarID = userRegistryVo.AvatarID
	return
}

func CreateUserToken(user model.User) (string, error) {
	tokenString, err := middleware.NewJwt().CreateToken(&middleware.CustomClaims{
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "boil",
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(global.G_CONFIG.Jwt.TokenTime)).Unix(),
		},
	})
	return tokenString, err
}

func CheckUserToken(userLoginVo model.UserLoginVo, salt string, password string) bool {
	pwd256Bytes := sha256.Sum256([]byte(userLoginVo.PassWord + salt))
	userLoginPwdToken := string(pwd256Bytes[:])
	return userLoginPwdToken == password
}

func CleanUser(user *model.User) {
	user.PassWord = ""
	user.Salt = ""
}

func GetUserInfoById(uid int, loginUserId int) (userInfoVo model.UserInfoVo, err error) {
	user, err := GetUserById(uid)
	if err != nil {
		return
	}
	userInfoVo.ID = uid
	userInfoVo.UserName = user.UserName
	userInfoVo.AvatarID = user.AvatarID
	userInfoVo.Bio = user.Bio
	userInfoVo.FollowerCount, _ = CountUserFollower(uid)
	userInfoVo.FollowingCount, _ = CountUserFollowing(uid)
	userInfoVo.IsFollow, _ = IsUserFollow(loginUserId, user.ID)
	userInfoVo.BoilCount, _ = CountUserBoil(uid)
	userInfoVo.CommentBoilCount, _ = CountUserCommentBoil(uid)
	userInfoVo.LikeBoilCount, _ = CountUserLikeBoil(uid)
	return
}

func InsertUserFollow(followerId, uid int) error {
	db := global.G_DB
	_, err := db.Exec("INSERT INTO boil_user_follow_user(follower_id,user_id) value (?,?)", followerId, uid)
	return err
}

func DeleteUserFollow(followerId, uid int) error {
	db := global.G_DB
	_, err := db.Exec("DELETE  FROM boil_user_follow_user WHERE follower_id=? AND user_id=?", followerId, uid)
	return err
}

func CountUserFollower(uid int) (count int, err error) {
	db := global.G_DB
	result, err := db.Query("SELECT COUNT(*) FROM boil_user_follow_user WHERE user_id=?", uid)
	if err != nil {
		return 0, err
	}
	result.Next()
	result.Scan(&count)
	result.Close()
	return
}

func CountUserFollowing(followerId int) (count int, err error) {
	db := global.G_DB
	result, err := db.Query("SELECT COUNT(*) FROM boil_user_follow_user WHERE follower_id=?", followerId)
	if err != nil {
		return 0, err
	}
	result.Next()
	result.Scan(&count)
	result.Close()
	return
}

func IsUserFollow(followerId, uid int) (bool, error) {
	db := global.G_DB
	result, err := db.Query("SELECT COUNT(*) FROM boil_user_follow_user WHERE follower_id=? AND user_id=?", followerId, uid)
	if err != nil {
		return false, err
	}
	result.Next()
	count := 0
	result.Scan(&count)
	result.Close()
	return count > 0, nil
}

func GetUserInfoByIds(ids []string, loginUserId int) ([]model.UserInfoVo, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	db := global.G_DB
	result, err := db.Query(
		"SELECT id,username,bio,avatar_id FROM boil_user WHERE " + "id in (" + strings.Join(ids, ",") + ")")
	if err != nil {
		return []model.UserInfoVo{}, err
	}
	userInfoVo := model.UserInfoVo{}
	userInfoVoArr := []model.UserInfoVo{}
	for result.Next() {
		result.Scan(&userInfoVo.ID, &userInfoVo.UserName, &userInfoVo.Bio, &userInfoVo.AvatarID)
		userInfoVo.FollowerCount, _ = CountUserFollower(userInfoVo.ID)
		userInfoVo.FollowingCount, _ = CountUserFollowing(userInfoVo.ID)
		userInfoVo.IsFollow, _ = IsUserFollow(loginUserId, userInfoVo.ID)
		userInfoVo.BoilCount, _ = CountUserBoil(userInfoVo.ID)
		userInfoVo.CommentBoilCount, _ = CountUserCommentBoil(userInfoVo.ID)
		userInfoVo.LikeBoilCount, _ = CountUserLikeBoil(userInfoVo.ID)
		userInfoVoArr = append(userInfoVoArr, userInfoVo)
	}
	result.Close()
	return userInfoVoArr, nil
}

func GetUserFollower(uid int, loginUserId int) ([]model.UserInfoVo, error) {
	db := global.G_DB
	result, err := db.Query("SELECT follower_id FROM boil_user_follow_user WHERE user_id=?", uid)
	if err != nil {
		return nil, err
	}
	followerIdArr := []string{}
	for result.Next() {
		id := "0"
		result.Scan(&id)
		followerIdArr = append(followerIdArr, id)
	}
	result.Close()
	userInfoArr, err := GetUserInfoByIds(followerIdArr, loginUserId)
	return userInfoArr, err
}

func GetUserFollowing(uid, loginUserId int) ([]model.UserInfoVo, error) {
	db := global.G_DB
	result, err := db.Query(
		"SELECT user_id FROM boil_user_follow_user WHERE follower_id=? AND user_id!=?",
		uid, loginUserId)
	if err != nil {
		return nil, err
	}
	followerIdArr := []string{}
	for result.Next() {
		id := "0"
		result.Scan(&id)
		followerIdArr = append(followerIdArr, id)
	}
	result.Close()
	userInfoArr, err := GetUserInfoByIds(followerIdArr, loginUserId)
	return userInfoArr, err
}

func GetUserFollowingIds(uid int) ([]string, error) {
	db := global.G_DB
	result, err := db.Query("SELECT user_id FROM boil_user_follow_user WHERE follower_id=?", uid)
	if err != nil {
		return nil, err
	}
	followerIdArr := []string{}
	for result.Next() {
		id := "0"
		result.Scan(&id)
		followerIdArr = append(followerIdArr, id)
	}
	result.Close()
	return followerIdArr, nil
}

func GetRecommendUsers(loginUserId int) ([]model.UserInfoVo, error) {
	ids, err := GetUserFollowingIds(loginUserId)
	userInfoVoArr := []model.UserInfoVo{}
	if err != nil {
		return userInfoVoArr, err
	}
	for _, sid := range ids {
		id, _ := strconv.Atoi(sid)
		db := global.G_DB
		result, err := db.Query(
			"SELECT user_id FROM boil_user_follow_user WHERE follower_id=? AND user_id!=? AND user_id NOT IN "+
				"(SELECT user_id FROM boil_user_follow_user WHERE follower_id=?)",
			id, loginUserId, loginUserId)
		if err != nil {
			return nil, err
		}
		followerIdArr := []string{}
		for result.Next() {
			id := "0"
			result.Scan(&id)
			followerIdArr = append(followerIdArr, id)
		}
		result.Close()
		followerUserInfoVoArr, err := GetUserInfoByIds(followerIdArr, loginUserId)
		userInfoVoArr = append(userInfoVoArr, followerUserInfoVoArr...)
	}
	return userInfoVoArr, nil
}

func GetAllUser(loginUserId int) ([]model.UserInfoVo, error) {
	db := global.G_DB
	result, err := db.Query(
		"SELECT id,username,bio,avatar_id FROM boil_user")
	if err != nil {
		return []model.UserInfoVo{}, err
	}
	userInfoVo := model.UserInfoVo{}
	userInfoVoArr := []model.UserInfoVo{}
	for result.Next() {
		result.Scan(&userInfoVo.ID, &userInfoVo.UserName, &userInfoVo.Bio, &userInfoVo.AvatarID)
		userInfoVo.FollowerCount, _ = CountUserFollower(userInfoVo.ID)
		userInfoVo.FollowingCount, _ = CountUserFollowing(userInfoVo.ID)
		userInfoVo.IsFollow, _ = IsUserFollow(loginUserId, userInfoVo.ID)
		userInfoVo.BoilCount, _ = CountUserBoil(userInfoVo.ID)
		userInfoVo.CommentBoilCount, _ = CountUserCommentBoil(userInfoVo.ID)
		userInfoVo.LikeBoilCount, _ = CountUserLikeBoil(userInfoVo.ID)
		userInfoVoArr = append(userInfoVoArr, userInfoVo)
	}
	result.Close()
	return userInfoVoArr, nil
}
