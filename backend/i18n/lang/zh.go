package lang

var Zh = map[string]map[string]string{
	"common": {
		"GenericError":              "发生了一个错误。请稍后重试。",
		"PermissionDenied":          "您无权执行此操作。",
		"RequestParameterIncorrect": "请求参数不正确。",
		"ModificationFailed":        "修改失败，请稍后重试。",
		"DeletionFailed":            "删除失败，请稍后重试。",
		"DoItToSelf":                "你无法对自己进行操作。",
		"OperationFailed":           "操作失败，请稍后再试。",
		"LinkExpiredTitle":          "链接已过期",
		"LinkExpired":               "链接已过期。",
		"SuccessfulDesc":            "3秒后系统将自动跳转至主页。如果没有，请点击这里。",
		"ImageUploadFailed":         "图片上传失败，请稍后重试。",
		"ImageTooLarge":             "图片文件过大，请选择较小的图片。",
		"EmailSendFailed":           "邮件发送失败，请稍后重试。",
		"TooManyOperations":         "您尝试的次数过多，请稍后再试。",
	},
	"category": {
		"DoesNotExist":     "分类不存在。",
		"CreationFailed":   "分类创建失败，请稍后重试。",
		"DeleteNonEmpty":   "无法删除非空的分类。",
		"IsNotCategory":    "这不是一个分类。",
		"CannotBeCopied":   "无法复制分类。",
		"CannotBeExported": "无法导出分类。",
		"FailedToGet":      "获取分类失败，请稍后重试。",
		"FailedToDelete":   "删除分类失败，请稍后重试。",
	},
	"share": {
		"FailedToGetStatus":           "获取共享状态失败，请稍后重试。",
		"SharingFailed":               "共享失败，请稍后重试。",
		"FailedToDisable":             "停止共享失败，请稍后重试。",
		"SharedKeyResetFailed":        "共享密钥重置失败，请稍后重试。",
		"SharedKeyVerificationFailed": "共享密钥验证失败，请稍后重试。",
		"SharedKeyError":              "共享密钥错误。",
		"PublicProjectShare":          "这是一个任何人都可以访问的公共项目。",
	},
	"user": {
		"LoginRequired":                    "请登录您的账号。",
		"OriginalPasswordWrong":            "原密码错误。",
		"PasswordUpdateFailed":             "密码修改失败，请稍后重试。",
		"PasswordResetFailed":              "密码重置失败，请稍后重试。",
		"PasswordResetSuccessfulTitle":     "密码重置成功",
		"EmailNotChanged":                  "这是您当前使用的邮箱。",
		"EmailHasBeenUsed":                 "该邮箱已被使用。",
		"EmailHasBeenVerified":             "该邮箱已被验证。",
		"EmailUpdateFailed":                "邮箱修改失败，请稍后重试。",
		"EmailUpdateSuccessfulTitle":       "邮箱修改成功",
		"EmailUpdateSuccessfulDesc":        "您将会使用新的邮箱进行登录。",
		"EmailVerificationSuccessfulTitle": "邮箱验证成功",
		"EmailVerificationFailedTitle":     "邮箱验证失败",
		"EmailHasVerifiedTitle":            "邮箱已验证",
		"EmailHasVerifiedDesc":             "您的邮箱已验证。",
		"ResendEmail":                      "请重新发送电子邮件。",
		"EmailHasRegistered":               "该邮箱已被注册。",
		"EmailVerificationFailed":          "邮箱验证失败。",
		"EmailDoesNotExist":                "邮箱不存在。",
		"LoginFailed":                      "登录失败，请稍后重试。",
		"IncorrectEmailOrPassword":         "错误的邮箱或密码。",
		"RegisterFailed":                   "注册失败，请稍后重试。",
		"OauthLoginFailed":                 "Oauth 登录失败，请稍后重试。",
		"LoginStatusExpired":               "登录状态已过期，请重新登录。",
		"InactiveEmail":                    "电子邮件无效，请在尝试前激活。",
		"NotSupportOauth":                  "不支持 %s。",
		"OauthConnectFailed":               "无法连接到 %s，请稍后再试。",
		"OauthConnectRepeat":               "此 %s 帐户已绑定。",
		"OauthDisconnectFailed":            "无法与 %s 解绑，请稍后再试。",
		"FailedToGetList":                  "获取用户列表失败，请稍后重试。",
		"DoesNotExist":                     "用户不存在。",
		"FailedToDelete":                   "删除用户失败，请稍后重试。",
	},
	"team": {
		"CreationFailed":             "创建团队失败，请稍后重试。",
		"FailedToGetList":            "获取团队列表失败，请稍后重试。",
		"FailedToGetCurrentTeam":     "获取当前团队失败，请稍后重试。",
		"CurrentTeamDoesNotExist":    "当前团队不存在。",
		"NoTeam":                     "您目前没有任何团队，创建一个团队立即开始协作。",
		"DoesNotExist":               "团队不存在。",
		"FailedToSwitch":             "切换团队失败，请稍后重试。",
		"InvitationTokenResetFailed": "邀请码重置失败，请稍后重试。",
		"InvitationTokenNotFound":    "未找到邀请码。",
		"InvalidInvitationToken":     "邀请码无效。",
		"FailedToJoinTeam":           "加入团队失败，请稍后重试。",
	},
	"teamMember": {
		"DoesNotExist":              "团队成员不存在。",
		"TeamTransferFailed":        "团队转移失败，请稍后重试。",
		"TeamTransferInvalidMember": "您只能将团队转移给管理员。",
		"FailedToGetList":           "获取团队成员列表失败，请稍后重试。",
		"RemoveFailed":              "团队成员删除失败，请稍后重试。",
		"CanNotQuitOwnTeam":         "您无法退出您自己的团队。",
		"TeamQuitFailed":            "退出失败，请稍后重试。",
		"JoinTeamRepeat":            "您已经是团队成员了。",
		"NotTeamMember":             "用户不是团队成员。",
		"Deactivated":               "团队成员已被停用。",
		"NotInTheTeam":              "用户不在团队中。",
	},
	"project": {
		"DoesNotExist":             "项目不存在。",
		"FileParseFailed":          "文件解析失败，请稍后重试。",
		"CreationFailed":           "项目创建失败，请稍后重试。",
		"FailedToGetList":          "获取项目列表失败，请稍后重试。",
		"FailedToDelete":           "删除项目失败，请稍后重试。",
		"TransferFailed":           "项目转移失败，请稍后重试。",
		"TransferToErrMember":      "项目只能转移给具有写入权限的成员。",
		"TransferToDisabledMember": "项目无法移交给被禁用的用户。",
		"CanNotQuitOwnProject":     "您无法退出您自己的项目。",
		"QuitFailed":               "退出失败，请稍后重试。",
		"FailedToUnfollowProject":  "无法取关项目，请稍后重试。",
		"ExportFailed":             "项目导出失败，请稍后重试。",
		"NotSupportFileType":       "不支持 %s。",
	},
	"projectGroup": {
		"DoesNotExist":          "项目分组不存在。",
		"GroupingFailed":        "更改分组失败，请稍后重试。",
		"FailedToFollowProject": "关注项目失败，请稍后重试。",
		"CreationFailed":        "创建项目分组失败，请稍后重试。",
		"NameHasBeenUsed":       "该分组名称已被使用。",
		"FailedToGetList":       "获取项目分组列表失败，请稍后重试。",
		"FailedToDelete":        "删除项目分组失败，请稍后重试。",
		"RenameFailed":          "重命名失败，请稍后重试。",
		"SortingFailed":         "项目分组排序失败，请稍后重试。",
	},
	"projectMember": {
		"DoesNotExist":             "项目成员不存在。",
		"CanNotAddProjectManager":  "一个项目只能有一个管理员。",
		"FailedToAddProjectMember": "添加项目成员失败，请稍后重试。",
		"FailedToGetList":          "获取项目成员列表失败，请稍后重试。",
		"RemoveFailed":             "项目成员删除失败，请稍后重试。",
		"NotInTheProject":          "成员不在项目中。",
	},
	"projectServer": {
		"CreationFailed":  "服务器 URL 创建失败，请稍后重试。",
		"HasBeenUsed":     "URL 已被使用。",
		"FailedToGetList": "获取服务器 URL 列表失败，请稍后重试。",
		"FailedToDelete":  "删除服务器 URL 失败，请稍后重试。",
		"SortingFailed":   "服务器 URL 排序失败，请稍后重试。",
	},
	"globalParameter": {
		"CreationFailed":  "全局参数创建失败，请稍后重试。",
		"HasBeenUsed":     "该参数已被使用。",
		"FailedToGetList": "获取参数列表失败，请稍后重试。",
		"DoesNotExist":    "全局参数不存在。",
		"FailedToDelete":  "删除参数失败，请稍后重试。",
		"FailedToSort":    "全局参数排序失败，请稍后重试。",
	},
	"definitionSchema": {
		"CreationFailed":   "模型创建失败，请稍后重试。",
		"GenerationFailed": "模型生成失败，请稍后重试。",
		"FailedToGetList":  "获取模型列表失败，请稍后重试。",
		"FailedToGet":      "获取模型失败，请稍后重试。",
		"DoesNotExist":     "模型不存在。",
		"FailedToDelete":   "删除模型失败，请稍后重试。",
		"FailedToMove":     "移动模型失败，请稍后重试。",
		"CopyFailed":       "模型复制失败，请稍后重试。",
	},
	"definitionSchemaHistory": {
		"FailedToGetList": "获取模型历史列表失败，请稍后重试。",
		"FailedToGet":     "获取模型历史记录失败，请稍后重试。",
		"DoesNotExist":    "模型历史不存在。",
		"RestoreFailed":   "模型恢复失败，请稍后重试。",
		"DiffFailed":      "模型对比失败，请稍后重试。",
	},
	"definitionResponse": {
		"CreationFailed":  "响应创建失败，请稍后重试。",
		"DoesNotExist":    "响应不存在。",
		"FailedToGetList": "获取响应列表失败，请稍后重试。",
		"FailedToGet":     "获取响应失败，请稍后再试。",
		"FailedToDelete":  "删除响应失败，请稍后重试。",
		"FailedToMove":    "移动响应失败，请稍后重试。",
		"CopyFailed":      "响应复制失败，请稍后重试。",
	},
	"iteration": {
		"DoesNotExist":    "迭代不存在。",
		"CreationFailed":  "迭代创建失败，请稍后重试。",
		"FailedToGetList": "获取迭代列表失败，请稍后重试。",
		"FailedToGet":     "获取迭代失败，请稍后重试。",
		"FailedToDelete":  "删除迭代失败，请稍后重试。",
	},
	"collection": {
		"FailedToGetList":  "获取 API 列表失败，请稍后重试。",
		"FailedToGet":      "获取 API 失败，请稍后重试。",
		"CreationFailed":   "API 创建失败，请稍后重试。",
		"GenerationFailed": "API 生成失败，请稍后重试。",
		"DoesNotExist":     "API 不存在。",
		"FailedToDelete":   "删除 API 失败，请稍后重试。",
		"FailedToMove":     "移动 API 失败，请稍后重试。",
		"CopyFailed":       "API 复制失败，请稍后重试。",
		"ExportFailed":     "API 导出失败，请稍后重试。",
	},
	"collectionHistory": {
		"FailedToGetList": "获取 API 历史列表失败，请稍后重试。",
		"FailedToGet":     "获取 API 历史记录失败，请稍后重试。",
		"DoesNotExist":    "API 历史记录不存在。",
		"RestoreFailed":   "API 恢复失败，请稍后重试。",
		"DiffFailed":      "API 比较失败，请稍后重试。",
	},
	"testCase": {
		"GenerationFailed":   "测试用例生成失败，请稍后重试。",
		"FailedToGet":        "获取测试用例失败，请稍后重试。",
		"DoesNotExist":       "测试用例不存在。",
		"RegenerationFailed": "重新生成测试用例失败，请稍后重试。",
		"FailedToDelete":     "删除测试用例失败，请稍后重试。",
	},
	"mock": {
		"FailedToMock": "Mock 失败，请稍后再试。",
	},
	"sysConfig": {
		"ServiceUpdateFailed":      "服务设置修改失败，请稍后重试。",
		"ServiceBindFailed":        "您要绑定的IP和端口不正确。",
		"ServiceBindPortSame":      "您的应用服务端口和 Mock 服务端口不能相同。",
		"OauthUpdateFailed":        "Oauth 设置失败，请稍后重试。",
		"FailedToGetStorageList":   "获取存储设置失败，请稍后重试。",
		"StorageUpdateFailed":      "存储设置失败，请稍后重试。",
		"LocalPathInvalid":         "存储路径设置有误。",
		"CloudflareConfigInvalid":  "Cloudflare 设置有误。",
		"QiniuConfigInvalid":       "七牛设置有误。",
		"FailedToGetCacheList":     "获取缓存设置失败，请稍后重试。",
		"CacheUpdateFailed":        "缓存设置失败，请稍后重试。",
		"RedisConfigInvalid":       "Redis 设置有误。",
		"FailedToGetEmailList":     "获取邮件设置失败，请稍后重试。",
		"EmailUpdateFailed":        "邮件设置失败，请稍后重试。",
		"SMTPConfigInvalid":        "SMTP 设置有误。",
		"FailedToGetModelList":     "获取模型设置失败，请稍后重试。",
		"ModelUpdateFailed":        "模型设置失败，请稍后重试。",
		"OpenAIConfigInvalid":      "OpenAI 设置有误。",
		"AzureOpenAIConfigInvalid": "Azure OpenAI 设置有误。",
	},
	"jsonschema": {
		"JsonSchemaIncorrect": "解析失败，请检查 JSON Schema 是否正确。",
		"FailedToParse":       "解析失败，请稍后重试。",
	},
}
