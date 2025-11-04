package utils

func GetPsInfo(code string) (string, string, string, string, string) {
	//电站id（3） 设备id（4） 数据类型（2） 设备类型（2）点号（3）
	if len(code) < 14 {
		return "", "", "", "", ""
	}
	return code[:3], code[3:7], code[7:9], code[9:11], code[11:14]
}
