package template

import "github.com/Yukio0315/reservation-backend/src/util"

const (
	// REGISTERSUB is title of register email
	REGISTERSUB = "【" + util.OFFICENAME + "】ご登録ありがとうございます"
	// REGISTERBODY is body of register email
	REGISTERBODY = util.OFFICENAME + "へのご登録が完了しました。"
	// CHANGEPASSWORDTITLE is title of change password mail
	CHANGEPASSWORDTITLE = "【" + util.OFFICENAME + "】パスワードの変更が完了しました"
	// CHANGEPASSWORDBODY is body of change password mail
	CHANGEPASSWORDBODY = "パスワードの変更が完了しました。"
	// RESETPASSWORDTITLE is title of password reset mail
	RESETPASSWORDTITLE = "【" + util.OFFICENAME + "】パスワードのリセットが完了しました"
	// RESETPASSWORDBODY is body of password reset mail
	RESETPASSWORDBODY = "パスワードのリセットが完了しました。"
	// CHANGEEMAILTITLE is title of change mail address mail
	CHANGEEMAILTITLE = "【" + util.OFFICENAME + "】Emailアドレスを変更しました"
	// CHANGEEMAILBODY is body of change email address mail
	CHANGEEMAILBODY = "Emailアドレスを変更しました。"
	// DELETEACCOUNTTITLE is title of delete account mail
	DELETEACCOUNTTITLE = "【" + util.OFFICENAME + "】アカウントを削除しました"
	// DELETEACCOUNTBODY is body of delete account mail
	DELETEACCOUNTBODY = "アカウントを削除しました。\nまたのご利用をお待ちしております。"
	// ONETIMEURLTITLE is title of delete account mail
	ONETIMEURLTITLE = "【" + util.OFFICENAME + "】パスワードリセット用URLを発行しました。"
)

// OneTimeURLBody is body of delete
func OneTimeURLBody(uuid string) string {
	return "1回のみアクセス可能なURLを発行しました。\n24時間以内にアクセスして、パスワードの再設定を行ってください。\n" + util.BASEURL + "/password-reset/" + uuid + "\n\n期限が切れた場合は再度アクセスしてください。"
}
