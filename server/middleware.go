//Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 11:03 上午
package server

type CustomClaims struct {
	Name     string `json:"name"`
	UserID   uint64 `json:"user_id"`   //用户Id
	UserName string `json:"user_name"` //用户名
	UserUUid string `json:"user_uuid"` //机构uuid
}
