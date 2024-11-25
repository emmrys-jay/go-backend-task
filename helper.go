package main

import (
	"log"
)

const (
	arialFont     = "arial"
	arialBoldFont = "arial-bold"

	bold int8 = 1
	bare int8 = 0

	fontSizeSmall float64 = 5
	fontSizeMid   float64 = 10
	fontSizeBig   float64 = 15

	large int8 = 0
	mid   int8 = 1
	small int8 = 2
)

var (
	fontMap = map[int8]string{
		bold: arialBoldFont,
		bare: arialFont,
	}

	sizeMap = map[int8]float64{
		large: fontSizeBig,
		mid:   fontSizeMid,
		small: fontSizeSmall,
	}

	col1, col2, col3, col4, col5 float64
)

// getAlignRightPosition uses the length of text and starting position from the right (startPos) to
// determine the position of text for it to be aligned at the right
// Example: The texts below are aligned to the right
//
//	hello world
//			hrl
func getAlignRightPosition(text string, startPos float64) float64 {
	textWidth, err := pdf.MeasureTextWidth(text)
	if err != nil {
		log.Fatal(err)
	}

	return startPos - textWidth - margins
}

// getAlignBottomPosition uses the lastPos, lastHeight, text and an optional ms (margins) parameter to
// determine the position of text for it to be aligned to the bottom of the page.
//
// use text if the last inserted element was a string, else use lastHeight. Set lastHeight to 0 if it is
// not used. lastHeight is prioritized above text
func getAlignBottomPosition(lastPos, lastHeight float64, text string, ms ...float64) float64 {
	if lastHeight != 0 {
		pos := lastPos - margins - lastHeight
		if len(ms) > 0 {
			pos -= ms[0]
		}
		return pos
	}

	textHeight, err := pdf.MeasureCellHeightByText(text)
	if err != nil {
		log.Fatal(err)
	}

	pos := lastPos - margins - textHeight
	if len(ms) > 0 {
		pos -= ms[0]
	}
	return pos
}

// getNextX uses the width of the last inserted text and its starting x position (startPos) to determine the
// starting position of the next x on the same line. margin defines the distance between the two texts
func getNextX(text string, margin float64, startPos float64) float64 {
	textWidth, err := pdf.MeasureTextWidth(text)
	if err != nil {
		log.Fatal(err)
	}

	return startPos + textWidth + margin
}

// getNextY uses the height of the last inserted text and its starting y position (startPos) to determine the
// starting position of the next y on the same line
func getNextY(lastLine string, startPos float64) float64 {
	textHeight, err := pdf.MeasureCellHeightByText(lastLine)
	if err != nil {
		log.Fatal(err)
	}

	return startPos + textHeight + lineSpace
}

// addFonts adds the available fonts
func addFonts() {
	err := pdf.AddTTFFont(arialFont, "./assets/font/Arial_Unicode.ttf")
	if err != nil {
		log.Fatal(err)
	}

	err = pdf.AddTTFFont(arialBoldFont, "./assets/font/Arial_Unicode_bold.ttf")
	if err != nil {
		log.Fatal(err)
	}
}

// setFont sets the current font based on font and size.
// font can either be 'bold' or 'bare'.
// size can either be 'large', 'mid' or 'small'
func setFont(font, size int8) {
	err := pdf.SetFont(fontMap[font], "", sizeMap[size])
	if err != nil {
		log.Fatal(err)
	}
}

// setColumnsXPos sets the x positions of the five colums for a table in the generated pdf.
// It must be calld before calling writeRow().
func setColumnsXPos(one, two, three, four, five float64) {
	col1, col2, col3, col4, col5 = one, two, three, four, five
}

// writeRow takes the current y position and at least 5 values to write a row as part of a
// 5-column table. values must contain at least 5 values or this function panics. An optional 6th
// string can be added to values to be placed under the 5th string for a multiline entry.
//
// setColumnsXPos() should be called before writeRow() to set the x position for values to be entered.
// You can call setColumsXPos() once and use it to write multiple rows.
func writeRow(yPos float64, values []string) (newY float64) {
	if len(values) < 5 {
		panic("values passed in must be at least five")
	}

	pdf.SetXY(col1, yPos)
	pdf.Text(values[0])

	pdf.SetXY(col2, yPos)
	pdf.Text(values[1])

	pdf.SetXY(col3, yPos)
	pdf.Text(values[2])

	pdf.SetXY(col4, yPos)
	pdf.Text(values[3])

	pdf.SetXY(col5, yPos)
	pdf.Text(values[4])
	text := values[4]

	if len(values) > 5 {
		yPos = getNextY(text, yPos)
		pdf.SetXY(col5, yPos)
		text = values[5]
		pdf.Text(text)
	}

	return yPos
}

// drawLine takes the current x position, y posion and an addon. addon is added to the y positon along with
// the default linespace to add a margin between the drawn line and coming texts
func drawLine(xPos, yPos, addon float64) (newY float64) {
	yPos += lineSpace + addon
	pdf.SetLineWidth(1)
	pdf.Line(xPos, yPos, getAlignRightPosition("", pageWidth), yPos)
	yPos += lineSpace + 10
	return yPos
}

// writeText abstracts setting coordinates and writing text to the pdf. xMargin is added to xPos, and
// yMargin is added to yPos before text is inserted into the pdf. The derived x and y values after
// addition  is returned.
func writeText(xPos, yPos, xMargin, yMargin float64, text string) (newX, newY float64) {
	newX = xPos + xMargin
	newY = yPos + yMargin
	pdf.SetXY(newX, newY)
	pdf.Text(text)

	return newX, newY
}
