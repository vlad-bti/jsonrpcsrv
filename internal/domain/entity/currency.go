package entity

var currency = map[string]string{
	"AED": "United Arab Emirates dirham",
	"ALL": "Albanian Lek",
	"AMD": "Armenian dram",
	"AOA": "Angolan Kwanza",
	"ARS": "Argentine peso",
	"AUD": "Australian dollar",
	"AZN": "Azerbaijani manat",
	"BAM": "Bosnia and Herzegovina convertible mark",
	"BDT": "Bangladeshi Taka",
	"BGN": "Bulgarian lev",
	"BHD": "Bahraini Dinar",
	"BLR": "Belarusian ruble",
	"BMD": "Bermudian Dollar",
	"BND": "Bruneian Dollar",
	"BOB": "Bolivian Boliviano",
	"BRL": "Brazilian real",
	"BTC": "Bitcoin",
	"BYN": "Belarusian ruble",
	"CAD": "Canadian dollar",
	"CDF": "Congolese franc",
	"CHF": "Swiss franc",
	"CLP": "Chilean peso",
	"CNY": "Chinese yuan",
	"COP": "Colombian peso",
	"CRC": "Costa Rican colon",
	"CZK": "Czech koruna",
	"DAS": "Dash Crypto",
	"DKK": "Danish krone",
	"DOG": "Dogecoin",
	"DOP": "Dominican peso",
	"DTC": "Datacoin",
	"DZD": "Algerian Dinari",
	"EGP": "Egyptian Pound",
	"EOS": "Cryptocurrency",
	"ETB": "Ethiopian birr",
	"ETH": "Ethereum",
	"EUR": "Euro",
	"GBP": "Pound sterling",
	"GEL": "Georgian lari",
	"GHS": "Ghanaian Cedi",
	"GMC": "Gridmaster",
	"GMD": "Gambian Dalasi",
	"GNF": "Guinean franc",
	"HKD": "Hong Kong dollar",
	"HRK": "Croatian Kuna",
	"HTG": "Haitian gourde",
	"HUF": "Hungarian Forint",
	"IDR": "Indonesian rupiah",
	"ILS": "Israeli new shekel",
	"INR": "Indian rupee",
	"IQD": "Iraqi Dinar",
	"IRR": "Iranian Rial",
	"NAN": "Not Available Now",
	"NGN": "Nigerian Naira",
	"NIO": "Nicaraguan córdoba",
	"NOK": "Norwegian krone",
	"NZD": "New Zealand dollar",
	"OMR": "Omani Rial",
	"PEN": "Peruvian Sol",
	"PHP": "Philippine peso",
	"PKR": "Pakistani rupee",
	"PLN": "Polish złoty",
	"PPT": "Populous",
	"PYG": "Paraguayan guaraní",
	"QAR": "Qatari Riyal",
	"RON": "Romanian Leu",
	"RSD": "Serbian dinar",
	"RUB": "Russian ruble",
	"SAR": "Saudi Arabian Riyal",
	"SCR": "Seychelles Rupee",
	"SDG": "Sudanese Pound",
	"SEK": "Swedish krona/kronor",
	"SGD": "Singapore Dollar",
	"SLL": "Sierra Leonean leone",
	"SRD": "Surinamese dollar",
	"SZL": "Swazi lilangeni",
	"THB": "Thai baht",
	"TJS": "Tajikistani Somoni",
	"TMT": "Turkmenistan manat",
	"TND": "Tunisian dinar",
	"TRY": "Turkish lira",
	"TWD": "Taiwan New Dollar",
	"TZS": "Tanzanian Shilling",
	"UAH": "Ukrainian hryvnia",
	"UBC": "μBTC",
	"UGX": "Ugandan shilling",
	"USD": "United States dollar",
	"UYU": "Uruguayan peso",
	"UZS": "Uzbekistani Som",
	"VEF": "Venezuelan bolívar",
	"VES": "Venezuelan bolívar",
	"VND": "Vietnamese Dong",
	"XAF": "Cameroon Franc",
	"XAU": "Gold Ounce",
	"XEM": "Cryptocurrency",
	"XMC": "mBCH mili Bitcoin Cash",
	"XME": "mETH mili Ethereum",
	"XML": "mLTC mili Litecoin",
	"XMR": "Monero",
	"XOF": "CFA franc",
	"XUT": "USDT USD Tether",
	"IRT": "Iranian touman",
	"ISK": "Icelandic Krona",
	"JOD": "Jordanian Dinar",
	"JPY": "Japanese yen",
	"KES": "Kenyan shilling",
	"KGS": "Kyrgyzstani Som",
	"KHR": "Cambodian Riel",
	"KRW": "South Korean Won",
	"KWD": "Kuwaiti Dinar",
	"KZT": "Kazakhstani tenge",
	"LBP": "Lebanese pound",
	"LRD": "Liberian dollar",
	"LTC": "Litecoin",
	"MAD": "Moroccan dirham",
	"MBC": "MicroBitcoin",
	"MBH": "Unknown",
	"MDL": "Moldovan leu",
	"MET": "Metronome",
	"MGA": "Malagasy Ariary",
	"MHC": "MetaHash",
	"MKD": "Macedonian denar",
	"MLT": "Maltese Lira",
	"MMK": "Burmese Kyat",
	"MNT": "Mongolian Tughrik",
	"MXM": "Maximine Coin",
	"MXN": "Mexican peso",
	"MYR": "Malaysian ringgit",
	"MZN": "Mozambican Metical",
	"NAD": "Namibian dollar",
	"ZAR": "South African rand",
	"ZEC": "ZCash",
	"ZMW": "Zambian kwacha",
	"ZWL": "Zimbabwean dollar",
}

func IsValidCurrency(name string) bool {
	_, ok := currency[name]
	return ok
}