package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/signintech/gopdf"
)

var (
	documentType = "EUR Statement"
	bank         = "Revolut Bank UAB"

	pageWidth  = gopdf.PageSizeA4.W
	pageHeight = gopdf.PageSizeA4.H

	margins   float64 = 30
	lineSpace float64 = 3

	pdf = &gopdf.GoPdf{}
)

var (
	accountStatement AccountStatement
)

func init() {
	file, err := os.Open("account_statement.json")
	if err != nil {
		log.Fatalln("Error opening account statement file, error: ", err)
	}
	defer file.Close()

	var account AccountStatement
	if err := json.NewDecoder(file).Decode(&account); err != nil {
		log.Fatalln("Error decoding file content, error: ", err)
	}

	accountStatement = account
}

func main() {
	var (
		bodyXPos, bodyYPos float64
	)

	// fmt.Printf("%+v", account)
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})

	fmt.Println("Generating PDF file...")

	addFonts()
	setFont(bare, mid)

	writeHeader()
	writeFooter()

	bodyXPos = getNextX("", 0, margins)
	bodyYPos = getNextY("", 60) // 60 was derived from the height of the logo

	writePage(bodyXPos, bodyYPos)

	// Write to file
	pdf.WritePdf("result.pdf")
	fmt.Println("Finished generating PDF file. Your file is ready!")
}

func writeHeader() {
	var (
		currentYPos float64

		text string
	)

	pdf.AddHeader(func() {
		setFont(bold, large)

		err := pdf.Image("./assets/image/logo.jpg", margins, margins, nil)
		if err != nil {
			log.Fatal(err)
			return
		}

		currentYPos = margins
		text = documentType
		pdf.SetXY(getAlignRightPosition(text, pageWidth), margins)
		pdf.Text(text)

		setFont(bare, mid)

		currentYPos = getNextY(text, currentYPos)
		text = fmt.Sprintf("Generated on %s", accountStatement.CreatedAt)
		pdf.SetXY(getAlignRightPosition(text, pageWidth), currentYPos)
		pdf.Text(text)

		currentYPos = getNextY(text, currentYPos)
		text = bank
		pdf.SetXY(getAlignRightPosition(text, pageWidth), currentYPos)
		pdf.Text(text)
	})
}

func writeFooter() {
	var (
		currentYPos float64
		currentXPos float64

		text string
	)

	pdf.AddFooter(func() {
		setFont(bold, mid)

		currentXPos = margins
		text = "© 2023 " + bank
		currentYPos = getAlignBottomPosition(pageHeight, 0, text)
		writeText(currentXPos, currentYPos, 0, 0, text)

		text = "Page 1 of 5"
		currentXPos = getAlignRightPosition(text, pageWidth) - 5
		writeText(currentXPos, currentYPos, 0, 0, text)

		currentXPos = margins
		currentYPos = getAlignBottomPosition(currentYPos, 0, text, 3)
		err := pdf.Image("./assets/image/randomqr.jpg", currentXPos, currentYPos, &gopdf.Rect{W: 30, H: 30})
		if err != nil {
			log.Fatal(err)
			return
		}

		setFont(bare, small)

		currentXPos += 30 + 3 // width of rectangle plus little margin
		text = "Report lost or stolen card"
		writeText(currentXPos, currentYPos, 0, 0, text)

		textWidth, err := pdf.MeasureTextWidth(text)
		if err != nil {
			log.Fatal(err)
			return
		}

		var bankDescription = "Revolut UAB bank is a credit institution licensed in the republik of Lithuania with company name given by "
		var bankDescriptionXPos float64
		bankDescriptionXPos, _ = writeText(currentXPos, currentYPos, textWidth+10, 0, bankDescription)

		currentYPos = getNextY(text, currentYPos)
		text = "+325 612 799"
		writeText(currentXPos, currentYPos, 0, 0, text)
		writeText(bankDescriptionXPos, currentYPos, 0, 0, bankDescription)

		currentYPos = getNextY(text, currentYPos)
		text = "Get help directly in the app"
		writeText(currentXPos, currentYPos, 0, 0, text)
		writeText(bankDescriptionXPos, currentYPos, 0, 0, bankDescription)

		currentYPos = getNextY(text, currentYPos)
		text = "Scan the QR code"
		writeText(currentXPos, currentYPos, 0, 0, text)
		writeText(bankDescriptionXPos, currentYPos, 0, 0, bankDescription)
	})

}

func writePage(xPos, yPos float64) {
	var (
		text        string
		currentXPos = xPos
		currentYPos = yPos
	)

	pdf.AddPage()

	setFont(bold, large)

	// Write name and address
	text = strings.ToUpper(accountStatement.Name)
	currentXPos, currentYPos = writeText(currentXPos, currentYPos, 0, 50, text)

	setFont(bold, mid)

	currentYPos = getNextY(text, currentYPos)
	text = fmt.Sprint(accountStatement.Address.HouseNo, " ", accountStatement.Address.Street)
	currentXPos, currentYPos = writeText(currentXPos, currentYPos, 0, 10, text)

	currentYPos = getNextY(text, currentYPos)
	text = accountStatement.Address.City
	writeText(currentXPos, currentYPos, 0, 0, text)

	currentYPos = getNextY(text, currentYPos)
	text = accountStatement.Address.State
	writeText(currentXPos, currentYPos, 0, 0, text)

	currentYPos = getNextY(text, currentYPos)
	text = accountStatement.Address.Country
	writeText(currentXPos, currentYPos, 0, 0, text)

	// Write IBAN details
	IBANDiclaimer1 := "(You cannot use this IBAN for bank transfers."
	IBANDiclaimer2 := "Please use the IBAN found in the app)"

	Iban := "IBAN"
	Bic := "BIC"
	currentXPos = getAlignRightPosition(IBANDiclaimer1, pageWidth)
	IbanXPos := getAlignRightPosition(Iban, currentXPos) - 10 // added value for the space between IBBAN and its value

	for i, v := range accountStatement.Iban {

		if i == 0 {
			currentYPos = getNextY(" ", currentYPos)
		} else {
			currentYPos = getNextY(text, currentYPos) + 30 // added value for the space between each IBAN and BIC entry
		}

		// Writing IBAN
		// IBAN
		setFont(bold, mid)
		writeText(IbanXPos, currentYPos, 0, 0, Iban)

		// IBAN No
		setFont(bare, mid)
		text = v.No
		writeText(currentXPos, currentYPos, 0, 0, text)

		// Writing BIC
		// BIC
		setFont(bold, mid)
		currentYPos = getNextY(text, currentYPos)
		writeText(IbanXPos, currentYPos, 0, 0, Bic)

		// BIC Value
		setFont(bare, mid)
		text = v.Bic
		writeText(currentXPos, currentYPos, 0, 0, text)

		if i > 0 {
			currentYPos = getNextY(text, currentYPos)
			text = IBANDiclaimer1
			writeText(currentXPos, currentYPos, 0, 0, text)

			currentYPos = getNextY(text, currentYPos)
			text = IBANDiclaimer2
			writeText(currentXPos, currentYPos, 0, 0, text)
		}
	}

	// Set new X and Y position
	currentXPos = getNextX("", margins, 0)
	currentYPos = getNextY(text, currentYPos)

	// ##########################################################################################
	// Balance Summary Table

	setFont(bold, large)
	currentYPos = getNextY(text, currentYPos)
	pdf.SetXY(currentXPos, currentYPos)
	text = "Balance summary"
	pdf.Text(text)

	setFont(bold, mid)
	var (
		productXPos        float64
		openingBalanceXPos float64
		moneyInXPos        float64
		moneyOutXPos       float64
		closingBalanceXPos float64
	)

	productXPos = currentXPos
	openingBalanceXPos = productXPos + 150
	moneyOutXPos = getNextX("Opening balance", 20, openingBalanceXPos) + 20
	moneyInXPos = getNextX("Money out", 20, moneyOutXPos) + 20
	closingBalanceXPos = getAlignRightPosition("Closing", pageWidth)

	currentYPos = getNextY(text, currentYPos) + 20

	setColumnsXPos(productXPos, openingBalanceXPos, moneyOutXPos, moneyInXPos, closingBalanceXPos)
	currentYPos = writeRow(currentYPos, []string{
		"Product",
		"Opening balance",
		"Money out",
		"Money in",
		"Closing",
		"balance",
	})
	currentYPos = drawLine(currentXPos, currentYPos, 0)
	setFont(bare, mid)

	var totalOb, totalMi, totalMo, totalCb float64
	for _, v := range accountStatement.BalanceSummary.Products {

		currentYPos = writeRow(currentYPos, []string{
			v.Product,
			fmt.Sprint(accountStatement.CurrencySymbol, v.OpeningBalance),
			fmt.Sprint(accountStatement.CurrencySymbol, v.MoneyOut),
			fmt.Sprint(accountStatement.CurrencySymbol, v.MoneyIn),
			fmt.Sprint(accountStatement.CurrencySymbol, v.ClosingBalance),
		})

		totalOb += v.OpeningBalance
		totalMo += v.MoneyOut
		totalMi += v.MoneyIn
		totalCb += v.ClosingBalance
		currentYPos = drawLine(currentXPos, currentYPos, 10)
	}

	// Total
	currentYPos = writeRow(currentYPos, []string{
		"Total",
		fmt.Sprint(accountStatement.CurrencySymbol, totalOb),
		fmt.Sprint(accountStatement.CurrencySymbol, totalMo),
		fmt.Sprint(accountStatement.CurrencySymbol, totalMi),
		fmt.Sprint(accountStatement.CurrencySymbol, totalCb),
	})
	currentYPos = getNextY(text, currentYPos) + 10

	setFont(bare, small)
	pdf.SetXY(currentXPos, currentYPos)
	text = "This statement in your statement might differ from the balance shown in your app"
	pdf.Text(text)

	// ######################################################################
	// Account transactions table
	setFont(bold, large)
	currentYPos = getNextY(text, currentYPos) + 20
	pdf.SetXY(currentXPos, currentYPos)
	text = "Account transactions from 1 February 2023 to 29 March 2023"
	pdf.Text(text)

	var (
		dateXPos        float64
		descriptionXPos float64
	)

	setFont(bold, mid)

	dateXPos = currentXPos
	descriptionXPos = dateXPos + 100
	currentYPos = getNextY(text, currentYPos) + 20

	setColumnsXPos(dateXPos, descriptionXPos, moneyOutXPos, moneyInXPos, closingBalanceXPos)
	currentYPos = writeRow(currentYPos, []string{
		"Date",
		"Description",
		"Money out",
		"Money in",
		"Balance",
	})
	currentYPos = drawLine(currentXPos, currentYPos, 0)

	setFont(bare, mid)

	for _, v := range accountStatement.Transactions {
		currentYPos = writeRow(currentYPos, []string{
			v.Date,
			v.Description,
			fmt.Sprint(accountStatement.CurrencySymbol, v.MoneyOut),
			fmt.Sprint(accountStatement.CurrencySymbol, v.MoneyIn),
			fmt.Sprint(accountStatement.CurrencySymbol, v.Balance),
		})
		currentYPos = drawLine(currentXPos, currentYPos, 10)
	}
}
