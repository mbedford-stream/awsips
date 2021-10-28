package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

//AWSIPs JSON format for AWS IPs Range
type AWSIPs struct {
	SyncToken  string `json:"syncToken"`
	CreateDate string `json:"createDate"`
	Prefixes   []struct {
		IPPrefix string `json:"ip_prefix"`
		Region   string `json:"region"`
		Service  string `json:"service"`
	} `json:"prefixes"`
}

// type AWSIPsV6 struct {
//   SyncToken  string `json:"syncToken"`
//   CreateDate string `json:"createDate"`
//   Prefixes   []struct {
//     IPPrefix string `json:"ipv6_prefix"`
//     Region   string `json:"region"`
//     Service  string `json:"service"`
//   } `json:"ipv6_prefixes"`
// }

type structAWSIPs struct {
	CreateDate   string `json:"createDate"`
	Ipv6Prefixes []struct {
		Ipv6Prefix string `json:"ipv6_prefix"`
		Region     string `json:"region"`
		Service    string `json:"service"`
	} `json:"ipv6_prefixes"`
	Prefixes []struct {
		IPPrefix string `json:"ip_prefix"`
		Region   string `json:"region"`
		Service  string `json:"service"`
	} `json:"prefixes"`
	SyncToken string `json:"syncToken"`
}

// getIPRanges function created/contributed by @amscotti
func getIPRanges() structAWSIPs {
	const url string = "https://ip-ranges.amazonaws.com/ip-ranges.json"

	res, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	// fmt.Println(body)

	var ipRanges structAWSIPs
	json.Unmarshal(body, &ipRanges)

	// fmt.Println(ipRanges)

	return ipRanges
}

func ipInNet(ipTest, netTest string) bool {
	// addr := "170.149.172.130/20"
	_, netA, err := net.ParseCIDR(netTest)
	if err != nil {
		log.Fatal("Could not parse network")
	}
	// fmt.Println(ipA, netA)
	addrTest := net.ParseIP(ipTest)
	if addrTest == nil {
		log.Fatal("Not a valid IP")
	}
	if netA.Contains(addrTest) == true {
		// fmt.Println("Yay")
		return true
	}
	return false

}

func main() {
	var awsService string
	var awsReg string
	var netSearch string
	var ipSearch string
	var vFour bool
	var vSix bool

	flag.StringVar(&awsService, "ser", "none", "AWS Service to generate addresses for (ec2, aws, cf, r53, or all")
	flag.StringVar(&awsReg, "reg", "none", "AWS region to return data for")
	flag.StringVar(&netSearch, "net", "none", "Search for a specific CIDR block")
	flag.StringVar(&ipSearch, "ip", "none", "Search for a specific IP within all networks")

	flag.BoolVar(&vFour, "v4", false, "Return only IPv4 values")
	flag.BoolVar(&vSix, "v6", false, "Return only IPv6 values")

	flag.Parse()

	if vFour && vSix {
		log.Fatal("You can't return ONLY v4 AND v6 at the same time..... ")
	}

	ipRanges := getIPRanges()

	ipTime, _ := time.Parse("2006-01-02-15-04-05", ipRanges.CreateDate)
	fmt.Printf("\nData accurate as of: %s\n", ipTime.Format(time.ANSIC))

	if netSearch != "none" {
		fmt.Printf("%-25s%-20s%-10s\n=========================================================================\n", "IP Prefix", "Region", "Service")

		if !vSix {
			for _, ip := range ipRanges.Prefixes {
				if ip.IPPrefix == netSearch {
					fmt.Printf("%-45s%-20s%-10s\n", ip.IPPrefix, ip.Region, ip.Service)
				}
			}
		}

		if !vFour {
			for _, ip := range ipRanges.Ipv6Prefixes {
				if ip.Ipv6Prefix == netSearch {
					fmt.Printf("%-45s%-20s%-10s\n", ip.Ipv6Prefix, ip.Region, ip.Service)
				}
			}
		}

		return
	}

	if awsReg == "none" && (awsService == "none" || awsService == "all") {
		fmt.Printf("%-25s%-20s%-10s\n=========================================================================\n", "IP Prefix", "Region", "Service")
		if !vSix {
			for _, ip := range ipRanges.Prefixes {
				if ipSearch != "none" {
					ipCheck := ipInNet(ipSearch, ip.IPPrefix)
					if ipCheck == true {
						fmt.Printf("%-45s%-20s%-10s\n", ip.IPPrefix, ip.Region, ip.Service)
					} else {
						continue
					}
				} else {
					fmt.Printf("%-45s%-20s%-10s\n", ip.IPPrefix, ip.Region, ip.Service)
				}
			}
		}

		if !vFour {
			for _, ip := range ipRanges.Ipv6Prefixes {
				if ipSearch != "none" {
					ipCheck := ipInNet(ipSearch, ip.Ipv6Prefix)
					if ipCheck == true {
						fmt.Printf("%-45s%-20s%-10s\n", ip.Ipv6Prefix, ip.Region, ip.Service)
					} else {
						continue
					}
				} else {
					fmt.Printf("%-45s%-20s%-10s\n", ip.Ipv6Prefix, ip.Region, ip.Service)
				}
			}
		}

		return
	}

	if awsService != "none" && awsReg == "none" {
		fmt.Printf("%-45s%-20s%-10s\n=========================================================================\n", "IP Prefix", "Region", "Service")

		if !vSix {
			for _, ip := range ipRanges.Prefixes {
				serMatch, sErr := regexp.MatchString(strings.ToLower(awsService), strings.ToLower(ip.Service))
				if sErr != nil {
					log.Fatal("Could not use -ser as a valid match term")
				}
				if serMatch {
					fmt.Printf("%-45s%-20s%-10s\n", ip.IPPrefix, ip.Region, ip.Service)
				}
			}
		}
		if !vFour {
			for _, ip := range ipRanges.Ipv6Prefixes {
				serMatch, sErr := regexp.MatchString(strings.ToLower(awsService), strings.ToLower(ip.Service))
				if sErr != nil {
					log.Fatal("Could not use -ser as a valid match term")
				}
				if serMatch {
					fmt.Printf("%-45s%-20s%-10s\n", ip.Ipv6Prefix, ip.Region, ip.Service)
				}
			}

		}
	}

	if awsService != "none" && awsReg != "none" {
		fmt.Printf("%-45s%-20s%-10s\n=========================================================================\n", "IP Prefix", "Region", "Service")

		if !vSix {
			for _, ip := range ipRanges.Prefixes {
				serMatch, sErr := regexp.MatchString(strings.ToLower(awsService), strings.ToLower(ip.Service))
				if sErr != nil {
					log.Fatal("Could not use -ser as a valid match term")
				}
				regMatch, rErr := regexp.MatchString(strings.ToLower(awsReg), strings.ToLower(ip.Region))
				if rErr != nil {
					log.Fatal("Could not use -reg as a valid match term")
				}

				if serMatch && regMatch {
					fmt.Printf("%-45s%-20s%-10s\n", ip.IPPrefix, ip.Region, ip.Service)
				}
			}
		}

		if !vFour {
			for _, ip := range ipRanges.Ipv6Prefixes {
				serMatch, sErr := regexp.MatchString(strings.ToLower(awsService), strings.ToLower(ip.Service))
				if sErr != nil {
					log.Fatal("Could not use -ser as a valid match term")
				}
				regMatch, rErr := regexp.MatchString(strings.ToLower(awsReg), strings.ToLower(ip.Region))
				if rErr != nil {
					log.Fatal("Could not use -reg as a valid match term")
				}

				if serMatch && regMatch {
					fmt.Printf("%-45s%-20s%-10s\n", ip.Ipv6Prefix, ip.Region, ip.Service)
				}
			}

		}
	}

	if awsService == "none" && awsReg != "none" {
		fmt.Printf("%-45s%-20s%-10s\n=========================================================================\n", "IP Prefix", "Region", "Service")

		if !vSix {
			for _, ip := range ipRanges.Prefixes {
				regMatch, rErr := regexp.MatchString(strings.ToLower(awsReg), strings.ToLower(ip.Region))
				if rErr != nil {
					log.Fatal("Could not use -reg as a valid match term")
				}

				if regMatch {
					fmt.Printf("%-45s%-20s%-10s\n", ip.IPPrefix, ip.Region, ip.Service)
				}
			}
		}

		if !vFour {
			for _, ip := range ipRanges.Ipv6Prefixes {
				regMatch, rErr := regexp.MatchString(strings.ToLower(awsReg), strings.ToLower(ip.Region))
				if rErr != nil {
					log.Fatal("Could not use -reg as a valid match term")
				}

				if regMatch {
					fmt.Printf("%-45s%-20s%-10s\n", ip.Ipv6Prefix, ip.Region, ip.Service)
				}
			}

		}
	}

	return
}
