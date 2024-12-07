package ecode

import "fmt"

type Errno struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *Errno) Error() string {
	return fmt.Sprintf("Ext - code: %d, msg: %s", e.Code, e.Msg)
}

var (
	Success = &Errno{Code: 0, Msg: "Success"}
	Error   = &Errno{Code: -1, Msg: "Fail"}

	ErrLogin                = &Errno{Code: 100001, Msg: "LOGIN ERROR"}
	ErrLoginIP              = &Errno{Code: 100002, Msg: "LOGIN ERROR,IP FROM ERROR"}
	ErrState                = &Errno{Code: 100003, Msg: "LOGIN STATE ERROR"}
	IpError                 = &Errno{Code: 100004, Msg: "CLIENT IP ERROR"}
	ErrPermissionNotAllow   = &Errno{Code: 100008, Msg: "DO NOT HAVE PERMISSION"}
	ErrSign                 = &Errno{Code: 100009, Msg: "SIGN ERROR"}
	ErrIllegal              = &Errno{Code: 100010, Msg: "ILLEGAL ERROR"}
	ErrBankUnderMaintenance = &Errno{Code: 100011, Msg: "BANK UNDER MAINTENANCE ERROR"}

	ErrReqParam   = &Errno{Code: 200001, Msg: "PARAMS ERROR"}
	ErrFrequently = &Errno{Code: 200002, Msg: "DO NOT SUBMIT FREQUENTLY ERROR"}

	ErrMerch               = &Errno{Code: 300001, Msg: "Merch ID ERROR"}
	ErrAccountInsufficient = &Errno{Code: 300003, Msg: "INSUFFICIENT ACCOUNT BALANCE"}
	ErrUpdateAccount       = &Errno{Code: 300002, Msg: "ACCOUNT BALANCE UPDATE ERROR "}

	ErrMerchLock     = &Errno{Code: 500001, Msg: "Merch BE LOCK"}
	ErrAccountLock   = &Errno{Code: 500004, Msg: "ACCOUNT BE LOCK"}
	ErrAccount       = &Errno{Code: 500005, Msg: "ACCOUNT CANNOT BE FIND"}
	ChannelErr       = &Errno{Code: 500003, Msg: "CHANNEL ID ERROR"}
	DuplicationError = &Errno{Code: 500004, Msg: "DUPLICATION ERROR"}

	CurrencyErr        = &Errno{Code: 600005, Msg: "CURRENCY ID ERROR"}
	CurrencyStateErr   = &Errno{Code: 600006, Msg: "CURRENCY BE LOCK"}
	RateErr            = &Errno{Code: 600007, Msg: "CURRENCY RATE ERROR"}
	BrokerageErr       = &Errno{Code: 600008, Msg: "CURRENCY BROKERAGE ERROR"}
	CurrencyUseErr     = &Errno{Code: 600009, Msg: "CURRENCY CANNOT USE"}
	ErrCurrencyDup     = &Errno{Code: 600020, Msg: "UP CHANNEL DUPLICATION"}
	ErrUpCurrency      = &Errno{Code: 600021, Msg: "UP CHANNEL CANNOT BE FIND"}
	ErrUpCurrencyState = &Errno{Code: 600022, Msg: "UP CHANNEL STATUS ERROR"}
	ErrBasicReqFail    = &Errno{Code: 600130, Msg: "BASIC REQ FAIL ERROR"}
	ErrDispatchCfgFail = &Errno{Code: 600131, Msg: "BASIC DISPATCH CONFIG SET ERROR"}

	OrderLowErr         = &Errno{Code: 6000030, Msg: "THE ORDER AMOUNT IS TOO LOW"}
	OrderHighErr        = &Errno{Code: 6000031, Msg: "THE ORDER AMOUNT IS TOO HIGH"}
	AmountLimitErr      = &Errno{Code: 6000031, Msg: "THE ORDER AMOUNT BE LIMIT"}
	OrderRepeatedErr    = &Errno{Code: 6000032, Msg: "THE ORDER_NO CANNOT BE REPEATED"}
	OrderErr            = &Errno{Code: 6000033, Msg: "THE ORDER CANNOT BE FIND"}
	OrderExtraErr       = &Errno{Code: 6000035, Msg: "THE ORDER EXTRA CANNOT BE FIND"}
	ErrOrderSync        = &Errno{Code: 6000036, Msg: "THE ORDER SYNC ERROR"}
	ErrOrderStateErr    = &Errno{Code: 6000037, Msg: "THE ORDER STATE ERROR"}
	ErrSrvWasBusyErr    = &Errno{Code: 6000038, Msg: "THE SERVICE WAS BUSY PLEASE WAIT"}
	ErrOrderCheckErr    = &Errno{Code: 6000039, Msg: "THE ORDER CHECK NOT FIND"}
	ErrOrderAmountError = &Errno{Code: 6000050, Msg: "AMOUNT ERROR"}

	ErrOrderClient  = &Errno{Code: 9000001, Msg: "THE ORDER CLIENT API CANNOT BE FIND"}
	ErrCurrencyCode = &Errno{Code: 9000002, Msg: "THE CURRENCY CODE CANNOT BE FIND"}

	ErrPixFormat            = &Errno{Code: 1200000, Msg: "Erro no formato da conta Pix!"}
	ErrCpfFormat            = &Errno{Code: 1200001, Msg: "Erro no formato da conta CPF!"}
	ErrEmailFormat          = &Errno{Code: 1200002, Msg: "Erro no formato da conta EMAIL!"}
	ErrPhoneFormat          = &Errno{Code: 1200003, Msg: "Erro no formato da conta PHONE!"}
	ErrEvpFormat            = &Errno{Code: 1200004, Msg: "Erro no formato da conta EVP!"}
	ErrCnpjCannotUse        = &Errno{Code: 1200005, Msg: "Erro Mantenimento do CNPJ!"}
	ErrAccountInvalid       = &Errno{Code: 1200006, Msg: "Erro A conta não existe ou é inválida!"}
	ErrActInvalid           = &Errno{Code: 1300001, Msg: "Err Account Invalid!"}
	ErrDepositImplNotFound  = &Errno{Code: 9900001, Msg: "DEPOSIT INTERFACE NOT IMPLEMENT!"}
	ErrDispatchImplNotFound = &Errno{Code: 9900002, Msg: "DISPATCH INTERFACE NOT IMPLEMENT!"}

	ErrGalaxyJobState = &Errno{Code: 12000001, Msg: "任务进行中，禁止修改。请先停止任务。"}
)
