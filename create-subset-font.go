package main

import (
    "bufio"
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "strconv"
    "os/exec"
    "path"
    "log"
    "flag"
)

// Struct which contains the complete definition of the XML
type SVG struct {
    XMLName xml.Name `xml:"svg"`
    XMLNS   string   `xml:"xmlns,attr"`
    Defs    Defs     `xml:"defs"`
}

type Defs struct {
    XMLName xml.Name `xml:"defs"`
    Font    Font     `xml:"font"`
}

type Font struct {
    XMLName       xml.Name      `xml:"font"`
    Id            string        `xml:"id,attr"`
    HorizAdvX     string        `xml:"horiz-adv-x,attr"`
    FontFace      FontFace      `xml:"font-face"`
    MissingGlyph  MissingGlyph  `xml:"missing-glyph"`
    Glyphs        []Glyph       `xml:"glyph"`
}

type FontFace struct {
    XMLName      xml.Name `xml:"font-face"`
    FontFamily   string   `xml:"font-family,attr"`
    UnitsPerEm   string   `xml:"units-per-em,attr"`
    Ascent       string   `xml:"ascent,attr"`
    Descent      string   `xml:"descent,attr"`
    FontWeight   string   `xml:"font-weight,attr"`
    FontStyle    string   `xml:"font-style,attr"`
}

type MissingGlyph struct {
    XMLName    xml.Name `xml:"missing-glyph"`
    HorizAdvX  string   `xml:"horiz-adv-x,attr"`
}

type Glyph struct {
    XMLName    xml.Name `xml:"glyph"`
    Name       string   `xml:"glyph-name,attr"`
    Unicode    string   `xml:"unicode,attr"`
    HorizAdvX  string   `xml:"horiz-adv-x,attr"`
    D          string   `xml:"d,attr"`
}

func main() {
    helpPtr := flag.Bool("help", false, "show help")
    var dryrun bool
    flag.BoolVar(&dryrun, "d", false, "dry run")
    flag.BoolVar(&dryrun, "dryrun", false, "dry run")
    flag.Parse()
    if *helpPtr || len(flag.Args()) == 0 {
        printHelp()
        os.Exit(1)
        return
    }
    var icons[]string = readIconsFIle(flag.Args()[0])
    var xml SVG = readXML(flag.Args()[1])
    generateFont(icons, xml, flag.Args()[2], dryrun)
}

func printHelp() {
    fmt.Println(`Subset Font Creator
Usage: go run create-subset-font.go [-d] iconslist svgfont font
Arguments:
    -d, --dryrun  does not generate any file
    iconslist     list of glyph names separated by comma or line breaks
    svgfont       source SVG font file
    font          source OpenType or TTF font file`)
}

func readIconsFIle(filename string)[]string {
    readFile, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    fileScanner := bufio.NewScanner(readFile)
    fileScanner.Split(bufio.ScanLines)
    var icons[]string
    for fileScanner.Scan() {
        tokens := strings.Split(fileScanner.Text(), ",")
        icons = append(icons, tokens...)
    }
    readFile.Close()
    fmt.Println("Icons in input file:", len(icons))
    return icons
}

func readXML(filename string) SVG {
    // Open our xmlFile
    xmlFile, err := os.Open(filename)
    // If os.Open returns an error then handle it
    if err != nil {
        log.Fatal(err)
    }
    // Defer the closing of our xmlFile so that we can parse it later on
    defer xmlFile.Close()
    // Read our opened xmlFile as a byte array.
    byteValue, _ := ioutil.ReadAll(xmlFile)
    // Initialize SVG struct
    var svg SVG
    // Unmarshal the byteArray which contains the content of xmlFile into svg which we defined above
    xml.Unmarshal(byteValue, &svg)
    fmt.Println("Icons in input SVG file:", len(svg.Defs.Font.Glyphs))
    return svg
}

func generateFont(icons []string, svg SVG, input string, dryrun bool) {
    var spans[]string
    // Unicode character to be used by pyftsubset
    var unicodes[]string
    var found[]string
    // Copy object
    newSvg := svg
    // Clear Glyphs collection
    newSvg.Defs.Font.Glyphs = nil
    glyphs := svg.Defs.Font.Glyphs
    for i := 0; i < len(glyphs); i++ {
        glyph := glyphs[i]
        if contains(icons, glyph.Name) {
            newSvg.Defs.Font.Glyphs = append(newSvg.Defs.Font.Glyphs, glyph)
            ascii := strconv.QuoteToASCII(glyph.Unicode) // Unicode to ascii "'\u554a'"
            ascii = ascii[3:len(ascii)-1]                // removes quotes and \u = "554a"
            unicodes = append(unicodes, "U+" + ascii)
            spans = append(spans, fmt.Sprintf("<span class='fab fa-%s'></span>", glyph.Name))
            if (dryrun) {
                fmt.Println(" +", glyph.Name, "(U+" + strings.ToUpper(ascii) + ")")
                found = append(found, glyph.Name)
            }
        }
    }
    fmt.Println("Icons in output file:", len(newSvg.Defs.Font.Glyphs))
    // fmt.Println("\n", strings.Join(spans[:], "\n"))
    // fmt.Println(strings.Join(unicodes[:], ","))
    if (!dryrun) {
        writeSVGFile(input, newSvg)
        createFontFile(input, unicodes, "woff")
        createFontFile(input, unicodes, "woff2")
    } else {
        if(len(found) < len(icons)) {
            fmt.Println("\nMissing icons:")
            for i := 0; i < len(icons); i++ {
                if !contains(found, icons[i]) {
                    fmt.Println(" -", icons[i])
                }
            }
        }
    }
}

func createFontFile(input string, unicodes[]string, flavour string) {
    err := os.Mkdir("subset", os.ModePerm);
    if err != nil && !os.IsExist(err) {  // If dir exists it will raise an error
        log.Fatal(err)
    }
    filename := "./subset/" + basename(input) + ".subset." + flavour
    // https://fonttools.readthedocs.io/en/latest/subset/index.html
    var cmdSrr [6]string
    cmdSrr[0] = "pyftsubset"
    cmdSrr[1] = input
    cmdSrr[2] = "--unicodes=" + strings.Join(unicodes[:], ",")
    cmdSrr[3] = "--flavor=" + flavour
    cmdSrr[4] = "--output-file=" + filename
    cmdSrr[5] = "--no-ignore-missing-unicodes"
    // fmt.Println(strings.Join(cmdSrr[:], " "))
    cmd := exec.Command(cmdSrr[0], cmdSrr[1], cmdSrr[2], cmdSrr[3], cmdSrr[4], cmdSrr[5])
    stderr, _ := cmd.StderrPipe()
    err = cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
    // Read output errors
    reader := bufio.NewReader(stderr)
    var errLine string
    line, err := reader.ReadString('\n')
    for err == nil {
        line, err = reader.ReadString('\n')
        if err == nil {
            errLine = line  // Get last error line
        }
    }
    cmd.Wait()
    if errLine != "" {
        fmt.Println("\nAn error was found:")
        fmt.Println(errLine)
    } else {
        fmt.Println("Created font file", filename)
    }
}

func writeSVGFile(input string, svg SVG) {
    // Save struct into an XML byte string
    xml, _ := xml.MarshalIndent(svg, "", "  ")
    header := `<?xml version="1.0" encoding="UTF-8"?>\n
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd" >\n`
    xml = []byte(header + string(xml))
    outputFile := "./subset/" + basename(input) + ".subset.svg"
    err := os.WriteFile(outputFile, xml, 0644)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Generated SVG font file:", outputFile)
}

// Go does not has a contains method for arrays or slices
func contains(elems []string, v string) bool {
    for _, s := range elems {
        if v == s {
            return true
        }
    }
    return false
}

// Returns a file basename (without the extension)
func basename(filename string) string {
	return strings.TrimSuffix(path.Base(filename), path.Ext(filename))
}