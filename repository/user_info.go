package repository

import (
	"gitee.com/cristiane/micro-mall-users-consumer/model/mysql"
	"gitee.com/kelvins-io/kelvins"
)

func GetUserInfoByUid(uid int) (*mysql.User, error) {
	var user mysql.User
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableUser).Select("user_name,email").Where("id = ?", uid).Get(&user)
	return &user, err
}

func GetUserByPhone(sqlSelect, countryCode, phone string) (*mysql.User, error) {
	var user mysql.User
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableUser).Select(sqlSelect).
		Where("country_code = ?", countryCode).
		Where("phone = ?", phone).
		Get(&user)
	return &user, err
}
