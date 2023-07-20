package main

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ValCurs struct {
	Date   string `xml:"Date,attr"`
	Valute []struct {
		CharCode string `xml:"CharCode"`
		Name     string `xml:"Name"`
		Value    string `xml:"Value"`
	} `xml:"Valute"`
}

func main() {
	date := "19/07/2023"
	currencyCode := "USD"
	
	getExchangeRate(date, currencyCode)
}

func getExchangeRate(date string, currencyCode string) {
	url := fmt.Sprintf("https://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date)

	body, err := getResponseBody(url)
	if err != nil {
		fmt.Println("Ошибка при получении ответа:", err)
		return
	}

	body, err = convertToUTF8(body)
	if err != nil {
		fmt.Println("Ошибка при преобразовании кодировки:", err)
		return
	}

	valCurs, err := parseXMLData(body)
	if err != nil {
		fmt.Println("Ошибка при разборе XML:", err)
		return
	}

	printCurrencyRates(valCurs, currencyCode)
}

func getResponseBody(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func convertToUTF8(input []byte) ([]byte, error) {
	reader, _ := charset.NewReaderLabel("windows-1251", strings.NewReader(string(input)))
	return ioutil.ReadAll(reader)
}

func parseXMLData(data []byte) (*ValCurs, error) {
	var valCurs ValCurs
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(&valCurs)
	if err != nil {
		return nil, err
	}
	return &valCurs, nil
}

func printCurrencyRates(valCurs *ValCurs, currencyCode string) {
	fmt.Println("Курс валют на", valCurs.Date)
	for _, valute := range valCurs.Valute {
		if valute.CharCode == currencyCode {
			value, err := strconv.ParseFloat(strings.ReplaceAll(valute.Value, ",", "."), 64)
			if err != nil {
				fmt.Println("Ошибка при преобразовании курса валюты:", err)
				return
			}

			encoder := charmap.Windows1251.NewEncoder()
			nameEncoded, err := encoder.String(valute.Name)
			if err != nil {
				fmt.Println("Ошибка при преобразовании имени валюты:", err)
				return
			}

			fmt.Printf("%s (%s): %.4f\n", valute.CharCode, nameEncoded, value)
			return
		}
	}

	fmt.Printf("Валюта с кодом %s не найдена\n", currencyCode)
}
