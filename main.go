package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"io"

	"github.com/getlantern/systray"
)

// Struct to parse JSON response
type IPInfo struct {
	IP         string `json:"ip"`
	Country    string `json:"country"`
	RegionCode string `json:"region_code"`
	ZipCode    string `json:"zip_code"`
	City       string `json:"city"`
	ASN        string `json:"asn"`
}

func main() {
	log.Println(`Starting locra...`)
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon())
	systray.SetTitle("locra")
	systray.SetTooltip("IP Information")

	// Menu items
	mIP := systray.AddMenuItem("IP: Fetching...", "Your IP will appear here")
	mCountry := systray.AddMenuItem("Country: Fetching...", "Your Country will appear here")
	mRegionCode := systray.AddMenuItem("Region Code: Fetching...", "Your Region Code will appear here")
	mZipCode := systray.AddMenuItem("Zip Code: Fetching...", "Your Zip Code will appear here")
	mCity := systray.AddMenuItem("City: Fetching...", "Your City will appear here")
	mASN := systray.AddMenuItem("ASN: Fetching...", "Your ASN will appear here")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Fetch IP info and update menu items
	go func() {
		for {
			select {
			case <-mIP.ClickedCh:
				info, err := fetchIPInfo()
				if err != nil {
					log.Println("Error fetching IP info:", err)
					mIP.SetTitle("IP: Error fetching IP info")
				} else {
					mIP.SetTitle(fmt.Sprintf("IP: %s", info.IP))
					mCountry.SetTitle(fmt.Sprintf("Country: %s", info.Country))
					mRegionCode.SetTitle(fmt.Sprintf("Region Code: %s", info.RegionCode))
					mZipCode.SetTitle(fmt.Sprintf("Zip Code: %s", info.ZipCode))
					mCity.SetTitle(fmt.Sprintf("City: %s", info.City))
					mASN.SetTitle(fmt.Sprintf("ASN: %s", info.ASN))
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

	// Fetch IP info initially
	go func() {
		info, err := fetchIPInfo()
		if err != nil {
			log.Println("Error fetching IP info:", err)
			mIP.SetTitle("IP: Error fetching IP info")
		} else {
			mIP.SetTitle(fmt.Sprintf("IP: %s", info.IP))
			mCountry.SetTitle(fmt.Sprintf("Country: %s", info.Country))
			mRegionCode.SetTitle(fmt.Sprintf("Region Code: %s", info.RegionCode))
			mZipCode.SetTitle(fmt.Sprintf("Zip Code: %s", info.ZipCode))
			mCity.SetTitle(fmt.Sprintf("City: %s", info.City))
			mASN.SetTitle(fmt.Sprintf("ASN: %s", info.ASN))
		}
	}()
}

func onExit() {
	// Cleanup here if needed
}

func fetchIPInfo() (*IPInfo, error) {
	resp, err := http.Get("https://ifconfig.co/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return nil, err
	}

	var ipInfo IPInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return nil, err
	}

	return &ipInfo, nil
}

func getIcon() []byte {
	// Default icon (display "loc" as text)
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG header
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10, // 16x16 pixels
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0xF3, 0xFF, // ... and more
		0x61, 0x00, 0x00, 0x00, 0x19, 0x74, 0x45, 0x58, // tEXt chunk
		0x74, 0x53, 0x6F, 0x66, 0x74, 0x77, 0x61, 0x72,
		0x65, 0x00, 0x41, 0x64, 0x6F, 0x62, 0x65, 0x20,
		0x49, 0x6D, 0x61, 0x67, 0x65, 0x52, 0x65, 0x61,
		0x64, 0x79, 0x71, 0xC9, 0x65, 0x3C, 0x00, 0x00,
		0x00, 0x0B, 0x49, 0x44, 0x41, 0x54, 0x78, 0xDA,
		0x63, 0xF8, 0xFF, 0xFF, 0x3F, 0x03, 0x00, 0x06,
		0x05, 0x02, 0xFE, 0xA7, 0x89, 0xE8, 0x47, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}
}
