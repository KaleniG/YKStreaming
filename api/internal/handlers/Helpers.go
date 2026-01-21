package handlers

import (
	"log"
	"net/http"

	"ykstreaming_api/internal/helpers"
)

type recordingAction int

const (
	StartRecording recordingAction = iota
	StopRecording
)

func requestStreamRecordingAction(key string, action recordingAction) error {
	controlURL := helpers.GetEnvDir("RTMP_CONTROL_URL")
	VODRecorderName := helpers.GetEnvDir("LIVE_APP_VOD_RECORDER_NAME")
	var actionString string
	if action == StartRecording {
		actionString = "start"
	} else if action == StopRecording {
		actionString = "stop"
	} else {
		log.Panic("unknown enum used in 'requestStreamRecordingAction'")
	}

	resp, err := http.Post(controlURL+"/record/"+actionString+"?app=live&name="+key+"&rec="+VODRecorderName, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func requestStreamStop(key string) error {
	controlURL := helpers.GetEnvDir("RTMP_CONTROL_URL")
	resp, err := http.Post(controlURL+"/drop/publisher?app=live&name="+key, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
