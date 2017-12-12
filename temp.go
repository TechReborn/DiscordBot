package main

import (
	"strconv"
	"strings"

	"fmt"
	"math"
)

const (
	absoluteZeroC = float64(273.15)
	absoluteZeroF = float64(459.67)
)

func CtoF(c float64) float64 {
	value := round(float64((c * 9 / 5) + 32))
	return value
}

func CtoK(c float64) float64 {
	value := round(float64(c + absoluteZeroC))
	return value
}

func FtoC(f float64) float64 {
	value := round(float64((f - 32) * 5 / 9))
	return value
}

func FtoK(f float64) float64 {
	value := round(float64((f + absoluteZeroF) * 5 / 9))
	return value
}

func KtoC(k float64) float64 {
	value := round(float64(k - absoluteZeroC))
	return value
}

func KtoF(k float64) float64 {
	value := round(float64((k * 9 / 5) - absoluteZeroF))
	return value
}

func round(val float64) float64 {
	var round float64
	digit := 10 * val
	if _, div := math.Modf(digit); div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / 10
}

func handleTempMessage(messageContent string) (string, bool) {
	split := strings.Split(messageContent, " ")

	result := ""

	for i := 0; i < len(split); i++ {
		str := split[i]
		if strings.HasSuffix(strings.ToLower(str), "c"){
			if len(str) != 1 {
				int, err := strconv.Atoi(TrimSuffix(strings.ToLower(str), "c"))
				if err == nil {
					result = result + handleC(int) + "\n"
				}
			} else if i > 1{
				int, err := strconv.Atoi(split[i -1])
				if err == nil {
					result = result + handleC(int) + "\n"
				}
			}
		} else if strings.HasSuffix(strings.ToLower(str), "f"){
			if len(str) != 1 {
				int, err := strconv.Atoi(TrimSuffix(strings.ToLower(str), "f"))
				if err == nil {
					result = result + handleF(int) + "\n"
				}
			} else if i > 1{
				int, err := strconv.Atoi(split[i -1])
				if err == nil {
					result = result + handleF(int) + "\n"
				}
			}
		} else if strings.HasSuffix(strings.ToLower(str), "k"){
			if len(str) != 1 {
				int, err := strconv.Atoi(TrimSuffix(strings.ToLower(str), "k"))
				if err == nil {
					result = result + handleK(int) + "\n"
				}
			} else if i > 1{
				int, err := strconv.Atoi(split[i -1])
				if err == nil {
					result = result + handleK(int) + "\n"
				}
			}
		}
	}

	if len(result) > 0 {
		return result, true
	}

	if strings.HasPrefix(messageContent, "!temp ") {
		unit := strings.ToLower(messageContent[6:][:1])
		int, err := strconv.Atoi(messageContent[8:])
		if err != nil {
			return "Failed to parse temperature", true
		}

		if unit == "c" {
			return handleC(int), true
		} else if unit == "f" {
			return handleF(int), true
		} else if unit == "k" {
			return handleK(int), true
		}

		return handleC(int), true
	}

	return "", false
}

func handleC(input int) string {
	return fmt.Sprintf("%dC is equal to: %sF and %sK", input, FloatToString(CtoF(float64(input))), FloatToString(CtoK(float64(input))))
}

func handleF(input int) string {
	return fmt.Sprintf("%dF is equal to: %sC and %sK", input, FloatToString(FtoC(float64(input))), FloatToString(FtoK(float64(input))))
}

func handleK(input int) string {
	return fmt.Sprintf("%dK is equal to: %sC and %sF", input, FloatToString(KtoC(float64(input))), FloatToString(KtoF(float64(input))))
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 1, 64)
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func main(){
	fmt.Println(handleTempMessage("The weather today is 20f and it was 6C yesterday. where as 6000 k is the heat of the sun"))
}
