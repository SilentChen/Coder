package utils

func GetGraphShapeSetting() []map[string]string {
	return []map[string]string{
		{"label":"Box",			"value":"box"},
		{"label":"Circle",		"value":"circle"},
		{"label":"Record", 		"value":"record"},
		{"label":"Egg", 		"value":"egg"},
		{"label":"Oval", 		"value":"oval"},
		{"label":"Diamond",		"value":"diamond"},
		{"label":"Triangle",	"value":"triangle"},
		{"label":"Invtriangle",	"value":"invtriangle"},
		{"label":"Rect",		"value":"rect"},
		{"label":"Square",		"value":"square"},
		{"label":"Rectangle",	"value":"rectangle"},
		{"label":"Star",		"value":"star"},
	}
}

func GetGraphColorSetting() []map[string]string {
	return []map[string]string {
		{"label":"白色",			"value":"white"},
		{"label":"红色",			"value":"red"},
		{"label":"蓝色",			"value":"blue"},
		{"label":"绿色",			"value":"green"},
		{"label":"金色",			"value":"gold"},
		{"label":"粉红",			"value":"pink"},
		{"label":"灰色",			"value":"gray"},
	}
}
