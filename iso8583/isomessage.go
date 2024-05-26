package utils_class

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/moov-io/iso8583"
)

func GenQueryMsg(
	requestType Request[GetNapasRecipientNameRequest],
	currentAccountDto CurrentAccountGetCurrentAccountResponse) (output IsoMessage) {
	request := requestType.Data
	// Get current time
	currentTime := NewTimeUtils()

	message := iso8583.NewMessage(Spec)

	// 0200 - request; 0210 – response
	message.MTI("0200")
	//DE #2: Primary Account Number
	if request.FromAccountNumber == "" {
		message.Field(2, "00"+currentTime.GetTimeYYYYMMDDHHMMSS())
	} else {
		message.Field(2, request.FromAccountNumber)
	}

	//DE #3: Processing Code
	var processingCode string = "43"
	if request.FromAccountNumber != "" {
		processingCode = processingCode + "20"
	} else {
		processingCode = processingCode + "00"
	}
	if strings.ToUpper(request.PaymentMethod) == "ACCOUNT" {
		processingCode = processingCode + "20"
	} else {
		processingCode = processingCode + "00"
	}
	message.Field(3, processingCode)
	//DE #4: Transaction Amount
	message.Field(4, PadLeft("", 12, '0'))
	//DE #7: Transmission Date and Time
	message.Field(7, currentTime.GetTimeGmtMMDDHHMMSS())
	//DE #11: System Trace Audit Number
	f11 := GetSystemTrace()
	message.Field(11, f11)
	//DE #12: Local Transaction Time
	message.Field(12, currentTime.GetTimeHHMMSS())
	//DE #13 Local Transaction Date
	message.Field(13, currentTime.GetTimeMMDD())
	//DE #15 Local Transaction Date
	message.Field(15, currentTime.GetTimeMMDD())
	//DE #18: Merchant Category Code
	message.Field(18, "7399")
	//DE #22: PAN mode
	message.Field(22, "000") //(Không biết chế độ PAN được nhập vào + Không biết được khả năng nhập số PIN của thiết bị đầu cuối)
	//DE #25: service condition
	message.Field(25, "05") //(Khách hàng có mặt tại nơi giao dịch nhưng không có thẻ)
	//DE #32: Acquiring Instititution Code
	message.Field(32, "970448") //970448 for OCB
	//DE #37: Retrieval reference number
	message.Field(37, currentTime.GetTimeYDDDHHNNNNNN(f11))
	//DE #38: Authorization identification response
	message.Field(38, "103698") //Napas trả về
	//DE #41 Card Acceptor Terminal Identification
	message.Field(41, "00000001") //Giá trị này ko có tài liệu cụ thể nên Huy-Napas kêu truyền như vậy !
	//DE #42 Card Acceptor Identification Code
	message.Field(42, "000000000000001") //Giá trị này ko có tài liệu cụ thể nên Huy-Napas kêu truyền như vậy !
	//DE #43 Card Acceptor Name/Location
	message.Field(43, padRight("OCB EBANKING", 22, ' ')+" "+padRight("TP HCM", 13, ' ')+" VNM")
	//DE #48: Additional private data
	//tmp = "Nguyen quang khanh" + (char)13 + "Ngan hang OCB";
	//message.Field(48, StringUtil.PadLeft(tmp.length() + "", 3, '0') + tmp);
	message.Field(48, fmt.Sprintf("%s\r", currentAccountDto.AccountName))
	//DE #49 (Currency Code) = 704 (VND)
	message.Field(49, "704")
	//DE #54: Additional amount
	//DE #60: Self – defined Field
	/*- 00: Không xác định
	- 01: ATM
	- 02: Counter (Quầy giao dịch)
	- 03: POS
	- 04: Internet Banking
	- 05: Mobile Application
	- 06: SMS Banking
	- 07: Kênh khác*/
	message.Field(60, "04")
	//DE #62 Service Code
	message.Field(62, "IF_INQ")
	//DE #63 Transaction Reference Number
	//message.Field(63, DateTimeUtil.genF63(dtDateTime.value));//Napas trả về
	//DE #100: Receiving Institution Identification Code
	message.Field(100, request.ToBankCode)
	//DE #90: Original Data Elements
	//message.Field(90, "");//Dùng cho g/d đảo
	//DE #102: From Account Identificati
	message.Field(102, request.FromAccountNumber)
	//DE #103: To Account Identification
	if strings.ToUpper(request.PaymentMethod) == "ACCOUNT" {
		message.Field(103, request.ToAccountNumber)
	} else {
		message.Field(103, request.ToCardNumber)
	}
	//DE #104 : Content transfer
	message.Field(104, "Truy Van Thong Tin")
	//DE #121 : cardMasking
	// message.Field(121, request.ToCreditCard.CardMasking)
	// //DE #122 : encrypted cardMasking
	// message.Field(122, request.ToCreditCard.EncryptedCardNo)
	// //DE #123 : encryptedKey
	// message.Field(123, request.ToCreditCard.EncryptedKey)
	//DE #128 : HMAC
	// TODO: read from env variable
	NAPAS_SOCKET_HMAC := "00112233445566778899AABBCCDDEEFF"
	message.Field(128, GenHMAC(message, NAPAS_SOCKET_HMAC))

	// iso8583.Describe(message, os.Stdout)
	// printIsoMsg(*message)

	return ParseEntity(message)
}

func ParseMessage(s string) *iso8583.Message {
	msg := iso8583.NewMessage(Spec)

	err := msg.Unpack([]byte(s))
	if err != nil {
		fmt.Println(err)
	}
	iso8583.Describe(msg, os.Stdout)

	entity := ParseEntity(msg)
	fmt.Println(entity)
	return msg
}

func ParseEntity(message *iso8583.Message) (output IsoMessage) {
	output.Field1, _ = message.GetField(1).String()
	output.Field2, _ = message.GetField(2).String()
	output.Field3, _ = message.GetField(3).String()
	output.Field4, _ = message.GetField(4).String()
	output.Field5, _ = message.GetField(5).String()
	output.Field6, _ = message.GetField(6).String()
	output.Field7, _ = message.GetField(7).String()
	output.Field8, _ = message.GetField(8).String()
	output.Field9, _ = message.GetField(9).String()
	output.Field10, _ = message.GetField(10).String()
	output.Field11, _ = message.GetField(11).String()
	output.Field12, _ = message.GetField(12).String()
	output.Field13, _ = message.GetField(13).String()
	output.Field14, _ = message.GetField(14).String()
	output.Field15, _ = message.GetField(15).String()
	output.Field16, _ = message.GetField(16).String()
	output.Field17, _ = message.GetField(17).String()
	output.Field18, _ = message.GetField(18).String()
	output.Field19, _ = message.GetField(19).String()
	output.Field20, _ = message.GetField(20).String()
	output.Field21, _ = message.GetField(21).String()
	output.Field22, _ = message.GetField(22).String()
	output.Field23, _ = message.GetField(23).String()
	output.Field24, _ = message.GetField(24).String()
	output.Field25, _ = message.GetField(25).String()
	output.Field26, _ = message.GetField(26).String()
	output.Field27, _ = message.GetField(27).String()
	output.Field28, _ = message.GetField(28).String()
	output.Field29, _ = message.GetField(29).String()
	output.Field30, _ = message.GetField(30).String()
	output.Field31, _ = message.GetField(31).String()
	output.Field32, _ = message.GetField(32).String()
	output.Field33, _ = message.GetField(33).String()
	output.Field34, _ = message.GetField(34).String()
	output.Field35, _ = message.GetField(35).String()
	output.Field36, _ = message.GetField(36).String()
	output.Field37, _ = message.GetField(37).String()
	output.Field38, _ = message.GetField(38).String()
	output.Field39, _ = message.GetField(39).String()
	output.Field40, _ = message.GetField(40).String()
	output.Field41, _ = message.GetField(41).String()
	output.Field42, _ = message.GetField(42).String()
	output.Field43, _ = message.GetField(43).String()
	output.Field44, _ = message.GetField(44).String()
	output.Field45, _ = message.GetField(45).String()
	output.Field46, _ = message.GetField(46).String()
	output.Field47, _ = message.GetField(47).String()
	output.Field48, _ = message.GetField(48).String()
	output.Field49, _ = message.GetField(49).String()
	output.Field50, _ = message.GetField(50).String()
	output.Field51, _ = message.GetField(51).String()
	output.Field52, _ = message.GetField(52).String()
	output.Field53, _ = message.GetField(53).String()
	output.Field54, _ = message.GetField(54).String()
	output.Field55, _ = message.GetField(55).String()
	output.Field56, _ = message.GetField(56).String()
	output.Field57, _ = message.GetField(57).String()
	output.Field58, _ = message.GetField(58).String()
	output.Field59, _ = message.GetField(59).String()
	output.Field60, _ = message.GetField(60).String()
	output.Field61, _ = message.GetField(61).String()
	output.Field62, _ = message.GetField(62).String()
	output.Field63, _ = message.GetField(63).String()
	output.Field64, _ = message.GetField(64).String()
	output.Field65, _ = message.GetField(65).String()
	output.Field66, _ = message.GetField(66).String()
	output.Field67, _ = message.GetField(67).String()
	output.Field68, _ = message.GetField(68).String()
	output.Field69, _ = message.GetField(69).String()
	output.Field70, _ = message.GetField(70).String()
	output.Field71, _ = message.GetField(71).String()
	output.Field72, _ = message.GetField(72).String()
	output.Field73, _ = message.GetField(73).String()
	output.Field74, _ = message.GetField(74).String()
	output.Field75, _ = message.GetField(75).String()
	output.Field76, _ = message.GetField(76).String()
	output.Field77, _ = message.GetField(77).String()
	output.Field78, _ = message.GetField(78).String()
	output.Field79, _ = message.GetField(79).String()
	output.Field80, _ = message.GetField(80).String()
	output.Field81, _ = message.GetField(81).String()
	output.Field82, _ = message.GetField(82).String()
	output.Field83, _ = message.GetField(83).String()
	output.Field84, _ = message.GetField(84).String()
	output.Field85, _ = message.GetField(85).String()
	output.Field86, _ = message.GetField(86).String()
	output.Field87, _ = message.GetField(87).String()
	output.Field88, _ = message.GetField(88).String()
	output.Field89, _ = message.GetField(89).String()
	output.Field90, _ = message.GetField(90).String()
	output.Field91, _ = message.GetField(91).String()
	output.Field92, _ = message.GetField(92).String()
	output.Field93, _ = message.GetField(93).String()
	output.Field94, _ = message.GetField(94).String()
	output.Field95, _ = message.GetField(95).String()
	output.Field96, _ = message.GetField(96).String()
	output.Field97, _ = message.GetField(97).String()
	output.Field98, _ = message.GetField(98).String()
	output.Field99, _ = message.GetField(99).String()
	output.Field100, _ = message.GetField(100).String()
	output.Field101, _ = message.GetField(101).String()
	output.Field102, _ = message.GetField(102).String()
	output.Field103, _ = message.GetField(103).String()
	output.Field104, _ = message.GetField(104).String()
	output.Field105, _ = message.GetField(105).String()
	output.Field106, _ = message.GetField(106).String()
	output.Field107, _ = message.GetField(107).String()
	output.Field108, _ = message.GetField(108).String()
	output.Field109, _ = message.GetField(109).String()
	output.Field110, _ = message.GetField(110).String()
	output.Field111, _ = message.GetField(111).String()
	output.Field112, _ = message.GetField(112).String()
	output.Field113, _ = message.GetField(113).String()
	output.Field114, _ = message.GetField(114).String()
	output.Field115, _ = message.GetField(115).String()
	output.Field116, _ = message.GetField(116).String()
	output.Field117, _ = message.GetField(117).String()
	output.Field118, _ = message.GetField(118).String()
	output.Field119, _ = message.GetField(119).String()
	output.Field120, _ = message.GetField(120).String()
	output.Field121, _ = message.GetField(121).String()
	output.Field122, _ = message.GetField(122).String()
	output.Field123, _ = message.GetField(123).String()
	output.Field124, _ = message.GetField(124).String()
	output.Field125, _ = message.GetField(125).String()
	output.Field126, _ = message.GetField(126).String()
	output.Field127, _ = message.GetField(127).String()
	output.Field128, _ = message.GetField(128).String()

	dataByte, _ := message.Pack()
	output.IsoMessage = string(dataByte)

	output.IsoMessage = PadLeft(fmt.Sprintf("%d", len(string(dataByte))), 4, '0') + string(dataByte)

	return output
}

func PadLeft(str string, length int, padChar byte) string {
	for len(str) < length {
		str = string(padChar) + str
	}
	return str
}

func padRight(str string, length int, padChar byte) string {
	for len(str) < length {
		str = str + string(padChar)
	}
	return str
}

func GenHMAC(isoMsg *iso8583.Message, macKey string) string {
	var macDefine = map[string]string{
		"2":   "LLVAR",
		"32":  "LLVAR",
		"48":  "LLLVAR",
		"63":  "LLLVAR",
		"102": "LLVAR",
		"103": "LLVAR",
		"120": "LLLVAR",
	}

	macList := []string{"2", "3", "4", "5", "6", "7", "11", "32", "37", "38", "39", "41", "42", "48", "63", "66", "90", "102", "103", "120"}

	data4Mac, err := isoMsg.GetMTI()
	if err != nil {
		log.Fatal(err)
	}

	for _, fieldId := range macList {
		fieldID, err := strconv.Atoi(fieldId)
		if err != nil {
			continue
			log.Fatal(err)
		}
		field, err := isoMsg.GetField(fieldID).String()
		if err != nil {
			log.Fatal(err)
		}
		if field != "" {
			l := len(field)
			var len string
			if definition, ok := macDefine[fieldId]; ok {
				if definition == "LLVAR" {
					len = fillZero(l, 2)
					data4Mac = data4Mac + len + field
				} else if definition == "LLLVAR" {
					len = fillZero(l, 3)
					data4Mac = data4Mac + len + field
				}

			} else {
				data4Mac = data4Mac + field
			}
		} else {
			data4Mac = data4Mac + ""
		}
	}

	return genHMac(data4Mac, macKey)
}

func genHMac(data string, macKey string) string {
	secretKey := []byte(macKey)
	algorithm := sha256.New
	mac := hmac.New(algorithm, secretKey)
	_, err := mac.Write([]byte(data))
	if err != nil {
		log.Printf("error writing data to HMAC: %+v", err)
	}
	hmacBytes := mac.Sum(nil)
	hmacHex := hex.EncodeToString(hmacBytes)
	return strings.ToUpper(hmacHex)
}

func fillZero(len int, length int) string {
	if length == 2 && len < 10 {
		return "0" + strconv.Itoa(len)
	} else if length == 3 && len < 10 {
		return "00" + strconv.Itoa(len)
	} else if length == 3 && len < 100 && len >= 10 {
		return "0" + strconv.Itoa(len)
	} else {
		return strconv.Itoa(len)
	}
}

func GetSystemTrace() string {
	// Define minimum and maximum values
	min := 1
	max := 999999

	// Generate random number between min and max (inclusive)
	randomValue := rand.Intn(max-min+1) + min
	stringValue := strconv.Itoa(randomValue)
	return PadLeft(stringValue+"", 6, '0')
}
