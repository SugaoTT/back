// Copyright 2016 CNI authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"os/exec"
	"strconv"

	"github.com/vishvananda/netlink"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	current "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/cni/pkg/version"

	"github.com/containernetworking/plugins/pkg/ns"
	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
)

type CNIData struct {
	CniVersion   string       `json:"cniVersion"`
	Type         string       `json:"type"`
	PodName      string       `json:"pod-name"`
	InterfaceSet InterfaceSet `json:"interface"`
}

type InterfaceSet struct {
	Items []Items `json:"items"`
}

type Items struct {
	Name           string `json:"name"`
	TargetPodName  string `json:"target-pod-name"`
	TargetPodNIC   string `json:"target-pod-nic"`
	SelfTunnelID   string `json:"self-tunnel-id"`
	TargetTunnelID string `json:"target-tunnel-id"`
	SessionID      string `json:"session-id"`
}

type TempNetwork struct {
	InterfaceName string `json: "interface-name"`
	TargetNode    string `json: "target-node"`
}

func parseNetConf(bytes []byte) (*types.NetConf, error) {
	conf := &types.NetConf{}

	return conf, nil
}

/*
*
サーバ内にあるl2tpethの中から最大のもの(自身が使うもの)の数値を取ってくる関数
*/
func calcL2tpEth() (int, error) {
	// l2tpethの最大値を検索する
	output, err := exec.Command("bash", "-c", "sudo ip link show | grep l2tpeth | awk '{print $2}' | sed 's/l2tpeth//g'").Output() //, "|", "grep", "l2tpeth", "|", "awk", "'{print $2}'", "|", "sed", "'s/l2tpeth//g'").Output()
	ioutil.WriteFile("l2tpdata.txt", []byte(output), 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to execute command: %s", err)
	}

	// 出力から最大値を抽出
	max := 0
	numbers := strings.Split(string(output), "\n")
	for _, number := range numbers {
		if len(number) > 0 {
			// インターフェース名の最後に":"が付いている場合，それを取り除く
			if strings.HasSuffix(number, ":") {
				number = number[:len(number)-1]
			}
			n, err := strconv.Atoi(number)
			if err == nil && n > max {
				max = n
			}
		}
	}
	return max, nil
}

func cmdAdd(args *skel.CmdArgs) error {

	//podマニフェストから送信されてきたネットワークトポロジを取得
	input := []byte(args.StdinData)

	//podマニフェストから送信されてきたネットワークトポロジを構造体に格納
	var ev CNIData
	json.Unmarshal(input, &ev)

	i := 0
	interfaces := ev.InterfaceSet.Items

	//ネットワークの接続先を一時的に管理するjsonファイルの作成
	// 空のスライスを作成
	data := []TempNetwork{}

	//podの各インタフェースに対して処理を実行
	for ; i < len(interfaces); i++ {

		//自身の機器が稼働するノード名と接続先機器が稼働するノード名を格納する変数
		var selfNode []byte
		var targetNode []byte

		//自身の機器が稼働するノード名を取得
		for {
			selfNode, _ = exec.Command("ssh", "-i", "/home/ubuntu/.ssh/id_ed25519", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "ubuntu@192.168.0.224", "kubectl", "get", "pods", "-o", "wide", "|", "grep", "-w", ev.PodName, "|", "awk", "'{print $7}'").Output()
			if string(selfNode) != "" {
				break
			}
			//取得できなかった場合1秒待った後再取得を試みる
			time.Sleep(time.Millisecond * 1000)
		}

		//接続先の機器が稼働するノード名を取得
		for {
			targetNode, _ = exec.Command("ssh", "-i", "/home/ubuntu/.ssh/id_ed25519", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "ubuntu@192.168.0.224", "kubectl", "get", "pods", "-o", "wide", "|", "grep", "-w", interfaces[i].TargetPodName, "|", "awk", "'{print $7}'").Output()
			if string(targetNode) != "" {
				break
			}
			//取得できなかった場合1秒待った後再取得を試みる
			time.Sleep(time.Millisecond * 1000)
		}

		//改行コードを取り除く
		selfNode = []byte(strings.TrimRight(string(selfNode), "\n"))
		ioutil.WriteFile("selfNode.txt", []byte(selfNode), 0644)
		targetNode = []byte(strings.TrimRight(string(targetNode), "\n"))
		ioutil.WriteFile("targetNode.txt", []byte(targetNode), 0644)

		// 新しい要素を追加
		newElement := TempNetwork{InterfaceName: "net" + strconv.Itoa(i+1), TargetNode: string(targetNode)}
		data = append(data, newElement)

		//仮想機器が稼働するノードに応じた処理を実行
		if string(selfNode) == string(targetNode) { //接続する2つの機器が同じノード上で稼働する場合
			//接続処理
			exec.Command("sudo", "ip", "link", "add", "veth-"+ev.PodName+"-"+interfaces[i].Name, "type", "veth", "peer", "name", "veth-"+interfaces[i].TargetPodName+"-"+interfaces[i].TargetPodNIC).Output()
			exec.Command("sudo", "brctl", "addif", ev.PodName+"-"+interfaces[i].Name, "veth-"+ev.PodName+"-"+interfaces[i].Name).Output()
			exec.Command("sudo", "ip", "link", "set", "veth-"+ev.PodName+"-"+interfaces[i].Name, "up").Output()
		} else { //接続する2つの機器が別のノード上で稼働する場合
			//自身の機器が稼働するノードのIPアドレスを取得
			selfNodeIP, _ := exec.Command("ssh", "-i", "/home/ubuntu/.ssh/id_ed25519", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "ubuntu@192.168.0.224", "kubectl", "get", "node", "-o", "wide", "|", "grep", "-w", string(selfNode), "|", "awk", "'{print $6}'").Output()
			//接続先の機器が稼働するノードのIPアドレスを取得
			targetNodeIP, _ := exec.Command("ssh", "-i", "/home/ubuntu/.ssh/id_ed25519", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "ubuntu@192.168.0.224", "kubectl", "get", "node", "-o", "wide", "|", "grep", "-w", string(targetNode), "|", "awk", "'{print $6}'").Output()

			//改行コードを取り除く
			selfNodeIP = []byte(strings.TrimRight(string(selfNodeIP), "\n"))
			targetNodeIP = []byte(strings.TrimRight(string(targetNodeIP), "\n"))

			//接続処理
			exec.Command("sudo", "ip", "l2tp", "add", "tunnel", "remote", string(targetNodeIP), "local", string(selfNodeIP), "tunnel_id", interfaces[i].SelfTunnelID, "peer_tunnel_id", interfaces[i].TargetTunnelID, "encap", "ip").Output()
			exec.Command("sudo", "ip", "l2tp", "add", "session", "tunnel_id", interfaces[i].SelfTunnelID, "session_id", interfaces[i].SessionID, "peer_session_id", interfaces[i].SessionID).Output()

			//正しいl2tpeth番号を取得
			l2tpEthNum, err := calcL2tpEth()
			if err != nil {
				fmt.Println("Error: ", err)
			}

			//接続処理
			exec.Command("sudo", "ip", "l", "set", "l2tpeth"+strconv.Itoa(l2tpEthNum), "up").Output()
			exec.Command("sudo", "ip", "a", "add", "10.1.0.44", "peer", "10.1.0.45", "dev", "l2tpeth"+strconv.Itoa(l2tpEthNum)).Output()
			exec.Command("sudo", "ip", "l", "set", "dev", "l2tpeth"+strconv.Itoa(l2tpEthNum), "master", ev.PodName+"-"+interfaces[i].Name).Output()
		}
	}

	//接続に用いたノード情報を一時的にサーバに格納．pod削除時に用いる．
	outputJson, err := json.Marshal(data)
	ioutil.WriteFile(ev.PodName+"-jsonData.json", []byte(outputJson), 0644)

	conf, err := parseNetConf(args.StdinData)
	if err != nil {
		return err
	}

	var v4Addr, v6Addr *net.IPNet

	args.IfName = "lo" // ignore config, this only works for loopback
	err = ns.WithNetNSPath(args.Netns, func(_ ns.NetNS) error {
		link, err := netlink.LinkByName(args.IfName)
		if err != nil {
			return err // not tested
		}

		err = netlink.LinkSetUp(link)
		if err != nil {
			return err // not tested
		}

		v4Addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
		if err != nil {
			return err // not tested
		}
		if len(v4Addrs) != 0 {
			v4Addr = v4Addrs[0].IPNet
			// sanity check that this is a loopback address
			for _, addr := range v4Addrs {
				if !addr.IP.IsLoopback() {
					return fmt.Errorf("loopback interface found with non-loopback address %q", addr.IP)
				}
			}
		}

		v6Addrs, err := netlink.AddrList(link, netlink.FAMILY_V6)
		if err != nil {
			return err // not tested
		}
		if len(v6Addrs) != 0 {
			v6Addr = v6Addrs[0].IPNet
			// sanity check that this is a loopback address
			for _, addr := range v6Addrs {
				if !addr.IP.IsLoopback() {
					return fmt.Errorf("loopback interface found with non-loopback address %q", addr.IP)
				}
			}
		}

		return nil
	})
	if err != nil {
		return err // not tested
	}

	var result types.Result
	if conf.PrevResult != nil {
		// If loopback has previous result which passes from previous CNI plugin,
		// loopback should pass it transparently
		result = conf.PrevResult
	} else {
		r := &current.Result{
			CNIVersion: conf.CNIVersion,
			Interfaces: []*current.Interface{
				&current.Interface{
					Name:    args.IfName,
					Mac:     "00:00:00:00:00:00",
					Sandbox: args.Netns,
				},
			},
		}

		if v4Addr != nil {
			r.IPs = append(r.IPs, &current.IPConfig{
				Interface: current.Int(0),
				Address:   *v4Addr,
			})
		}

		if v6Addr != nil {
			r.IPs = append(r.IPs, &current.IPConfig{
				Interface: current.Int(0),
				Address:   *v6Addr,
			})
		}

		result = r
	}

	return types.PrintResult(result, conf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	//podマニフェストから送信されてきたネットワークトポロジを取得
	input := []byte(args.StdinData)

	//podマニフェストから送信されてきたネットワークトポロジを構造体に格納
	var ev CNIData
	json.Unmarshal(input, &ev)

	//pod作成時に保存したノード情報を記述したファイルの読み込み
	filePath := ev.PodName + "-jsonData.json"
	file, _ := ioutil.ReadFile(filePath)
	var ev2 []TempNetwork
	json.Unmarshal(file, &ev2)

	i := 0
	interfaces := ev.InterfaceSet.Items

	//podの各インタフェースに対して処理を実行
	for ; i < len(interfaces); i++ {

		//自身の機器が稼働するノード名と接続先機器が稼働するノード名を格納する変数
		var selfNode []byte

		//自身の機器が稼働するノード名を取得
		for {
			selfNode, _ = exec.Command("ssh", "-i", "/home/ubuntu/.ssh/id_ed25519", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "ubuntu@192.168.0.224", "kubectl", "get", "pods", "-o", "wide", "|", "grep", "-w", ev.PodName, "|", "awk", "'{print $7}'").Output()
			if string(selfNode) != "" {
				break
			}
			//取得できなかった場合1秒待った後再取得を試みる
			time.Sleep(time.Millisecond * 1000)
		}

		//改行コードを取り除く
		selfNode = []byte(strings.TrimRight(string(selfNode), "\n"))

		//仮想機器が稼働するノードに応じた処理を実行
		if string(selfNode) == string(ev2[i].TargetNode) { //接続する2つの機器が同じノード上で稼働する場合
			//後片付け処理
			exec.Command("sudo", "ip", "link", "set", "veth-"+ev.PodName+"-"+interfaces[i].Name, "down").Output()
			exec.Command("sudo", "ip", "link", "del", "veth-"+ev.PodName+"-"+interfaces[i].Name, "type", "veth", "peer", "name", "veth-"+interfaces[i].TargetPodName+"-"+interfaces[i].TargetPodNIC).Output()
		} else { //接続する2つの機器が別のノード上で稼働する場合
			//後片付け処理
			exec.Command("sudo", "ip", "l2tp", "del", "session", "tunnel_id", interfaces[i].SelfTunnelID, "session_id", interfaces[i].SessionID).Output()
		}
	}
	// ファイルを削除する
	os.Remove(filePath)

	if args.Netns == "" {
		return nil
	}
	args.IfName = "lo" // ignore config, this only works for loopback
	err := ns.WithNetNSPath(args.Netns, func(ns.NetNS) error {
		link, err := netlink.LinkByName(args.IfName)
		if err != nil {
			return err // not tested
		}

		err = netlink.LinkSetDown(link)
		if err != nil {
			return err // not tested
		}

		return nil
	})
	if err != nil {
		//  if NetNs is passed down by the Cloud Orchestration Engine, or if it called multiple times
		// so don't return an error if the device is already removed.
		// https://github.com/kubernetes/kubernetes/issues/43014#issuecomment-287164444
		_, ok := err.(ns.NSPathNotExistErr)
		if ok {
			return nil
		}
		return err
	}

	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, bv.BuildString("loopback"))
}

func cmdCheck(args *skel.CmdArgs) error {
	args.IfName = "lo" // ignore config, this only works for loopback

	return ns.WithNetNSPath(args.Netns, func(_ ns.NetNS) error {
		link, err := netlink.LinkByName(args.IfName)
		if err != nil {
			return err
		}

		if link.Attrs().Flags&net.FlagUp != net.FlagUp {
			return errors.New("loopback interface is down")
		}

		return nil
	})
}
