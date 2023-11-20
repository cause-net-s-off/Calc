package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Calc struct {
	operand_1        int
	operand_1_Arabic bool
	operand_2        int
	operand_2_Arabic bool
	operation        byte
	Result           string
}

type Empty struct{}

type Roman struct {
	dec int
	rom string
}

var transTable = []Roman{
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{10, "X"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

var operations = map[byte]Empty{'+': {}, '-': {}, '*': {}, '/': {}}

var (
	ErrIncorrectNumber          = errors.New("Err: Incorrect number or out of range [I..X]!")
	ErrIncorrectOperation       = errors.New("Err: Incorrect operation!")
	ErrNotMathematicalOperation = errors.New("Err: Input is not a mathematical operation!")
	ErrIncorrectFormat          = errors.New("Err: Incorrect task format! (need: <operand_1> <+, -, *, /> <operand_2>)")
	ErrIncorrectRoman           = errors.New("Err: Roman number can not be negative or zero!")
	ErrMismatchOperands         = errors.New("Err: Both operands must be of the same number (Arabian or Roman)!")
	ErrOutOfNumberRange         = errors.New("Err: Number must be in range [1..10]!")
	ErrDivideByZero             = errors.New("Err: Divide by zero!")
)

func getOperation(operStr string, oper *byte) error { //			Определение арифм. операции
	if len(operStr) != 1 {
		return fmt.Errorf("%q %w", operStr, ErrIncorrectOperation)
	}
	operation := []byte(operStr)[0]
	_, val := operations[operation]
	if val {
		*oper = operation
		return nil
	}
	return fmt.Errorf("%q %w", operStr, ErrNotMathematicalOperation)
}

func decToRoman(num int) string { //					Перевод десятичного числа в римское
	res := ""
	for _, record := range transTable {
		for num >= record.dec {
			res += record.rom
			num -= record.dec
		}
	}
	return res
}

/*												// Перевод римского в десятичное (некорректно)
func romanToDec2(num string) (int, error) {
	num = "lxl"
	num = strings.ToUpper(num)
	fmt.Println("Num:", num)
	res := 0
	tableIndex := 0
	for i := 0; i < len(num); i++ {
		oneDigit := 0
		tryFind2D := ""
		find1D := ""
		// 0 1 2 3 4	5
		//   *
		find1D = string(num[i])
		if i+1 < len(num) {
			tryFind2D = string(num[i : i+2])
		} else {
			tryFind2D = string(num[i])
		}
		//fmt.Println(tryFind2D, i, find1D)

		for ind, record := range transTable {
			if tryFind2D == record.rom {
				if ind < tableIndex {
					return 0, ErrIncorrectRomanFormat
				}
				res += record.dec
				oneDigit = 0
				fmt.Println("ind", ind)
				tableIndex = ind + 4
				i++
				break
			}
			if find1D == record.rom {
				if ind < tableIndex {
					return 0, ErrIncorrectRomanFormat
				}
				oneDigit = record.dec
				tableIndex = ind
			}
		}
		fmt.Println(tableIndex)
		res += oneDigit
	}
	return res, nil
}
*/

func romanToDec(num string) (int, error) { //			Конвертация римского числа в десятичное
	res := 0
	err := errors.New("")
	num = strings.ToUpper(num)
	switch num {
	case "I":
		res = 1
	case "II":
		res = 2
	case "III":
		res = 3
	case "IV":
		res = 4
	case "V":
		res = 5
	case "VI":
		res = 6
	case "VII":
		res = 7
	case "VIII":
		res = 8
	case "IX":
		res = 9
	case "X":
		res = 10
	default:
		err = ErrIncorrectNumber
	}
	if res == 0 {
		return 0, err
	}
	return res, nil
}

func (c Calc) checkEqual() error { //						Проверка на идентичность исчисления операндов
	if c.operand_1_Arabic != c.operand_2_Arabic {
		return ErrMismatchOperands
	}
	return nil
}

func decodeOperand(operand string, op *int, arabic *bool) error { //	Извлечение операнда
	isArabic := true
	oper, err := strconv.Atoi(operand)
	if err != nil {
		isArabic = false
		oper, err = romanToDec(operand)
	}
	if err != nil {
		return err
	}
	if oper < 1 || oper > 10 {
		return ErrOutOfNumberRange
	}
	*op = oper
	*arabic = isArabic
	return nil
}

func analyzeFill(message string, calc *Calc) error { //		Парсинг введенной сроки и заполнение структуры
	message = RemoveSpaces(message)
	el := strings.Split(message, " ")
	if len(el) != 3 {
		return ErrIncorrectFormat
	}
	err := decodeOperand(el[0], &calc.operand_1, &calc.operand_1_Arabic)
	if err != nil {
		return err
	}
	err = decodeOperand(el[2], &calc.operand_2, &calc.operand_2_Arabic)
	if err != nil {
		return err
	}
	err = getOperation(el[1], &calc.operation)
	if err != nil {
		return err
	}
	err = calc.checkEqual()
	if err != nil {
		return err
	}

	return nil
}

func RemoveSpaces(str string) string { // 						Удаление лишних пробелов
	str = strings.TrimSpace(str)
	strArr := []byte(str)
	j := 0
	for i := 1; i < len(strArr); i++ {
		if strArr[j] != strArr[i] || strArr[j] != ' ' {
			j++
			strArr[j] = strArr[i]
		}
	}
	return string(strArr[:j+1])
}

func (c *Calc) getSolution(solution int) error { //		запись решения в Арабской/Римской, с проверкой
	if !c.operand_1_Arabic && !c.operand_2_Arabic {
		if solution < 1 {
			c.Result = ""
			return ErrIncorrectRoman
		}
		c.Result = decToRoman(solution)
	} else {
		c.Result = strconv.Itoa(solution)
	}
	return nil
}

func (c *Calc) add() error { //								Сложение (43)
	res := c.operand_1 + c.operand_2
	return c.getSolution(res)
}

func (c *Calc) substract() error { //						Вычитание (45)
	res := c.operand_1 - c.operand_2
	return c.getSolution(res)
}

func (c *Calc) multiply() error { //						Умножение (42)
	res := c.operand_1 * c.operand_2
	return c.getSolution(res)
}

func (c *Calc) divide() error { //							Деление (47)
	if c.operand_2 == 0 {
		c.Result = ""
		return ErrDivideByZero
	}
	res := c.operand_1 / c.operand_2
	return c.getSolution(res)
}

func (c Calc) Perform(message string, calc *Calc) error { //	Основной метод, запускает анализ строки,
	err := analyzeFill(message, calc) //						и метод арифм. операции
	if err != nil {
		return err
	}

	switch string(calc.operation) {
	case "+":
		err = calc.add()
	case "-":
		err = calc.substract()
	case "*":
		err = calc.multiply()
	case "/":
		err = calc.divide()
	}
	if err != nil {
		return err
	}
	return nil
}

func main() { //										Начало
	calculations := Calc{}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Input expression: (q for exit)")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message == "q" {
			break
		}
		err := calculations.Perform(message, &calculations)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(calculations.Result)
		}
	}
	fmt.Println("Exit")

}
