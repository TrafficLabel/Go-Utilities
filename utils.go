// Copyright 2018 Traffic Label. All rights reserved.
// Source Code Written for public use.
// @author Noy Hillel

package utils

import (
	"bytes"
	"encoding/json"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ConvertToFloat(s string) float64 {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("Could not convert %v to a float!", s)
		return 0
	}
	return num
}

func PrintEmoji(emoji string, emojiMap map[string]string) string {
	for key, value := range emojiMap {
		if ":"+emoji+":" == key {
			return value
		}
	}
	return ""
}

func Commaf(v float64) string {
	return humanize.Commaf(v)
}

func Comma(v int64) string {
	return humanize.Comma(v)
}

func BubbleSortDesc(arr []string) []string {
	temp := ""
	for i := 0; i < len(arr); i++ {
		for j := 1; j < len(arr)-i; j++ {
			if arr[j-1] < arr[j] {
				temp = arr[j-1]
				arr[j-1] = arr[j]
				arr[j] = temp
			}
		}
	}
	return arr
}

func String(n int) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func ProperlyFormatDate(date string, headerMap map[string]interface{}) (string, error) {
	properDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		headerMap["Err"] = "There was an error loading the date!"
		return "", err
	}
	return properDate.Format("January 2 2006"), err
}

func FormatFloat(num float64) string {
	return strconv.FormatFloat(num, 'f', 2, 64)
}

func SendHTTPError(reason string, rw http.ResponseWriter) {
	http.Error(rw, reason, http.StatusBadRequest)
}

func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	var result []string
	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func TrimCompletelyAfter(s, remover string) string {
	if idx := strings.Index(s, remover); idx != -1 {
		return s[:idx]
	}
	return s
}

func Interface(f interface{}) string {
	return f.(string)
}

func MonthInSlice(a interface{}, list []time.Month) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func DaysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func GetRealAddr(r *http.Request) string {
	remoteIP := ""
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}
	return remoteIP
}

func DenyAccess(w http.ResponseWriter, data string) {
	color.HiRed("[%v] Access Denied From: %v", time.Now().Format(time.RFC850), color.HiBlueString(data))
	http.Error(w, "You can't access this from "+data, 401)
}

func RedirectToHome(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/", http.StatusFound)
}

func OpenFile(path string) (file io.ReadWriteCloser, err error) {
	file, err = os.Open(path)
	return
}

func Mode(mode []string) string {
	if len(mode) == 0 {
		return ""
	}
	var modeMap = map[string]int{}
	var maxEl = mode[0]
	var maxCount = 1
	for i := 0; i < len(mode); i++ {
		var el = mode[i]
		if modeMap[el] == 0 {
			modeMap[el] = 1
		} else {
			modeMap[el]++
		}
		if modeMap[el] > maxCount {
			maxEl = el
			maxCount = modeMap[el]
		}
	}
	return maxEl
}

func GetExchangeRates(currency string, fallback float64) float64 {
	type Rates struct {
		GBP float64 `json:"GBP"`
	}
	type CurrencyResult struct {
		Rates Rates `json:"rates"`
	}
	// FREE API!
	re, err := http.Get("https://api.exchangeratesapi.io/latest?base=" + currency)
	if err != nil {
		log.Printf("Something went wrong getting the API: %v", err.Error())
		return fallback
	}
	defer re.Body.Close()
	body, err := ioutil.ReadAll(re.Body)
	if err != nil {
		log.Printf("Error with reading body.. falling back on %v. Error: %v", fallback, err.Error())
		return fallback
	}
	var cR CurrencyResult
	if err = json.Unmarshal(body, &cR); err != nil {
		if err != nil {
			log.Printf("Error with unmarshal, falling back on %v. Error: %v", fallback, err.Error())
			return fallback
		}
	}
	return cR.Rates.GBP
}

func JsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func GetDaysInMonth(month string, year int) int {
	switch month {
	case time.January.String():
		return 31
	case time.February.String():
		daysInYear := time.Date(year, time.December, 31, 0, 0, 0, 0, time.Local).YearDay()
		if daysInYear > 365 {
			return 29
		} else {
			return 28
		}
	case time.March.String():
		return 31
	case time.April.String():
		return 30
	case time.May.String():
		return 31
	case time.June.String():
		return 30
	case time.July.String():
		return 31
	case time.August.String():
		return 31
	case time.September.String():
		return 30
	case time.October.String():
		return 31
	case time.November.String():
		return 30
	case time.December.String():
		return 31
	default:
		break
	}
	return 31
}

func GetMonthFromName(month string) (t time.Month, err error) {
	parsedMonth, err := time.Parse("January", month)
	if err != nil {
		log.Println(err.Error())
		return t, err
	} else {
		return parsedMonth.Month(), nil
	}
}

func IntArrayToString(A []int, delim string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(A); i++ {
		buffer.WriteString(strconv.Itoa(A[i]))
		if i != len(A)-1 {
			buffer.WriteString(delim)
		}
	}
	return buffer.String()
}

func ReverseList(reversed []interface{}) []interface{} {
	sort.SliceStable(reversed, func(i, j int) bool {
		return true
	})
	return reversed
}

func ParseCSVFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("err opening file, %v", err.Error())
	}
	return file
}

func FormatDateWithSuffix(t time.Time) string {
	suffix := "th"
	switch t.Day() {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return t.Format("January 2" + suffix + ", 2006")
}

func DownloadAndSaveFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
