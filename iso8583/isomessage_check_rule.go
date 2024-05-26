package utils_class

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/moov-io/iso8583"
)

func CheckHmac(msg *iso8583.Message) bool {
	actual, _ := msg.GetField(128).String()

	NAPAS_SOCKET_HMAC := "00112233445566778899AABBCCDDEEFF"

	expected := GenHMAC(msg, NAPAS_SOCKET_HMAC)

	if actual != expected {
		fmt.Println("HMAC is not valid")
	}

	return actual == expected
}

func CheckRule(msg *iso8583.Message, ls []NapasRuleEntity) (isValid bool) {
	// check rule
	m := checkFormat(ls, parseIsoMessageToMap(msg))
	return m["result"] == "1"
}

func parseIsoMessageToMap(msg *iso8583.Message) map[string]string {
	isomsgMap := make(map[string]string)
	isomsgMap["MTI"], _ = msg.GetField(0).String()

	for i := 2; i <= 128; i++ {
		field, _ := msg.GetField(i).String()
		if i < 10 {
			isomsgMap["F0"+fmt.Sprintf("%d", i)] = field
		} else {
			isomsgMap["F"+fmt.Sprintf("%d", i)] = field
		}

		if i == 48 {
			f48Array := strings.Split(field, "\r")
			count := 0
			for _, f48Item := range f48Array {
				if count == 0 {
					isomsgMap["F481"] = f48Item
				}
				if count == 1 {
					isomsgMap["F482"] = f48Item
				}
			}
		}
	}

	return isomsgMap
}

func checkFormat(ls []NapasRuleEntity, isomsgMap map[string]string) map[string]string {
	transType := ""
	route := ""
	checkResult := ""
	resultMap := make(map[string]string)

	if strings.HasPrefix(isomsgMap["F03"][:2], "43") {
		transType = "QUERY"
	}
	if strings.HasPrefix(isomsgMap["F03"][:2], "91") {
		transType = "PAYMENT"
	}

	if isomsgMap["MTI"] == "0210" {
		route = "NHTH"
	}
	if isomsgMap["MTI"] == "0200" {
		route = "NHTH"
	}

	for _, rule := range ls {
		ruleTransType := rule.TransType
		ruleRoute := rule.Route
		ruleFieldName := rule.FieldName
		checkResult += "-" + ruleFieldName

		if (transType == ruleTransType || ruleTransType == "ALL") && (route == ruleRoute || ruleRoute == "ALL") {
			result := checkFormatByValue(rule, isomsgMap[ruleFieldName])
			checkResult += ":" + result
			if result == "30" {
				resultMap["result"] = ruleFieldName
				resultMap["resultlog"] = checkResult
				return resultMap
			}
		}
	}

	resultMap["result"] = "1"
	resultMap["resultlog"] = checkResult
	return resultMap
}

func checkFormatByValue(napasRule NapasRuleEntity, value string) string {
	checkNull := napasRule.CheckNull
	lengthType := napasRule.LengthType
	length := napasRule.Length
	fieldName := napasRule.FieldName
	pattern := napasRule.Pattern
	route := napasRule.Route

	if checkNull == "1" && value == "" {
		return "30"
	}

	if value != "" {
		if lengthType == "EQUAL" && len(value) != int(length) {
			return "30"
		}
		if lengthType == "MAX" && len(value) > int(length) {
			return "30"
		}
	}

	if fieldName == "F104" && pattern != "" {
		if route == "NHTH" || route == "ALL" {
			match, err := regexp.MatchString(pattern, value)
			if err != nil || match {
				return "30"
			}
		}
	}

	if fieldName == "F07" && pattern == "checkMMddHHmmss" {
		if !checkMMddHHmmss(value) {
			return "30"
		}
	}

	if fieldName == "F13" && pattern == "checkMMDD" {
		if !checkMMDD(value) {
			return "30"
		}
	}

	if fieldName == "F12" && pattern == "checkHHmmss" {
		if !checkHHmmss(value) {
			return "30"
		}
	}

	return "1"
}

// checkMMddHHmmss checks if the input string matches the pattern MMddHHmmss
func checkMMddHHmmss(value string) bool {
	// Define the regular expression pattern for MMddHHmmss
	pattern := `^(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])(0[0-9]|1[0-9]|2[0-3])([0-5][0-9]){2}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(value)
}

// checkMMDD checks if the input string matches the pattern MMDD
func checkMMDD(value string) bool {
	// Define the regular expression pattern for MMDD
	pattern := `^(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(value)
}

// checkHHmmss checks if the input string matches the pattern HHmmss
func checkHHmmss(value string) bool {
	// Define the regular expression pattern for HHmmss
	pattern := `^(0[0-9]|1[0-9]|2[0-3])([0-5][0-9]){2}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(value)
}
