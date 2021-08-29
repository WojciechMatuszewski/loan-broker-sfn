package main_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestReplace(t *testing.T) {
	str := "BrokerAPIDeployment$TIMESTAMP$"

	nowStr := strconv.Itoa(int(time.Now().Unix()))

	replacedStr := strings.ReplaceAll(str, "$TIMESTAMP$", nowStr)

	fmt.Println(replacedStr)
}
