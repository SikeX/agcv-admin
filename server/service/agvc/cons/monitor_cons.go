package cons

var NBPointLabelMap = map[int]string{
	1:   "totalPowerGeration",   //总发电量(kWh)"
	2:   "dailyPowerGeration",   //日发电量(kWh)"
	3:   "monthlyPowerGeration", //月发电量(kWh)"
	4:   "annualPowerGeration",  //年发电量(kWh)"
	10:  "ACPower",              //交流功率(kW)"
	202: "DCPower",              //直流功率(kW)"
	27:  "reactivePower",        //无功功率(kVar)"
	28:  "apparentPower",        //视在功率(kVa)"
	15:  "gridFrequency",        //电网频率(Hz)"
	94:  "phaseAVoltage",        //A相电压Ua(V)"
	95:  "phaseBVoltage",        //B相电压Ub(V)"
	96:  "phaseCVoltage",        //C相电压Uc(V)"
	7:   "lineABVoltage",        //AB线电压Uab(V)"
	8:   "lineBCVoltage",        //BC线电压Ubc(V)"
	9:   "lineCAVoltage",        //CA线电压Uca(V)"
	12:  "phaseACurrent",        //A相电流Ia(A)"
	13:  "phaseBCurrent",        //B相电流Ib(A)"
	14:  "phaseCCurrent",        //C相电流Ic(A)"
	21:  "powerFactor",          //功率因数"
	22:  "deviceTemperature",    //设备温度(℃)"
	102: "conversionEfficiency", //转换效率(%)"
	501: "deviceStatusCode",     //设备状态码"
	502: "insulationImpedance",  //绝缘阻抗"
	300: "totalDCVoltage",       //总直流电压",
	301: "DCVoltage1",           //直流电压1"
	302: "DCVoltage2",           //直流电压2"
	303: "DCVoltage3",           //直流电压3"
	304: "DCVoltage4",           //直流电压4"
	305: "DCVoltage5",           //直流电压5"
	306: "DCVoltage6",           //直流电压6"
	307: "DCVoltage7",           //直流电压7"
	308: "DCVoltage8",           //直流电压8"
	309: "DCVoltage9",           //直流电压9"
	310: "DCVoltage10",          //直流电压10"
	311: "DCVoltage11",          //直流电压11"
	312: "DCVoltage12",          //直流电压12"
	313: "DCVoltage13",          //直流电压13"
	314: "DCVoltage14",          //直流电压14"
	315: "DCVoltage15",          //直流电压15"
	316: "DCVoltage16",          //直流电压16"
	317: "DCVoltage17",          //直流电压17"
	318: "DCVoltage18",          //直流电压18"
	319: "DCVoltage19",          //直流电压19"
	320: "DCVoltage20",          //直流电压20"
	16:  "totalDCCurrent",       //总直流电流(A)"
	36:  "DCCurrent1",           //直流电流1"
	37:  "DCCurrent2",           //直流电流2"
	38:  "DCCurrent3",           //直流电流3"
	39:  "DCCurrent4",           //直流电流4"
	40:  "DCCurrent5",           //直流电流5"
	41:  "DCCurrent6",           //直流电流6"
	42:  "DCCurrent7",           //直流电流7"
	43:  "DCCurrent8",           //直流电流8"
	44:  "DCCurrent9",           //直流电流9"
	45:  "DCCurrent10",          //直流电流10"
	46:  "DCCurrent11",          //直流电流11"
	47:  "DCCurrent12",          //直流电流12"
	48:  "DCCurrent13",          //直流电流13"
	49:  "DCCurrent14",          //直流电流14"
	50:  "DCCurrent15",          //直流电流15"
	51:  "DCCurrent16",          //直流电流16"
	52:  "DCCurrent17",          //直流电流17"
	53:  "DCCurrent18",          //直流电流18"
	54:  "DCCurrent19",          //直流电流19"
	55:  "DCCurrent20",          //直流电流20"
	401: "phaseADCPower",        //A相直流功率(Pa)"
	402: "phaseBDCPower",        //B相直流功率(Pb)"
	403: "phaseCDCPower",        //C相直流功率(Pc)"
	404: "phaseAReactivePower",  //A相无功功率(Qa)"
	405: "phaseBReactivePower",  //B相无功功率(Qb)"
	406: "phaseCReactivePower",  //C相无功功率(Qc)"
	407: "phaseAApparentPower",  //A相视在功率(Sa)"
	408: "phaseBApparentPower",  //B相视在功率(Sb)"
	409: "phaseCApparentPower",  //C相视在功率(Sc)"
	410: "phaseAPowerFactor",    //A相功率因数"
	411: "phaseBPowerFactor",    //B相功率因数"
	412: "phaseCPowerFactor",    //C相功率因数"
	512: "shutdownTime",         //关机时间"
	251: "phaseATemperature",    //A相温度"
	252: "phaseBTemperature",    //B相温度"
	253: "phaseCTemperature",    //C相温度"
}

var NBLabelPointMap = map[string]int{
	"totalPowerGeration":   1,   // 总发电量(kWh)
	"dailyPowerGeration":   2,   // 日发电量(kWh)
	"monthlyPowerGeration": 3,   // 月发电量(kWh)
	"annualPowerGeration":  4,   // 年发电量(kWh)
	"ACPower":              10,  // 交流功率(kW)
	"DCPower":              202, // 直流功率(kW)
	"reactivePower":        27,  // 无功功率(kVar)
	"apparentPower":        28,  // 视在功率(kVa)
	"gridFrequency":        15,  // 电网频率(Hz)
	"phaseAVoltage":        94,  // A相电压Ua(V)
	"phaseBVoltage":        95,  // B相电压Ub(V)
	"phaseCVoltage":        96,  // C相电压Uc(V)
	"lineABVoltage":        7,   // AB线电压Uab(V)
	"lineBCVoltage":        8,   // BC线电压Ubc(V)
	"lineCAVoltage":        9,   // CA线电压Uca(V)
	"phaseACurrent":        12,  // A相电流Ia(A)
	"phaseBCurrent":        13,  // B相电流Ib(A)
	"phaseCCurrent":        14,  // C相电流Ic(A)
	"powerFactor":          21,  // 功率因数
	"deviceTemperature":    22,  // 设备温度(℃)
	"conversionEfficiency": 102, // 转换效率(%)
	"deviceStatusCode":     501, // 设备状态码
	"insulationImpedance":  502, // 绝缘阻抗
	"totalDCVoltage":       300, // 总直流电压
	"DCVoltage1":           301, // 直流电压1
	"DCVoltage2":           302, // 直流电压2
	"DCVoltage3":           303, // 直流电压3
	"DCVoltage4":           304, // 直流电压4
	"DCVoltage5":           305, // 直流电压5
	"DCVoltage6":           306, // 直流电压6
	"DCVoltage7":           307, // 直流电压7
	"DCVoltage8":           308, // 直流电压8
	"DCVoltage9":           309, // 直流电压9
	"DCVoltage10":          310, // 直流电压10
	"DCVoltage11":          311, // 直流电压11
	"DCVoltage12":          312, // 直流电压12
	"DCVoltage13":          313, // 直流电压13
	"DCVoltage14":          314, // 直流电压14
	"DCVoltage15":          315, // 直流电压15
	"DCVoltage16":          316, // 直流电压16
	"DCVoltage17":          317, // 直流电压17
	"DCVoltage18":          318, // 直流电压18
	"DCVoltage19":          319, // 直流电压19
	"DCVoltage20":          320, // 直流电压20
	"totalDCCurrent":       16,  // 总直流电流(A)
	"DCCurrent1":           36,  // 直流电流1
	"DCCurrent2":           37,  // 直流电流2
	"DCCurrent3":           38,  // 直流电流3
	"DCCurrent4":           39,  // 直流电流4
	"DCCurrent5":           40,  // 直流电流5
	"DCCurrent6":           41,  // 直流电流6
	"DCCurrent7":           42,  // 直流电流7
	"DCCurrent8":           43,  // 直流电流8
	"DCCurrent9":           44,  // 直流电流9
	"DCCurrent10":          45,  // 直流电流10
	"DCCurrent11":          46,  // 直流电流11
	"DCCurrent12":          47,  // 直流电流12
	"DCCurrent13":          48,  // 直流电流13
	"DCCurrent14":          49,  // 直流电流14
	"DCCurrent15":          50,  // 直流电流15
	"DCCurrent16":          51,  // 直流电流16
	"DCCurrent17":          52,  // 直流电流17
	"DCCurrent18":          53,  // 直流电流18
	"DCCurrent19":          54,  // 直流电流19
	"DCCurrent20":          55,  // 直流电流20
	"phaseADCPower":        401, // A相直流功率(Pa)
	"phaseBDCPower":        402, // B相直流功率(Pb)
	"phaseCDCPower":        403, // C相直流功率(Pc)
	"phaseAReactivePower":  404, // A相无功功率(Qa)
	"phaseBReactivePower":  405, // B相无功功率(Qb)
	"phaseCReactivePower":  406, // C相无功功率(Qc)
	"phaseAApparentPower":  407, // A相视在功率(Sa)
	"phaseBApparentPower":  408, // B相视在功率(Sb)
	"phaseCApparentPower":  409, // C相视在功率(Sc)
	"phaseAPowerFactor":    410, // A相功率因数
	"phaseBPowerFactor":    411, // B相功率因数
	"phaseCPowerFactor":    412, // C相功率因数
	"shutdownTime":         512, // 关机时间
	"phaseATemperature":    251, // A相温度
	"phaseBTemperature":    252, // B相温度
	"phaseCTemperature":    253, // C相温度
}
