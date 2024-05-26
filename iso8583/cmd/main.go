package main

import (
	"fmt"
	utils_class "iso8583"
	"time"

	"github.com/google/uuid"
	"github.com/moov-io/iso8583"
)

func main() {
	// testGenIsoMessage()
	testParseIsoMessage()
}

func testGenIsoMessage() {
	request := utils_class.Request[utils_class.GetNapasRecipientNameRequest]{
		Trace: utils_class.Trace{
			From:     uuid.NewString(),
			To:       uuid.NewString(),
			Cid:      uuid.NewString(),
			Sid:      uuid.NewString(),
			Cts:      time.Now().Unix(),
			Sts:      time.Now().Unix(),
			Username: uuid.NewString(),
		},
		Data: utils_class.GetNapasRecipientNameRequest{
			PaymentMethod:     "ACCOUNT",
			Channel:           uuid.NewString(),
			ClientTransId:     uuid.NewString(),
			FromAccountNumber: "0037100018823009",
			ToBankCode:        "1234",
			ToAccountNumber:   "0037100018823009",
			ToCardNumber:      uuid.NewString(),
			ToCreditCard: utils_class.GetNapasRecipientNameRequestCreditCartInfo{
				EncryptedCardNo: uuid.NewString(),
				EncryptedKey:    uuid.NewString(),
				CardMasking:     uuid.NewString(),
			},
		},
	}
	currentAccountDto := utils_class.CurrentAccountGetCurrentAccountResponse{
		CIFNum:                     uuid.NewString(),
		AccountNum:                 "0037100018823009",
		AccountName:                uuid.NewString(),
		AccountCurrency:            uuid.NewString(),
		AccountPostingRestrictions: uuid.NewString(),
		MappingBankAccount:         uuid.NewString(),
		AccountCategoryID:          uuid.NewString(),
		VirtualAccountNum:          uuid.NewString(),
		VaAccountName:              uuid.NewString(),
		AccountMaturityDate:        uuid.NewString(),
	}
	output := utils_class.GenQueryMsg(request, currentAccountDto)
	fmt.Printf("%+v", output)
}

func testParseIsoMessage() {
	var reqExecMsg string = "03930200F23A44810CE1800600000000170001E11600371000241520084320200000040000000516092545294410162545051605167399000050697044841370929441010369800000001000000000000001OCB EBANKING           TP HCM        VNM015PHAM NGOC LALA\r70406IF_DEP0184129OCBB880430610606970406160037100024152008100129837294017chuyen tien napas015NGUYEN VAN TEST0000000005276737BAE47D278B92910074D58A203D97C131B6406A3D32A4E715FD46480CD"
	var reqQueryMsg string = "03930200F23A44810CE1801400000000170000011600371000188230094320200000000000000516101818860314171818051605167399000050697044841371086031410369800000001000000000000001OCB EBANKING           TP HCM        VNM021SHORTNAME-977239test\r7040020406IF_INQ06970406160037100018823009100129837294018Truy Van Thong Tin9498512CA7BE618F59E89B0C25E5082D9AF0A473BB266A8782A2647843B72930"

	data := [][]interface{}{
		{"F50", "1", 3, "", "PAYMENT", "NHPL", "EQUAL"},
		{"F43", "1", 40, "", "ALL", "NHTH", "EQUAL"},
		{"F48", "0", 999, "", "QUERY", "NHPL", "MAX"},
		{"F12", "1", 6, "", "ALL", "ALL", "EQUAL"},
		{"F60", "1", 60, "", "ALL", "ALL", "MAX"},
		{"F22", "1", 3, "", "ALL", "NHTH", "EQUAL"},
		{"F11", "1", 6, "", "ALL", "ALL", "EQUAL"},
		{"F05", "1", 12, "", "PAYMENT", "NHPL", "EQUAL"},
		{"F03", "1", 6, "", "ALL", "ALL", "EQUAL"},
		{"F62", "1", 10, "", "ALL", "ALL", "MAX"},
		{"F102", "1", 28, "", "ALL", "ALL", "MAX"},
		{"F41", "1", 8, "", "ALL", "ALL", "EQUAL"},
		{"F25", "1", 2, "", "ALL", "NHTH", "EQUAL"},
		{"F63", "1", 16, "", "ALL", "NHPL", "MAX"},
		{"F120", "0", 70, "", "QUERY", "NHPL", "MAX"},
		{"F07", "1", 10, "checkMMddHHmmss", "ALL", "ALL", "EQUAL"},
		{"F13", "1", 4, "checkMMDD", "ALL", "ALL", "EQUAL"},
		{"F09", "1", 8, "", "PAYMENT", "NHPL", "EQUAL"},
		{"F37", "1", 12, "", "ALL", "ALL", "EQUAL"},
		{"F104", "1", 210, "", "ALL", "ALL", "MAX"},
		{"F128", "1", 64, "", "ALL", "ALL", "EQUAL"},
		{"F481", "0", 300, "", "ALL", "ALL", "MAX"},
		{"F49", "1", 3, "", "ALL", "ALL", "EQUAL"},
		{"F04", "1", 12, "", "ALL", "ALL", "EQUAL"},
		{"F39", "1", 2, "", "ALL", "NHPL", "EQUAL"},
		{"F32", "1", 11, "", "ALL", "ALL", "MAX"},
		{"F15", "1", 4, "", "ALL", "NHPL", "EQUAL"},
		{"F02", "1", 19, "", "ALL", "ALL", "MAX"},
		{"F482", "0", 300, "", "ALL", "ALL", "MAX"},
		{"F103", "1", 28, "", "ALL", "ALL", "MAX"},
	}

	var isoMessage *iso8583.Message
	isoMessage = utils_class.ParseMessage(reqExecMsg[4:])
	isoMessage = utils_class.ParseMessage(reqQueryMsg[4:])

	var ls []utils_class.NapasRuleEntity

	for _, item := range data {
		ls = append(ls, utils_class.NapasRuleEntity{
			FieldName:  item[0].(string),
			CheckNull:  item[1].(string),
			Length:     int64(item[2].(int)),
			Pattern:    item[3].(string),
			TransType:  item[4].(string),
			Route:      item[5].(string),
			LengthType: item[6].(string),
		})
	}

	fmt.Println(utils_class.CheckHmac(isoMessage))
	fmt.Println(utils_class.CheckRule(isoMessage, ls))

}
