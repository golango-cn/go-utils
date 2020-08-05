package ip

import (
	"fmt"
	"os"
	"testing"
)

func TestIP(t *testing.T)  {

	pwd, _ := os.Getwd()
	ipdata := pwd + "/qqwry.dat"
	fmt.Println(ipdata)

	data, _ := NewIPData(ipdata)

	ips := []string{"163.177.65.160", "114.114.114.114"}

	qwry := NewQQwry(data)
	for _, v := range ips {
		a := qwry.Find(v)
		fmt.Println(v, a)
	}

}
