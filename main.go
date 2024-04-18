package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type MemoItem struct {
	ID   int    `json:"id"`   // id
	Text string `json:"text"` // 内容
	Time string `json:"time"` // 时间
}

func main() {
	app := cli.NewApp()
	app.Name = "memo"

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a new memo item",
			Action:  addMemo,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all memo items",
			Action:  listMemo,
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete a memo item",
			Action:  delMemo,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("启动程序失败!")
		log.Fatal(err)
	}

}

// 追加单条内容
func addMemo(c *cli.Context) error {
	// 读取命令 输入 内容
	text := c.Args().First()

	// 读取 记事本 现有内容
	items, err := readMemo()
	if err != nil {
		fmt.Println("读取内容失败")
		return err
	}
	// 实例化结构体
	item := MemoItem{
		ID:   len(items) + 1,
		Text: text,
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}
	// 追加内容
	items = append(items, item)
	// 将所有内容 重新写入 记事本
	if err = writeMemo(items); err != nil {
		fmt.Println("追加内容后写入记事本出错")
		return err
	}

	return nil
}

// 删除 某条 记录
func delMemo(c *cli.Context) error {
	// 获取命令 后 对应的 id值
	idd := c.Args().First()
	id, err := strconv.Atoi(idd)
	if err != nil {
		fmt.Println("整数转换失败！")
		return err
	}

	// 读取 记事本 目前所有内容
	items, err := readMemo()
	if err != nil {
		return err
	}
	// 在原来 的 基础上 进行 切除
	items = append(items[:id-1], items[id:]...)

	// 删除后，格式化序号
	if err != nil {
		return err
	}
	var redoItems []MemoItem

	for i, text := range items {
		item := MemoItem{
			ID:   i + 1,
			Text: text.Text,
			Time: text.Time,
		}
		redoItems = append(redoItems, item)
	}
	err = writeMemo(redoItems)
	if err != nil {
		fmt.Println("格式化重新写入操作失败！")
		return err
	}
	return nil

}

// 查看 所有内容
func listMemo(c *cli.Context) error {
	// 读取 记事本 现有内容
	items, err := readMemo()
	if err != nil {
		fmt.Println("读取记事本内容失败")
		return err
	}
	// 遍历 所有内容 格式化
	for _, item := range items {
		fmt.Printf("%d: %s  %s\n", item.ID, item.Text, item.Time)
	}
	return nil
}

// 所有内容 写入到 备忘录
func writeMemo(items []MemoItem) error {
	fileData, err := json.MarshalIndent(items, "", "     ")
	if err != nil {
		return err
	}
	// ioutil.WriteFile  如果没有目标文件，进行创建
	err = ioutil.WriteFile("data.json", fileData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// 获取 所有内容
func readMemo() ([]MemoItem, error) {
	fileData, err := ioutil.ReadFile("data.json")

	if err != nil {
		/*
			如果此时 没有 data.json文件，那就没办法读取内容，会导致报错
			所以 返回一个 空集合
			让后续写入的时候，进行创建文件
		*/
		if os.IsNotExist(err) {
			return []MemoItem{}, nil
		}
		return nil, err
	}

	var items []MemoItem
	// 将读取到的内容 填充到 items结构体中
	err = json.Unmarshal(fileData, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
