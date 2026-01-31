package models

type GaveDoneStatistic struct {
	Gave     int    `json:"gave"`
	Got      int    `json:"got"`
	UserName string `json:"user_name"`
}

type RecordsProMonthStatistic struct {
	Records  ProMonth `json:"records"`
	Services ProMonth `json:"services"`
}

type ProMonth struct {
	January   int `json:"january"`
	February  int `json:"february"`
	March     int `json:"march"`
	April     int `json:"april"`
	May       int `json:"may"`
	June      int `json:"june"`
	July      int `json:"july"`
	August    int `json:"august"`
	September int `json:"september"`
	October   int `json:"october"`
	November  int `json:"november"`
	December  int `json:"december"`
}

type ByType struct {
	Type  string `json:"type"`
	Total int    `json:"total"`
}

type StatisticByType struct {
	Types  []string `json:"types"`
	Amount []int    `json:"amount"`
}

type MainStatistic struct {
	UserServices  int `json:"services"`
	Recieved      int `json:"recieved"`
	Provided      int `json:"provided"`
	Done          int `json:"done"`
	Confirm       int `json:"confirm"`
	Promoservises int `json:"promoservices"`
	Comments      int `json:"comments"`
}
