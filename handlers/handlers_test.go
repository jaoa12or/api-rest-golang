package handlers

import (
	"challenge-backend/models"
	"testing"
)

func TestGetPageInfo(t *testing.T) {
	expected := models.ScrapingResponse{Icon: "", Title: "Google"}
	got, err := GetPageInfo("google.com")
	if err != nil {
		t.Error("Error: ", err.Error())
	}
	if got != expected {
		t.Error("GetPageInfo(google.com) = ", got, "want: ", models.ScrapingResponse{Icon: "", Title: "Google"})
	}
}

func TestGetPageInfoFail(t *testing.T) {
	expected := models.ScrapingResponse{}
	got, err := GetPageInfo("")
	if err == nil {
		t.Error("Error: ", err.Error())
	}
	if got != expected {
		t.Error("GetPageInfo(google.com) = ", got, "want: ", models.ScrapingResponse{Icon: "", Title: "Google"})
	}
}

func TestGetOwnerData(t *testing.T) {
	expected := "Google LLC"
	got := GetOwnerData("216.58.194.174", "OrgName")
	if got != expected {
		t.Error("GetOwnerData(216.58.194.174, OrgNam) = ", got, "want: Google LLC")
	}
}