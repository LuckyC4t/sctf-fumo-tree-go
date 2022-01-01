package data

var Sanitizer = map[string]bool{
	"crypt":         true,
	"md5":           true,
	"sha1":          true,
	"base64_encode": true, // 懒得写计数器了, 能找出路径就行
}
