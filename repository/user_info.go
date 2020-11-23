package repository

import (
	"gitee.com/cristiane/micro-mall-users-consumer/model/mysql"
	"gitee.com/kelvins-io/kelvins"
)

func GetUserNameByUid(uid int) (*mysql.User, error) {
	var user mysql.User
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableUser).Select("user_name").Where("id = ?", uid).Get(&user)
	return &user, err
}

func GetUserByPhone(sqlSelect, countryCode, phone string) (*mysql.User, error) {
	var user mysql.User
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableUser).Select(sqlSelect).Where("country_code = ? and phone = ?", countryCode, phone).Get(&user)
	return &user, err
}
