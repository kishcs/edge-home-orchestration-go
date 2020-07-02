/*******************************************************************************
 * Copyright 2019 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/

// Package resourceutil provides the information of resource usage of local device
package resourceutil

import (
	"fmt"
	"time"

	"restinterface/resthelper"

	netDB "db/bolt/network"
)

const (
	pingAPI            = "/api/v1/ping"
	internalPort       = 56002
	defaultRttDuration = 5
)

var (
	helper        resthelper.RestHelper
	netDBExecutor netDB.DBInterface
)

func init() {
	helper = resthelper.GetHelper()
	netDBExecutor = netDB.Query{}
}

func processRTT() {
	go func() {
		for {
			netInfos, err := netDBExecutor.GetList()
			if err != nil {
				return
			}

			for _, netInfo := range netInfos {
				totalCount := len(netInfo.IPv4)
				ch := make(chan float64, totalCount)
				for _, ip := range netInfo.IPv4 {
					go func(targetIP string) {
						ch <- checkRTT(targetIP)
					}(ip)
				}
				go func(info netDB.NetworkInfo) {
					result := selectMinRTT(ch, totalCount)
					info.RTT = result
					netDBExecutor.Update(info)
				}(netInfo)
			}
			time.Sleep(time.Duration(defaultRttDuration) * time.Second)
		}
	}()
}

func processTotal() {
    log.Printf(" ==== processTotal() ==== ")
    go func() {
        for {
            calculateScore()
			//time.Sleep(time.Duration(1) * time.Second)
			time.Sleep(time.Duration(defaultRttDuration) * time.Second)
        }
    }()
}

func calculateScore() {

	ips, _ = networkhelper.GetInstance().GetOutboundIP()
	//ips := netInfo.GetIP()
	log.Printf(" ==== processTotal() ==== LINUX NATIVE ==> [%s] ips", ips)
	//var st uint64
	//st = androidexecutor.GetInstance().GetStatus()
	//log.Printf(" ANDROID ==> [%d] executionStatus", st)

	type Scores struct {
		Scpu   string `json:"cpu"`
		Smem   string `json:"mem"`
		Snet   string `json:"net"`
		Sren   string `json:"ren"`
		SIp    string `json:"ip"`
		Status string `json:"status"`
		Score  string `json:"score"`
	}

	var out Scores
	var netVal float64
	out.SIp = ips

	out.Scpu = strconv.FormatFloat(cpuScores, 'f', 6, 64)
	out.Smem = strconv.FormatFloat(mems, 'f', 6, 64)
	netVal = 1 / (8770 * math.Pow(nets, -0.9))
	out.Snet = strconv.FormatFloat(netVal, 'f', 6, 64)
	out.Sren = strconv.FormatFloat(rtts, 'f', 6, 64)
	finaScore := float64(netVal + (cpuScores / 2) + rtts)
	//out.Status = strconv.FormatUint(st, 16)
	out.Score = strconv.FormatFloat(finaScore, 'f', 6, 64)

	// message := "cpu" + cScore + "," + "MemoryAvail" + mem + "," + "NetScore : " + netscore + "," + "RenderScore : " + rendscore + "," + "Score : " + Score

	in, err := json.Marshal(out)

	//ioutil.WriteFile("/storage/emulated/0/Android/data/com.samsung.orchestration.service/files/score.json", []byte(in), 0644)
	ioutil.WriteFile("/tmp/score.json", []byte(in), 0644)

	//androidexecutor.GetInstance().SetStatus(0)
	if err != nil {
		return
	}
}

func calculateScore1() {

	ips, _ = networkhelper.GetInstance().GetOutboundIP()
	 //ips := netInfo.GetIP()
	 //log.Printf("------------[%s]---------------ips",ips)
	/*var st uint64
	st = nativeexecutor.GetInstance().GetStatus()
	if st == 1 {
	   log.Printf("------------[%d]--------------- executionStatus",st)
	}*/

	type Scores struct{
	   Scpu string `json:"cpu"`
	   Smem string `json:"mem"`
	   Snet string `json:"net"`
	   Sren string `json:"ren"`
	   SIp string `json:"ip"`
	   Status string `json:"status"`
	   Score string `json:"score"`
	}

	var out Scores
	var netVal float64
	
	out.SIp = ips
	out.Scpu = strconv.FormatFloat(cpuScores, 'f', 6, 64)
	out.Smem = strconv.FormatFloat(mems, 'f', 6, 64)
	netVal = 1 / (8770 * math.Pow(nets, -0.9))
	out.Snet = strconv.FormatFloat(netVal, 'f', 6, 64)
	out.Sren = strconv.FormatFloat(rtts, 'f', 6, 64)
	finaScore := float64(netVal + (cpuScores / 2) + rtts)
	out.Status = strconv.FormatUint(st, 16)
	out.Score = strconv.FormatFloat(finaScore, 'f', 6, 64)

   // message := "cpu" + cScore + "," + "MemoryAvail" + mem + "," + "NetScore : " + netscore + "," + "RenderScore : " + rendscore + "," + "Score : " + Score

	in, err := json.Marshal(out)

	//ioutil.WriteFile("/tmp/score.json", []byte(in), 0644)

	//nativeexecutor.GetInstance().SetStatus(0)

	service1, err := ioutil.ReadFile("serverip.txt")

	if err != nil {
		log.Println(logPrefix, "Dashboard error : ", err.Error())
		return
	}

	service := string(service1)
	if len(service) <= 0 {
		log.Printf("---service NULL ")
		return
	}
	service = strings.TrimSuffix(service, "\n")
	if len(service) <= 0 {
		log.Printf("service NULL ")
		return
	}

	//service := "107.108.87.9:1046"
	tcpAddr, err := network.ResolveTCPAddr("tcp4", service)
	if err != nil {
		log.Println(logPrefix, "KKK error : ", err.Error())
		return
	}
	conn, err := network.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(logPrefix, "Dashboard error : ", err.Error())
		return
	}
	_, err = conn.Write([]byte(in))
	if err != nil {
		log.Println(logPrefix, "KKK error : ", err.Error())
		return
	}
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Println(logPrefix, "KKK error : ", err.Error())
		return
	}

	fmt.Println(string(result))

	//nativeexecutor.GetInstance().SetStatus(0)

	if err != nil {
			return
	}
}

func checkRTT(ip string) (rtt float64) {
	targetURL := helper.MakeTargetURL(ip, internalPort, pingAPI)

	reqTime := time.Now()
	_, _, err := helper.DoGet(targetURL)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return time.Now().Sub(reqTime).Seconds()
}

func selectMinRTT(ch chan float64, totalCount int) (minRTT float64) {
	for i := 0; i < totalCount; i++ {
		select {
		case rtt := <-ch:
			if (rtt != 0 && rtt < minRTT) || minRTT == 0 {
				minRTT = rtt
			}
		}
	}
	return
}
