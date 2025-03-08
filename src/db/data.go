package db

import (
	"io"
	"os"

	"github.com/koyo-os/tokino/src/models"

	"github.com/bytedance/sonic"
)

type Data struct{
	File *os.File
}

func New() (*Data, error) {
	file, err := os.Open("storage/blocks.json")
	if err != nil{
		return nil,err
	}

	return &Data{
		File: file,
	}, nil
}

func (d *Data) Add(b *models.Block) error {
	data,err := io.ReadAll(d.File)
	if err != nil{
		return err
	}

	if len(data) < 1 {
		newdata,err := sonic.Marshal(b)
		if err != nil{
			return err
		}

		_, err = io.WriteString(d.File, string(newdata))
	}

	var blocks []models.Block

	if err = sonic.Unmarshal(data, &blocks);err != nil{
		return err
	}

	blocks = append(blocks, *b)

	newdata,err := sonic.Marshal(blocks)
	if err != nil{
		return err
	}

	_, err = io.WriteString(d.File, string(newdata))
	return err
}

func (d *Data) GetAll() ([]models.Block, error) {
	var blocks []models.Block

	data, err := io.ReadAll(d.File)
	if err != nil{
		return nil,err
	}

	if err = sonic.Unmarshal(data, &blocks);err != nil{
		return nil,err
	}

	return blocks, nil
}

func (d *Data) GetLast() (*models.Block, error) {
	blocks, err := d.GetAll()
	if err != nil{
		return nil,err
	}
	return &blocks[len(blocks)-1], nil
}