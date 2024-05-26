package utils_class

var GetNapasRecipientNameTag = "GetNapasRecipientName"

// REQUEST
type GetNapasRecipientNameRequest struct {
	PaymentMethod     string                                     `json:"paymentMethod" validate:"required"`
	Channel           string                                     `json:"channel" validate:"required"`
	ClientTransId     string                                     `json:"clientTransId" validate:"required"`
	FromAccountNumber string                                     `json:"fromAccountNumber",omitempty`
	ToBankCode        string                                     `json:"toBankCode,omitempty",omitempty`
	ToAccountNumber   string                                     `json:"toAccountNumber,omitempty"`
	ToCardNumber      string                                     `json:"toCardNumber,omitempty"`
	ToCreditCard      GetNapasRecipientNameRequestCreditCartInfo `json:"toCreditCard,omitempty"`
}

type GetNapasRecipientNameRequestCreditCartInfo struct {
	EncryptedCardNo string `json:"encryptedCardNo",omitempty`
	EncryptedKey    string `json:"encryptedKey",omitempty`
	CardMasking     string `json:"cardMasking",omitempty`
}

// --------------------------------------------------------------------------------------------------------------------------------
// RESPONSE
type GetNapasRecipientNameResponse struct {
	BenefitCustomerName     string                             `json:"benefitCustomerName"`
	NapasRefNumber          string                             `json:"napasRefNumber"`
	MerchantCategoryCode    string                             `json:"merchantCategoryCode"`
	AcceptorNameAndLocation string                             `json:"acceptorNameAndLocation"`
	NapasResponse           GetNapasRecipientNameNapasResponse `json:"NapasResponse"`
}
type GetNapasRecipientNameNapasResponse struct {
	DetailCode string `json:"detailCode"`
}

type CurrentAccountGetCurrentAccountRequest struct {
	AccountNumber string `json:"accountNumber" validate:"required"`
}

type CurrentAccountGetCurrentAccountResponse struct {
	CIFNum                     string `json:"CIFNum" db:"CIFNum"`
	AccountNum                 string `json:"accountNum" db:"accountNum"`
	AccountName                string `json:"accountName" db:"accountName"`
	AccountCurrency            string `json:"accountCurrency" db:"accountCurrency"`
	AccountPostingRestrictions string `json:"accountPostingRestrictions" db:"accountPostingRestrictions"`
	MappingBankAccount         string `json:"mappingBankAccount" db:"mappingBankAccount"`
	AccountCategoryID          string `json:"accountCategoryID" db:"accountCategoryID"`
	VirtualAccountNum          string `json:"virtualAccountNum" db:"virtualAccountNum"`
	VaAccountName              string `json:"vaAccountName" db:"vaAccountName"`
	AccountMaturityDate        string `json:"accountMaturityDate" db:"accountMaturityDate"`
}

type IsoMessage struct {
	IsoMessage string
	Field0     string `json:"field0"`
	Field1     string `json:"field1"`
	Field2     string `json:"field2"`
	Field3     string `json:"field3"`
	Field4     string `json:"field4"`
	Field5     string `json:"field5"`
	Field6     string `json:"field6"`
	Field7     string `json:"field7"`
	Field8     string `json:"field8"`
	Field9     string `json:"field9"`
	Field10    string `json:"field10"`
	Field11    string `json:"field11"`
	Field12    string `json:"field12"`
	Field13    string `json:"field13"`
	Field14    string `json:"field14"`
	Field15    string `json:"field15"`
	Field16    string `json:"field16"`
	Field17    string `json:"field17"`
	Field18    string `json:"field18"`
	Field19    string `json:"field19"`
	Field20    string `json:"field20"`
	Field21    string `json:"field21"`
	Field22    string `json:"field22"`
	Field23    string `json:"field23"`
	Field24    string `json:"field24"`
	Field25    string `json:"field25"`
	Field26    string `json:"field26"`
	Field27    string `json:"field27"`
	Field28    string `json:"field28"`
	Field29    string `json:"field29"`
	Field30    string `json:"field30"`
	Field31    string `json:"field31"`
	Field32    string `json:"field32"`
	Field33    string `json:"field33"`
	Field34    string `json:"field34"`
	Field35    string `json:"field35"`
	Field36    string `json:"field36"`
	Field37    string `json:"field37"`
	Field38    string `json:"field38"`
	Field39    string `json:"field39"`
	Field40    string `json:"field40"`
	Field41    string `json:"field41"`
	Field42    string `json:"field42"`
	Field43    string `json:"field43"`
	Field44    string `json:"field44"`
	Field45    string `json:"field45"`
	Field46    string `json:"field46"`
	Field47    string `json:"field47"`
	Field48    string `json:"field48"`
	Field49    string `json:"field49"`
	Field50    string `json:"field50"`
	Field51    string `json:"field51"`
	Field52    string `json:"field52"`
	Field53    string `json:"field53"`
	Field54    string `json:"field54"`
	Field55    string `json:"field55"`
	Field56    string `json:"field56"`
	Field57    string `json:"field57"`
	Field58    string `json:"field58"`
	Field59    string `json:"field59"`
	Field60    string `json:"field60"`
	Field61    string `json:"field61"`
	Field62    string `json:"field62"`
	Field63    string `json:"field63"`
	Field64    string `json:"field64"`
	Field65    string `json:"field65"`
	Field66    string `json:"field66"`
	Field67    string `json:"field67"`
	Field68    string `json:"field68"`
	Field69    string `json:"field69"`
	Field70    string `json:"field70"`
	Field71    string `json:"field71"`
	Field72    string `json:"field72"`
	Field73    string `json:"field73"`
	Field74    string `json:"field74"`
	Field75    string `json:"field75"`
	Field76    string `json:"field76"`
	Field77    string `json:"field77"`
	Field78    string `json:"field78"`
	Field79    string `json:"field79"`
	Field80    string `json:"field80"`
	Field81    string `json:"field81"`
	Field82    string `json:"field82"`
	Field83    string `json:"field83"`
	Field84    string `json:"field84"`
	Field85    string `json:"field85"`
	Field86    string `json:"field86"`
	Field87    string `json:"field87"`
	Field88    string `json:"field88"`
	Field89    string `json:"field89"`
	Field90    string `json:"field90"`
	Field91    string `json:"field91"`
	Field92    string `json:"field92"`
	Field93    string `json:"field93"`
	Field94    string `json:"field94"`
	Field95    string `json:"field95"`
	Field96    string `json:"field96"`
	Field97    string `json:"field97"`
	Field98    string `json:"field98"`
	Field99    string `json:"field99"`
	Field100   string `json:"field100"`
	Field101   string `json:"field101"`
	Field102   string `json:"field102"`
	Field103   string `json:"field103"`
	Field104   string `json:"field104"`
	Field105   string `json:"field105"`
	Field106   string `json:"field106"`
	Field107   string `json:"field107"`
	Field108   string `json:"field108"`
	Field109   string `json:"field109"`
	Field110   string `json:"field110"`
	Field111   string `json:"field111"`
	Field112   string `json:"field112"`
	Field113   string `json:"field113"`
	Field114   string `json:"field114"`
	Field115   string `json:"field115"`
	Field116   string `json:"field116"`
	Field117   string `json:"field117"`
	Field118   string `json:"field118"`
	Field119   string `json:"field119"`
	Field120   string `json:"field120"`
	Field121   string `json:"field121"`
	Field122   string `json:"field122"`
	Field123   string `json:"field123"`
	Field124   string `json:"field124"`
	Field125   string `json:"field125"`
	Field126   string `json:"field126"`
	Field127   string `json:"field127"`
	Field128   string `json:"field128"`
}

type ListNapasRuleEntity struct {
	NapasRuleEntityLst []NapasRuleEntity `json:"napasRuleEntityList"`
}

type NapasRuleEntity struct {
	FieldName  string `json:"fieldName" db:"FIELDNAME"`
	CheckNull  string `json:"checkNull" db:"CHECKNULL"`
	Length     int64  `json:"length" db:"LENGTH"`
	Pattern    string `json:"pattern" db:"PATTERN"`
	TransType  string `json:"transtype" db:"TRANSTYPE"`
	Route      string `json:"route" db:"ROUTE"`
	LengthType string `json:"lengthtype" db:"LENGTHTYPE"`
}
