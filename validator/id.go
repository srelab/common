package validator

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

var (
	ErrFormatInvalid      = errors.New("格式错误")
	ErrAddressInvalid     = errors.New("地址码错误")
	ErrBirthFormatInvalid = errors.New("出生日期格式错误")
	ErrBirthRangeInvalid  = errors.New("出生日期范围错误")
	ErrSumInvalid         = errors.New("校验和错误")

	reg  = regexp.MustCompile("^(\\d{6})(18|19|20)?(\\d{2})(0\\d|10|11|12)([012]\\d|3[01])(\\d{3})(\\d|X)?$")
	area = map[string]string{
		"11": "北京", "12": "天津", "13": "河北", "14": "山西", "15": "内蒙",
		"21": "辽宁", "22": "吉林", "23": "黑龙", " 31": "上海", "32": "江苏",
		"33": "浙江", "34": "安徽", "35": "福建", "36": "江西", "37": "山东",
		"41": "河南", "42": "湖北", "43": "湖南", "44": "广东", "45": "广西",
		"46": "海南", "50": "重庆", "51": "四川", "52": "贵州", "53": "云南",
		"54": "西藏", "61": "陕西", "62": "甘肃", "63": "青海", "64": "宁夏",
		"65": "新疆", "71": "台湾", "81": "香港", "82": "澳门", "91": "国外",
	}

	min_date = time.Date(1890, 0, 0, 0, 0, 0, 0, time.Local)
	max_date = time.Now()

	//十七位数字本体码权重
	weight = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	code   = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
)

type IDCard struct {
	Number string
}

//整体校验格式
func (i *IDCard) validateReg() error {
	if reg.MatchString(i.Number) {
		return nil
	}
	return ErrFormatInvalid
}

//校验地区码
func (i *IDCard) validateArea() error {
	if _, ok := area[i.Number[0:2]]; ok {
		return nil
	}

	return ErrAddressInvalid
}

//校验生日,包括格式和范围
func (i *IDCard) validateBirth() error {
	birth := i.Number[6:14]
	if date, err := time.Parse("20060102", birth); err != nil {
		return ErrBirthFormatInvalid
	} else if date.After(max_date) && date.Before(min_date) {
		return ErrBirthRangeInvalid
	}
	return nil
}

//校验和
func (i *IDCard) validateSum() error {
	sum := 0
	for i, char := range i.Number[:len(i.Number)-1] {
		cf, _ := strconv.ParseFloat(string(char), 64)
		sum += int(cf) * weight[i]
	}
	if code[sum%11] == i.Number[len(i.Number)-1] {
		return nil
	}
	return ErrSumInvalid
}

// 校验
func (i *IDCard) Validate() (flag bool, err error) {
	if err = i.validateReg(); err != nil {
		return false, err
	}

	if err = i.validateArea(); err != nil {
		return false, err
	}

	if err = i.validateBirth(); err != nil {
		return false, err
	}

	if err = i.validateSum(); err != nil {
		return false, err
	}

	return true, nil
}
