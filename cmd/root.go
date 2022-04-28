package cmd

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gerow/go-color"
	"github.com/spf13/cobra"
)

type HEX string

type HEX_HSL struct {
	HEX HEX
	HSL color.HSL
}

const FILE_WRITE_PERM = 0644

func HEXToHSL(hex HEX) (color.HSL, error) {
	rgb, err := color.HTMLToRGB(string(hex))
	return rgb.ToHSL(), err
}

func distance(c1 color.HSL, c2 color.HSL) float64 {
	hDist := math.Abs(c1.H - c2.H)
	sDist := c1.S - c2.S
	lDist := c1.L - c2.L
	if hDist > 0.5 {
		hDist = 1.0 - hDist
	}
	return math.Sqrt(hDist*hDist + sDist*sDist + lDist*lDist)
}

func getClosestColor(target HEX, colors *[]HEX_HSL, threshold float64) (HEX, error) {
	targetHSL, err := HEXToHSL(target)
	if err != nil {
		return target, err
	}

	result := &HEX_HSL{HEX: target, HSL: targetHSL}
	minD := math.Inf(0)
	for _, color := range *colors {
		d := distance(targetHSL, color.HSL)
		if d < minD && d <= threshold {
			result = &color
			minD = d
		}
	}
	return result.HEX, nil
}

func genColorList(colors *[]HEX) (*[]HEX_HSL, error) {
	result := []HEX_HSL{}
	for _, c := range *colors {
		colorHSL, err := HEXToHSL(c)
		if err != nil {
			return &result, err
		}
		result = append(result, HEX_HSL{
			HEX: c,
			HSL: colorHSL,
		})
	}
	return &result, nil
}

func getColors(s string) []HEX {
	var result []HEX
	r := regexp.MustCompile(`#[\dABCDEFabcdef]{6}`)
	codes := r.FindAllStringSubmatch(s, -1)
	for _, code := range codes {
		result = append(result, HEX(code[0]))
	}
	return result
}

func exec(pathname string, colors *[]HEX, threshold float64, isDry bool) error {
	var colorsWithHls *[]HEX_HSL
	var fileNames []string
	var err error
	changedFileNames := []string{}
	replacedColorMap := map[HEX]HEX{}

	threshold = threshold * math.Sqrt(2.5)

	if colorsWithHls, err = genColorList(colors); err != nil {
		return err
	}
	if fileNames, err = filepath.Glob(pathname); err != nil {
		return err
	}

	for _, fileName := range fileNames {
		bytes, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		fileIn := string(bytes)
		fileOut := fileIn
		targets := getColors(fileIn)
		for _, target := range targets {
			var replaceHEX HEX
			if color, have := replacedColorMap[target]; have {
				replaceHEX = color
			} else {
				if replaceHEX, err = getClosestColor(target, colorsWithHls, threshold); err != nil {
					return err
				}
				replacedColorMap[target] = replaceHEX
			}
			fileOut = strings.ReplaceAll(fileOut, string(target), string(replaceHEX))
		}
		if fileOut != fileIn {
			changedFileNames = append(changedFileNames, fileName)
			if !isDry {
				os.WriteFile(fileName, []byte(fileOut), FILE_WRITE_PERM)
			}
		}
	}
	fmt.Println(changedFileNames)
	fmt.Println(replacedColorMap)
	return nil
}

var RootCmd = &cobra.Command{
	Use:   "replace",
	Short: "replace colors",
	Run: func(cmd *cobra.Command, args []string) {
		dry, _ := cmd.Flags().GetBool("dry")

		pathname, _ := cmd.Flags().GetString("pathname")
		colors_str, _ := cmd.Flags().GetString("colors")

		if pathname == "" || colors_str == "" {
			fmt.Println("path and colors are required")
			return
		}

		colors := []HEX{}

		threshold, _ := cmd.Flags().GetFloat64("threshold")

		for _, color := range strings.Split(colors_str, ",") {
			colors = append(colors, HEX(color))
		}
		fmt.Println(pathname, colors, threshold, dry)
		err := exec(pathname, &colors, threshold, dry)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.Flags().StringP("pathname", "p", "", "file pathname name (glob)")
	RootCmd.Flags().StringP("colors", "c", "", "hex colors (#000000,#111111...)")
	RootCmd.Flags().Float64("threshold", 0.1, "threshold (default: 0.1)")
	RootCmd.Flags().BoolP("dry", "d", false, "dry run")
}
