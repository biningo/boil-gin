package service

import (
	"crypto/sha256"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/middleware"
	"github.com/biningo/boil-gin/model"
	"github.com/biningo/boil-gin/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

/**
*@Author lyer
*@Date 4/16/21 23:10
*@Describe
**/

func GetUserById(uid int) (user model.User, err error) {
	err = global.G_DB.Get(&user, "SELECT username,bio,avatar_id FROM boil_user WHERE id=?", uid)
	return
}

func GetUserByName(username string) (user model.User, err error) {
	err = global.G_DB.Get(&user, "SELECT id,username,password,bio,avatar_id,salt FROM boil_user WHERE username=?", username)
	return
}

func UpdateUserBio(uid int, bio string) error {
	_, err := global.G_DB.Exec("UPDATE boil_user SET bio=? WHERE id=?", bio, uid)
	return err
}

func InsertUser(user model.User) error {
	_, err := global.G_DB.Exec(
		"INSERT INTO boil_user(username,password,avatar_id,salt,bio) VALUE(:username,:password,:avatar_id,:salt,:bio)", user)
	return err
}

func CountUserByName(username string) (count int, err error) {
	r := global.G_DB.QueryRowx("SELECT COUNT(*) FROM boil_user WHERE username=?", username)
	err = r.Scan(&count)
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
	_, err := global.G_DB.Exec("INSERT INTO boil_user_follow_user(follower_id,user_id) value (?,?)", followerId, uid)
	return err
}

func DeleteUserFollow(followerId, uid int) error {
	_, err := global.G_DB.Exec("DELETE  FROM boil_user_follow_user WHERE follower_id=? AND user_id=?", followerId, uid)
	return err
}

func CountUserFollower(uid int) (count int, err error) {
	r := global.G_DB.QueryRowx("SELECT COUNT(*) FROM boil_user_follow_user WHERE user_id=?", uid)
	err = r.Scan(&count)
	return
}

func CountUserFollowing(followerId int) (count int, err error) {
	r := global.G_DB.QueryRowx("SELECT COUNT(*) FROM boil_user_follow_user WHERE follower_id=?", followerId)
	err = r.Scan(&count)
	return
}

func IsUserFollow(followerId, uid int) (bool, error) {
	result := global.G_DB.QueryRowx("SELECT COUNT(*) FROM boil_user_follow_user WHERE follower_id=? AND user_id=?", followerId, uid)
	count := 0
	if err := result.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserInfoByIds(ids []string, loginUserId int) (userInfoVoArr []model.UserInfoVo, err error) {
	if ids == nil || len(ids) == 0 {
		return
	}
	userInfoSql, args, err := sqlx.In("SELECT id,username,bio,avatar_id FROM boil_user WHERE id in (?)", ids)
	if err != nil {
		return nil, err
	}
	if err = global.G_DB.Select(&userInfoVoArr, userInfoSql, args); err != nil {
		return nil, err
	}
	for _, userInfoVo := range userInfoVoArr {
		userInfoVo.FollowerCount, _ = CountUserFollower(userInfoVo.ID)
		userInfoVo.FollowingCount, _ = CountUserFollowing(userInfoVo.ID)
		userInfoVo.IsFollow, _ = IsUserFollow(loginUserId, userInfoVo.ID)
		userInfoVo.BoilCount, _ = CountUserBoil(userInfoVo.ID)
		userInfoVo.CommentBoilCount, _ = CountUserCommentBoil(userInfoVo.ID)
		userInfoVo.LikeBoilCount, _ = CountUserLikeBoil(userInfoVo.ID)
		userInfoVoArr = append(userInfoVoArr, userInfoVo)
	}
	return
}

func GetUserFollowerIds(uid int) (followerIdArr []string, err error) {
	if err = global.G_DB.Select(&followerIdArr, "SELECT follower_id FROM boil_user_follow_user WHERE user_id=?", uid); err != nil {
		return
	}
	return
}
func GetUserFollowingIds(uid int) (followingIdArr []string, err error) {
	if err = global.G_DB.Select(&followingIdArr, "SELECT user_id FROM boil_user_follow_user WHERE follower_id=?", uid); err != nil {
		return
	}
	return
}

func GetUserFollower(uid int, loginUserId int) (userInfoVoArr []model.UserInfoVo, err error) {
	followerIdArr, _ := GetUserFollowerIds(uid)
	userInfoVoArr, err = GetUserInfoByIds(followerIdArr, loginUserId)
	return
}
func GetUserFollowing(uid int, loginUserId int) (userInfoVoArr []model.UserInfoVo, err error) {
	followingIdArr, _ := GetUserFollowingIds(uid)
	userInfoVoArr, err = GetUserInfoByIds(followingIdArr, loginUserId)
	return
}

func GetRecommendUsers(loginUserId int) (userInfoVoArr []model.UserInfoVo, err error) {
	ids, err := GetUserFollowingIds(loginUserId)
	if err != nil {
		return
	}
	for _, sid := range ids {
		id, _ := strconv.Atoi(sid)
		var followerIdArr []string
		err = global.G_DB.Select(&followerIdArr,
			"SELECT user_id FROM boil_user_follow_user WHERE follower_id=? AND user_id!=? AND user_id NOT IN "+
				"(SELECT user_id FROM boil_user_follow_user WHERE follower_id=?)",
			id, loginUserId, loginUserId)
		if err != nil {
			return nil, err
		}
		followerUserInfoVoArr, err := GetUserInfoByIds(followerIdArr, loginUserId)
		if err != nil {
			return userInfoVoArr, err
		}
		userInfoVoArr = append(userInfoVoArr, followerUserInfoVoArr...)
	}
	return userInfoVoArr, nil
}

func GetAllUser(loginUserId int) (userInfoVoArr []model.UserInfoVo, err error) {
	if err = global.G_DB.Select(&userInfoVoArr, "SELECT id,username,bio,avatar_id FROM boil_user"); err != nil {
		return []model.UserInfoVo{}, err
	}
	return userInfoVoArr, nil
}
