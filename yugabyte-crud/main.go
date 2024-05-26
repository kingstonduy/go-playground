package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
	configuration "yugabyte-crud/config"
	data_yugabyte "yugabyte-crud/data"
	"yugabyte-crud/domain"

	"github.com/google/uuid"
	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
)

func main() {
	db := configuration.GetYugabyteDatabase(configuration.GetYugabyteConfig())

	repository := data_yugabyte.NewCRUDRepository(db)

	isoMsg := GenQueryMsg()

	entity, _ := parseIsoMessageToEntity(isoMsg)

	repository.Insert(context.Background(), entity)

}

func parseIsoMessageToEntity(isoMessage iso8583.Message) (entity domain.SmartlinkRetrieveTraceEntity, err error) {
	entity.Field0, err = isoMessage.GetField(0).String()
	if err != nil {
		return entity, err
	}

	entity.Field2, err = isoMessage.GetField(2).String()
	if err != nil {
		return entity, err
	}

	entity.Field3, err = isoMessage.GetField(3).String()
	if err != nil {
		return entity, err
	}

	entity.Field4, err = isoMessage.GetField(4).String()
	if err != nil {
		return entity, err
	}

	entity.Field7, err = isoMessage.GetField(7).String()
	if err != nil {
		return entity, err
	}

	entity.Field11, err = isoMessage.GetField(11).String()
	if err != nil {
		return entity, err
	}

	entity.Field12, err = isoMessage.GetField(12).String()
	if err != nil {
		return entity, err
	}

	entity.Field13, err = isoMessage.GetField(13).String()
	if err != nil {
		return entity, err
	}

	entity.Field15, err = isoMessage.GetField(15).String()
	if err != nil {
		return entity, err
	}

	entity.Field18, err = isoMessage.GetField(18).String()
	if err != nil {
		return entity, err
	}

	entity.Field19, err = isoMessage.GetField(19).String()
	if err != nil {
		return entity, err
	}

	entity.Field22, err = isoMessage.GetField(22).String()
	if err != nil {
		return entity, err
	}

	entity.Field25, err = isoMessage.GetField(25).String()
	if err != nil {
		return entity, err
	}

	entity.Field32, err = isoMessage.GetField(32).String()
	if err != nil {
		return entity, err
	}

	entity.Field37, err = isoMessage.GetField(37).String()
	if err != nil {
		return entity, err
	}

	entity.Field38, err = isoMessage.GetField(38).String()
	if err != nil {
		return entity, err
	}

	// fixed hard code
	entity.Field39 = "O2"
	if err != nil {
		return entity, err
	}

	entity.Field41, err = isoMessage.GetField(41).String()
	if err != nil {
		return entity, err
	}

	entity.Field42, err = isoMessage.GetField(42).String()
	if err != nil {
		return entity, err
	}

	entity.Field43, err = isoMessage.GetField(43).String()
	if err != nil {
		return entity, err
	}

	entity.Field48, err = isoMessage.GetField(48).String()
	if err != nil {
		return entity, err
	}

	entity.Field49, err = isoMessage.GetField(49).String()
	if err != nil {
		return entity, err
	}

	entity.Field60, err = isoMessage.GetField(60).String()
	if err != nil {
		return entity, err
	}

	entity.Field62, err = isoMessage.GetField(62).String()
	if err != nil {
		return entity, err
	}

	entity.Field100, err = isoMessage.GetField(100).String()
	if err != nil {
		return entity, err
	}

	entity.Field102, err = isoMessage.GetField(102).String()
	if err != nil {
		return entity, err
	}

	entity.Field103, err = isoMessage.GetField(103).String()
	if err != nil {
		return entity, err
	}

	entity.Field104, err = isoMessage.GetField(104).String()
	if err != nil {
		return entity, err
	}

	entity.Field128OCB, err = isoMessage.GetField(128).String()
	if err != nil {
		return entity, err
	}

	// channel from input
	entity.Channel = ""
	if err != nil {
		return entity, err
	}

	entity.TraceType = "OCB2SL"
	entity.TransID = uuid.NewString()
	entity.ClientID = uuid.NewString()

	entity.CardMasking, err = isoMessage.GetField(121).String()
	if err != nil {
		return entity, err
	}

	entity.EncryptedKey, err = isoMessage.GetField(122).String()
	if err != nil {
		return entity, err
	}

	entity.EncryptedCardNo, err = isoMessage.GetField(123).String()
	if err != nil {
		return entity, err
	}

	entity.WorkflowID = uuid.NewString()
	if err != nil {
		return entity, err
	}

	entity.Status = "0"
	entity.PrevStatus = "0"
	entity.ProcessingStatus = "0"

	now := time.Now().Unix()

	entity.CreatedDate = fmt.Sprintf("%d", now)
	entity.LastUpdated = fmt.Sprintf("%d", now)
	return entity, nil
}

func GenQueryMsg() iso8583.Message {
	// Get current time
	currentTime := NewTimeUtils()

	message := iso8583.NewMessage(Spec)

	// 0200 - request; 0210 – response
	message.MTI("0200")
	//DE #2: Primary Account Number
	message.Field(2, "00"+currentTime.GetTimeYYYYMMDDHHMMSS())
	//DE #3: Processing Code
	var processingCode string = "430000"
	message.Field(3, processingCode)
	//DE #4: Transaction Amount
	message.Field(4, padLeft("", 12, '0'))
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
	//message.Field(48, StringUtil.padLeft(tmp.length() + "", 3, '0') + tmp);
	message.Field(48, "")
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
	message.Field(100, "")
	//DE #90: Original Data Elements
	//message.Field(90, "");//Dùng cho g/d đảo
	//DE #102: From Account Identificati
	message.Field(102, "")
	//DE #103: To Account Identification
	message.Field(103, "")
	//DE #104 : Content transfer
	message.Field(104, "Truy Van Thong Tin")
	//DE #121 : cardMasking
	message.Field(121, "")
	//DE #122 : encrypted cardMasking
	message.Field(122, "")
	//DE #123 : encryptedKey
	message.Field(123, "")
	//DE #128 : HMAC
	// TODO: read from env variable
	NAPAS_SOCKET_HMAC := "00112233445566778899AABBCCDDEEFF"
	message.Field(128, genHMAC(*message, NAPAS_SOCKET_HMAC))

	iso8583.Describe(message, os.Stdout)

	printIsoMsg(*message)
	return *message
}

func genHMAC(isoMsg iso8583.Message, macKey string) string {
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
	return hmacHex
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
	return padLeft(stringValue+"", 6, '0')
}

func printIsoMsg(isoMsg iso8583.Message) {
	msg := "<isomsg>" + "<!-- " + isoMsg.GetSpec().Name + " -->"
	keys := make([]int, 0)
	for k, _ := range isoMsg.GetFields() {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		field, err := isoMsg.GetField(k).String()
		if err == nil {
			msg += "<field id=\"" + strconv.Itoa(k) + "\" value=\"" + field + "\"/>"
		}
	}
	msg += "</isomsg>"
	fmt.Println(msg)
	data, err := isoMsg.Pack()
	if err != nil {
		log.Fatal(err)
	}
	iso8583.Describe(&isoMsg, os.Stdout)
	fmt.Println("parseIsoMsg: " + string(data))
}

var example = iso8583.Spec87
var Spec = &iso8583.MessageSpec{
	Fields: map[int]field.Field{
		0: field.NewString(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		// field 1 ok
		1: field.NewBitmap(&field.Spec{
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		}),
		2: field.NewString(&field.Spec{
			Length:      19,
			Description: "PAN - PRIMARY ACCOUNT NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		3: field.NewString(&field.Spec{
			Length:      6,
			Description: "PROCESSING CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		4: field.NewString(&field.Spec{
			Length:      12,
			Description: "AMOUNT,TRANSACTION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		5: field.NewString(&field.Spec{
			Length:      12,
			Description: "AMOUNT,SETTLEMENT",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		6: field.NewString(&field.Spec{
			Length:      12,
			Description: "AMOUNT,CARDHOLDER BILLING",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		7: field.NewString(&field.Spec{
			Length:      10,
			Description: "TRANSMISSION DATE AND TIME",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		8: field.NewString(&field.Spec{
			Length:      8,
			Description: "AMOUNT,CARDHOLDER BILLING FEE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		9: field.NewString(&field.Spec{
			Length:      8,
			Description: "CONVERSION RATE,SETTLEMENT,",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		10: field.NewString(&field.Spec{
			Length:      8,
			Description: "CONVERSION RATE,CARDHOLDER BILLING",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		11: field.NewString(&field.Spec{
			Length:      6,
			Description: "SYSTEM TRACE AUDIT NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		12: field.NewString(&field.Spec{
			Length:      6,
			Description: "TIME,LOCAL TRANSACTION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		13: field.NewString(&field.Spec{
			Length:      4,
			Description: "DATE,LOCAL TRANSACTION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		14: field.NewString(&field.Spec{
			Length:      4,
			Description: "DATE,EXPIRATION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		15: field.NewString(&field.Spec{
			Length:      4,
			Description: "DATE,SETTLEMENT",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		16: field.NewString(&field.Spec{
			Length:      4,
			Description: "DATE,CONVERSION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		17: field.NewString(&field.Spec{
			Length:      4,
			Description: "DATE,CAPTURE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		18: field.NewString(&field.Spec{
			Length:      4,
			Description: "MERCHANTS TYPE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		19: field.NewString(&field.Spec{
			Length:      3,
			Description: "ACQUIRING INSTITUTION COUNTRY CODE,",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		20: field.NewString(&field.Spec{
			Length:      3,
			Description: "PAN EXTENDED COUNTRY CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		21: field.NewString(&field.Spec{
			Length:      3,
			Description: "FORWARDING INSTITUTION COUNTRY CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		22: field.NewString(&field.Spec{
			Length:      3,
			Description: "POINT OF SERVICE ENTRY MODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		23: field.NewString(&field.Spec{
			Length:      3,
			Description: "CARD SEQUENCE NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		24: field.NewString(&field.Spec{
			Length:      3,
			Description: "NETWORK INTERNATIONAL IDENTIFIEER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		25: field.NewString(&field.Spec{
			Length:      2,
			Description: "POINT OF SERVICE CONDITION CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		26: field.NewString(&field.Spec{
			Length:      2,
			Description: "POINT OF SERVICE PIN CAPTURE CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		27: field.NewString(&field.Spec{
			Length:      1,
			Description: "AUTHORIZATION IDENTIFICATION RESP LEN",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		28: field.NewString(&field.Spec{
			Length:      9,
			Description: "AMOUNT,TRANSACTION FEE",
			Enc:         encoding.ASCII,
			Pad:         padding.Left('0'),
		}),
		29: field.NewString(&field.Spec{
			Length:      9,
			Description: "AMOUNT,SETTLEMENT FEE,",
			Enc:         encoding.ASCII,
			Pad:         padding.Left('0'),
		}),
		30: field.NewString(&field.Spec{
			Length:      9,
			Description: "AMOUNT,TRANSACTION PROCESSING FEE",
			Enc:         encoding.ASCII,
			Pad:         padding.Left('0'),
		}),
		31: field.NewString(&field.Spec{
			Length:      9,
			Description: "AMOUNT,SETTLEMENT PROCESSING FEE",
			Enc:         encoding.ASCII,
			Pad:         padding.Left('0'),
		}),
		32: field.NewString(&field.Spec{
			Length:      11,
			Description: "ACQUIRING INSTITUTION IDENT CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		33: field.NewString(&field.Spec{
			Length:      11,
			Description: "FORWARDING INSTITUTION IDENT CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		34: field.NewString(&field.Spec{
			Length:      28,
			Description: "PAN EXTENDED",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		35: field.NewString(&field.Spec{
			Length:      37,
			Description: "TRACK 2 DATA",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		36: field.NewString(&field.Spec{
			Length:      104,
			Description: "TRACK 3 DATA",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		37: field.NewString(&field.Spec{
			Length:      12,
			Description: "RETRIEVAL REFERENCE NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		38: field.NewString(&field.Spec{
			Length:      6,
			Description: "AUTHORIZATION IDENTIFICATION RESPONSE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		39: field.NewString(&field.Spec{
			Length:      2,
			Description: "RESPONSE CODE,",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		40: field.NewString(&field.Spec{
			Length:      3,
			Description: "SERVICE RESTRICTION CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		41: field.NewString(&field.Spec{
			Length:      8,
			Description: "CARD ACCEPTOR TERMINAL IDENTIFICACION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		42: field.NewString(&field.Spec{
			Length:      15,
			Description: "CARD ACCEPTOR IDENTIFICATION CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		43: field.NewString(&field.Spec{
			Length:      40,
			Description: "CARD ACCEPTOR NAME/LOCATION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		44: field.NewString(&field.Spec{
			Length:      25,
			Description: "ADITIONAL RESPONSE DATA",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		45: field.NewString(&field.Spec{
			Length:      76,
			Description: "TRACK 1 DATA",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		46: field.NewString(&field.Spec{
			Length:      999,
			Description: "ADITIONAL DATA - ISO",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		47: field.NewString(&field.Spec{
			Length:      999,
			Description: "ADITIONAL DATA - NATIONAL",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		48: field.NewString(&field.Spec{
			Length:      999,
			Description: "ADITIONAL DATA - PRIVATE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		49: field.NewString(&field.Spec{
			Length:      3,
			Description: "CURRENCY CODE,TRANSACTION,",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		50: field.NewString(&field.Spec{
			Length:      3,
			Description: "CURRENCY CODE,SETTLEMENT",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		51: field.NewString(&field.Spec{
			Length:      3,
			Description: "CURRENCY CODE,CARDHOLDER BILLING",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		52: field.NewString(&field.Spec{
			Length:      16,
			Description: "PIN DATA",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		53: field.NewString(&field.Spec{
			Length:      8,
			Description: "SECURITY RELATED CONTROL INFORMATION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		54: field.NewString(&field.Spec{
			Length:      120,
			Description: "ADDITIONAL AMOUNTS",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		55: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED ISO",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		56: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED ISO",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		57: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED NATIONAL",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		58: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED NATIONAL",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		59: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED NATIONAL,",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		60: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED PRIVATE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		61: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED PRIVATE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		62: field.NewString(&field.Spec{
			Length:      10,
			Description: "RESERVED PRIVATE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		63: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED PRIVATE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		64: field.NewBinary(&field.Spec{
			Length:      8,
			Description: "MESSAGE AUTHENTICATION CODE FIELD",
			Enc:         encoding.Binary,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		65: field.NewBinary(&field.Spec{
			Length:      1,
			Description: "BITMAP,EXTENDED",
			Enc:         encoding.Binary,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		66: field.NewString(&field.Spec{
			Length:      1,
			Description: "SETTLEMENT CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		67: field.NewString(&field.Spec{
			Length:      2,
			Description: "EXTENDED PAYMENT CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		68: field.NewString(&field.Spec{
			Length:      3,
			Description: "RECEIVING INSTITUTION COUNTRY CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		69: field.NewString(&field.Spec{
			Length:      3,
			Description: "SETTLEMENT INSTITUTION COUNTRY CODE,",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		70: field.NewString(&field.Spec{
			Length:      3,
			Description: "NETWORK MANAGEMENT INFORMATION CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		71: field.NewString(&field.Spec{
			Length:      4,
			Description: "MESSAGE NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		72: field.NewString(&field.Spec{
			Length:      4,
			Description: "MESSAGE NUMBER LAST",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		73: field.NewString(&field.Spec{
			Length:      6,
			Description: "DATE ACTION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		74: field.NewString(&field.Spec{
			Length:      10,
			Description: "CREDITS NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		75: field.NewString(&field.Spec{
			Length:      10,
			Description: "CREDITS REVERSAL NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		76: field.NewString(&field.Spec{
			Length:      10,
			Description: "DEBITS NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		77: field.NewString(&field.Spec{
			Length:      10,
			Description: "DEBITS REVERSAL NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		78: field.NewString(&field.Spec{
			Length:      10,
			Description: "TRANSFER NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		79: field.NewString(&field.Spec{
			Length:      10,
			Description: "TRANSFER REVERSAL NUMBER,",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		80: field.NewString(&field.Spec{
			Length:      10,
			Description: "INQUIRIES NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		81: field.NewString(&field.Spec{
			Length:      10,
			Description: "AUTHORIZATION NUMBER",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		82: field.NewString(&field.Spec{
			Length:      12,
			Description: "CREDITS,PROCESSING FEE AMOUNT",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		83: field.NewString(&field.Spec{
			Length:      12,
			Description: "CREDITS,TRANSACTION FEE AMOUNT",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		84: field.NewString(&field.Spec{
			Length:      12,
			Description: "DEBITS,PROCESSING FEE AMOUNT",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		85: field.NewString(&field.Spec{
			Length:      12,
			Description: "DEBITS,TRANSACTION FEE AMOUNT",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		86: field.NewString(&field.Spec{
			Length:      16,
			Description: "CREDITS,AMOUNT",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		87: field.NewString(&field.Spec{
			Length:      16,
			Description: "CREDITS,REVERSAL AMOUNT",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		88: field.NewString(&field.Spec{
			Length:      16,
			Description: "DEBITS,AMOUNT",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		89: field.NewString(&field.Spec{
			Length:      16,
			Description: "DEBITS,REVERSAL AMOUNT,",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		90: field.NewString(&field.Spec{
			Length:      42,
			Description: "ORIGINAL DATA ELEMENTS",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		91: field.NewString(&field.Spec{
			Length:      1,
			Description: "FILE UPDATE CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		92: field.NewString(&field.Spec{
			Length:      2,
			Description: "FILE SECURITY CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		93: field.NewString(&field.Spec{
			Length:      6,
			Description: "RESPONSE INDICATOR",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		94: field.NewString(&field.Spec{
			Length:      7,
			Description: "SERVICE INDICATOR",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		95: field.NewString(&field.Spec{
			Length:      42,
			Description: "REPLACEMENT AMOUNTS",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		96: field.NewBinary(&field.Spec{
			Length:      16,
			Description: "MESSAGE SECURITY CODE",
			Enc:         encoding.Binary,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		97: field.NewString(&field.Spec{
			Length:      17,
			Description: "AMOUNT,NET SETTLEMENT",
			Enc:         encoding.ASCII,
			Pad:         padding.Left('0'),
		}),
		98: field.NewString(&field.Spec{
			Length:      25,
			Description: "PAYEE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
		99: field.NewString(&field.Spec{
			Length:      11,
			Description: "SETTLEMENT INSTITUTION IDENT CODE,",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		100: field.NewString(&field.Spec{
			Length:      11,
			Description: "RECEIVING INSTITUTION IDENT CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		101: field.NewString(&field.Spec{
			Length:      17,
			Description: "FILE NAME",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		102: field.NewString(&field.Spec{
			Length:      28,
			Description: "ACCOUNT IDENTIFICATION 1",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		103: field.NewString(&field.Spec{
			Length:      28,
			Description: "ACCOUNT IDENTIFICATION 2",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		104: field.NewString(&field.Spec{
			Length:      210,
			Description: "TRANSACTION DESCRIPTION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		105: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED ISO USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		106: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED ISO USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		107: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED ISO USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		108: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED ISO USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		109: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED ISO USE,",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		110: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED ISO USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		111: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED ISO USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		112: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED NATIONAL USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		113: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED NATIONAL USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		114: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED NATIONAL USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		115: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED NATIONAL USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		116: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED NATIONAL USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		117: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED NATIONAL USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		118: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED NATIONAL USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		119: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED NATIONAL USE,",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		120: field.NewString(&field.Spec{
			Length:      70,
			Description: "RESERVED PRIVATE USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		121: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED PRIVATE USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		122: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED PRIVATE USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		123: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED PRIVATE USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		124: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED PRIVATE USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		125: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED PRIVATE USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		126: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED PRIVATE USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		127: field.NewString(&field.Spec{
			Length:      999,
			Description: "RESERVED PRIVATE USE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		128: field.NewString(&field.Spec{
			Length:      64,
			Description: "MAC 2 };",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Right(' '),
		}),
	},
}

type TimeUtils struct {
	now time.Time
}

func NewTimeUtils() TimeUtils {
	return TimeUtils{
		now: time.Now(),
	}
}

func (t *TimeUtils) GetCurrentTime() {
	t.now = time.Now()
}

func (t *TimeUtils) Format(now time.Time, format string) string {
	return now.Format(format)
}

func (t *TimeUtils) GetTimeYYYYMMDDHHMMSS() string {
	return t.now.Format("20060102150405") // YYYYMMDDHHMMSS
}

func (t *TimeUtils) GetTimeGmtMMDDHHMMSS() string {
	return t.now.In(time.FixedZone("GMT", 0)).Format("0102150405") // MMDDHHMMSS
}

func (t *TimeUtils) GetTimeHHMMSS() string {
	return t.now.Format("150405") // HHMMSS
}

func (t *TimeUtils) GetTimeMMDD() string {
	return t.now.Format("0102") // MMDD
}

func (t *TimeUtils) GetTimeYDDDHHNNNNNN(f11 string) string {
	lastLetterYear := t.now.Year() % 10
	dayOfYear := t.now.YearDay()
	fmt.Println(strconv.Itoa(lastLetterYear) + strconv.Itoa(dayOfYear) + f11)
	return strconv.Itoa(lastLetterYear) + strconv.Itoa(dayOfYear) + f11
}

func padLeft(str string, length int, padChar byte) string {
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
