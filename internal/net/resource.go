package net

import (
	"log"
	"syscall"
)

func configureRLimit() {
	var rLimit syscall.Rlimit

	/// RLIMIT_NOFILE은 소켓 파일 최대 연결 수 자원
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Println("[ERROR] System Call Get RLIMIT : ", err)
	}

	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Println("[ERROR] System Call Set RLIMIT : ", err)

	}
}
