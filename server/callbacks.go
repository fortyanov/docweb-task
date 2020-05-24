package server

import "docweb-task/log"

func uploadPreCallback() {
	log.Info("called uploadPreCallback")
}

func uploadPostCallback() {
	log.Info("called uploadPostCallback")
}

func downloadPreCallback() {
	log.Info("called downloadPreCallback")
}

func downloadPostCallback() {
	log.Info("called downloadPostCallback")
}

func deletePreCallback() {
	log.Info("called deletePreCallback")
}

func deletePostCallback() {
	log.Info("called deletePostCallback")
}
