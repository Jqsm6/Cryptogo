package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"Cryptogo/internal/models"
)

func infoETH(pirp *models.PaymentInfoResponse) (bool, error) {
	var ethModel *models.ETHBalance
	url := fmt.Sprintf("https://api.ethplorer.io/getAddressInfo/%s?apiKey=freekey", pirp.ToAddress)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, errors.New("api did not respond with a 200 code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(body, &ethModel)
	if err != nil {
		return false, err
	}

	balance := strconv.Itoa(ethModel.ETH.Balance)

	if balance == pirp.Amount {
		return true, nil
	}

	return false, nil
}
