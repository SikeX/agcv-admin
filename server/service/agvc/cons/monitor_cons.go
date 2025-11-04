package cons

var NBPointLabelMap = map[string]string{
	"001":  "totalPowerGeration",   //总发电量(kWh)"
	"002":  "dailyPowerGeration",   //日发电量(kWh)"
	"003":  "monthlyPowerGeration", //月发电量(kWh)"
	"004":  "annualPowerGeration",  //年发电量(kWh)"
	"010":  "ACPower",              //交流功率(kW)"
	"0202": "DCPower",              //直流功率(kW)"
	"027":  "reactivePower",        //无功功率(kVar)"
	"028":  "apparentPower",        //视在功率(kVa)"
	"015":  "gridFrequency",        //电网频率(Hz)"
	"094":  "phaseAVoltage",        //A相电压Ua(V)"
	"095":  "phaseBVoltage",        //B相电压Ub(V)"
	"096":  "phaseCVoltage",        //C相电压Uc(V)"
	"007":  "lineABVoltage",        //AB线电压Uab(V)"
	"008":  "lineBCVoltage",        //BC线电压Ubc(V)"
	"009":  "lineCAVoltage",        //CA线电压Uca(V)"
	"012":  "phaseACurrent",        //A相电流Ia(A)"
	"013":  "phaseBCurrent",        //B相电流Ib(A)"
	"014":  "phaseCCurrent",        //C相电流Ic(A)"
	"021":  "powerFactor",          //功率因数"
	"022":  "deviceTemperature",    //设备温度(℃)"
	"102":  "conversionEfficiency", //转换效率(%)"
	"501":  "deviceStatusCode",     //设备状态码"
	"502":  "insulationImpedance",  //绝缘阻抗"
	"300":  "totalDCVoltage",       //总直流电压",
	"301":  "DCVoltage1",           //直流电压1"
	"302":  "DCVoltage2",           //直流电压2"
	"303":  "DCVoltage3",           //直流电压3"
	"304":  "DCVoltage4",           //直流电压4"
	"305":  "DCVoltage5",           //直流电压5"
	"306":  "DCVoltage6",           //直流电压6"
	"307":  "DCVoltage7",           //直流电压7"
	"308":  "DCVoltage8",           //直流电压8"
	"309":  "DCVoltage9",           //直流电压9"
	"310":  "DCVoltage10",          //直流电压10"
	"311":  "DCVoltage11",          //直流电压11"
	"312":  "DCVoltage12",          //直流电压12"
	"313":  "DCVoltage13",          //直流电压13"
	"314":  "DCVoltage14",          //直流电压14"
	"315":  "DCVoltage15",          //直流电压15"
	"316":  "DCVoltage16",          //直流电压16"
	"317":  "DCVoltage17",          //直流电压17"
	"318":  "DCVoltage18",          //直流电压18"
	"319":  "DCVoltage19",          //直流电压19"
	"320":  "DCVoltage20",          //直流电压20"
	"016":  "totalDCCurrent",       //总直流电流(A)"
	"036":  "DCCurrent1",           //直流电流1"
	"037":  "DCCurrent2",           //直流电流2"
	"038":  "DCCurrent3",           //直流电流3"
	"039":  "DCCurrent4",           //直流电流4"
	"040":  "DCCurrent5",           //直流电流5"
	"041":  "DCCurrent6",           //直流电流6"
	"042":  "DCCurrent7",           //直流电流7"
	"043":  "DCCurrent8",           //直流电流8"
	"044":  "DCCurrent9",           //直流电流9"
	"045":  "DCCurrent10",          //直流电流10"
	"046":  "DCCurrent11",          //直流电流11"
	"047":  "DCCurrent12",          //直流电流12"
	"048":  "DCCurrent13",          //直流电流13"
	"049":  "DCCurrent14",          //直流电流14"
	"050":  "DCCurrent15",          //直流电流15"
	"051":  "DCCurrent16",          //直流电流16"
	"052":  "DCCurrent17",          //直流电流17"
	"053":  "DCCurrent18",          //直流电流18"
	"054":  "DCCurrent19",          //直流电流19"
	"055":  "DCCurrent20",          //直流电流20"
	"401":  "phaseADCPower",        //A相直流功率(Pa)"
	"402":  "phaseBDCPower",        //B相直流功率(Pb)"
	"403":  "phaseCDCPower",        //C相直流功率(Pc)"
	"404":  "phaseAReactivePower",  //A相无功功率(Qa)"
	"405":  "phaseBReactivePower",  //B相无功功率(Qb)"
	"406":  "phaseCReactivePower",  //C相无功功率(Qc)"
	"407":  "phaseAApparentPower",  //A相视在功率(Sa)"
	"408":  "phaseBApparentPower",  //B相视在功率(Sb)"
	"409":  "phaseCApparentPower",  //C相视在功率(Sc)"
	"410":  "phaseAPowerFactor",    //A相功率因数"
	"411":  "phaseBPowerFactor",    //B相功率因数"
	"412":  "phaseCPowerFactor",    //C相功率因数"
	"512":  "shutdownTime",         //关机时间"
	"251":  "phaseATemperature",    //A相温度"
	"252":  "phaseBTemperature",    //B相温度"
	"253":  "phaseCTemperature",    //C相温度"
}

var NBLabelPointMap = map[string]string{
	"totalPowerGeration":   "001", // 总发电量(kWh)
	"dailyPowerGeration":   "002", // 日发电量(kWh)
	"monthlyPowerGeration": "003", // 月发电量(kWh)
	"annualPowerGeration":  "004", // 年发电量(kWh)
	"ACPower":              "010", // 交流功率(kW)
	"DCPower":              "202", // 直流功率(kW)
	"reactivePower":        "027", // 无功功率(kVar)
	"apparentPower":        "028", // 视在功率(kVa)
	"gridFrequency":        "015", // 电网频率(Hz)
	"phaseAVoltage":        "094", // A相电压Ua(V)
	"phaseBVoltage":        "095", // B相电压Ub(V)
	"phaseCVoltage":        "096", // C相电压Uc(V)
	"lineABVoltage":        "007", // AB线电压Uab(V)
	"lineBCVoltage":        "008", // BC线电压Ubc(V)
	"lineCAVoltage":        "109", // CA线电压Uca(V)
	"phaseACurrent":        "012", // A相电流Ia(A)
	"phaseBCurrent":        "013", // B相电流Ib(A)
	"phaseCCurrent":        "014", // C相电流Ic(A)
	"powerFactor":          "021", // 功率因数
	"deviceTemperature":    "022", // 设备温度(℃)
	"conversionEfficiency": "102", // 转换效率(%)
	"deviceStatusCode":     "501", // 设备状态码
	"insulationImpedance":  "502", // 绝缘阻抗
	"totalDCVoltage":       "300", // 总直流电压
	"DCVoltage1":           "301", // 直流电压1
	"DCVoltage2":           "302", // 直流电压2
	"DCVoltage3":           "303", // 直流电压3
	"DCVoltage4":           "304", // 直流电压4
	"DCVoltage5":           "305", // 直流电压5
	"DCVoltage6":           "306", // 直流电压6
	"DCVoltage7":           "307", // 直流电压7
	"DCVoltage8":           "308", // 直流电压8
	"DCVoltage9":           "309", // 直流电压9
	"DCVoltage10":          "310", // 直流电压10
	"DCVoltage11":          "311", // 直流电压11
	"DCVoltage12":          "312", // 直流电压12
	"DCVoltage13":          "313", // 直流电压13
	"DCVoltage14":          "314", // 直流电压14
	"DCVoltage15":          "315", // 直流电压15
	"DCVoltage16":          "316", // 直流电压16
	"DCVoltage17":          "317", // 直流电压17
	"DCVoltage18":          "318", // 直流电压18
	"DCVoltage19":          "319", // 直流电压19
	"DCVoltage20":          "320", // 直流电压20
	"totalDCCurrent":       "016", // 总直流电流(A)
	"DCCurrent1":           "036", // 直流电流1
	"DCCurrent2":           "037", // 直流电流2
	"DCCurrent3":           "038", // 直流电流3
	"DCCurrent4":           "039", // 直流电流4
	"DCCurrent5":           "040", // 直流电流5
	"DCCurrent6":           "041", // 直流电流6
	"DCCurrent7":           "042", // 直流电流7
	"DCCurrent8":           "043", // 直流电流8
	"DCCurrent9":           "044", // 直流电流9
	"DCCurrent10":          "045", // 直流电流10
	"DCCurrent11":          "046", // 直流电流11
	"DCCurrent12":          "047", // 直流电流12
	"DCCurrent13":          "048", // 直流电流13
	"DCCurrent14":          "049", // 直流电流14
	"DCCurrent15":          "050", // 直流电流15
	"DCCurrent16":          "051", // 直流电流16
	"DCCurrent17":          "052", // 直流电流17
	"DCCurrent18":          "053", // 直流电流18
	"DCCurrent19":          "054", // 直流电流19
	"DCCurrent20":          "055", // 直流电流20
	"phaseADCPower":        "401", // A相直流功率(Pa)
	"phaseBDCPower":        "402", // B相直流功率(Pb)
	"phaseCDCPower":        "403", // C相直流功率(Pc)
	"phaseAReactivePower":  "404", // A相无功功率(Qa)
	"phaseBReactivePower":  "405", // B相无功功率(Qb)
	"phaseCReactivePower":  "406", // C相无功功率(Qc)
	"phaseAApparentPower":  "407", // A相视在功率(Sa)
	"phaseBApparentPower":  "408", // B相视在功率(Sb)
	"phaseCApparentPower":  "409", // C相视在功率(Sc)
	"phaseAPowerFactor":    "410", // A相功率因数
	"phaseBPowerFactor":    "411", // B相功率因数
	"phaseCPowerFactor":    "412", // C相功率因数
	"shutdownTime":         "512", // 关机时间
	"phaseATemperature":    "251", // A相温度
	"phaseBTemperature":    "252", // B相温度
	"phaseCTemperature":    "253", // C相温度
}
